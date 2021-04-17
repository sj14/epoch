# epoch

![Action](https://github.com/sj14/epoch/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/sj14/epoch)](https://goreportcard.com/report/github.com/sj14/epoch)
[![GoDoc](https://godoc.org/github.com/sj14/epoch/epoch?status.png)](https://godoc.org/github.com/sj14/epoch/epoch)

`epoch` converts unix timestamps to human readable formats and vice versa.

**Why?**  
To convert timestamps to dates, you have to run different commands for different `date` implementations, such as GNU's, BSD's or Busybox's date implementation. For example, `date -d @1267619929` (GNU) vs `date -r 1267619929` (BSD), and what about handling nanosecond timestamps? Furthermore, have you ever tried converting a time formatted string such as `"2019-01-25 21:51:38 +0100 CET"` to a timestamp? Of course, you can do all this somehow, but all ways I've found so far were too cumbersome. This tool tries to solve all this with ease:

```bash
$ epoch "2019-01-25 21:51:38 +0100 CET"
using seconds as unit
1548449498
```

```bash
$ epoch 1548449498
guessed unit: seconds
2019-01-25 21:51:38 +0100 CET
```

Convert between timezones:

```bash
$ epoch -tz "Europe/Berlin" "Sat Jul 18 15:46:45 UTC 2020"
Sat Jul 18 17:46:45 CEST 2020
```

Timestamp to formatted string of specific timezone:

```bash
$ epoch -tz "America/New_York" 1595088886
guessed unit: seconds
2020-07-18 12:14:46 -0400 EDT
```

The functionallity is implemented as a package and can be used in other programs.

## Installation

### Precompiled Binaries

Binaries are available for all major platforms. See the [releases](https://github.com/sj14/epoch/releases) page. Usually, `epoch` uses the timezone data from the operating system. When the operating system has no timezone data installed, you can use the 'full' binaries which have this information embedded.

### Homebrew

Using the [Homebrew](https://brew.sh/) package manager for macOS:

```bash
brew install sj14/tap/epoch
```

### Manually

It's also possible to install the current development snapshot with `go get`:

```bash
go get -u github.com/sj14/epoch/cmd/epoch
```

## Usage

```text
Usage of epoch:
  -format string
        human readable output format, see readme for details
  -quiet
        don't output guessed units
  -tz string
        the timezone to use, e.g. 'Local', 'UTC', or a name corresponding to the IANA Time Zone database, such as 'America/New_York' (default "Local")
  -unit string
        unit for timestamps: s, ms, us, ns (default "guess")
```

## Examples

### Timestamps to human readable format

#### conversion based on given unit

Use the `-unit` flag or append the unit at a suffix to the input.

```bash
$ epoch -unit s 1548449513
Fri Jan 25 21:51:53 CET 2019
```

```bash
$ epoch 1548449513s
Fri Jan 25 21:51:53 CET 2019
```

---

```bash
$ epoch -unit ms 1548449513
Sun Jan 18 23:07:29 CET 1970
```

```bash
$ epoch 1548449513ms
Sun Jan 18 23:07:29 CET 1970
```

---

```bash
$ epoch -unit us 1548449513
Thu Jan  1 01:25:48 CET 1970
```

```bash
$ epoch 1548449513us
Thu Jan  1 01:25:48 CET 1970
```

---

```bash
$ epoch -unit ns 1548449513
Thu Jan  1 01:00:01 CET 1970
```

```bash
$ epoch 1548449513ns
Thu Jan  1 01:00:01 CET 1970
```

#### set the output format

```bash
$ epoch -unit ms -format rfc850 1548449513
Sunday, 18-Jan-70 23:07:29 CET
```

```bash
$ epoch -unit ms -format ruby 1548449513
Sun Jan 18 23:07:29 +0100 1970
```

```bash
$ epoch -unit ms -format ansic 1548449513
Sun Jan 18 23:07:29 1970
```

#### guess the unit

Guess the unit. Internally, the guess is done by comparing the number of digits with the current epoch timestamps (in `s`, `ms`, `us`, `ns`) of your machine. The smallest difference wins.

seconds:

```bash
$ epoch 1548449513
guessed unit seconds
2019-01-25 21:51:53 +0100 CET
```

milliseconds:

```bash
$ epoch 1548449513940
guessed unit milliseconds
2019-01-25 21:51:53 +0100 CET
```

microseconds:

```bash
$ epoch 1548449513940562
guessed unit microseconds
2019-01-25 21:51:53.940562 +0100 CET
```

nanoseconds:

```bash
$ epoch 1548449513940562000
guessed unit nanoseconds
2019-01-25 21:51:53.940562 +0100 CET
```

#### negative timestamp

```bash
$ epoch -- -15484495
guessed unit: seconds
1969-07-05 19:45:05 +0100 CET
```

#### using the pipe

```bash
$ echo -15484495 | epoch
guessed unit: seconds
1969-07-05 19:45:05 +0100 CET
```

### Formatted input to epoch timestamps

seconds (default when no `unit` flag given):

```bash
$ epoch -unit s "2019-01-25 21:51:38.272173 +0100 CET"
1548449498
```

milliseconds:

```bash
$ epoch -unit ms "2019-01-25 21:51:38.272173 +0100 CET"
1548449498272
```

microseconds:

```bash
$ epoch -unit us "2019-01-25 21:51:38.272173 +0100 CET"
1548449498272173
```

nanoseconds:

```bash
$ epoch -unit ns "2019-01-25 21:51:38.272173 +0100 CET"
1548449498272173000
```

## Supported Formats

All current Go formats as of 2019-01-26 (https://golang.org/pkg/time/#pkg-constants):

```go
ANSIC       = "Mon Jan _2 15:04:05 2006"
UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
RFC822      = "02 Jan 06 15:04 MST"
RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
RFC3339     = "2006-01-02T15:04:05Z07:00"
RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
Kitchen     = "3:04PM"
// Handy time stamps.
Stamp      = "Jan _2 15:04:05"
StampMilli = "Jan _2 15:04:05.000"
StampMicro = "Jan _2 15:04:05.000000"
StampNano  = "Jan _2 15:04:05.000000000"
// HTTP Timestamp time.RFC1123 but hard-codes GMT as the time zone.
HTTP = "Mon, 02 Jan 2006 15:04:05 GMT"
```
