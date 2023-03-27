package commands

import (
	"bufio"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/signal"
	models "pragmatic-zac/goModeS/models"
	"pragmatic-zac/goModeS/streaming"
	"sync"
	"syscall"
	"time"
)

// probably move these so that other stuff can use them

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
	connectCmd.Flags().StringVarP(&mode, "lat", "l", "", "receiver latitude")
	connectCmd.Flags().StringVarP(&mode, "lon", "o", "", "receiver longitude")

	connectCmd.MarkFlagRequired("address")
	connectCmd.MarkFlagRequired("mode")
	connectCmd.MarkFlagRequired("latitude")
	connectCmd.MarkFlagRequired("longitude")

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

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// TODO: render the table
			for _, f := range flightsState {
				fmt.Printf("ICAO: %s, Call: %s, Alt: %d \n", f.Icao, f.Callsign, f.Altitude)
			}

			// do this once a second
			time.Sleep(1 * time.Second)
		}
	}
}
