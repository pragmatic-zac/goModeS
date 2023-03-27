package streaming

import (
	"pragmatic-zac/goModeS/decode"
	"pragmatic-zac/goModeS/internal"
	models "pragmatic-zac/goModeS/models"
	"time"
)

func DecodeAdsB(msg string, flightsState map[string]models.Flight) {
	cleanedMsg := internal.CleanMessage(msg)

	icao, _ := decode.Icao(cleanedMsg)
	tc, _ := decode.Typecode(cleanedMsg)
	timestamp := time.Now()

	f := flightsState[icao]
	f.LastSeen = timestamp

	if tc >= 1 && tc <= 4 {
		// identification
		ident, _ := decode.Callsign(cleanedMsg)
		f.Callsign = ident
	}

	if tc >= 5 && tc <= 8 || tc == 19 {
		// velocity (either type)
		vel, _ := decode.CombinedVelocity(cleanedMsg)
		f.Velocity = vel
	}

	if tc >= 5 && tc <= 18 {
		oddEven := decode.OddEvenFlag(cleanedMsg)
		println("odd or even -> ", oddEven)

		// for now use position with ref, but will need a way to determine position with ref or position with odd/even message pair

		if tc >= 5 && tc <= 8 {
			// surface position
			pos, _ := decode.SurfacePositionWithRef(cleanedMsg, 36.04863, -86.95218)
			f.Position = pos

			alt, _ := decode.Altitude(cleanedMsg)
			f.Altitude = alt
		} else {
			// airborne position
			pos, _ := decode.AirbornePositionWithRef(cleanedMsg, 36.04863, -86.95218)
			f.Position = pos

			alt, _ := decode.Altitude(cleanedMsg)
			f.Altitude = alt
		}
	}

	// update the flight in the cache
	flightsState[icao] = f

	expireCache(flightsState, timestamp)

}

func expireCache(flightsState map[string]models.Flight, t time.Time) {
	for _, flight := range flightsState {
		diff := t.Sub(flight.LastSeen)
		if diff > 60*time.Second {
			delete(flightsState, flight.Icao)
		}
	}
}
