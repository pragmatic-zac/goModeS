package internal

import (
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
		println("Error converting DF to binary: ", err)
	}

	df, err := strconv.ParseInt(bin[0:5], 2, 32)
	if err != nil {
		return 0, err
	}

	return df, nil
}
