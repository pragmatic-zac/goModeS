package decode

import (
	"testing"
	"time"
)

func TestCategory(t *testing.T) {
	msg := "8D406B902015A678D4D220AA4BDA"

	actual, _ := Category(msg)

	want := 0

	if actual != int64(want) {
		t.Fatalf("Category incorrect, wanted %v got %v", want, actual)
	}
}

func TestCallsign(t *testing.T) {
	msg := "8D4840D6202CC371C32CE0576098"

	actual, _ := Callsign(msg)

	want := "KLM1023 "

	if actual != want {
		t.Fatalf("Category incorrect, wanted %v got %v", want, actual)
	}
}

func TestAirbornePosition(t *testing.T) {
	msg0 := "8D40621D58C382D690C8AC2863A7"
	msg1 := "8D40621D58C386435CC412692AD6"
	t0 := time.Unix(int64(1457996402), 0)
	t1 := time.Unix(int64(1457996400), 0)

	actualLat, actualLon, _ := AirbornePosition(msg0, msg1, t0, t1)

	wantedLat := 52.2572
	wantedLon := 3.91937

	if actualLat != wantedLat {
		t.Fatalf("Latitude incorrect, wanted %v got %v", wantedLat, actualLat)
	}

	if actualLon != wantedLon {
		t.Fatalf("Longitude incorrect, wanted %v got %v", wantedLon, actualLon)
	}
}

func TestSurfacePosition(t *testing.T) {
	msg0 := "8C4841753AAB238733C8CD4020B1"
	msg1 := "8C4841753A8A35323FAEBDAC702D"
	t0 := time.Unix(int64(1457996410), 0)
	t1 := time.Unix(int64(1457996412), 0)
	latRef := 51.990
	lonRef := 4.375

	actualLat, actualLon, _ := SurfacePosition(msg0, msg1, t0, t1, latRef, lonRef)

	wantedLat := 52.32061
	wantedLon := 4.73473

	if actualLat != wantedLat {
		t.Fatalf("Latitude incorrect, wanted %v got %v", wantedLat, actualLat)
	}

	if actualLon != wantedLon {
		t.Fatalf("Longitude incorrect, wanted %v got %v", wantedLon, actualLon)
	}
}

func TestSurfacePositionWithRef(t *testing.T) {
	msg0 := "8C4841753A9A153237AEF0F275BE"
	latRef := 51.990
	lonRef := 4.375

	actualLat, actualLon, _ := SurfacePositionWithRef(msg0, latRef, lonRef)

	wantedLat := 52.320561
	wantedLon := 4.735735

	if actualLat != wantedLat {
		t.Fatalf("Latitude incorrect, wanted %v got %v", wantedLat, actualLat)
	}

	if actualLon != wantedLon {
		t.Fatalf("Longitude incorrect, wanted %v got %v", wantedLon, actualLon)
	}
}

func TestSurfaceVelocity(t *testing.T) {
	msg := "8C4841753A9A153237AEF0F275BE"

	spd, trk, vertRate, spdType, _ := SurfaceVelocity(msg)

	wantedSpd := 17.0
	wantedTrk := 92.8
	wantedVertRate := 0
	wantedSpdType := "GS"

	if spd != wantedSpd {
		t.Fatalf("Speed incorrect, wanted %v got %v", wantedSpd, spd)
	}

	if trk != wantedTrk {
		t.Fatalf("Track incorrect, wanted %v got %v", wantedTrk, trk)
	}

	if vertRate != int32(wantedVertRate) {
		t.Fatalf("Vertical rate incorrect, wanted %v got %v", wantedVertRate, vertRate)
	}

	if spdType != wantedSpdType {
		t.Fatalf("Speed type incorrect, wanted %v got %v", wantedSpdType, spdType)
	}
}

func TestAirborneVelocityGS(t *testing.T) {
	msg := "8D485020994409940838175B284F"

	vel, _ := AirborneVelocity(msg)

	wantedSpd := int64(159)
	wantedTrk := 182.88
	wantedVertRate := int64(-832)
	wantedSpdType := "GS"

	if vel.speed != wantedSpd {
		t.Fatalf("Speed incorrect, wanted %v got %v", wantedSpd, vel.speed)
	}

	if vel.angle != wantedTrk {
		t.Fatalf("Track incorrect, wanted %v got %v", wantedTrk, vel.angle)
	}

	if vel.vertRate != wantedVertRate {
		t.Fatalf("Vertical rate incorrect, wanted %v got %v", wantedVertRate, vel.vertRate)
	}

	if vel.speedType != wantedSpdType {
		t.Fatalf("Speed type incorrect, wanted %v got %v", wantedSpdType, vel.speedType)
	}
}

func TestAirborneVelocityTAS(t *testing.T) {
	msg := "8DA05F219B06B6AF189400CBC33F"

	vel, _ := AirborneVelocity(msg)

	wantedSpd := int64(375)
	wantedTrk := 243.98
	wantedVertRate := int64(-2304)
	wantedSpdType := "TAS"

	if vel.speed != wantedSpd {
		t.Fatalf("Speed incorrect, wanted %v got %v", wantedSpd, vel.speed)
	}

	if vel.angle != wantedTrk {
		t.Fatalf("Track incorrect, wanted %v got %v", wantedTrk, vel.angle)
	}

	if vel.vertRate != wantedVertRate {
		t.Fatalf("Vertical rate incorrect, wanted %v got %v", wantedVertRate, vel.vertRate)
	}

	if vel.speedType != wantedSpdType {
		t.Fatalf("Speed type incorrect, wanted %v got %v", wantedSpdType, vel.speedType)
	}
}

func BenchmarkAirbornePosition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msg0 := "8D40621D58C382D690C8AC2863A7"
		msg1 := "8D40621D58C386435CC412692AD6"
		t0 := time.Unix(int64(1457996402), 0)
		t1 := time.Unix(int64(1457996400), 0)

		AirbornePosition(msg0, msg1, t0, t1)
	}
}

func BenchmarkSurfacePosition(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msg0 := "8C4841753AAB238733C8CD4020B1"
		msg1 := "8C4841753A8A35323FAEBDAC702D"
		t0 := time.Unix(int64(1457996410), 0)
		t1 := time.Unix(int64(1457996412), 0)
		latRef := 51.990
		lonRef := 4.375

		SurfacePosition(msg0, msg1, t0, t1, latRef, lonRef)
	}
}

func BenchmarkAirborneVelocity(b *testing.B) {
	for i := 0; i < b.N; i++ {
		msg := "8DA05F219B06B6AF189400CBC33F"

		AirborneVelocity(msg)
	}
}
