package internal

import "testing"

func TestHexToBinary(t *testing.T) {
	msg := "8D4840D6202CC371C32CE0576098"

	actual, err := HexToBinary(msg)
	if err != nil {
		println(err)
	}

	want := "1000110101001000010000001101011000100000001011001100001101110001110000110010110011100000010101110110000010011000"

	if actual != want {
		t.Fatalf("Binary incorrect, wanted %v got %v", want, actual)
	}
}

func BenchmarkHexToBinary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HexToBinary("8D4840D6202CC371C32CE0576098")
	}
}
