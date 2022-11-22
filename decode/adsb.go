package decode

import (
	"errors"
	"pragmatic-zac/goModeS/internal"
	"strconv"
	"strings"
)

func Category(msg string) (int64, error) {
	tc, err := internal.Typecode(msg)
	if err != nil {
		return 0, err
	}

	if tc < 1 || tc > 4 {
		err = errors.New("not an identification message")
	}

	msgBin, err := internal.HexToBinary(msg)
	if err != nil {
		return 0, err
	}

	bin := msgBin[32:87]

	return strconv.ParseInt(bin[5:8], 2, 32)
}

func Callsign(msg string) (string, error) {
	lookup := "#ABCDEFGHIJKLMNOPQRSTUVWXYZ##### ###############0123456789######"

	bin, err := internal.HexToBinary(msg[8:22])
	if err != nil {
		println(err)
	}

	var callsign strings.Builder

	for i := 8; i < len(bin); i += 6 {
		output, err := strconv.ParseInt(bin[i:i+6], 2, 32)
		if err != nil {
			println(err)
		}

		callsign.WriteString(string(lookup[output]))
	}

	return callsign.String(), nil
}
