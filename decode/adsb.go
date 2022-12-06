package decode

import (
	"errors"
	"math"
	"pragmatic-zac/goModeS/internal"
	"strconv"
	"strings"
	"time"
)

func Df(msg string) (int64, error) {
	res, err := internal.Df(msg)
	if err != nil {
		return 0, err
	}

	return res, nil
}

func Icao(msg string) (string, error) {
	res, err := internal.Icao(msg)
	if err != nil {
		return "", err
	}

	return res, nil
}

func Typecode(msg string) (int64, error) {
	res, err := internal.Typecode(msg)
	if err != nil {
		return 0, err
	}

	return res, nil
}

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

func AirbornePosition(msg0 string, msg1 string, t0 time.Time, t1 time.Time) (float64, float64, error) {
	bin0, err := internal.HexToBinary(msg0)
	if err != nil {
		return 0, 0, err
	}
	bin1, err := internal.HexToBinary(msg1)
	if err != nil {
		return 0, 0, err
	}

	mb0 := bin0[32:]
	mb1 := bin1[32:]

	// add oe1 and oe2 checks

	// fix these names
	// we need the first one just to parse the int value
	// we need the second one to be a float and that's what we're really working with
	// maybe split this into another function and use a pointer?
	// this seems like a ridiculous amount of allocation
	cprLatEven, err := strconv.ParseInt(mb0[22:39], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	latCprE := float64(cprLatEven) / 131072

	cprLonEven, err := strconv.ParseInt(mb0[39:56], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	lonCprE := float64(cprLonEven) / 131072

	cprLatOdd, err := strconv.ParseInt(mb1[22:39], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	latCprO := float64(cprLatOdd) / 131072

	cprLonOdd, err := strconv.ParseInt(mb1[39:56], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	lonCprO := float64(cprLonOdd) / 131072

	airDLatEven := 360.0 / 60.0
	airDLatOdd := 360.0 / 59.0

	j := int(math.Floor(59*latCprE - 60*latCprO + 0.5))

	t := float64(j % 59)

	latEven := airDLatEven * (internal.Modulo(float64(j), 60) + latCprE)
	latOdd := airDLatOdd * (t + latCprO)

	if latEven >= 270 {
		latEven = latEven - 360
	}

	if latOdd >= 270 {
		latOdd = latOdd - 360
	}

	if internal.CprNL(latEven) != internal.CprNL(latOdd) {
		return 0, 0, nil
	}

	var lat float64
	var lon float64

	if t0.After(t1) {
		lat = latEven

		var nl float64 = internal.CprNL(lat)

		ni := math.Max(nl, 1)

		m := math.Floor(lonCprE*(nl-1) - lonCprO*nl + 0.5)

		lon = (360 / ni) * (internal.Modulo(m, ni) + lonCprE)
	} else {
		lat = float64(latOdd)

		nl := internal.CprNL(lat)

		ni := math.Max(float64(nl)-1.0, 1)

		m := math.Floor(lonCprE*(nl-1) - lonCprO*nl + 0.5)

		lon = (360 / ni) * (internal.Modulo(m, ni) + float64(lonCprO))
	}

	if lon > 180.0 {
		lon = lon - 360
	}

	return internal.RoundFloat(lat, 5), internal.RoundFloat(lon, 5), nil
}

func SurfacePosition(msg0 string, msg1 string, t0 time.Time, t1 time.Time, latRef float64, lonRef float64) (float64, float64, error) {
	bin0, err := internal.HexToBinary(msg0)
	if err != nil {
		return 0, 0, err
	}
	bin1, err := internal.HexToBinary(msg1)
	if err != nil {
		return 0, 0, err
	}

	cprLatEven, err := strconv.ParseInt(bin0[54:71], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	latCprE := float64(cprLatEven) / 131072

	cprLonEven, err := strconv.ParseInt(bin0[71:88], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	lonCprE := float64(cprLonEven) / 131072

	cprLatOdd, err := strconv.ParseInt(bin1[54:71], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	latCprO := float64(cprLatOdd) / 131072

	cprLonOdd, err := strconv.ParseInt(bin1[71:88], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	lonCprO := float64(cprLonOdd) / 131072

	airDLatEven := 90.0 / 60.0
	airDLatOdd := 90.0 / 59.0

	j := int(math.Floor(59*latCprE - 60*latCprO + 0.5))

	// north hemisphere
	latEvenN := airDLatEven * (float64(j%60) + latCprE)
	latOddN := airDLatOdd * (float64(j%59) + latCprO)

	// south hemisphere
	latEvenS := latEvenN - 90
	latOddS := latOddN - 90

	var latE float64
	var latO float64
	if latRef > 0 {
		latE = latEvenN
		latO = latOddN
	} else {
		latE = latEvenS
		latO = latOddS
	}

	// check if both are in same lat zone
	if internal.CprNL(latE) != internal.CprNL(latO) {
		return 0, 0, nil
	}

	var lat float64
	var lon float64
	if t0.After(t1) {
		lat = latE
		nl := internal.CprNL(latE)
		ni := math.Max(nl, 1)
		m := math.Floor(lonCprE*(nl-1.0) - lonCprO*nl + 0.5)
		lon = (90 / ni) * (math.Mod(m, ni) + lonCprE)
	} else {
		lat = latO
		nl := internal.CprNL(latO)
		ni := math.Max(nl-1, 1)
		m := math.Floor(lonCprE*(nl-1.0) - lonCprO*nl + 0.5)
		lon = (90 / ni) * (math.Mod(m, ni) + lonCprO)
	}

	// there are four possible solutions
	lons := []float64{lon, lon + 90, lon + 180, lon + 270}

	// make sure all lon values are valid, between -180 and 180
	for i, f := range lons {
		lons[i] = math.Mod(f+180, 360) - 180
	}

	// we want the one closest to the receiver
	var closest float64
	for _, f := range lons {
		abs := math.Abs(lonRef - f)
		if abs < closest {
			closest = f
		}
	}

	return internal.RoundFloat(lat, 5), internal.RoundFloat(lon, 5), nil
}

func SurfacePositionWithRef(msg string, latRef float64, lonRef float64) (float64, float64, error) {
	msgBin, err := internal.HexToBinary(msg)
	if err != nil {
		return 0, 0, err
	}

	bin := msgBin[32:]

	cprLatInt, err := strconv.ParseInt(bin[22:39], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	cprLat := float64(cprLatInt) / 131072

	cprLonInt, err := strconv.ParseInt(bin[39:56], 2, 64)
	if err != nil {
		return 0, 0, err
	}
	cprLon := float64(cprLonInt) / 131072

	i, _ := strconv.Atoi(bin[21:22])
	var dLat float64
	if i != 0 {
		dLat = 90.0 / 59.0
	} else {
		dLat = 90.0 / 60.0
	}

	j := math.Floor(latRef/dLat) + math.Floor(0.5+((math.Mod(latRef, dLat)/dLat)-cprLat))

	lat := dLat * (j + cprLat)

	ni := internal.CprNL(lat) - float64(i)

	var dLon float64
	if ni > 0 {
		dLon = 90.0 / ni
	} else {
		dLon = 90.0
	}

	m := math.Floor(lonRef/dLon) + math.Floor(0.5+((math.Mod(lonRef, dLon)/dLon)-cprLon))

	lon := dLon * (m + cprLon)

	return internal.RoundFloat(lat, 6), internal.RoundFloat(lon, 6), nil
}
