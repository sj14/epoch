package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sj14/epoch/epoch"
)

func main() {
	var (
		input      string
		unitFlag   = flag.String("unit", "guess", "unit for timestamp output: s, ms, us, ns")
		formatFlag = flag.String("format", "", "human readable output format, see readme for details")
		utcFlag    = flag.Bool("utc", false, "use UTC instead of local zone")
		quieteFlag = flag.Bool("quiete", false, "don't ouput guessed units")
	)
	flag.Parse()

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

	if input == "" {
		input = time.Now().String()
	}

	// if the input can be parsed as an int, we assume it's an epoch timestamp
	if i, err := strconv.ParseInt(input, 10, 64); err == nil {
		t := outputTimestamp(*unitFlag, i, *quieteFlag)
		printFormatted(t, *formatFlag, *utcFlag)
		return
	}

	// output formatted time
	outputFormatted(input, *unitFlag, *quieteFlag)
}

func outputFormatted(input, unitFlag string, quieteFlag bool) {
	// convert fromatted string to time type
	t, err := epoch.ParseFormatted(input)
	if err != nil {
		log.Fatalf("failed to convert input: %v", err)
	}

	unit, err := epoch.ParseUnit(unitFlag)
	if err != nil {
		// use seconds as default unit
		unit = epoch.UnitSeconds
		if !quieteFlag {
			fmt.Println("using seconds as unit")
		}
	}

	// convert time to timestamp
	timestamp, err := epoch.ToTimestamp(t, unit)
	if err != nil {
		log.Fatalf("failed to convert timestamp: %v", err)
	}
	fmt.Println(timestamp)
}

func outputTimestamp(unitFlag string, i int64, quieteFlag bool) time.Time {
	unit, err := epoch.ParseUnit(unitFlag)
	if err != nil {
		unit = epoch.GuessUnit(i, time.Now())

		if !quieteFlag {
			switch unit {
			case epoch.UnitSeconds:
				fmt.Fprintln(os.Stderr, "guessed unit: seconds")
			case epoch.UnitMilliseconds:
				fmt.Fprintln(os.Stderr, "guessed unit: milliseconds")
			case epoch.UnitMicroseconds:
				fmt.Fprintln(os.Stderr, "guessed unit: microseconds")
			case epoch.UnitNanoseconds:
				fmt.Fprintln(os.Stderr, "guessed unit: nanoseconds")
			}
		}
	}

	t, err := epoch.ParseTimestamp(i, unit)
	if err != nil {
		log.Fatalf("failed to convert from timestamp: %v", err)
	}
	return t
}

func printFormatted(t time.Time, format string, utc bool) {
	if utc {
		t = t.In(time.UTC)
	}

	format = strings.ToLower(format)

	switch format {
	case "":
		fmt.Println(t)
	case "unix":
		fmt.Println(t.Format(time.UnixDate))
	case "ruby":
		fmt.Println(t.Format(time.RubyDate))
	case "ansic":
		fmt.Println(t.Format(time.ANSIC))
	case "rfc822":
		fmt.Println(t.Format(time.RFC822))
	case "rfc822z":
		fmt.Println(t.Format(time.RFC822Z))
	case "rfc850":
		fmt.Println(t.Format(time.RFC850))
	case "rfc1123":
		fmt.Println(t.Format(time.RFC1123))
	case "rfc1123z":
		fmt.Println(t.Format(time.RFC1123Z))
	case "rfc3339":
		fmt.Println(t.Format(time.RFC3339))
	case "rfc3339nano":
		fmt.Println(t.Format(time.RFC3339Nano))
	case "kitchen":
		fmt.Println(t.Format(time.Kitchen))
	case "stamp":
		fmt.Println(t.Format(time.Stamp))
	case "stampms":
		fmt.Println(t.Format(time.StampMilli))
	case "stampus":
		fmt.Println(t.Format(time.StampMicro))
	case "stampns":
		fmt.Println(t.Format(time.StampNano))
	case "http":
		fmt.Println(t.Format(http.TimeFormat))
	default:
		fmt.Fprintf(os.Stderr, "failed to parse format '%v'\n", format)
		fmt.Println(t)
	}
}
