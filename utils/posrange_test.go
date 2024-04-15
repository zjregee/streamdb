package utils

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestPositionRange(t *testing.T) {
	testCases := []struct {
		p          PositionRange
		query      string
		lineOffset int
		expected   string
	}{
		{
			p: PositionRange{ Start: 5 },
			query: "",
			lineOffset: 0,
			expected: "unknown position",
		},
		{
			p: PositionRange{ Start: 12 },
			query: "test query",
			lineOffset: 0,
			expected: "invalid position",
		},
		{
			p: PositionRange{ Start: 5 },
			query: "test query",
			lineOffset: 0,
			expected: "1:6",
		},
		{
			p: PositionRange{ Start: 5 },
			query: "test\nquery",
			lineOffset: 0,
			expected: "2:1",
		},
		{
			p: PositionRange{ Start: 5 },
			query: "test\nquery",
			lineOffset: 2,
			expected: "4:1",
		},
	}
	for _, tc := range testCases {
		require.Equal(t, tc.expected, tc.p.StartPosInput(tc.query, tc.lineOffset))
	}
}
