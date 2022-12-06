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
