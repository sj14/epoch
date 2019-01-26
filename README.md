# epoch

`epoch` converts unix timestamps to human readable formats and vice-versa. Support input of timestamps with second and nanosecond accurary.

## Installation

``` text 
go get -u github.com/sj14/epoch
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