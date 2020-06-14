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
		unitFlag   = flag.String("unit", "guess", "unit for timestamps: s, ms, us, ns")
		formatFlag = flag.String("format", "", "human readable output format, see readme for details")
		tzFlag     = flag.String("tz", "Local", "the timezone to use, e.g. 'Local', 'UTC', or a name corresponding to the IANA Time Zone database, such as 'America/New_York'")
		quietFlag  = flag.Bool("quiet", false, "don't output guessed units")
	)
	flag.Parse()

	input := readInput()

	if input == "" {
		if *tzFlag != "Local" {
			log.Fatalln("can't use empty input with specific timezone")
		}
		input = time.Now().String()
	}

	// use suffix of input as unit, e.g.
	// "1234567890s" -> unit: "s"; input: "1234567890"
	//
	// keep "s" as last element in slice, otherwise,
	// it will match all other units as they end with an "s", too.
	for _, unit := range []string{"ns", "us", "ms", "s"} {
		if !strings.HasSuffix(input, unit) {
			continue
		}

		// check if remaining input is an integer, if not
		// it might be a time zone ending with 'unit'.
		// (I'm currently not aware of any, but let's be sure)
		inputTrim := strings.TrimSuffix(input, unit)
		if _, err := strconv.ParseInt(inputTrim, 10, 64); err != nil {
			continue
		}

		if *unitFlag != "guess" && *unitFlag != unit {
			log.Fatalf("mismatch between unit flag (%v) and input unit (%v)\n", *unitFlag, unit)
		}
		*unitFlag = unit
		input = inputTrim
		break
	}

	// if the input can be parsed as an int, we assume it's an epoch timestamp
	if i, err := strconv.ParseInt(input, 10, 64); err == nil {
		t := outputTimestamp(*unitFlag, i, *quietFlag)
		tz := *tzFlag

		if strings.ToLower(tz) == "local" {
			tz = "Local" // capital is important
		}

		loc, err := time.LoadLocation(tz)
		if err != nil {
			log.Fatalf("failed loading timezone '%v': %v\n", tz, err)
		}

		t = t.In(loc)
		printFormatted(t, *formatFlag)
		return
	}

	// likely not an epoch timestamp, output formatted time
	outputFormatted(input, *unitFlag, *quietFlag)
}

// read program input from stdin or argument
func readInput() string {
	// from stdin/pipe
	if flag.NArg() == 0 {

		// check if it's piped or from empty stdin
		// https://stackoverflow.com/a/26567513
		stat, err := os.Stdin.Stat()
		if err != nil {
			log.Fatalf("failed to get stdin stats: %v\n", err)
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return ""
		}

		// read the input from the pipe
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("failed to read input: %v\n", err)
		}
		return strings.TrimSpace(input)
	}

	// from argument
	if flag.NArg() > 1 {
		log.Fatalln("takes at most one input")
	}
	return flag.Arg(0)
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

func printFormatted(t time.Time, format string) {
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
