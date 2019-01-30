package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type timeUnit byte

const (
	unitSeconds timeUnit = iota
	unitMilliseconds
	unitMicroseconds
	unitNanoseconds
)

func toUnit(input string) (timeUnit, error) {
	switch input {
	case "s":
		return unitSeconds, nil
	case "ms":
		return unitMilliseconds, nil
	case "us":
		return unitMicroseconds, nil
	case "ns":
		return unitNanoseconds, nil
	}
	return 255, fmt.Errorf("failed to convert %v to time unit", input)
}

func main() {
	var (
		input      string
		unitFlag   = flag.String("unit", "s", "unit for timestamp output: s, ms, us, ns")
		formatFlag = flag.String("format", "UnixDate", "TODO")
		localFal   = flag.Bool("local", false, "use local time instead of UTC")
	)
	flag.Parse()

	unit, err := toUnit(*unitFlag)
	if err != nil {
		log.Fatal(err)
	}

	// read program input
	if flag.NArg() == 0 { // from stdin/pipe
		reader := bufio.NewReader(os.Stdin)
		var err error
		input, err = reader.ReadString('\n')
		if err != nil {
			log.Fatalln("failed to read input")
		}
		input = strings.TrimSpace(input)
	} else { // from argument
		if flag.NArg() > 1 {
			log.Fatalln("takes at most one input")
		}
		input = flag.Arg(0)
	}

	// if the input can be parsed as an int, we assume it's an epoch timestamp
	if i, err := strconv.ParseInt(input, 10, 64); err == nil {
		printFormatted(fromTimestamp(i), *formatFlag, *localFal)
		return
	}

	// output unix timestamp
	if t, err := fromFormatted(input); err == nil {
		printTimestamp(t, unit)
		return
	}

	log.Fatalln("failed to convert input")
}

func abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

func printTimestamp(t time.Time, unit timeUnit) {
	epoch := t.Unix()

	switch unit {
	case unitSeconds:
		// calculated as default value,
		// nothing to to
	case unitMilliseconds:
		epoch = t.UnixNano() / (1000 * 1000)
	case unitMicroseconds:
		epoch = t.UnixNano() / 1000
	case unitNanoseconds:
		epoch = t.UnixNano()
	default:
		panic("forgot to add a unit")
	}

	fmt.Println(epoch)
}

func printFormatted(t time.Time, format string, local bool) {
	// TODO: local not working
	if local {
		t = t.Local()
	}

	// TODO: use format parameter
	fmt.Println(t.Format(time.UnixDate))
}

func fromTimestamp(timestamp int64) time.Time {
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
		return time.Unix(timestamp, 0)
	} else {
		// number of digits is closer to current nanoseconds timestamp
		return time.Unix(0, timestamp)
	}
}

func fromFormatted(input string) (time.Time, error) {
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
