# epoch

`epoch` converts unix timestamps to human readable formats and vice-versa.

Why?  
To convert timestamps to dates, you have to run different commands for Linux and MacOS: `date -d @1267619929` vs `date -r 1267619929` which might still be ok, but have you ever tried converting a formatted timestamp such as `2019-01-25 21:51:38 +0100 CET` to a timestamp? Of course you can do this, but it's cumbersome. This tool tries to solve the task with ease:

```text
$ epoch "2019-01-25 21:51:38 +0100 CET"
1548449498
```

Currently, supports input of timestamps with second and nanosecond accuracy.

## Installation

``` text
go get -u github.com/sj14/epoch
```

## Supported Formats

All current Go formats as of 26.01.2019 (https://golang.org/pkg/time/#pkg-constants):

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

## Examples

### Timestamps to human readable format

seconds:

``` text
$ epoch 1548449513
2019-01-25 21:51:53 +0100 CET
```

nanoseconds:

``` text
$ epoch 1548449513940562000
2019-01-25 21:51:53.940562 +0100 CET
```

Negative timestamp:

``` text
$ epoch -- -15484495
1969-07-05 19:45:05 +0100 CET
```

### Formatted input to epoch timestamps

seconds:

``` text
$ epoch "2019-01-25 21:51:38.272173 +0100 CET"
1548449498
```

nanoseconds:

```text
$ epoch -nsec "2019-01-25 21:51:38.272173 +0100 CET"
1548449498272173000
```