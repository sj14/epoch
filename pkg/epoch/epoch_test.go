package epoch

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func equalError(t *testing.T, got, want error) {
	t.Helper()

	if (got == nil && want != nil) ||
		(got != nil && want == nil) ||
		(got.Error() != want.Error()) {
		t.Fatalf("got %q, want %q\n", got, want)
	}
}

func equal[T any](t *testing.T, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got '%#v', want '%#v\n", got, want)
	}
}

func TestParseUnit(t *testing.T) {
	testCases := []struct {
		description string
		given       string
		expected    TimeUnit
		expectedErr error
	}{
		{
			description: "empty",
			given:       "",
			expectedErr: errors.New("failed to parse input '' to unit"),
		},
		{
			description: "seconds",
			given:       "s",
			expected:    UnitSeconds,
		},
		{
			description: "milliseconds",
			given:       "ms",
			expected:    UnitMilliseconds,
		},
		{
			description: "microseconds",
			given:       "us",
			expected:    UnitMicroseconds,
		},
		{
			description: "nanosecods",
			given:       "ns",
			expected:    UnitNanoseconds,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			unit, err := ParseUnit(tt.given)
			if err != nil {
				equalError(t, err, tt.expectedErr)
				return
			} else if tt.expectedErr != nil {
				equalError(t, err, tt.expectedErr)
				return
			}
			equal(t, unit, tt.expected)
		})
	}
}

func TestToTimestamp(t *testing.T) {
	type givenType struct {
		time time.Time
		unit TimeUnit
	}

	type expecedType struct {
		timestamp int64
		err       error
	}

	testCases := []struct {
		description string
		given       givenType
		expected    expecedType
	}{
		{
			description: "empty",
			expected:    expecedType{timestamp: -62135596800},
		},
		{
			description: "wrong unit",
			given:       givenType{time: time.Unix(0, 1549727875568573000), unit: 42},
			expected:    expecedType{err: errors.New("unknown unit '42'")},
		},
		{
			description: "seconds",
			given:       givenType{time: time.Unix(0, 1549727875568573000), unit: UnitSeconds},
			expected:    expecedType{timestamp: 1549727875},
		},
		{
			description: "milliseconds",
			given:       givenType{time: time.Unix(0, 1549727875568573000), unit: UnitMilliseconds},
			expected:    expecedType{timestamp: 1549727875568},
		},
		{
			description: "microseconds",
			given:       givenType{time: time.Unix(0, 1549727875568573000), unit: UnitMicroseconds},
			expected:    expecedType{timestamp: 1549727875568573},
		},
		{
			description: "nanoseconds",
			given:       givenType{time: time.Unix(0, 1549727875568573000), unit: UnitNanoseconds},
			expected:    expecedType{timestamp: 1549727875568573000},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			timestamp, err := ToTimestamp(tt.given.time, tt.given.unit)
			if err != nil {
				equalError(t, err, tt.expected.err)
				return
			} else if tt.expected.err != nil {
				equalError(t, err, tt.expected.err)
				return
			}

			equal(t, timestamp, tt.expected.timestamp)
		})
	}
}

func TestParseTimestamp(t *testing.T) {
	type givenType struct {
		timestamp int64
		unit      TimeUnit
	}

	type expecedType struct {
		time time.Time
		err  error
	}

	testCases := []struct {
		description string
		given       givenType
		expected    expecedType
	}{
		{
			description: "empty",
			expected:    expecedType{time: time.Date(1970, 1, 1, 0, 0, 0, 00000000, time.UTC)},
		},
		{
			description: "wrong unit",
			given:       givenType{timestamp: 1549741094, unit: 42},
			expected:    expecedType{err: errors.New("unknown unit '42'")},
		},
		{
			description: "seconds",
			given:       givenType{timestamp: 1549741094, unit: UnitSeconds},
			expected:    expecedType{time: time.Date(2019, 2, 9, 19, 38, 14, 00000000, time.UTC)},
		},
		{
			description: "milliseconds",
			given:       givenType{timestamp: 1549741094065, unit: UnitMilliseconds},
			expected:    expecedType{time: time.Date(2019, 2, 9, 19, 38, 14, 65000000, time.UTC)},
		},
		{
			description: "microseconds",
			given:       givenType{timestamp: 1549741094065178, unit: UnitMicroseconds},
			expected:    expecedType{time: time.Date(2019, 2, 9, 19, 38, 14, 65178000, time.UTC)},
		},
		{
			description: "nanoseconds",
			given:       givenType{timestamp: 1549741094065178000, unit: UnitNanoseconds},
			expected:    expecedType{time: time.Date(2019, 2, 9, 19, 38, 14, 65178000, time.UTC)},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			timestamp, err := ParseTimestamp(tt.given.timestamp, tt.given.unit)
			if err != nil {
				equalError(t, err, tt.expected.err)
				return
			} else if tt.expected.err != nil {
				equalError(t, err, tt.expected.err)
				return
			}

			equal(t, timestamp.UTC(), tt.expected.time.UTC())
		})
	}
}

func TestGuessUnit(t *testing.T) {
	type givenType struct {
		timestamp int64
		ref       time.Time
	}

	ref := time.Unix(0, 1549777538844829000)

	type expecedType struct {
		unit TimeUnit
	}

	testCases := []struct {
		description string
		given       givenType
		expected    expecedType
	}{
		{
			description: "empty",
			expected:    expecedType{unit: UnitSeconds},
		},
		{
			description: "seconds/exactly",
			given:       givenType{timestamp: 1549777538, ref: ref},
			expected:    expecedType{unit: UnitSeconds},
		},
		{
			description: "milliseconds/exactly",
			given:       givenType{timestamp: 1549777538844, ref: ref},
			expected:    expecedType{unit: UnitMilliseconds},
		},
		{
			description: "microseconds/exactly",
			given:       givenType{timestamp: 1549777538844829, ref: ref},
			expected:    expecedType{unit: UnitMicroseconds},
		},
		{
			description: "nanoseconds/exactly",
			given:       givenType{timestamp: 1549777538844829000, ref: ref},
			expected:    expecedType{unit: UnitNanoseconds},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			unit := GuessUnit(tt.given.timestamp, tt.given.ref)
			equal(t, unit, tt.expected.unit)
		})
	}
}

func TestParseFormatted(t *testing.T) {
	type givenType struct {
		formatted string
	}

	type expecedType struct {
		time   time.Time
		layout string
		err    error
	}

	testCases := []struct {
		description string
		given       givenType
		expected    expecedType
	}{
		{
			description: "empty",
			given:       givenType{formatted: ""},
			expected:    expecedType{err: ErrParseFormatted, layout: ""},
		},
		{
			description: "rfc1123",
			given:       givenType{formatted: "Mon, 02 Jan 2006 15:04:05 UTC"},
			expected:    expecedType{time: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC), layout: time.RFC1123},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.description, func(t *testing.T) {
			parsed, layout, err := ParseFormatted(tt.given.formatted, time.Local)
			equal(t, layout, tt.expected.layout)
			if err != nil {
				equalError(t, err, tt.expected.err)
				return
			} else if tt.expected.err != nil {
				equalError(t, err, tt.expected.err)
				return
			}

			equal(t, parsed, tt.expected.time)
		})
	}
}

func TestFormatSimple(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		format string
		want   string
	}{
		{
			name:   "year-month-day long",
			format: "{YYYY}-{MM}-{DD}",
			want:   "2022-09-08",
		},
		{
			name:   "month/day/year short",
			format: "{M}/{D}/{YY}",
			want:   "9/8/22",
		},
		{
			name:   "hour:minute:second long",
			format: "{hh}:{mm}:{ss}",
			want:   "07:06:05",
		},
		{
			name:   "second-minute short",
			format: "{s}-{m}",
			want:   "5-6",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			goFormat := FormatSimple(tt.format)
			fixedDate := time.Date(
				2022, // year
				9,    // month
				8,    // day
				7,    // hour
				6,    // minute
				5,    // second
				4,    // nanosecond
				time.UTC,
			)

			if got := fixedDate.Format(goFormat); got != tt.want {
				t.Errorf("ParseFormat() got %q, want %q", got, tt.want)
			}
		})
	}
}
