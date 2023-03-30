package decode

import (
	"testing"
)

var binTests = []struct {
	msg  string
	want string
}{
	{"8D4840D6202CC371C32CE0576098", "1000110101001000010000001101011000100000001011001100001101110001110000110010110011100000010101110110000010011000"},
	{"6E", "01101110"},
	{"6e", "01101110"},
}

func TestHexToBinary(t *testing.T) {
	for _, test := range binTests {
		t.Run(test.msg, func(t *testing.T) {
			actual, _ := hexToBinary(test.msg)
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

	if actual != want {
		t.Fatalf("DF incorrect, wanted %v got %v", want, actual)
	}
}

var icaoTests = []struct {
	msg  string
	want string
}{
	{"8D406B902015A678D4D220AA4BDA", "406B90"},
	{"A0001839CA3800315800007448D9", "400940"},
	{"A000139381951536E024D4CCF6B5", "3C4DD2"},
	{"A000029CFFBAA11E2004727281F1", "4243D0"},
}

func TestIcao(t *testing.T) {
	for _, test := range icaoTests {
		t.Run(test.msg, func(t *testing.T) {
			actual, _ := Icao(test.msg)
			if actual != test.want {
				t.Errorf("ICAO incorrect, wanted %v got %v", test.want, actual)
			}
		})
	}
}

func TestTypecode(t *testing.T) {
	msg := "8D4840D6202CC371C32CE0576098"

	actual, _ := Typecode(msg)

	want := int64(4)

	if actual != want {
		t.Fatalf("Typecode incorrect, wanted %v got %v", want, actual)
	}
}

var crcTests = []struct {
	msg  string
	want int
}{
	{"8D406B902015A678D4D220AA4BDA", 0},
	{"8d8960ed58bf053cf11bc5932b7d", 0},
	{"8d45cab390c39509496ca9a32912", 0},
	{"8d49d3d4e1089d00000000744c3b", 0},
	{"8d4400cd9b0000b4f87000e71a10", 0},
	{"8d4065de58a1054a7ef0218e226a", 0},
	{"c80b2dca34aa21dd821a04cb64d4", 10719924},
	{"8d4ca251204994b1c36e60a5343d", 16},
}

func TestCrc(t *testing.T) {
	for _, test := range crcTests {
		t.Run(test.msg, func(t *testing.T) {
			actual, _ := crc(test.msg, false)
			if actual != test.want {
				t.Errorf("CRC incorrect, wanted %v got %v", test.want, actual)
			}
		})
	}
}

func BenchmarkHexToBinary(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hexToBinary("8D4840D6202CC371C32CE0576098")
	}
}

func BenchmarkCrc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		crc("8D406B902015A678D4D220AA4BDA", false)
	}
}
