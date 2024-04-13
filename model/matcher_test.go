package model

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestMatcherMatches(t *testing.T) {
	require.Equal(t, NewMatcher(MatchEqual, "test", "bar").Matches("bar"), true)
	require.Equal(t, NewMatcher(MatchEqual, "test", "foo-bar").Matches("bar"), false)
	require.Equal(t, NewMatcher(MatchNotEqual, "test", "bar").Matches("bar"), false)
	require.Equal(t, NewMatcher(MatchNotEqual, "test", "foo-bar").Matches("bar"), true)
}

func TestMatcherInverse(t *testing.T) {
	require.Equal(t, NewMatcher(MatchEqual, "test", "bar").Inverse().Matches("bar"), false)
	require.Equal(t, NewMatcher(MatchEqual, "test", "foo-bar").Inverse().Matches("bar"), true)
	require.Equal(t, NewMatcher(MatchNotEqual, "test", "bar").Inverse().Matches("bar"), true)
	require.Equal(t, NewMatcher(MatchNotEqual, "test", "foo-bar").Inverse().Matches("bar"), false)
}
