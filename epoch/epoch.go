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
	UnitGuess TimeUnit = iota
	UnitSeconds
	UnitMilliseconds
	UnitMicroseconds
	UnitNanoseconds
)

func ParseUnit(input string) (TimeUnit, error) {
	switch input {
	case "":
		return UnitGuess, nil
	case "s":
		return UnitSeconds, nil
	case "ms":
		return UnitMilliseconds, nil
	case "us":
		return UnitMicroseconds, nil
	case "ns":
		return UnitNanoseconds, nil
	}
	return 255, fmt.Errorf("failed to convert %v to time unit", input)
}

func Timestamp(t time.Time, unit TimeUnit) (int64, error) {
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

// TODO: add milli and microseconds
func FromTimestamp(timestamp int64, unit TimeUnit) (time.Time, error) {

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
	case UnitGuess:
		// guess if it's seconds or nanoseconds
		var (
			now     = time.Now()
			lenIn   = len(fmt.Sprintf("%v", timestamp))          // number of digits of timestamp to guess
			lenSec  = len(strconv.FormatInt(now.Unix(), 10))     // number of digits in current seconds timestamp
			lenNano = len(strconv.FormatInt(now.UnixNano(), 10)) // number of digits in current nanoseconds timestamp
		)
		// guessing if the input is seconds or nanoseconds based on
		// the difference of the length of the current epoch times
		if abs(lenSec-lenIn) < abs(lenNano-lenIn) {
			// number of digits is closer to current seconds timestamp
			return time.Unix(timestamp, 0), nil
		} else {
			// number of digits is closer to current nanoseconds timestamp
			return time.Unix(0, timestamp), nil
		}
	default:
		return time.Time{}, fmt.Errorf("unknown unit '%v'", unit)
	}
}

func FromFormatted(input string) (time.Time, error) {
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
