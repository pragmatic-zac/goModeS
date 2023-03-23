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

func Altitude(binString string) (int, error) {
	if len(binString) != 13 {
		// also check to make sure it's binary?
		return 0, errors.New("binary string must be 13 bits long")
	}

	_, err := strconv.ParseInt(binString, 2, 64)
	if err != nil {
		return 0, errors.New("input must be a binary string")
	}

	Mbit := string(binString[6])
	Qbit := string(binString[8])

	r, _ := strconv.ParseInt(binString, 2, 64)
	if r == 0 {
		return 0, nil // altitude unknown or invalid
	}

	var alt int

	if Mbit == "0" { // unit in ft
		if Qbit == "1" { // 25ft interval
			vbin := binString[:6] + binString[7:8] + binString[9:]
			vint, _ := strconv.ParseInt(vbin, 2, 64)
			alt = int(vint)*25 - 1000
		}
		if Qbit == "0" { // 100ft interval, above 50187.5ft
			C1 := string(binString[0])
			A1 := string(binString[1])
			C2 := string(binString[2])
			A2 := string(binString[3])
			C4 := string(binString[4])
			A4 := string(binString[5])
			B1 := string(binString[7])
			B2 := string(binString[9])
			D2 := string(binString[10])
			B4 := string(binString[11])
			D4 := string(binString[12])

			graystr := D2 + D4 + A1 + A2 + A4 + B1 + B2 + B4 + C1 + C2 + C4
			alt = gray2alt(graystr)
		}
	}

	if Mbit == "1" { // unit in meter
		vbin := binString[:6] + binString[7:]
		vint, _ := strconv.ParseInt(vbin, 2, 64)
		alt = int(float64(vint) * 3.28084) // convert to ft
	}

	return alt, nil
}

func gray2int(binString string) int {
	num, _ := strconv.ParseInt(binString, 2, 64)
	num ^= num >> 8
	num ^= num >> 4
	num ^= num >> 2
	num ^= num >> 1
	return int(num)
}

func gray2alt(binString string) int {
	gc500 := binString[:8]
	n500 := gray2int(gc500)

	// 100-ft step must be converted first
	gc100 := binString[8:]
	n100 := gray2int(gc100)

	if n100 == 0 || n100 == 5 || n100 == 6 {
		return 0
	}

	if n100 == 7 {
		n100 = 5
	}

	if n500%2 == 1 {
		n100 = 6 - n100
	}

	alt := (n500*500 + n100*100) - 1300
	return alt
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
