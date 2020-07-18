package main

import (
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		input      string
		now        string
		unitFlag   string
		formatFlag string
		tzFlag     string
		quietFlag  bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "empty input", args: args{now: "2020-07-18 17:46:45.215239 +0200 CEST"}, want: "1595087205"},
		{name: "empty input/utc", args: args{now: "2020-07-18 17:46:45.215239 +0200 CEST", tzFlag: "UTC"}, want: "2020-07-18 15:46:45.215239 +0000 UTC"},

		{name: "timedate", args: args{input: "2020-07-18 17:46:45.215239 +0200 CEST"}, want: "1595087205"},
		{name: "timedate/unit", args: args{input: "2020-07-18 17:46:45.215239 +0200 CEST", unitFlag: "ms"}, want: "1595087205215"},
		{name: "timedate/timezone", args: args{input: "2020-07-18 17:46:45.215239 +0200 CEST", tzFlag: "MST"}, want: "2020-07-18 08:46:45.215239 -0700 MST"},
		{name: "timedate/timezone/format", args: args{input: "2020-07-18 17:46:45.215239 +0200 CEST", formatFlag: "unix", tzFlag: "UTC"}, want: "Sat Jul 18 15:46:45 UTC 2020"},

		{name: "timestamp", args: args{input: "1595087205"}, want: "2020-07-18 17:46:45 +0200 CEST"},
		{name: "timestamp/unit", args: args{input: "1595087205", unitFlag: "ms"}, want: "1970-01-19 12:04:47.205 +0100 CET"},
		{name: "timestamp/timezone", args: args{input: "1595087205", tzFlag: "UTC"}, want: "2020-07-18 15:46:45 +0000 UTC"},
		{name: "timestamp/timezone/format", args: args{input: "1595087205", formatFlag: "ruby", tzFlag: "UTC"}, want: "Sat Jul 18 15:46:45 +0000 2020"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := run(tt.args.input, tt.args.now, tt.args.unitFlag, tt.args.formatFlag, tt.args.tzFlag, tt.args.quietFlag); got != tt.want {
				t.Errorf("run() = %v, want %v", got, tt.want)
			}
		})
	}
}
