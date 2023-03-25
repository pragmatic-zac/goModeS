package streaming

import (
	"fmt"
	"pragmatic-zac/goModeS/decode"
)

func DecodeAdsB(msg string) {
	icao, _ := decode.Icao(msg)
	println(icao)

	tc, _ := decode.Typecode(msg)

	println("message typecode: ", tc)

	// set the time, will need to do this on the time that each msg is received
	// timestamp := time.Now()
	// println("current time is ", timestamp.)

	if tc >= 1 && tc <= 4 {
		// identification
		ident, _ := decode.Callsign(msg)
		println(ident)
	}

	if tc >= 5 && tc <= 8 || tc == 19 {
		// velocity (either type)
		vel, _ := decode.CombinedVelocity(msg)
		println(int(vel.Speed))
	}

	if tc >= 5 && tc <= 18 {
		oddEven := decode.OddEvenFlag(msg)
		println("odd or even -> ", oddEven)

		// for now use position with ref, but will need a way to determine position with ref or position with odd/even message pair

		if tc >= 5 && tc <= 8 {
			// surface position
			pos, _ := decode.SurfacePositionWithRef(msg, 36.04863, -86.95218)
			fmt.Printf("Lat: %f", pos.Latitude)
			fmt.Printf("Lon: %f", pos.Longitude)

			alt, _ := decode.Altitude(msg)
			fmt.Printf("altitude: %d\n", alt)
		} else {
			// airborne position
			pos, _ := decode.AirbornePositionWithRef(msg, 36.04863, -86.95218)
			fmt.Printf("Lat: %f\n", pos.Latitude)
			fmt.Printf("Lon: %f\n", pos.Longitude)

			// not sure if belongs here
			alt, _ := decode.Altitude(msg)
			fmt.Printf("altitude: %d\n", alt)
		}

		vel, _ := decode.SurfaceVelocity(msg)
		fmt.Printf("Surface velocity: %f\n", vel.Speed)
	}
}
