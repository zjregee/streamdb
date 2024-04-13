package model

import "fmt"

type MatchType int

const (
	MatchEqual    MatchType = iota
	MatchNotEqual
)

var matchTypeToStr = [...]string{
	MatchEqual:    "=",
	MatchNotEqual: "!=",
}

func (m MatchType) String() string {
	if m < MatchEqual || m > MatchNotEqual {
		panic("model.MatchType.String: invalid match type")
	}
	return matchTypeToStr[m]
}

func NewMatcher(mt MatchType, name, value string) *Matcher {
	return &Matcher{
		Type: mt,
		Name: name,
		Value: value,
	}
}

type Matcher struct {
	Type  MatchType
	Name  string
	Value string
}

func (m *Matcher) String() string {
	return fmt.Sprintf("%s%s%q", m.Name, m.Type, m.Value)
}

func (m *Matcher) Matches(value string) bool {
	switch m.Type {
	case MatchEqual:
		return m.Value == value
	case MatchNotEqual:
		return m.Value != value
	}
	panic("model.Matcher.Matches: invalid match type")
}

func (m *Matcher) Inverse() *Matcher {
	switch m.Type {
	case MatchEqual:
		return NewMatcher(MatchNotEqual, m.Name, m.Value)
	case MatchNotEqual:
		return NewMatcher(MatchEqual, m.Name, m.Value)
	}
	panic("model.Matcher.Inverse: invalid match type")
}
