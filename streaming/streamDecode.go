package streaming

import (
	"fmt"
	"pragmatic-zac/goModeS/decode"
	"pragmatic-zac/goModeS/internal"
)

// TODO: also going to need a pointer to a data structure that stores flights
func DecodeAdsB(msg string) {
	cleanedMsg := internal.CleanMessage(msg)

	icao, _ := decode.Icao(cleanedMsg)
	println(icao)

	tc, _ := decode.Typecode(cleanedMsg)

	println("message typecode: ", tc)

	// set the time, will need to do this on the time that each msg is received
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

		vel, _ := decode.SurfaceVelocity(cleanedMsg)
		fmt.Printf("Surface velocity: %f\n", vel.Speed)
	}
}
