package streaming

import (
	"fmt"
	"pragmatic-zac/goModeS/decode"
	"pragmatic-zac/goModeS/internal"
	models "pragmatic-zac/goModeS/models"
)

// TODO: does this need its own state to keep other recent messages? for example, previous lat/long pairs
func DecodeAdsB(msg string, flightsState map[string]models.Flight) {
	cleanedMsg := internal.CleanMessage(msg)

	icao, _ := decode.Icao(cleanedMsg)
	println(icao)
	flightsState[icao] = models.Flight{Icao: icao, Altitude: "fixme"}

	// place this into flights for testing
	for _, flight := range flightsState {
		if flight.Icao == icao {
			println("this is already in state")
		}
	}

	tc, _ := decode.Typecode(cleanedMsg)

	println("message typecode: ", tc)

	// TODO: set the time, will need to do this on the time that each msg is received
	// timestamp := time.Now()
	// println("current time is ", timestamp.)

	if tc >= 1 && tc <= 4 {
		// identification
		ident, _ := decode.Callsign(cleanedMsg)
		println(ident)
	}

	if tc >= 5 && tc <= 8 || tc == 19 {
		// velocity (either type)
		vel, _ := decode.CombinedVelocity(cleanedMsg)
		println(int(vel.Speed))
	}

	if tc >= 5 && tc <= 18 {
		oddEven := decode.OddEvenFlag(cleanedMsg)
		println("odd or even -> ", oddEven)

		// for now use position with ref, but will need a way to determine position with ref or position with odd/even message pair

		if tc >= 5 && tc <= 8 {
			// surface position
			pos, _ := decode.SurfacePositionWithRef(cleanedMsg, 36.04863, -86.95218)
			fmt.Printf("Lat: %f", pos.Latitude)
			fmt.Printf("Lon: %f", pos.Longitude)

			alt, _ := decode.Altitude(cleanedMsg)
			fmt.Printf("altitude: %d\n", alt)
		} else {
			// airborne position
			pos, _ := decode.AirbornePositionWithRef(cleanedMsg, 36.04863, -86.95218)
			fmt.Printf("Lat: %f\n", pos.Latitude)
			fmt.Printf("Lon: %f\n", pos.Longitude)

			// not sure if belongs here
			alt, _ := decode.Altitude(cleanedMsg)
			fmt.Printf("altitude: %d\n", alt)
		}
	}

	// TODO: add all this data to the flight map
}
