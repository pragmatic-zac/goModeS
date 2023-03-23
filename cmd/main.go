package main

import (
	"fmt"
	"pragmatic-zac/goModeS/decode"
)

// 8DACD0CF990C2D3250041196585C - 19
// 8DA41C50E11A3400000000160B30 - 28
// 8DA41C5099103297782403FCA1FC - 19
// 8DA3FB6F990CAD2CB804120E1E2A - 19
// 8DAC0BFB99953D88109412D0E7BB - 19
// 8DA57788221102B4D71820B37C63 - 4
// 8DA57788990C7A10A0C00A1BB1FF - 19
// 8DA57788990C7910C0C00B2EDEDF - 19
// 8DA3FB6F589B909697401DC40C0C - 11
// 8DA3FB6F589B942D65BC1825EA07 - 11
// 8DA1701BEA44785EE75C08240817 - 29

// typecodes
// 1- 4		: identification
// 5 - 8	: surface position
// 9 - 18	: airborne position with baro altitude
// 19 		: airborne velocity
// 20 - 22	: airborne position with gnss height
// 29		: target state and status info

// user needs to pass in lat/lon for reference position

func main() {
	msg := "8DA3FB6F990CAD2CB804120E1E2A"

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

	println("done")
}
