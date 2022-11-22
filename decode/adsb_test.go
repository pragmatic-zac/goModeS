package decode

import "testing"

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
