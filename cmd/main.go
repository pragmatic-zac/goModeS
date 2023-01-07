package main

import (
	"pragmatic-zac/goModeS/decode"
	"time"
)

// 8DACD0CF990C2D3250041196585C
// 8DA41C50E11A3400000000160B30
// 8DA41C5099103297782403FCA1FC

// user needs to pass in lat/lon for reference position

func main() {
	msg := "8DA41C5099103297782403FCA1FC"

	tc, _ := decode.Typecode(msg)

	println("message typecode: ", tc)

	// set the time, will need to do this on the time that each msg is received
	timestamp := time.Now()
	println("current time is ", timestamp)

	if tc >= 1 && tc <= 4 {
		// identification
		ident, _ := decode.Callsign(msg)
		println(ident)
	}

	if tc >= 5 && tc <= 8 {
		vel, _ := decode.SurfaceVelocity(msg)
		println(int(vel.Speed))
	}

	if tc >= 5 && tc <= 18 {
		oddEven := decode.OddEvenFlag(msg)
		println("odd or even -> ", oddEven)

		// for now use position with ref, but will need a way to determine position with ref or position with odd/even message pair

		if tc >= 5 && tc <= 8 {
			// surface position
			//input := decode.PositionInput{}
		} else {
			// airborne position
		}

		vel, _ := decode.SurfaceVelocity(msg)
		println(int(vel.Speed))
	}

	if tc >= 9 && tc <= 18 {
		// airborne position BARO
	}

	if tc == 19 {
		// airborne velocity
		println("airborne velocity")
		vel, _ := decode.AirborneVelocity(msg)
		println(int(vel.Speed))
	}

	if tc >= 20 && tc <= 27 {
		// airborne position GNSS
		println("airborne position GNSS")
	}

	println("done")
}
