package internal

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func HexToBinary(hex string) (string, error) {
	m := map[string]string{
		"0": "0000",
		"1": "0001",
		"2": "0010",
		"3": "0011",
		"4": "0100",
		"5": "0101",
		"6": "0110",
		"7": "0111",
		"8": "1000",
		"9": "1001",
		"A": "1010",
		"B": "1011",
		"C": "1100",
		"D": "1101",
		"E": "1110",
		"F": "1111",
	}

	var bin strings.Builder

	for i := 0; i < len(hex); i++ {
		bin.WriteString(m[string(hex[i])])
	}

	return bin.String(), nil
}

func Df(msg string) (int64, error) {
	bin, err := HexToBinary(msg[0:2])
	if err != nil {
		return 0, err
	}

	df, err := strconv.ParseInt(bin[0:5], 2, 32)
	if err != nil {
		return 0, err
	}

	return df, nil
}

func Icao(msg string) (string, error) {
	df, err := Df(msg)
	if err != nil {
		return "", err
	}

	var addr string

	// currently only supporting DF17
	if df != 17 {
		return "", errors.New("currently only DF17 messages are supported")
	}

	addr = msg[2:8]

	return addr, nil
}

func Typecode(msg string) (int64, error) {
	df, err := Df(msg)
	if err != nil {
		return 0, err
	}

	if df != 17 {
		return 0, nil
	}

	bin, err := HexToBinary(msg[8:10])
	if err != nil {
		return 0, nil
	}

	typecode, err := strconv.ParseInt(bin[0:5], 2, 32)
	if err != nil {
		return 0, err
	}

	return typecode, nil
}

func Modulo(x float64, y float64) float64 {
	if y == 0.0 {
		panic("Y may not be zero.") // panic or error?
	}

	return x - y*math.Floor(x/y)
}

func CprNL(lat float64) float64 {
	if lat == 0 {
		return 59
	}

	if math.Abs(lat) == 87 {
		return 2
	}

	if lat > 87 || lat < -87 {
		return 1
	}

	nz := 15.0

	denom := math.Acos(1 - ((1 - math.Cos(math.Pi/(2*nz))) / math.Pow(math.Cos((math.Pi/180)*lat), 2)))

	nl := math.Floor(2 * math.Pi / denom)

	return nl
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
