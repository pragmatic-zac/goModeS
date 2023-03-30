// Package decode provides methods for decoding ADS-B messages.
package decode

import (
	"errors"
	"math"
	"strconv"
	"strings"
	"time"
)

// Position is a struct that represents the calculated airborne position information, including the latitude and longitude.
//
// Fields:
//   - Latitude: a float64 that represents the latitude of the airborne position.
//   - Longitude: a float64 that represents the longitude of the airborne position.
type Position struct {
	Latitude  float64
	Longitude float64
}

// Velocity is a struct that represents the calculated airborne velocity information, including the speed, angle,
// vertical rate, and speed type.
//
// Fields:
//   - Speed: a float64 that represents the speed of the airborne velocity in knots.
//   - Angle: a float64 that represents the angle of the airborne velocity in degrees.
//   - VertRate: an int32 that represents the vertical rate of the airborne velocity in feet per minute.
//   - SpeedType: a string that represents the type of speed, either "GS" for ground speed, "IAS" for indicated air speed, or "TAS" for true air speed.
//   - RateSource: a string that represents the source of the vertical rate, either "GNSS" or "BARO".
type Velocity struct {
	Speed      float64
	Angle      float64
	VertRate   int32
	SpeedType  string
	RateSource string
}

// PositionInput is a struct that represents the necessary information to calculate the airborne position, including the
// message strings, time stamps, and latitude and longitude reference points.
//
// Fields:
//   - Msg0: 28 character hexadecimal string message, even frame.
//   - Msg1: 28 character hexadecimal string message, odd frame.
//   - T0: a time.Time that represents the time stamp of the even message.
//   - T1: a time.Time that represents the time stamp of the odd message.
//   - LatRef: a pointer to a float64 that represents the latitude reference point used to calculate the airborne position.
//     This field may be nil if no reference point is needed.
//   - LonRef: a pointer to a float64 that represents the longitude reference point used to calculate the airborne position.
//     This field may be nil if no reference point is needed.
type PositionInput struct {
	Msg0   string
	Msg1   string
	T0     time.Time
	T1     time.Time
	LatRef *float64
	LonRef *float64
}

// Df is a function that decodes the Downlink Format value.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - int: an integer that represents the result of the `df` function if successful.
//   - error: an error that indicates whether an error occurred during the processing of the message.
func Df(msg string) (int, error) {
	res, err := df(msg)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// Icao is a function that decodes the ICAO value. Usable with DF4, DF5, DF20, DF21 messages.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - string: a string that represents the ICAO value if successful.
//   - error: an error that indicates whether an error occurred during the processing of the message.
func Icao(msg string) (string, error) {
	res, err := icao(msg)
	if err != nil {
		return "", err
	}

	return res, nil
}

// Typecode is a function that decodes the typecode value of an ADS-B message.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - int: an integer that represents the message typecode if successful.
//   - error: an error that indicates whether an error occurred during the processing of the message.
func Typecode(msg string) (int64, error) {
	res, err := typecode(msg)
	if err != nil {
		return 0, err
	}

	return res, nil
}

// Category is a function that decodes the category value of an ADS-B message.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - int: an integer that represents the message category if successful.
//   - error: an error that indicates whether an error occurred during the processing of the message.
func Category(msg string) (int64, error) {
	tc, err := typecode(msg)
	if err != nil {
		return 0, err
	}

	if tc < 1 || tc > 4 {
		err = errors.New("not an identification message")
	}

	msgBin, err := hexToBinary(msg)
	if err != nil {
		return 0, err
	}

	bin := msgBin[32:87]

	return strconv.ParseInt(bin[5:8], 2, 32)
}

// Callsign is a function that decodes the callsign value in an ADS-B message.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - string: a string that represents the aircraft's callsign if successful.
//   - error: an error that indicates whether an error occurred during the processing of the message.
func Callsign(msg string) (string, error) {
	lookup := "#ABCDEFGHIJKLMNOPQRSTUVWXYZ##### ###############0123456789######"

	bin, err := hexToBinary(msg[8:22])
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

// AirbornePosition is a function that takes a PositionInput as input and returns a Position and an error.
// It calculates the airborne position based on the given input and returns the result if successful.
//
// Parameters:
//   - input: a struct that contains the necessary information to calculate the airborne position
//
// Returns:
//   - Position: a struct that contains the calculated airborne position information, including the latitude
//     and longitude.
//   - error: an error that indicates whether an error occurred during the calculation of the airborne position.
func AirbornePosition(input PositionInput) (Position, error) {
	bin0, err := hexToBinary(input.Msg0)
	if err != nil {
		return Position{}, err
	}
	bin1, err := hexToBinary(input.Msg1)
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
		input.Msg0, input.Msg1 = input.Msg1, input.Msg0
		input.T0, input.T1 = input.T1, input.T0
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

	latEven := airDLatEven * (modulo(float64(j), 60) + latCprE)
	latOdd := airDLatOdd * (t + latCprO)

	if latEven >= 270 {
		latEven = latEven - 360
	}

	if latOdd >= 270 {
		latOdd = latOdd - 360
	}

	if cprNL(latEven) != cprNL(latOdd) {
		return Position{}, err
	}

	var lat float64
	var lon float64

	if input.T0.After(input.T1) {
		lat = latEven

		var nl = cprNL(lat)

		ni := math.Max(nl, 1)

		m := math.Floor(lonCprE*(nl-1) - lonCprO*nl + 0.5)

		lon = (360 / ni) * (modulo(m, ni) + lonCprE)
	} else {
		lat = latOdd

		nl := cprNL(lat)

		ni := math.Max(float64(nl)-1.0, 1)

		m := math.Floor(lonCprE*(nl-1) - lonCprO*nl + 0.5)

		lon = (360 / ni) * (modulo(m, ni) + lonCprO)
	}

	if lon > 180.0 {
		lon = lon - 360
	}

	pos := Position{
		Latitude:  roundFloat(lat, 5),
		Longitude: roundFloat(lon, 5),
	}

	return pos, nil
}

// AirbornePositionWithRef is a function that takes a message, a latitude reference, and a longitude reference as input,
// and returns a Position and an error. It calculates the airborne position based on the given message and reference points,
// and returns the result if successful.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//   - latRef: a float64 that represents the latitude reference point used to calculate the airborne position.
//   - lonRef: a float64 that represents the longitude reference point used to calculate the airborne position.
//
// Returns:
//   - Position: a struct that contains the calculated airborne position information, including the latitude
//     and longitude.
//   - error: an error that indicates whether an error occurred during the calculation of the airborne position.
func AirbornePositionWithRef(msg string, latRef float64, lonRef float64) (Position, error) {
	msgBin, err := hexToBinary(msg)
	if err != nil {
		return Position{}, err
	}

	bin := msgBin[32:]

	cprLatInt, err := strconv.ParseInt(bin[22:39], 2, 64)
	if err != nil {
		return Position{}, err
	}
	cprLat := float64(cprLatInt) / 131072

	cprLonInt, err := strconv.ParseInt(bin[39:56], 2, 64)
	if err != nil {
		return Position{}, err
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

	ni := cprNL(lat) - float64(i)

	var dLon float64
	if ni > 0 {
		dLon = 90.0 / ni
	} else {
		dLon = 90.0
	}

	m := math.Floor(lonRef/dLon) + math.Floor(0.5+((math.Mod(lonRef, dLon)/dLon)-cprLon))

	lon := dLon * (m + cprLon)

	p := Position{
		Latitude:  roundFloat(lat, 6),
		Longitude: roundFloat(lon, 6),
	}

	return p, nil
}

// SurfacePosition is a function that takes a PositionInput as input and returns a Position and an error.
// It calculates the surface position based on the given input and returns the result if successful.
//
// Parameters:
//   - input: a struct that contains the necessary information to calculate the surface position, including
//     the message strings, time stamps, and latitude and longitude reference points.
//
// Returns:
//   - Position: a struct that contains the calculated surface position information, including the latitude and
//     longitude.
//   - error: an error that indicates whether an error occurred during the calculation of the surface position.
func SurfacePosition(input PositionInput) (Position, error) {
	bin0, err := hexToBinary(input.Msg0)
	if err != nil {
		return Position{}, err
	}
	bin1, err := hexToBinary(input.Msg1)
	if err != nil {
		return Position{}, err
	}

	cprLatEven, err := strconv.ParseInt(bin0[54:71], 2, 64)
	if err != nil {
		return Position{}, err
	}
	latCprE := float64(cprLatEven) / 131072

	cprLonEven, err := strconv.ParseInt(bin0[71:88], 2, 64)
	if err != nil {
		return Position{}, err
	}
	lonCprE := float64(cprLonEven) / 131072

	cprLatOdd, err := strconv.ParseInt(bin1[54:71], 2, 64)
	if err != nil {
		return Position{}, err
	}
	latCprO := float64(cprLatOdd) / 131072

	cprLonOdd, err := strconv.ParseInt(bin1[71:88], 2, 64)
	if err != nil {
		return Position{}, err
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
	if *input.LatRef > 0 {
		latE = latEvenN
		latO = latOddN
	} else {
		latE = latEvenS
		latO = latOddS
	}

	// check if both are in same lat zone
	if cprNL(latE) != cprNL(latO) {
		return Position{}, err
	}

	var lat float64
	var lon float64
	if input.T0.After(input.T1) {
		lat = latE
		nl := cprNL(latE)
		ni := math.Max(nl, 1)
		m := math.Floor(lonCprE*(nl-1.0) - lonCprO*nl + 0.5)
		lon = (90 / ni) * (math.Mod(m, ni) + lonCprE)
	} else {
		lat = latO
		nl := cprNL(latO)
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
		abs := math.Abs(*input.LonRef - f)
		if abs < closest {
			closest = f
		}
	}

	pos := Position{
		Latitude:  roundFloat(lat, 5),
		Longitude: roundFloat(lon, 5),
	}

	return pos, nil
}

// SurfacePositionWithRef is a function that takes a single message, a latitude reference, and a longitude reference as input,
// and returns a Position and an error. It calculates the surface position based on the given message and reference points,
// and returns the result if successful.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//   - latRef: a float64 that represents the latitude reference point used to calculate the surface position.
//   - lonRef: a float64 that represents the longitude reference point used to calculate the surface position.
//
// Returns:
//   - Position: a struct that contains the calculated surface position information, including the latitude and
//     longitude.
//   - error: an error that indicates whether an error occurred during the calculation of the surface position.
func SurfacePositionWithRef(msg string, latRef float64, lonRef float64) (Position, error) {
	msgBin, err := hexToBinary(msg)
	if err != nil {
		return Position{}, err
	}

	bin := msgBin[32:]

	cprLatInt, err := strconv.ParseInt(bin[22:39], 2, 64)
	if err != nil {
		return Position{}, err
	}
	cprLat := float64(cprLatInt) / 131072

	cprLonInt, err := strconv.ParseInt(bin[39:56], 2, 64)
	if err != nil {
		return Position{}, err
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

	ni := cprNL(lat) - float64(i)

	var dLon float64
	if ni > 0 {
		dLon = 90.0 / ni
	} else {
		dLon = 90.0
	}

	m := math.Floor(lonRef/dLon) + math.Floor(0.5+((math.Mod(lonRef, dLon)/dLon)-cprLon))

	lon := dLon * (m + cprLon)

	p := Position{
		Latitude:  roundFloat(lat, 6),
		Longitude: roundFloat(lon, 6),
	}

	return p, nil
}

// SurfaceVelocity is a function that takes a message string as input and returns a Velocity and an error.
// It calculates the surface velocity based on the given message and returns the result if successful.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - Velocity: a struct that contains the calculated surface velocity information, including the speed, angle,
//     vertical rate, speed type, and rate source.
//   - error: an error that indicates whether an error occurred during the calculation of the surface velocity.
func SurfaceVelocity(msg string) (Velocity, error) {
	tc, err := Typecode(msg)
	if err != nil {
		return Velocity{}, err
	}

	if tc < 5 || tc > 8 {
		err = errors.New("not a surface message, expecting a Typecode between 5 and 8")
		return Velocity{}, err
	}

	msgBin, err := hexToBinary(msg)
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
		trk = roundFloat(trk, 1)
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
		Speed:      spd,
		Angle:      trk,
		VertRate:   0,
		SpeedType:  "GS",
		RateSource: "",
	}
	return v, nil
}

// AirborneVelocity is a function that takes a message string as input and returns a Velocity and an error.
// It calculates the airborne velocity based on the given message and returns the result if successful.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - Velocity: a struct that contains the calculated airborne velocity information, including the speed, angle,
//     vertical rate, speed type, and rate source.
//   - error: an error that indicates whether an error occurred during the calculation of the airborne velocity.
func AirborneVelocity(msg string) (Velocity, error) {
	tc, err := Typecode(msg)
	if err != nil {
		return Velocity{}, err
	}

	if tc != 19 {
		err = errors.New("not an airborne velocity message, expecting typecode 19")
	}

	msgBin, err := hexToBinary(msg)
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
		trk = roundFloat(trk, 2)
		spdType = "GS"
	} else {
		status, _ := strconv.Atoi(bin[13:14])
		if status == 0 {
			trk = 0
		} else {
			trk = float64(ew) / 1024.0 * 360.0
			trk = roundFloat(trk, 2)
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
		Speed:      float64(spd),
		Angle:      trk,
		VertRate:   vs,
		SpeedType:  spdType,
		RateSource: vrSource,
	}

	return v, nil
}

// CombinedVelocity is a function that takes a message string as input and returns a Velocity and an error.
// It calculates the airborne OR surface velocity based on the given message and returns the result if successful.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - Velocity: a struct that contains the calculated airborne velocity information, including the speed, angle,
//     vertical rate, speed type, and rate source.
//   - error: an error that indicates whether an error occurred during the calculation of the airborne velocity.
func CombinedVelocity(msg string) (Velocity, error) {
	tc, err := Typecode(msg)
	if err != nil {
		return Velocity{}, err
	}

	if tc >= 5 && tc <= 8 {
		v, err := SurfaceVelocity(msg)
		if err != nil {
			return Velocity{}, err
		}
		return v, nil
	} else if tc == 19 {
		v, err := AirborneVelocity(msg)
		if err != nil {
			return Velocity{}, err
		}
		return v, nil
	} else {
		return Velocity{}, errors.New("incorrect message type, expecting 5 thru 8 or 19")
	}
}

// Altitude is a function that takes a message string as input and returns an integer altitude value and an error.
// It calculates the altitude based on the given message and returns the result if successful.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - int: an integer that represents the calculated altitude in feet.
//   - error: an error that indicates whether an error occurred during the calculation of the altitude.
func Altitude(msg string) (int, error) {
	tc, err := Typecode(msg)
	if err != nil {
		return 0, err
	}

	// check for surface position and return 0

	if tc < 9 || tc == 19 || tc > 22 {
		return 0, errors.New("cannot decode altitude, not an airborne position message")
	}

	bin, err := hexToBinary(msg)
	if err != nil {
		return 0, err
	}

	msgBin := bin[32:]

	var alt int

	altBin := msgBin[8:20]
	if tc < 19 {
		altCode := altBin[0:6] + "0" + altBin[6:]
		alt, err = altitude(altCode)
		if err != nil {
			return 0, err
		}
	} else {
		n, _ := strconv.ParseInt(altBin, 2, 64)
		alt = int(float64(n) * 3.28084)
	}

	return alt, nil
}

// OddEvenFlag is a function that takes a message string as input and returns an integer.
// It decodes whether the message is an odd or even frame.
//
// Parameters:
//   - msg: 28 character hexadecimal string message.
//
// Returns:
//   - int: an integer that represents the calculated odd/even flag for the message. The value can be either 0 or 1.
func OddEvenFlag(msg string) int {
	bin, _ := hexToBinary(msg)
	res, _ := strconv.Atoi(bin[53:54])
	return res
}
