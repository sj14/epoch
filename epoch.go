package epoch

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// TimeUnit represents a time unit.
type TimeUnit byte

const (
	// UnitSeconds represents seconds.
	UnitSeconds TimeUnit = iota
	// UnitMilliseconds represents milliseconds.
	UnitMilliseconds
	// UnitMicroseconds represents microseconds.
	UnitMicroseconds
	// UnitNanoseconds represents nanoseconds.
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
	}
	// number of digits is closer to current nanoseconds timestamp
	return UnitNanoseconds
}

// ErrParseFormatted is used when parsing the formatted string failed.
var ErrParseFormatted = errors.New("failed to convert string to time")

const (
	// FormatGo handles Go's default time.Now() format (e.g. 2019-01-26 09:43:57.377055 +0100 CET m=+0.644739467)
	FormatGo = "2006-01-02 15:04:05.999999999 -0700 MST"
	// FormatSimple handles "2019-01-25 21:51:38"
	FormatSimple = "2006-01-02 15:04:05.999999999"
	// FormatHTTP instead of importing main with http.TimeFormat which would increase the binary size significantly.
	FormatHTTP = "Mon, 02 Jan 2006 15:04:05 GMT"
)

// ParseFormatted takes a human readable time string and returns Go's default time type and the layout it recognized.
// Example input: "Mon, 02 Jan 2006 15:04:05 MST".
func ParseFormatted(input string) (time.Time, string, error) {
	// "Mon, 02 Jan 2006 15:04:05 MST"
	if t, err := time.Parse(time.RFC1123, input); err == nil {
		return t, time.RFC1123, nil
	}

	// "Mon, 02 Jan 2006 15:04:05 -0700"
	if t, err := time.Parse(time.RFC1123Z, input); err == nil {
		return t, time.RFC1123Z, nil
	}

	// "2006-01-02T15:04:05Z07:00"
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, time.RFC3339, nil
	}

	// "2006-01-02T15:04:05.999999999Z07:00"
	if t, err := time.Parse(time.RFC3339Nano, input); err == nil {
		return t, time.RFC3339Nano, nil
	}

	// "02 Jan 06 15:04 MST"
	if t, err := time.Parse(time.RFC822, input); err == nil {
		return t, time.RFC822, nil
	}

	// "02 Jan 06 15:04 -0700"
	if t, err := time.Parse(time.RFC822Z, input); err == nil {
		return t, time.RFC822Z, nil
	}

	// "Monday, 02-Jan-06 15:04:05 MST"
	if t, err := time.Parse(time.RFC850, input); err == nil {
		return t, time.RFC850, nil
	}

	// "Mon Jan _2 15:04:05 2006"
	if t, err := time.Parse(time.ANSIC, input); err == nil {
		return t, time.ANSIC, nil
	}

	// "Mon Jan _2 15:04:05 MST 2006"
	if t, err := time.Parse(time.UnixDate, input); err == nil {
		return t, time.UnixDate, nil
	}

	// "Mon Jan 02 15:04:05 -0700 2006"
	if t, err := time.Parse(time.RubyDate, input); err == nil {
		return t, time.RubyDate, nil
	}

	// "3:04PM"
	if t, err := time.Parse(time.Kitchen, input); err == nil {
		return t, time.Kitchen, nil
	}

	// "Jan _2 15:04:05"
	if t, err := time.Parse(time.Stamp, input); err == nil {
		return t, time.Stamp, nil
	}

	// "Jan _2 15:04:05.000"
	if t, err := time.Parse(time.StampMilli, input); err == nil {
		return t, time.StampMilli, nil
	}

	// "Jan _2 15:04:05.000000"
	if t, err := time.Parse(time.StampMicro, input); err == nil {
		return t, time.StampMicro, nil
	}

	// "Jan _2 15:04:05.000000000"
	if t, err := time.Parse(time.StampNano, input); err == nil {
		return t, time.StampNano, nil
	}

	// "Mon, 02 Jan 2006 15:04:05 GMT"
	if t, err := time.Parse(FormatHTTP, input); err == nil {
		return t, FormatHTTP, nil
	}

	if t, err := time.Parse(FormatGo, strings.Split(input, " m=")[0]); err == nil {
		return t, FormatGo, nil
	}

	// "2019-01-25 21:51:38"
	if t, err := time.Parse(FormatSimple, input); err == nil {
		return t, FormatSimple, nil
	}

	return time.Time{}, "", ErrParseFormatted
}

// Operator for arithemtic operation.
type Operator uint

const (
	// Undefined operator.
	Undefined Operator = iota
	// Add operation.
	Add
	// Sub operation.
	Sub
)

// ErrUnkownOperator is returned when no matching operator was found.
var ErrUnkownOperator = errors.New("unkown operator")

// ToOperator return the matching operator to the given string.
func ToOperator(s string) (Operator, error) {
	switch s {
	case "+", "add", "plus":
		return Add, nil
	case "-", "sub", "minus":
		return Sub, nil
	}
	return Undefined, fmt.Errorf("%w: '%v'", ErrUnkownOperator, s)
}

// Arithmetics does basic add/sub calculations on the given input.
func Arithmetics(input time.Time, op Operator, duration time.Duration) time.Time {
	switch op {
	case Undefined:
		return input
	case Add:
		return input.Add(duration)
	case Sub:
		return input.Add(-duration)
	default:
		log.Fatalf("failed parsing operator: '%v'\n", op)
	}

	return time.Time{}
}