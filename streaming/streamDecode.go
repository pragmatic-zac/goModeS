package streaming

import (
	"github.com/pragmatic-zac/goModeS/decode"
	"github.com/pragmatic-zac/goModeS/internal"
	models "github.com/pragmatic-zac/goModeS/models"
	"time"
)

func DecodeAdsB(msg string, flightsState map[string]models.Flight, latRef float64, lonRef float64) {
	cleanedMsg := internal.CleanMessage(msg)

	icao, _ := decode.Icao(cleanedMsg)
	tc, _ := decode.Typecode(cleanedMsg)
	timestamp := time.Now()

	f := flightsState[icao]
	f.Icao = icao
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
		if oddEven == 0 {
			f.EvenMessage = cleanedMsg
			f.EvenMessageTime = timestamp
		} else {
			f.OddMessage = cleanedMsg
			f.OddMessageTime = timestamp
		}

		// for now use position with ref, later on add support for position with odd/even message pair

		if tc >= 5 && tc <= 8 {
			// surface position
			pos, _ := decode.SurfacePositionWithRef(cleanedMsg, latRef, lonRef)
			f.Position = pos

			alt, _ := decode.Altitude(cleanedMsg)
			if alt != 0 {
				f.Altitude = alt
			}
		} else {
			// airborne position
			pos, _ := decode.AirbornePositionWithRef(cleanedMsg, latRef, lonRef)
			f.Position = pos

			alt, _ := decode.Altitude(cleanedMsg)
			if alt != 0 {
				f.Altitude = alt
			}
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
