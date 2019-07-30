# epoch

[![Build Status](https://travis-ci.org/sj14/epoch.svg?branch=master)](https://travis-ci.org/sj14/epoch)
[![Go Report Card](https://goreportcard.com/badge/github.com/sj14/epoch)](https://goreportcard.com/report/github.com/sj14/epoch)
[![GoDoc](https://godoc.org/github.com/sj14/epoch/epoch?status.png)](https://godoc.org/github.com/sj14/epoch/epoch)

`epoch` converts unix timestamps to human readable formats and vice versa.

**Why?**  
To convert timestamps to dates, you have to run different commands for Linux and macOS: `date -d @1267619929` vs `date -r 1267619929`, and what about handling nanosecond timestamps? Seriously, I don't know how to do this with `date`.  
Furthermore, have you ever tried converting a time formatted string such as `"2019-01-25 21:51:38 +0100 CET"` to a timestamp? Of course, you can do this somehow, but all ways I've found so far were too cumbersome. This tool tries to solve all this with ease:

```text
$ epoch "2019-01-25 21:51:38 +0100 CET"
using seconds as unit
1548449498
```

```text
$ epoch 1548449498
guessed unit: seconds
2019-01-25 21:51:38 +0100 CET
```

The functionallity is implemented as a package and can be used in other programs.

## Installation

### Precompiled Binaries

Binaries are available for all major platforms. See the [releases](https://github.com/sj14/epoch/releases) page.

### Homebrew

Using the [Homebrew](https://brew.sh/) package manager for macOS:

``` text
brew install sj14/tap/epoch
```

### Manually

It's also possible to install the current development snapshot with `go get`:

``` text
go get -u github.com/sj14/epoch
```

## Usage

```text
Usage of epoch:
  -format string
        human readable output format, see readme for details
  -quiet
        don't output guessed units
  -unit string
        unit for timestamp output: s, ms, us, ns (default "guess")
  -utc
        use UTC instead of local zone
```

## Examples

### Timestamps to human readable format

#### conversion based on given unit

Use the `-unit` flag or append the unit at a suffix to the input.

``` text
$ epoch -unit s 1548449513
Fri Jan 25 21:51:53 CET 2019
```

``` text
$ epoch 1548449513s
Fri Jan 25 21:51:53 CET 2019
```

---

``` text
$ epoch -unit ms 1548449513
Sun Jan 18 23:07:29 CET 1970
```

``` text
$ epoch 1548449513ms
Sun Jan 18 23:07:29 CET 1970
```

---

``` text
$ epoch -unit us 1548449513
Thu Jan  1 01:25:48 CET 1970
```

``` text
$ epoch 1548449513us
Thu Jan  1 01:25:48 CET 1970
```

---

``` text
$ epoch -unit ns 1548449513
Thu Jan  1 01:00:01 CET 1970
```

``` text
$ epoch 1548449513ns
Thu Jan  1 01:00:01 CET 1970
```

#### set the output format

``` text
$ epoch -unit ms -format rfc850 1548449513
Sunday, 18-Jan-70 23:07:29 CET
```

``` text
$ epoch -unit ms -format ruby 1548449513
Sun Jan 18 23:07:29 +0100 1970
```

``` text
$ epoch -unit ms -format ansic 1548449513
Sun Jan 18 23:07:29 1970
```

#### guess the unit

Guess the unit. Internally, the guess is done by comparing the number of digits with the current epoch timestamps (in `s`, `ms`, `us`, `ns`) of your machine. The smallest difference wins.

seconds:

``` text
$ epoch 1548449513
guessed unit seconds
2019-01-25 21:51:53 +0100 CET
```

milliseconds:

``` text
$ epoch 1548449513940
guessed unit milliseconds
2019-01-25 21:51:53 +0100 CET
```

microseconds:

``` text
$ epoch 1548449513940562
guessed unit microseconds
2019-01-25 21:51:53.940562 +0100 CET
```

nanoseconds:

``` text
$ epoch 1548449513940562000
guessed unit nanoseconds
2019-01-25 21:51:53.940562 +0100 CET
```

#### negative timestamp

``` text
$ epoch -- -15484495
guessed unit: seconds
1969-07-05 19:45:05 +0100 CET
```

#### using the pipe

``` text
$ echo -15484495 | epoch
guessed unit: seconds
1969-07-05 19:45:05 +0100 CET
```

### Formatted input to epoch timestamps

seconds (default when no `unit` flag given):

``` text
$ epoch -unit s "2019-01-25 21:51:38.272173 +0100 CET"
1548449498
```

milliseconds:

```text
$ epoch -unit ms "2019-01-25 21:51:38.272173 +0100 CET"
1548449498272
```

microseconds:

```text
$ epoch -unit us "2019-01-25 21:51:38.272173 +0100 CET"
1548449498272173
```

nanoseconds:

```text
$ epoch -unit ns "2019-01-25 21:51:38.272173 +0100 CET"
1548449498272173000
```

## Supported Formats

All current Go formats as of 2019-01-26 (https://golang.org/pkg/time/#pkg-constants):

``` go
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
