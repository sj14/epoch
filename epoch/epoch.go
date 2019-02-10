package epoch

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TimeUnit byte

const (
	UnitSeconds TimeUnit = iota
	UnitMilliseconds
	UnitMicroseconds
	UnitNanoseconds
)

// ParseUnit takes a string and returns the corresponding unit.
func ParseUnit(input string) (TimeUnit, error) {
	switch input {
	case "s", "sec":
		return UnitSeconds, nil
	case "ms", "milli":
		return UnitMilliseconds, nil
	case "us", "micro":
		return UnitMicroseconds, nil
	case "ns", "nano":
		return UnitNanoseconds, nil
	}
	return UnitSeconds, fmt.Errorf("failed to parse input '%v' to unit", input)
}

// ToTimestamp takes Go's default time type returns a timestamp of the given unit.
func ToTimestamp(t time.Time, unit TimeUnit) (int64, error) {
	epoch := t.Unix()

	switch unit {
	case UnitSeconds:
		// calculated as default value,
		// nothing to to
	case UnitMilliseconds:
		epoch = t.UnixNano() / (1000 * 1000)
	case UnitMicroseconds:
		epoch = t.UnixNano() / 1000
	case UnitNanoseconds:
		epoch = t.UnixNano()
	default:
		return 255, fmt.Errorf("unknown unit '%v'", unit)
	}

	return epoch, nil
}

func abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

// ParseTimestamp takes a timestamp of the given unit and returns Go's default time type.
func ParseTimestamp(timestamp int64, unit TimeUnit) (time.Time, error) {
	switch unit {
	case UnitSeconds:
		return time.Unix(timestamp, 0), nil
	case UnitMilliseconds:
		// add digits to match nanosecond accuracy
		timestamp *= 1000 * 1000
		return time.Unix(0, timestamp), nil
	case UnitMicroseconds:
		// add digits to match nanosecond accuracy
		timestamp *= 1000
		return time.Unix(0, timestamp), nil
	case UnitNanoseconds:
		return time.Unix(0, timestamp), nil
	default:
		return time.Time{}, fmt.Errorf("unknown unit '%v'", unit)
	}
}

// GuessUnit guesses if the input is sec, ms, us or ns based on
// the difference of the length (number of digits) of the 'ref' epoch times.
func GuessUnit(timestamp int64, ref time.Time) TimeUnit {
	var (
		lenIn    = len(fmt.Sprintf("%v", timestamp))                      // number of digits of timestamp to guess
		lenSec   = len(strconv.FormatInt(ref.Unix(), 10))                 // number of digits in current seconds timestamp
		lenMill  = len(strconv.FormatInt(ref.UnixNano()/(1000*1000), 10)) // number of digits in current milliseconds timestamp
		lenMicro = len(strconv.FormatInt(ref.UnixNano()/1000, 10))        // number of digits in current microseconds timestamp
		lenNano  = len(strconv.FormatInt(ref.UnixNano(), 10))             // number of digits in current nanoseconds timestamp

		diffSec   = abs(lenSec - lenIn)
		diffMill  = abs(lenMill - lenIn)
		diffMicro = abs(lenMicro - lenIn)
		diffNano  = abs(lenNano - lenIn)
	)

	// TODO: maybe there is a better way to do this guessing.
	if diffSec <= diffMill &&
		diffSec <= diffMicro &&
		diffSec <= diffNano {
		// number of digits is closer to current seconds timestamp
		return UnitSeconds
	} else if diffMill <= diffSec &&
		diffMill <= diffMicro &&
		diffMill <= diffNano {
		// number of digits is closer to current milliseconds timestamp
		return UnitMilliseconds
	} else if diffMicro <= diffSec &&
		diffMicro <= diffMill &&
		diffMicro <= diffNano {
		// number of digits is closer to current microseconds timestamp
		return UnitMicroseconds
	} else {
		// number of digits is closer to current nanoseconds timestamp
		return UnitNanoseconds
	}
}

// ParseFormatted takes a human readable time string and returns Go's default time type.
// Example input: "Mon, 02 Jan 2006 15:04:05 MST".
func ParseFormatted(input string) (time.Time, error) {
	// "Mon, 02 Jan 2006 15:04:05 MST"
	if t, err := time.Parse(time.RFC1123, input); err == nil {
		return t, nil
	}

	// "Mon, 02 Jan 2006 15:04:05 -0700"
	if t, err := time.Parse(time.RFC1123Z, input); err == nil {
		return t, nil
	}

	// "2006-01-02T15:04:05Z07:00"
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, nil
	}

	// "2006-01-02T15:04:05.999999999Z07:00"
	if t, err := time.Parse(time.RFC3339Nano, input); err == nil {
		return t, nil
	}

	// "02 Jan 06 15:04 MST"
	if t, err := time.Parse(time.RFC822, input); err == nil {
		return t, nil
	}

	// "02 Jan 06 15:04 -0700"
	if t, err := time.Parse(time.RFC822Z, input); err == nil {
		return t, nil
	}

	// "Monday, 02-Jan-06 15:04:05 MST"
	if t, err := time.Parse(time.RFC850, input); err == nil {
		return t, nil
	}

	// "Mon Jan _2 15:04:05 2006"
	if t, err := time.Parse(time.ANSIC, input); err == nil {
		return t, nil
	}

	// "Mon Jan _2 15:04:05 MST 2006"
	if t, err := time.Parse(time.UnixDate, input); err == nil {
		return t, nil
	}

	// "Mon Jan 02 15:04:05 -0700 2006"
	if t, err := time.Parse(time.RubyDate, input); err == nil {
		return t, nil
	}

	// "3:04PM"
	if t, err := time.Parse(time.Kitchen, input); err == nil {
		return t, nil
	}

	// "Jan _2 15:04:05"
	if t, err := time.Parse(time.Stamp, input); err == nil {
		return t, nil
	}

	// "Jan _2 15:04:05.000"
	if t, err := time.Parse(time.StampMilli, input); err == nil {
		return t, nil
	}

	// "Jan _2 15:04:05.000000"
	if t, err := time.Parse(time.StampMicro, input); err == nil {
		return t, nil
	}

	// "Jan _2 15:04:05.000000000"
	if t, err := time.Parse(time.StampNano, input); err == nil {
		return t, nil
	}

	// "Mon, 02 Jan 2006 15:04:05 GMT"
	if t, err := time.Parse(http.TimeFormat, input); err == nil {
		return t, nil
	}

	// handle Go's default time.Now() format (e.g. 2019-01-26 09:43:57.377055 +0100 CET m=+0.644739467)
	if t, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", strings.Split(input, " m=")[0]); err == nil {
		return t, nil
	}

	// "2019-01-25 21:51:38"
	if t, err := time.Parse("2006-01-02 15:04:05.999999999", input); err == nil {
		return t, nil
	}

	return time.Time{}, errors.New("failed to convert string to time")
}
