# epoch

`epoch` converts unix timestamps to human readable formats and vice versa.

**Why?**  
To convert timestamps to dates, you have to run different commands for Linux and macOS: `date -d @1267619929` vs `date -r 1267619929`, and what about handling nanosecond timestamps? Seriously, I don't know how to do this with `date`.  
Furthermore, have you ever tried converting a time formatted string such as `"2019-01-25 21:51:38 +0100 CET"` to a timestamp? Of course, you can do this somehow, but all ways I've found so far were too cumbersome. This tool tries to solve all this with ease:

```text
$ epoch "2019-01-25 21:51:38 +0100 CET"
1548449498
```

Currently, it's a quickly hacked tool, which solves my needs, but I thought it might be useful for other people, too.

## Installation

``` text
go get -u github.com/sj14/epoch
```

## Examples

### Timestamps to human readable format

#### conversion based on given unit

``` text
$ epoch --unit s 1548449513
Fri Jan 25 21:51:53 CET 2019
```

``` text
$ epoch --unit ms 1548449513
Sun Jan 18 23:07:29 CET 1970
```

``` text
$ epoch --unit us 1548449513
Thu Jan  1 01:25:48 CET 1970
```

``` text
$ epoch --unit ns 1548449513
Thu Jan  1 01:00:01 CET 1970
```

#### guess the unit

seconds (auto guessed):

``` text
$ epoch 1548449513
2019-01-25 21:51:53 +0100 CET
```

nanoseconds (auto guessed):

``` text
$ epoch 1548449513940562000
2019-01-25 21:51:53.940562 +0100 CET
```

#### negative timestamp

``` text
$ epoch -- -15484495
1969-07-05 19:45:05 +0100 CET
```

#### using the pipe

``` text
$ echo -15484495 | epoch
1969-07-05 19:45:05 +0100 CET
```

### Formatted input to epoch timestamps

seconds (default when no `unit` flag given):

``` text
$ epoch --unit s "2019-01-25 21:51:38.272173 +0100 CET"
1548449498
```

milliseconds:

```text
$ epoch --unit ms "2019-01-25 21:51:38.272173 +0100 CET"
1548449498272
```

microseconds:

```text
$ epoch --unit us "2019-01-25 21:51:38.272173 +0100 CET"
1548449498272173
```

nanoseconds:

```text
$ epoch --unit ns "2019-01-25 21:51:38.272173 +0100 CET"
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
