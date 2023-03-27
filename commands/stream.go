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
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Display table of aircraft tracked by receiver running on provided port",
	Long:  `Connects to a receiver running on a provided port. Decodes messages and displays tracked aircraft in a table.`,
	Run: func(cmd *cobra.Command, args []string) {
		println("connect")
		println(address)
		println(mode)

		// some parent level state
		trackedFlights := make(map[string]models.Flight)

		// network stuff
		tcpAddr, err := net.ResolveTCPAddr("tcp", "192.168.1.190:30002") // replace this with address
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
		fmt.Println("Shutting down...")

		cancel()
		wg.Wait()

		fmt.Println("Successfully shut down.")
	},
}

func init() {
	connectCmd.Flags().StringVarP(&address, "address", "a", "", "address to connect to (include port)")
	connectCmd.Flags().StringVarP(&mode, "mode", "m", "", "mode of source (currently only raw supported)")

	connectCmd.MarkFlagRequired("address")
	connectCmd.MarkFlagRequired("mode")

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
			println("TCP SHUTTING DOWN")
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
			println("MESSAGE PROCESSOR SHUTTING DOWN")
			return
		case msg := <-msgChan:
			// ignore other messages for now
			if len(msg) == 31 {
				streaming.DecodeAdsB(msg, flightsState)
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
			println("RENDER LOOP SHUTTING DOWN")
			return
		default:
			// println("-- render --")
			fmt.Printf("Planes in cache: %d\n", len(flightsState))

			// TODO: render the table

			// do this once a second
			time.Sleep(1 * time.Second)
		}
	}
}
