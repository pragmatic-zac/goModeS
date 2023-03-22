package internal

import (
	"errors"
	"fmt"
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
		bin.WriteString(m[strings.ToUpper(string(hex[i]))])
	}

	return bin.String(), nil
}

func Df(msg string) (int, error) {
	bin, err := HexToBinary(msg[0:2])
	if err != nil {
		return 0, err
	}

	df, err := strconv.ParseInt(bin[0:5], 2, 32)
	if err != nil {
		return 0, err
	}

	return int(df), nil
}

func Icao(msg string) (string, error) {
	df, err := Df(msg)
	if err != nil {
		return "", err
	}

	var addr string
	filter := []int{0, 4, 5, 16, 20, 21}

	if df == 11 || df == 17 || df == 18 {
		addr = msg[2:8]
	} else if contains(&filter, &df) {
		c0, _ := Crc(msg, true)
		c1, _ := strconv.ParseInt(msg[len(msg)-6:], 16, 32)
		result := c0 ^ int(c1)
		addr = fmt.Sprintf("%06X", result)
	}

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

func Crc(msg string, encode bool) (int, error) {
	if len(msg) != 28 {
		return 0, errors.New("message should be exactly 28 characters long")
	}

	G := []int{255, 250, 4, 128}

	if encode {
		msg = msg[:len(msg)-6] + "000000"
	}

	msgBin, err := HexToBinary(msg)
	if err != nil {
		return 0, err
	}
	msgBinSplit := wrap(msgBin, 8)
	mBytes := make([]int, 0, len(msgBinSplit))
	for _, s := range msgBinSplit {
		i, _ := strconv.ParseInt(s, 2, 32)
		mBytes = append(mBytes, int(i))
	}

	//for iByte := 0; iByte < len(mBytes)-3; iByte++ {
	//	for ibit := 0; ibit < 8; ibit++ {
	//		mask := 0x80 >> int(ibit)
	//		bits := mBytes[iByte] & mask
	//
	//		if bits > 0 {
	//			mBytes[iByte] = mBytes[iByte] ^ (G[0] >> ibit)
	//			mBytes[iByte+1] = mBytes[iByte+1] ^ (0xFF & ((G[0]<<8 - ibit) | (G[1] >> ibit)))
	//			mBytes[iByte+2] = mBytes[iByte+2] ^ (0xFF & ((G[1]<<8 - ibit) | (G[2] >> ibit)))
	//			mBytes[iByte+3] = mBytes[iByte+3] ^ (0xFF & ((G[2]<<8 - ibit) | (G[3] >> ibit)))
	//		}
	//	}
	//}

	for i := 0; i < len(mBytes)-3; i++ {
		for j := 0; j < 8; j++ {
			mask := 0x80 >> uint(j)
			bits := mBytes[i] & mask

			if bits > 0 {
				mBytes[i] ^= (G[0] >> uint(j))
				mBytes[i+1] ^= 0xFF & ((G[0] << (8 - uint(j))) | (G[1] >> uint(j)))
				mBytes[i+2] ^= 0xFF & ((G[1] << (8 - uint(j))) | (G[2] >> uint(j)))
				mBytes[i+3] ^= 0xFF & ((G[2] << (8 - uint(j))) | (G[3] >> uint(j)))
			}
		}
	}

	result := (mBytes[len(mBytes)-3] << 16) | (mBytes[len(mBytes)-2] << 8) | mBytes[len(mBytes)-1]

	return result, nil
}

func wrap(s string, length int) []string {
	var lines []string

	for i := 0; i < len(s); i = i + length {
		lines = append(lines, s[i:i+length])
	}

	return lines
}

func contains(s *[]int, e *int) bool {
	for _, a := range *s {
		if a == *e {
			return true
		}
	}
	return false
}
