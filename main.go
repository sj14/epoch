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

	"github.com/sj14/epoch/epoch"
)

func main() {
	var (
		input      string
		unitFlag   = flag.String("unit", "guess", "unit for timestamp output: s, ms, us, ns")
		formatFlag = flag.String("format", "UnixDate", "TODO")
		localFal   = flag.Bool("local", false, "use local time instead of UTC")
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

	// if the input can be parsed as an int, we assume it's an epoch timestamp
	if i, err := strconv.ParseInt(input, 10, 64); err == nil {
		unit, err := epoch.ParseUnit(*unitFlag)
		if err != nil {
			unit = epoch.GuessUnit(i)
			switch unit {
			case epoch.UnitSeconds:
				fmt.Println("guessed unit seconds")
			case epoch.UnitMilliseconds:
				fmt.Println("guessed unit milliseconds")
			case epoch.UnitMicroseconds:
				fmt.Println("guessed unit microseconds")
			case epoch.UnitNanoseconds:
				fmt.Println("guessed unit nanoseconds")
			}
		}

		t, err := epoch.ParseTimestamp(i, unit)
		if err != nil {
			log.Fatalf("failed to convert from timestamp: %v", err)
		}
		printFormatted(t, *formatFlag, *localFal)
		return
	}

	// output unix timestamp

	// convert fromatted string to time type
	t, err := epoch.ParseFormatted(input)
	if err != nil {
		log.Fatalf("failed to convert input: %v", err)
	}

	unit, err := epoch.ParseUnit(*unitFlag)
	if err != nil {
		// use seconds as default unit
		fmt.Println("using seconds as unit")
		unit = epoch.UnitSeconds
	}

	// convert time to timestamp
	timestamp, err := epoch.ToTimestamp(t, unit)
	if err != nil {
		log.Fatalf("failed to convert timestamp: %v", err)
	}
	fmt.Println(timestamp)
}

func printFormatted(t time.Time, format string, local bool) {
	// TODO: local not working
	if local {
		t = t.Local()
	}

	// TODO: use format parameter
	fmt.Println(t.Format(time.UnixDate))
}
