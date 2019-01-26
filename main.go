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

func main() {
	var (
		input string
		nano  = flag.Bool("nsec", false, "use nanoseconds instead of seconds")
		local = flag.Bool("local", false, "use local time instead of UTC")
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
		printHuman(toHuman(i, guess), *local)
		return
	}

	// output unix timestamp
	if t, err := toUnix(input); err == nil {
		printEpoch(t, *nano)
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

func printEpoch(t time.Time, nano bool) {
	epoch := t.Unix()
	if nano {
		epoch = t.UnixNano()
	}

	fmt.Println(epoch)
}

func printHuman(t time.Time, local bool) {
	// TODO: local not working
	if local {
		t = t.Local()
	}

	fmt.Println(t)
}

type tsType byte

const (
	guess tsType = iota
	sec
	nsec
)

func toHuman(timestamp int64, typ tsType) time.Time {
	switch typ {
	case sec:
		return time.Unix(timestamp, 0)
	case nsec:
		return time.Unix(0, timestamp)
	}

	// neither seconds nor nanoseconds, guess type
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

func toUnix(input string) (time.Time, error) {
	if t, err := time.Parse(time.RFC1123, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC1123Z, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC3339Nano, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC822, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC822Z, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RFC850, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.ANSIC, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.UnixDate, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.RubyDate, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.Kitchen, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.Stamp, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.StampMicro, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.StampMilli, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(time.StampNano, input); err == nil {
		return t, nil
	}
	if t, err := time.Parse(http.TimeFormat, input); err == nil {
		return t, nil
	}

	// handle Go's default time.Now() format (e.g. 2019-01-26 09:43:57.377055 +0100 CET m=+0.644739467)
	const defaultNow = "2006-01-02 15:04:05.999999999 -0700 MST"
	if t, err := time.Parse(defaultNow, strings.Split(input, " m=")[0]); err == nil {
		return t, nil
	}

	return time.Time{}, errors.New("failed to convert string to time")
}
