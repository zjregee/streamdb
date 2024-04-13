package model

import (
	"testing"
	"github.com/stretchr/testify/require"
)

func TestFromMap(t *testing.T) {
	testCases := []struct {
		m        map[string]string
		expected Labels
	}{
		{
			m: map[string]string{
				"t1": "v1",
				"t2": "v2",
			},
			expected: Labels{
				Label{Name: "t1", Value: "v1"},
				Label{Name: "t2", Value: "v2"},
			},
		},
		{
			m: map[string]string{},
			expected: Labels{},
		},
	}
	for _, tc := range testCases {
		result := FromMap(tc.m)
		require.Equal(t, tc.expected, result)
	}
}

func TestFromStrings(t *testing.T) {
	testCases := []struct {
		ss       []string
		expected Labels
	}{
		{
			ss: []string{"t1", "v1", "t2", "v2"},
			expected: Labels{
				Label{Name: "t1", Value: "v1"},
				Label{Name: "t2", Value: "v2"},
			},
		},
		{
			ss: []string{},
			expected: Labels{},
		},
	}
	for _, tc := range testCases {
		result := FromStrings(tc.ss...)
		require.Equal(t, tc.expected, result)
	}
}

func TestLabelsMap(t *testing.T) {
	require.Equal(t, map[string]string{}, Labels{}.Map())
	require.Equal(t, map[string]string{
		"t1": "v1",
		"t2": "v2",
	}, FromMap(map[string]string{
		"t1": "v1",
		"t2": "v2",
	}).Map())
}

func TestLabelsString(t *testing.T) {
	require.Equal(t, "{}", Labels{}.String())
	require.Equal(t, "{t1=\"v1\", t2=\"v2\"}", FromStrings("t1", "v1", "t2", "v2").String())
}

func TestLabelsContains(t *testing.T) {
	require.Equal(t, false, FromStrings("aaa", "111", "bbb", "222").Contains("foo"))
	require.Equal(t, true, FromStrings("aaaa", "111", "bbb", "222").Contains("aaaa"))
	require.Equal(t, true, FromStrings("aaaa", "111", "bbb", "222").Contains("bbb"))
}

func TestLabelsGet(t *testing.T) {
	require.Equal(t, "", FromStrings("aaa", "111", "bbb", "222").Get("foo"))
	require.Equal(t, "111", FromStrings("aaaa", "111", "bbb", "222").Get("aaaa"))
	require.Equal(t, "222", FromStrings("aaaa", "111", "bbb", "222").Get("bbb"))
}

func TestLabelsCopy(t *testing.T) {
	require.Equal(t, Labels{}, Labels{}.Copy())
	require.Equal(t, FromStrings("t1", "v1", "t2", "v2"), FromStrings("t1", "v1", "t2", "v2").Copy())
}

func TestLabelsIsEmpty(t *testing.T) {
	require.Equal(t, true, Labels{}.IsEmpty())
	require.Equal(t, false, FromStrings("t1", "v1", "t2", "v2").IsEmpty())
}

func TestLabelsIsValid(t *testing.T) {
	require.Equal(t, true, FromStrings("aaa", "111", "bbb", "222").IsValid())
	require.Equal(t, false, FromStrings("", "111", "bbb", "222").IsValid())
	require.Equal(t, false, FromStrings("aaaa", "", "bbb", "222").IsValid())
}

func TestLabelsBuilderReset(t *testing.T) {
	lb := NewLabelsBuilder(Labels{})
	require.Equal(t, Labels{}, lb.Reset(Labels{}).Labels())
	require.Equal(t, FromStrings("t1", "v1", "t2", "v2"), lb.Reset(FromStrings("t1", "v1", "t2", "v2")).Labels())
}

func TestLabelsBuilderLabels(t *testing.T) {
	require.Equal(t, Labels{}, NewLabelsBuilder(Labels{}).Labels())
	require.Equal(t, FromStrings("t1", "v1", "t2", "v2"), NewLabelsBuilder(FromStrings("t1", "v1", "t2", "v2")).Labels())
}

func TestLabelsBuilderGet(t *testing.T) {
	require.Equal(t, "", NewLabelsBuilder(FromStrings("aaa", "111", "bbb", "222")).Get("foo"))
	require.Equal(t, "111", NewLabelsBuilder(FromStrings("aaa", "111", "bbb", "222")).Get("aaa"))
	require.Equal(t, "222", NewLabelsBuilder(FromStrings("aaa", "111", "bbb", "222")).Get("bbb"))
}

func TestLabelsBuilderSet(t *testing.T) {
	lb := NewLabelsBuilder(FromStrings("t1", "v1", "t2", "v2"))
	require.Equal(t, FromStrings("t1", "v2", "t2", "v2"), lb.Set("t1", "v2").Labels())
	require.Equal(t, FromStrings("t1", "v2", "t2", "v1"), lb.Set("t2", "v1").Labels())
	require.Equal(t, FromStrings("t1", "v2", "t2", "v1", "t3", "v3"), lb.Set("t3", "v3").Labels())
}

func TestLabelsBuilderDelete(t *testing.T) {
	lb := NewLabelsBuilder(FromStrings("t1", "v1", "t2", "v2", "t3", "v3"))
	require.Equal(t, FromStrings("t2", "v2", "t3", "v3"), lb.Delete("t1").Labels())
	require.Equal(t, Labels{}, lb.Delete("t2", "t3").Labels())
}
