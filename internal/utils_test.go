package internal

import "testing"

var binTests = []struct {
	msg  string
	want string
}{
	{"8D4840D6202CC371C32CE0576098", "1000110101001000010000001101011000100000001011001100001101110001110000110010110011100000010101110110000010011000"},
	{"6E", "01101110"},
	{"6e", "01101110"},
}

func TestHexToBin(t *testing.T) {
	for _, test := range binTests {
		t.Run(test.msg, func(t *testing.T) {
			actual, _ := HexToBinary(test.msg)
			if actual != test.want {
				t.Errorf("Binary incorrect, wanted %v got %v", test.want, actual)
			}
		})
	}
}

func TestDf(t *testing.T) {
	msg := "8D4840D6202CC371C32CE0576098"

	actual, _ := Df(msg)

	want := 17

	if actual != int64(want) {
		t.Fatalf("DF incorrect, wanted %v got %v", want, actual)
	}
}

func TestIcao(t *testing.T) {
	msg := "8D406B902015A678D4D220AA4BDA"

	actual, _ := Icao(msg)

	want := "406B90"

	if actual != want {
		t.Fatalf("ICAO incorrect, wanted %v got %v", want, actual)
	}
}

func TestTypecode(t *testing.T) {
	msg := "8D4840D6202CC371C32CE0576098"

	actual, _ := Typecode(msg)

	want := int64(4)

	if actual != want {
		t.Fatalf("ICAO incorrect, wanted %v got %v", want, actual)
	}
}

func TestCrc(t *testing.T) {
	msg := "8D406B902015A678D4D220AA4BDA"

	actual, _ := Crc(msg, false)

	want := 0

	if actual != want {
		t.Fatalf("CRC incorrect, wanted %v got %v", want, actual)
	}
}

func BenchmarkHexToBinary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HexToBinary("8D4840D6202CC371C32CE0576098")
	}
}
