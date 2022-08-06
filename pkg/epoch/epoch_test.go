package epoch

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

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
				require.EqualError(t, tt.expectedErr, err.Error(), err)
				return
			} else if tt.expectedErr != nil {
				require.EqualError(t, err, tt.expectedErr.Error(), tt.expectedErr)
				return
			}

			require.Equal(t, tt.expected, unit)
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
				require.EqualError(t, tt.expected.err, err.Error(), err)
				return
			} else if tt.expected.err != nil {
				require.EqualError(t, err, tt.expected.err.Error(), tt.expected.err)
				return
			}

			require.Equal(t, tt.expected.timestamp, timestamp)
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
				require.EqualError(t, tt.expected.err, err.Error(), err)
				return
			} else if tt.expected.err != nil {
				require.EqualError(t, err, tt.expected.err.Error(), tt.expected.err)
				return
			}

			require.Equal(t, tt.expected.time.UTC(), timestamp.UTC())
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
			require.Equal(t, tt.expected.unit, unit)
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
			parsed, layout, err := ParseFormatted(tt.given.formatted)
			require.Equal(t, tt.expected.layout, layout)
			if err != nil {
				require.EqualError(t, tt.expected.err, err.Error(), err)
				return
			} else if tt.expected.err != nil {
				require.EqualError(t, err, tt.expected.err.Error(), tt.expected.err)
				return
			}

			require.Equal(t, tt.expected.time, parsed)
		})
	}
}
