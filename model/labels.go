package model

import (
	"bytes"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"
)

type Label struct {
	Name, Value string
}

func FromMap(m map[string]string) Labels {
	res := make(Labels, 0, len(m))
	for k, v := range m {
		res = append(res, Label{Name: k, Value: v})
	}
	slices.SortFunc(res, func(a, b Label) int { return strings.Compare(a.Name, b.Name) })
	return res
}

func FromStrings(ss ...string) Labels {
	if len(ss) % 2 != 0 {
		panic("model.FromStrings: invalid number of strings")
	}
	res := make(Labels, 0, len(ss) / 2)
	for i := 0; i < len(ss); i += 2 {
		res = append(res, Label{Name: ss[i], Value: ss[i + 1]})
	}
	slices.SortFunc(res, func(a, b Label) int { return strings.Compare(a.Name, b.Name) })
	return res
}

type Labels []Label

func (ls Labels) Len() int           { return len(ls) }
func (ls Labels) Swap(i, j int)      { ls[i], ls[j] = ls[j], ls[i] }
func (ls Labels) Less(i, j int) bool { return ls[i].Name < ls[j].Name }

func (ls Labels) Range(f func(l Label)) {
	for _, l := range ls {
		f(l)
	}
}

func (ls Labels) Validate(f func(l Label) error) error {
	for _, l := range ls {
		if err := f(l); err != nil {
			return err
		}
	}
	return nil
}

func (ls Labels) Map() map[string]string {
	m := make(map[string]string)
	ls.Range(func(l Label) {
		m[l.Name] = l.Value
	})
	return m
}

func (ls Labels) String() string {
	var bytea [1024]byte // avoid memory allocation
	b := bytes.NewBuffer(bytea[:0])
	b.WriteByte('{')
	i := 0
	ls.Range(func(l Label) {
		if i > 0 {
			b.WriteByte(',')
			b.WriteByte(' ')
		}
		b.WriteString(l.Name)
		b.WriteByte('=')
		b.Write(strconv.AppendQuote(b.AvailableBuffer(), l.Value))
		i++
	})
	b.WriteByte('}')
	return b.String()
}

func (ls Labels) Contains(name string) bool {
	for _, l := range ls {
		if l.Name == name {
			return true
		}
	}
	return false
}

func (ls Labels) Get(name string) string {
	for _, l := range ls {
		if l.Name == name {
			return l.Value
		}
	}
	return ""
}

func (ls Labels) Copy() Labels {
	res := make(Labels, len(ls))
	copy(res, ls)
	return res
}

func (ls Labels) IsEmpty() bool {
	return len(ls) == 0
}

func (ls Labels) IsValid() bool {
	err := ls.Validate(func(l Label) error {
		if len(l.Name) == 0 || len(l.Value) == 0 {
			return strconv.ErrSyntax
		}
		if !utf8.ValidString(l.Name) || !utf8.ValidString(l.Value) {
			return strconv.ErrSyntax
		}
		return nil
	})
	return err == nil
}

func NewLabelsBuilder(base Labels) *LabelsBuilder {
	b := &LabelsBuilder{
		del: make([]string, 0, 5),
		add: make([]Label, 0, 5),
	}
	return b.Reset(base)
}

type LabelsBuilder struct {
	base Labels
	add  Labels
	del  []string
}

func (lb *LabelsBuilder) Reset(base Labels) *LabelsBuilder {
	lb.base = base
	lb.del = lb.del[:0]
	lb.add = lb.add[:0]
	return lb
}

func (lb *LabelsBuilder) Labels() Labels {
	if len(lb.del) == 0 && len(lb.add) == 0 {
		return lb.base
	}
	expectSize := len(lb.base) + len(lb.add) - len(lb.del)
	res := make(Labels, 0, expectSize)
	for _, l := range lb.base {
		if slices.Contains(lb.del, l.Name) || lb.add.Contains(l.Name) {
			continue
		}
		res = append(res, l)
	}
	if len(lb.add) > 0 {
		res = append(res, lb.add...)
		slices.SortFunc(res, func(a, b Label) int { return strings.Compare(a.Name, b.Name) })
	}
	return res
}

func (lb *LabelsBuilder) Get(name string) string {
	for _, add := range lb.add {
		if add.Name == name {
			return add.Value
		}
	}
	if slices.Contains(lb.del, name) {
		return ""
	}
	return lb.base.Get(name)
}

func (lb *LabelsBuilder) Set(name, value string) *LabelsBuilder {
	for i, add := range lb.add {
		if add.Name == name {
			lb.add[i].Value = value
			return lb
		}
	}
	lb.add = append(lb.add, Label{Name: name, Value: value})
	return lb
}

func (lb *LabelsBuilder) Delete(names ...string) *LabelsBuilder {
	for _, name := range names {
		for i, add := range lb.add {
			if add.Name == name {
				lb.add = append(lb.add[:i], lb.add[i + 1:]...)
			}
		}
		lb.del = append(lb.del, name)
	}
	return lb
}
