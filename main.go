package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"

	"github.com/sj14/epoch/epoch"
)

var version = "undefined" // will be replaced during the build process

func main() {
	var (
		unit        = flag.String("unit", "guess", "unit for timestamps: s, ms, us, ns")
		format      = flag.String("format", "", "human readable output format, see readme for details")
		tz          = flag.String("tz", "", `the timezone to use, e.g. 'Local', 'UTC', or a name corresponding to the IANA Time Zone database, such as 'America/New_York' (default "Local")`)
		quiet       = flag.Bool("quiet", false, "don't output guessed units")
		versionFlag = flag.Bool("version", false, fmt.Sprintf("print version (%v)", version))
	)
	flag.Parse()

	if *versionFlag {
		fmt.Println(version)
		os.Exit(0)
	}

	input, err := readInput()
	if err != nil {
		log.Fatalln(err)
	}
	result, err := run(input, time.Now().String(), *unit, *format, *tz, *quiet)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
}

func run(input, now, unit, format, tz string, quiet bool) (string, error) {
	if input == "" {
		input = now
	}

	input, unit, err := parseUnit(input, unit)
	if err != nil {
		return "", err
	}

	// If the input can be parsed as an int, we assume it's an epoch timestamp. Convert to formatted string.
	if i, err := strconv.ParseInt(input, 10, 64); err == nil {
		t := parseTimestamp(unit, i, quiet)
		return formattedString(t, format, tz), nil
	}

	// Likely not an epoch timestamp as input. But a timezone and/or format was specified. Convert formatted input to another timezone and/or format.
	if tz != "" || format != "" {
		if unit != "guess" {
			return "", fmt.Errorf("can't use unit flag together with timezone or format flag on a formatted string (omit -unit flag)")
		}

		t, _, err := epoch.ParseFormatted(input)
		if err != nil {
			return "", fmt.Errorf("failed to convert input: %v", err)
		}

		return formattedString(t, format, tz), nil
	}

	// Likely not an epoch timestamp as input, output formatted input time to timestamp.
	if format != "" {
		return "", fmt.Errorf("can't use specific format when converting to timestamp (omit -format flag)")
	}

	return strconv.FormatInt(timestamp(input, unit, quiet), 10), nil
}

// read program input from stdin or argument
func readInput() (string, error) {
	// from stdin/pipe
	if flag.NArg() == 0 {

		// check if it's piped or from empty stdin
		// https://stackoverflow.com/a/26567513
		stat, err := os.Stdin.Stat()
		if err != nil {
			return "", fmt.Errorf("failed to get stdin stats: %v", err)
		}
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			return "", nil
		}

		// read the input from the pipe
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %v", err)
		}
		return strings.TrimSpace(input), nil
	}

	// from argument
	if flag.NArg() > 1 {
		return "", fmt.Errorf("takes at most one input")
	}
	return flag.Arg(0), nil
}

func parseUnit(input, unitFlag string) (string, string, error) {
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

		if unitFlag != "guess" && unitFlag != unit {
			return "", "", fmt.Errorf("mismatch between unit flag (%v) and input unit (%v)", unitFlag, unit)
		}

		return inputTrim, unit, nil
	}
	return input, unitFlag, nil
}

func location(tz string) *time.Location {
	if strings.ToLower(tz) == "local" || tz == "" {
		tz = "Local" // capital is important
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatalf("failed loading timezone '%v': %v\n", tz, err)
	}
	return loc
}

func timestamp(input, unitFlag string, quiete bool) int64 {
	// convert fromatted string to time type
	t, _, err := epoch.ParseFormatted(input)
	if err != nil {
		log.Fatalf("failed to convert input: %v", err)
	}

	unit, err := epoch.ParseUnit(unitFlag)
	if err != nil {
		// use seconds as default unit
		unit = epoch.UnitSeconds
		if !quiete {
			fmt.Println("using seconds as unit")
		}
	}

	// convert time to timestamp
	timestamp, err := epoch.ToTimestamp(t, unit)
	if err != nil {
		log.Fatalf("failed to convert timestamp: %v", err)
	}
	return timestamp
}

func parseTimestamp(unitFlag string, i int64, quiete bool) time.Time {
	unit, err := epoch.ParseUnit(unitFlag)
	if err != nil {
		unit = epoch.GuessUnit(i, time.Now())

		if !quiete {
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

func formattedString(t time.Time, format, tz string) string {
	t = t.In(location(tz))

	format = strings.ToLower(format)

	switch format {
	case "":
		return t.String()
	case "unix":
		return t.Format(time.UnixDate)
	case "ruby":
		return t.Format(time.RubyDate)
	case "ansic":
		return t.Format(time.ANSIC)
	case "rfc822":
		return t.Format(time.RFC822)
	case "rfc822z":
		return t.Format(time.RFC822Z)
	case "rfc850":
		return t.Format(time.RFC850)
	case "rfc1123":
		return t.Format(time.RFC1123)
	case "rfc1123z":
		return t.Format(time.RFC1123Z)
	case "rfc3339":
		return t.Format(time.RFC3339)
	case "rfc3339nano":
		return t.Format(time.RFC3339Nano)
	case "kitchen":
		return t.Format(time.Kitchen)
	case "stamp":
		return t.Format(time.Stamp)
	case "stampms":
		return t.Format(time.StampMilli)
	case "stampus":
		return t.Format(time.StampMicro)
	case "stampns":
		return t.Format(time.StampNano)
	case "http":
		return t.Format(epoch.FormatHTTP)
	default:
		fmt.Fprintf(os.Stderr, "failed to parse format '%v'\n", format)
		return t.String()
	}
}
