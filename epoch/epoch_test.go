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
