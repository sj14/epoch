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

	"github.com/sj14/epoch"
)

var version = "undefined" // will be replaced during the build process

func main() {
	var (
		unit        = flag.String("unit", "guess", "unit for timestamps: s, ms, us, ns")
		format      = flag.String("format", "", "human readable output format, such as 'rfc3339' (see readme for details)")
		tz          = flag.String("tz", "", `the timezone to use, e.g. 'Local' (default), 'UTC', or a name corresponding to the IANA Time Zone database, such as 'America/New_York'`)
		quiet       = flag.Bool("quiet", false, "don't output guessed units")
		versionFlag = flag.Bool("version", false, fmt.Sprintf("print version (%v)", version))
		calc        = flag.String("calc", "", "apply basic time calculations, e.g. '+30m -5h +3M -10Y'")
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
	result, err := run(input, time.Now().String(), *calc, *unit, *format, *tz, *quiet)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
}

type calculation struct {
	operator epoch.Operator
	amount   int
	unit     string
}

func run(input string, now, calc string, unit, format, tz string, quiet bool) (string, error) {
	var (
		err          error
		calculations []calculation
	)

	calcInpufStrings := strings.Split(calc, " ")
	for _, calcInputString := range calcInpufStrings {
		if strings.TrimSpace(calcInputString) == "" {
			continue
		}
		operator, err := epoch.ToOperator(calcInputString[:1])
		if err != nil {
			return "", err
		}

		amount, err := strconv.Atoi(calcInputString[1 : len(calcInputString)-1])
		if err != nil {
			return "", err
		}

		suffix := calcInputString[len(calcInputString)-1:]
		calculations = append(calculations, calculation{operator: operator, amount: amount, unit: suffix})
	}

	if input == "" {
		input = now
	}

	input, unit, err = parseUnit(input, unit)
	if err != nil {
		return "", err
	}

	// If the input can be parsed as a number, we assume it's an epoch timestamp. Convert to formatted string.
	if i, err := strconv.ParseFloat(input, 64); err == nil {
		t := parseTimestamp(unit, int64(i), quiet).In(location(tz))

		if len(calculations) > 0 {
			// when applying arithmetics here, return as timestamp again
			for _, calc := range calculations {
				t = epoch.Calculate(t, calc.operator, calc.amount, calc.unit)
			}
			// always quite as we already output unit above in parseTimestmap
			return strconv.FormatInt(timestamp(t, unit, true), 10), nil
		}

		return epoch.FormattedString(t, format)
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
		t = t.In(location(tz))

		for _, calc := range calculations {
			t = epoch.Calculate(t, calc.operator, calc.amount, calc.unit)
		}

		return epoch.FormattedString(t, format)
	}

	// Likely not an epoch timestamp as input, output formatted input time to timestamp.
	if format != "" {
		return "", fmt.Errorf("can't use specific format when converting to timestamp (omit -format flag)")
	}

	// convert formatted string to time type
	t, _, err := epoch.ParseFormatted(input)
	if err != nil {
		log.Fatalf("failed to convert input: %v", err)
	}
	t = t.In(location(tz))

	for _, calc := range calculations {
		t = epoch.Calculate(t, calc.operator, calc.amount, calc.unit)
	}

	return strconv.FormatInt(timestamp(t, unit, quiet), 10), nil
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

func timestamp(t time.Time, unitFlag string, quiete bool) int64 {
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
