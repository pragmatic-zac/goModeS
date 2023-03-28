package models

import (
	"github.com/pragmatic-zac/goModeS/decode"
	"time"
)

type Flight struct {
	Icao            string
	Callsign        string
	Altitude        int
	Position        decode.Position
	Velocity        decode.Velocity
	LastSeen        time.Time
	OddMessage      string
	OddMessageTime  time.Time
	EvenMessage     string
	EvenMessageTime time.Time
}
