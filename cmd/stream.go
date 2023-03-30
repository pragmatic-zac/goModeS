package main

import (
	"bufio"
	"context"
	"fmt"
	tm "github.com/buger/goterm"
	models "github.com/pragmatic-zac/goModeS/models"
	"github.com/pragmatic-zac/goModeS/streaming"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var address string
var mode string
var latRef float64
var lonRef float64
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Display table of aircraft tracked by receiver running on provided port",
	Long:  `Connects to a receiver running on a provided port. Decodes messages and displays tracked aircraft in a table.`,
	Run: func(cmd *cobra.Command, args []string) {
		if mode != "raw" {
			println("only raw format is currently supported!")
		}

		trackedFlights := make(map[string]models.Flight)

		// network stuff
		tcpAddr, err := net.ResolveTCPAddr("tcp", address)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer conn.Close()

		// set up channels and such
		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup

		msgChan := make(chan string)

		go handleConnection(ctx, conn, msgChan, &wg)
		go processMessages(ctx, msgChan, &wg, trackedFlights)
		go renderLoop(ctx, &wg, trackedFlights)

		// Wait for SIGINT or SIGTERM to trigger a graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		cancel()
		wg.Wait()

		fmt.Println("Successfully shut down.")
	},
}

func init() {
	connectCmd.Flags().StringVarP(&address, "address", "a", "", "address to connect to (include port)")
	connectCmd.Flags().StringVarP(&mode, "mode", "m", "", "mode of source (currently only raw supported)")
	connectCmd.Flags().Float64VarP(&latRef, "lat", "l", 0, "receiver latitude")
	connectCmd.Flags().Float64VarP(&lonRef, "lon", "o", 0, "receiver longitude")

	connectCmd.MarkFlagRequired("address")
	connectCmd.MarkFlagRequired("mode")
	connectCmd.MarkFlagRequired("lat")
	connectCmd.MarkFlagRequired("lon")

	rootCmd.AddCommand(connectCmd)
}

func handleConnection(ctx context.Context, conn net.Conn, msgChan chan<- string, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	readChan := make(chan string)
	errChan := make(chan error)

	// start a goroutine to read data from the connection
	// bufio blocks if we don't do this, and we never get graceful shutdown
	go func() {
		for {
			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				errChan <- err
				return
			}
			readChan <- msg
		}
	}()

	for {
		select {
		case <-ctx.Done():
			close(msgChan)
			return
		case msg := <-readChan:
			msgChan <- msg
		case err := <-errChan:
			fmt.Println("Error reading message:", err)
			return
		}
	}
}

func processMessages(ctx context.Context, msgChan <-chan string, wg *sync.WaitGroup, flightsState map[string]models.Flight) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgChan:
			// ignore other messages for now
			if len(msg) == 31 {
				streaming.DecodeAdsB(msg, flightsState, latRef, lonRef)
			}
		}
	}
}

func renderLoop(ctx context.Context, wg *sync.WaitGroup, flightsState map[string]models.Flight) {
	wg.Add(1)
	defer wg.Done()

	tm.Clear()

	for {
		select {
		case <-ctx.Done():
			tm.Clear()
			return
		default:
			tm.MoveCursor(1, 1)

			tbl := tm.NewTable(0, 10, 5, ' ', 0)
			fmt.Fprintf(tbl, "ICAO\t Callsign \t Altitude \t Speed \tHeading \t VertRate \t Lat \t Lon \n")

			for _, f := range flightsState {
				fmt.Fprintf(tbl, "%s \t %s \t %d \t %f \t %f \t %d \t %f \t %f \n", f.Icao, f.Callsign, f.Altitude, f.Velocity.Speed, f.Velocity.Angle, f.Velocity.VertRate, f.Position.Latitude, f.Position.Longitude)
			}

			tm.Println(tbl)

			tm.Flush()

			// do this once a second
			time.Sleep(1 * time.Second)
		}
	}
}
