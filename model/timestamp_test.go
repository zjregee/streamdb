package model

import (
	"time"
	"testing"
	"github.com/stretchr/testify/require"
)

func TestFromTime(t *testing.T) {
	testCases := []struct {
		input    time.Time
		expected int64
	}{
		{time.Date(2024, time.April, 1, 12, 0, 0, 0, time.UTC), 1711972800000},
		{time.Date(2024, time.April, 1, 14, 0, 0, 123456789, time.UTC), 1711980000123},
	}

	for _, tc := range testCases {
		result := FromTime(tc.input)
		require.Equal(t, tc.expected, result)
	}
}

func TestTime(t *testing.T) {
	testCases := []struct {
		input    int64
		expected time.Time
	}{
		{1711972800000, time.Date(2024, time.April, 1, 12, 0, 0, 0, time.UTC)},
		{1711980000123, time.Date(2024, time.April, 1, 14, 0, 0, 123000000, time.UTC)},
	}

	for _, tc := range testCases {
		result := Time(tc.input)
		require.Equal(t, tc.expected, result)
	}
}

func TestFromFloatSeconds(t *testing.T) {
	testCases := []struct {
		input    float64
		expected int64
	}{
		{1.234, 1234},
		{1.555, 1555},
	}

	for _, tc := range testCases {
		result := FromFloatSeconds(tc.input)
		require.Equal(t, tc.expected, result)
	}
}
