package decode

import (
	"errors"
	"math"
	"pragmatic-zac/goModeS/internal"
	"strconv"
	"strings"
	"time"
)

type Position struct {
	latitude  float64
	longitude float64
}

type Velocity struct {
	speed      float64
	angle      float64
	vertRate   int32
	speedType  string
	rateSource string
}

type PositionInput struct {
	msg0   string
	msg1   string
	t0     time.Time
	t1     time.Time
	latRef *float64
	lonRef *float64
}

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

func AirbornePosition(input *PositionInput) (Position, error) {
	bin0, err := internal.HexToBinary(input.msg0)
	if err != nil {
		return Position{}, err
	}
	bin1, err := internal.HexToBinary(input.msg1)
	if err != nil {
		return Position{}, err
	}

	mb0 := bin0[32:]
	mb1 := bin1[32:]

	// check if the user mixed up odd/even messages
	oddEven0, _ := strconv.ParseInt(mb0[21:22], 2, 32)
	oddEven1, _ := strconv.ParseInt(mb1[21:22], 2, 32)
	if oddEven0 == 0 && oddEven1 == 1 {

	} else if oddEven0 == 1 && oddEven1 == 0 {
		input.msg0, input.msg1 = input.msg1, input.msg0
		input.t0, input.t1 = input.t1, input.t0
	} else {
		return Position{}, errors.New("both an even + odd message are required")
	}

	cprLatEven, err := strconv.ParseInt(mb0[22:39], 2, 64)
	if err != nil {
		return Position{}, err
	}
	latCprE := float64(cprLatEven) / 131072

	cprLonEven, err := strconv.ParseInt(mb0[39:56], 2, 64)
	if err != nil {
		return Position{}, err
	}
	lonCprE := float64(cprLonEven) / 131072

	cprLatOdd, err := strconv.ParseInt(mb1[22:39], 2, 64)
	if err != nil {
		return Position{}, err
	}
	latCprO := float64(cprLatOdd) / 131072

	cprLonOdd, err := strconv.ParseInt(mb1[39:56], 2, 64)
	if err != nil {
		return Position{}, err
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
		return Position{}, err
	}

	var lat float64
	var lon float64

	if input.t0.After(input.t1) {
		lat = latEven

		var nl = internal.CprNL(lat)

		ni := math.Max(nl, 1)

		m := math.Floor(lonCprE*(nl-1) - lonCprO*nl + 0.5)

		lon = (360 / ni) * (internal.Modulo(m, ni) + lonCprE)
	} else {
		lat = latOdd

		nl := internal.CprNL(lat)

		ni := math.Max(float64(nl)-1.0, 1)

		m := math.Floor(lonCprE*(nl-1) - lonCprO*nl + 0.5)

		lon = (360 / ni) * (internal.Modulo(m, ni) + lonCprO)
	}

	if lon > 180.0 {
		lon = lon - 360
	}

	pos := Position{
		latitude:  internal.RoundFloat(lat, 5),
		longitude: internal.RoundFloat(lon, 5),
	}

	return pos, nil
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

func SurfaceVelocity(msg string) (Velocity, error) {
	tc, err := internal.Typecode(msg)
	if err != nil {
		return Velocity{}, err
	}

	if tc < 5 || tc > 8 {
		err = errors.New("not a surface message, expecting a Typecode between 5 and 8")
		return Velocity{}, err
	}

	msgBin, err := internal.HexToBinary(msg)
	if err != nil {
		return Velocity{}, err
	}

	bin := msgBin[32:]

	// ground track
	var trk float64
	trkStatus, _ := strconv.Atoi(bin[12:13])
	if trkStatus == 1 {
		tmp, err := strconv.ParseInt(bin[13:20], 2, 64)
		if err != nil {
			return Velocity{}, err
		}

		trk = float64(tmp) * 360 / 128
		trk = internal.RoundFloat(trk, 1)
	} else {
		trk = 0
	}

	// ground speed
	mov, err := strconv.ParseInt(bin[5:12], 2, 64)
	if err != nil {
		return Velocity{}, err
	}

	var spd float64
	if mov == 0 || mov > 124 {
		spd = 0
	} else if mov == 1 {
		spd = 0
	} else if mov == 124 {
		spd = 175.0
	} else {
		mvmt := []int64{2, 9, 13, 39, 94, 109, 124}
		kts := []float64{0.125, 1, 2, 15, 70, 100, 175}
		step := []float64{0.125, 0.25, 0.5, 1, 2, 5}

		var idx int

		for i := 0; i < len(mvmt); i++ {
			if mov >= mvmt[i] && mov <= mvmt[i+1] {
				idx = i + 1
			}
		}

		spd = kts[idx-1] + float64(mov-mvmt[idx-1])*step[idx-1]
	}

	v := Velocity{
		speed:      spd,
		angle:      trk,
		vertRate:   0,
		speedType:  "GS",
		rateSource: "",
	}
	return v, nil
}

func AirborneVelocity(msg string) (Velocity, error) {
	tc, err := internal.Typecode(msg)
	if err != nil {
		return Velocity{}, err
	}

	if tc != 19 {
		err = errors.New("not an airborne velocity message, expecting typecode 19")
	}

	msgBin, err := internal.HexToBinary(msg)
	if err != nil {
		return Velocity{}, err
	}

	bin := msgBin[32:]

	subtype, err := strconv.ParseInt(bin[5:8], 2, 64)
	if err != nil {
		return Velocity{}, err
	}

	// check velocity components
	ew, err := strconv.ParseInt(bin[14:24], 2, 64)
	if err != nil {
		return Velocity{}, err
	}

	ns, err := strconv.ParseInt(bin[25:35], 2, 64)
	if err != nil {
		return Velocity{}, err
	}

	if ew == 0 || ns == 0 {
		return Velocity{}, err
	}

	var trk float64
	var spd int64
	var spdType string
	var vrSource string
	var vs int32

	if subtype == 1 || subtype == 2 {
		ewBit, _ := strconv.Atoi(bin[13:14]) // direction EW
		nsBit, _ := strconv.Atoi(bin[24:25]) // direction NS
		if ewBit == 1 {
			ewBit = -1
		}
		if nsBit == 1 {
			nsBit = -1
		}

		// check if velocity is supersonic
		if subtype == 2 {
			ewBit = ewBit * 4
			nsBit = nsBit * 4
		}

		ew = ew - 1
		ns = ns - 1

		vwe := int64(ewBit) * ew
		vsn := int64(nsBit) * ns

		spd = int64(math.Sqrt(float64(vsn*vsn + vwe*vwe)))

		trk = math.Atan2(float64(vwe), float64(vsn))
		trk = trk * (180 / math.Pi)
		if trk < 0 {
			trk = trk + 360
		}
		trk = internal.RoundFloat(trk, 2)
		spdType = "GS"
	} else {
		status, _ := strconv.Atoi(bin[13:14])
		if status == 0 {
			trk = 0
		} else {
			trk = float64(ew) / 1024.0 * 360.0
			trk = internal.RoundFloat(trk, 2)
		}

		if ns == 0 {
			spd = 0
		} else {
			spd = ns - 1
		}

		// supersonic check
		if subtype == 4 && spd != 0 {
			spd = spd * 4
		}

		tBit, _ := strconv.Atoi(bin[24:25])
		if tBit == 0 {
			spdType = "IAS"
		} else {
			spdType = "TAS"
		}
	}

	srcBit, _ := strconv.Atoi(bin[35:36])
	if srcBit == 0 {
		vrSource = "GNSS"
	} else {
		vrSource = "BARO"
	}

	var vrSign int64
	vrSignBit, _ := strconv.Atoi(bin[36:37])
	if vrSignBit == 1 {
		vrSign = -1
	}

	vr, err := strconv.ParseInt(bin[37:46], 2, 64)
	if err != nil {
		return Velocity{}, err
	}

	if vr == 0 {
		vs = 0
	} else {
		vs = int32(vrSign * (vr - 1) * 64)
	}

	v := Velocity{
		speed:      float64(spd),
		angle:      trk,
		vertRate:   vs,
		speedType:  spdType,
		rateSource: vrSource,
	}

	return v, nil
}
