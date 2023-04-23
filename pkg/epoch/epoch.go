package epoch

import (
	"errors"
	"fmt"
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

func abs(i int64) int64 {
	if i < 0 {
		return -i
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
// the difference to the 'ref' epoch times.
func GuessUnit(timestamp int64, ref time.Time) TimeUnit {
	var (
		asSec   = ref.Unix()
		asMill  = ref.UnixNano() / (1000 * 1000)
		asMicro = ref.UnixNano() / 1000
		asNano  = ref.UnixNano()

		diffSec   = abs(asSec - timestamp)
		diffMill  = abs(asMill - timestamp)
		diffMicro = abs(asMicro - timestamp)
		diffNano  = abs(asNano - timestamp)
	)

	// TODO: maybe there is a better way to do this guessing.
	if diffSec <= diffMill &&
		diffSec <= diffMicro &&
		diffSec <= diffNano {
		// difference is closer to current seconds timestamp
		return UnitSeconds
	} else if diffMill <= diffSec &&
		diffMill <= diffMicro &&
		diffMill <= diffNano {
		// difference is closer to current milliseconds timestamp
		return UnitMilliseconds
	} else if diffMicro <= diffSec &&
		diffMicro <= diffMill &&
		diffMicro <= diffNano {
		// difference is closer to current microseconds timestamp
		return UnitMicroseconds
	}
	// difference is closer to current nanoseconds timestamp
	return UnitNanoseconds
}

// ErrParseFormatted is used when parsing the formatted string failed.
var ErrParseFormatted = errors.New("failed to convert string to time")

const (
	// TimeFormatGo handles Go's default time.Now() format (e.g. 2019-01-26 09:43:57.377055 +0100 CET m=+0.644739467)
	TimeFormatGo = "2006-01-02 15:04:05.999999999 -0700 MST"
	// TimeFormatSimple handles "2019-01-25 21:51:38"
	TimeFormatSimple = "2006-01-02 15:04:05.999999999"
	// TimeFormatHTTP instead of importing main with http.TimeFormat which would increase the binary size significantly.
	TimeFormatHTTP = "Mon, 02 Jan 2006 15:04:05 GMT"
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
	if t, err := time.Parse(TimeFormatHTTP, input); err == nil {
		return t, TimeFormatHTTP, nil
	}

	if t, err := time.Parse(TimeFormatGo, strings.Split(input, " m=")[0]); err == nil {
		return t, TimeFormatGo, nil
	}

	// "2019-01-25 21:51:38"
	if t, err := time.Parse(TimeFormatSimple, input); err == nil {
		return t, TimeFormatSimple, nil
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
	case "+": //, "add", "plus":
		return Add, nil
	case "-": //, "sub", "minus":
		return Sub, nil
	}
	return Undefined, fmt.Errorf("%w: '%v'", ErrUnkownOperator, s)
}

// Calculate does basic add/sub calculations on the given input.
func Calculate(input time.Time, op Operator, amount int, unit string) time.Time {
	switch op {
	case Sub:
		amount = -amount
	}

	var duration time.Duration

	switch unit {
	case "ns":
		duration = time.Duration(amount) * time.Nanosecond
		return input.Add(duration)
	case "us":
		duration = time.Duration(amount) * time.Microsecond
		return input.Add(duration)
	case "ms":
		duration = time.Duration(amount) * time.Millisecond
		return input.Add(duration)
	case "s":
		duration = time.Duration(amount) * time.Second
		return input.Add(duration)
	case "m":
		duration = time.Duration(amount) * time.Minute
		return input.Add(duration)
	case "h":
		duration = time.Duration(amount) * time.Hour
		return input.Add(duration)
	case "D":
		return input.AddDate(0, 0, amount)
	case "W":
		return input.AddDate(0, 0, amount*7)
	case "M":
		return input.AddDate(0, amount, 0)
	case "Y":
		return input.AddDate(amount, 0, 0)
	}

	return time.Time{}
}

// FormatName returns the formatting to the given name (e.g. 'unix' or 'rfc3339').
// When 'format' is not recognized, it will return Go's default format and an error.
func FormatName(format string) (string, error) {
	format = strings.ToLower(format)

	switch format {
	case "":
		return TimeFormatGo, nil
	case "unix":
		return time.UnixDate, nil
	case "ruby":
		return time.RubyDate, nil
	case "ansic":
		return time.ANSIC, nil
	case "rfc822":
		return time.RFC822, nil
	case "rfc822z":
		return time.RFC822Z, nil
	case "rfc850":
		return time.RFC850, nil
	case "rfc1123":
		return time.RFC1123, nil
	case "rfc1123z":
		return time.RFC1123Z, nil
	case "rfc3339":
		return time.RFC3339, nil
	case "rfc3339nano":
		return time.RFC3339Nano, nil
	case "kitchen":
		return time.Kitchen, nil
	case "stamp":
		return time.Stamp, nil
	case "stampmilli":
		return time.StampMilli, nil
	case "stampmicro":
		return time.StampMicro, nil
	case "stampnano":
		return time.StampNano, nil
	case "http":
		return TimeFormatHTTP, nil
	default:
		return TimeFormatGo, fmt.Errorf("failed to parse format %q", format)
	}
}

// FormatSimple converts a simple format to Go's native formatting.
func FormatSimple(format string) string {
	format = strings.ReplaceAll(format, "{YYYY}", "2006") // Long year
	format = strings.ReplaceAll(format, "{YY}", "06")     // Short year

	format = strings.ReplaceAll(format, "{MMMM}", "January")
	format = strings.ReplaceAll(format, "{MMM}", "Jan")
	format = strings.ReplaceAll(format, "{MM}", "01") // Month (2-digit)
	format = strings.ReplaceAll(format, "{M}", "1")   // Month (1-digit)

	format = strings.ReplaceAll(format, "{DDDD}", "002") // Day of year
	format = strings.ReplaceAll(format, "{DDDD}", "__2") // Day of year

	format = strings.ReplaceAll(format, "{DD}", "02") // Day of month (2-digit)
	format = strings.ReplaceAll(format, "{D}", "2")   // Day of month (1-digit)

	format = strings.ReplaceAll(format, "{dddd}", "Monday") // Day of week
	format = strings.ReplaceAll(format, "{ddd}", "Mon")     // Day of week

	format = strings.ReplaceAll(format, "{HH}", "15") // Hour 24 (2-digit)
	format = strings.ReplaceAll(format, "{hh}", "03") // Hour 12 (2-digit)
	format = strings.ReplaceAll(format, "{h}", "3")   // Hour 12 (1-digit)

	format = strings.ReplaceAll(format, "{A}", "PM")
	format = strings.ReplaceAll(format, "{a}", "pm")

	format = strings.ReplaceAll(format, "{mm}", "04") // Minute (2-digit)
	format = strings.ReplaceAll(format, "{m}", "4")   // Minute (1-digit)

	format = strings.ReplaceAll(format, "{ss}", "05") // Second (2-digit)
	format = strings.ReplaceAll(format, "{s}", "5")   // Second (1-digit)

	// Arbitrary precision of fractional seconds.
	// The dot won't be removed from the output :-/
	format = strings.ReplaceAll(format, "{F}", ".0")
	format = strings.ReplaceAll(format, "{FF}", ".00")
	format = strings.ReplaceAll(format, "{FFF}", ".000")
	format = strings.ReplaceAll(format, "{FFFF}", ".0000")
	format = strings.ReplaceAll(format, "{FFFFF}", ".00000")
	format = strings.ReplaceAll(format, "{FFFFFF}", ".000000")
	format = strings.ReplaceAll(format, "{FFFFFFF}", ".0000000")
	format = strings.ReplaceAll(format, "{FFFFFFFF}", ".00000000")
	format = strings.ReplaceAll(format, "{FFFFFFFFF}", ".000000000")

	format = strings.ReplaceAll(format, "{f}", ".9")
	format = strings.ReplaceAll(format, "{ff}", ".99")
	format = strings.ReplaceAll(format, "{fff}", ".999")
	format = strings.ReplaceAll(format, "{ffff}", ".9999")
	format = strings.ReplaceAll(format, "{fffff}", ".99999")
	format = strings.ReplaceAll(format, "{ffffff}", ".999999")
	format = strings.ReplaceAll(format, "{fffffff}", ".9999999")
	format = strings.ReplaceAll(format, "{ffffffff}", ".99999999")
	format = strings.ReplaceAll(format, "{fffffffff}", ".999999999")

	// Timezone / Offset
	format = strings.ReplaceAll(format, "{ZZZ}", "-07:00")
	format = strings.ReplaceAll(format, "{ZZ}", "-0700")
	format = strings.ReplaceAll(format, "{Z}", "-07")
	format = strings.ReplaceAll(format, "{z}", "MST")

	return format
}
