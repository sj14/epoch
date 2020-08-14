package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	type args struct {
		input      string
		now        string
		unitFlag   string
		formatFlag string
		tzFlag     string
		quietFlag  bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// Remember that timezone "Local" can't be used as the CI might have a different timezone!
		{name: "empty input", args: args{now: "2020-07-18 17:46:45.215239 +0200 CEST", unitFlag: "guess"}, want: "1595087205"},
		{name: "empty input/utc", args: args{now: "2020-07-18 17:46:45.215239 +0200 CEST", tzFlag: "UTC", unitFlag: "guess"}, want: "2020-07-18 15:46:45.215239 +0000 UTC"},

		{name: "timedate", args: args{input: "2020-07-18 17:46:45.215239 +0200 CEST", unitFlag: "guess"}, want: "1595087205"},
		{name: "timedate/unit", args: args{input: "2020-07-18 17:46:45.215239 +0200 CEST", unitFlag: "ms"}, want: "1595087205215"},
		{name: "timedate/timezone", args: args{input: "2020-07-18 17:46:45.215239 +0200 CEST", tzFlag: "MST", unitFlag: "guess"}, want: "2020-07-18 08:46:45.215239 -0700 MST"},
		{name: "timedate/timezone/format", args: args{input: "2020-07-18 17:46:45.215239 +0200 CEST", formatFlag: "unix", tzFlag: "UTC", unitFlag: "guess"}, want: "Sat Jul 18 15:46:45 UTC 2020"},
		{name: "timedate/timezone/unit/FAIL", args: args{input: "2020-07-18 17:46:45.215239 +0200 CEST", unitFlag: "ms", formatFlag: "unix", tzFlag: "UTC"}, wantErr: true},

		{name: "timestamp/timezone/unitsuffix", args: args{input: "1595087205us", tzFlag: "Europe/Berlin", unitFlag: "guess"}, want: "1970-01-01 01:26:35.087205 +0100 CET"},
		{name: "timestamp/timezone/unit", args: args{input: "1595087205", tzFlag: "Europe/Berlin", unitFlag: "ms"}, want: "1970-01-19 12:04:47.205 +0100 CET"},
		{name: "timestamp/timezone", args: args{input: "1595087205", tzFlag: "UTC", unitFlag: "guess"}, want: "2020-07-18 15:46:45 +0000 UTC"},
		{name: "timestamp/timezone/format", args: args{input: "1595087205", formatFlag: "ruby", tzFlag: "UTC", unitFlag: "guess"}, want: "Sat Jul 18 15:46:45 +0000 2020"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := run(tt.args.input, tt.args.now, tt.args.unitFlag, tt.args.formatFlag, tt.args.tzFlag, tt.args.quietFlag)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("run() = %v, want %v", got, tt.want)
			}

		})
	}
}