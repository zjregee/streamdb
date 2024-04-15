package utils

import (
	"fmt"
	"errors"
)

func NewAnnotations() *Annotations {
	return &Annotations{}
}

type Annotations map[string]error

func (a *Annotations) Add(err error) Annotations {
	(*a)[err.Error()] = err
	return *a
}

func (a *Annotations) Merge(aa Annotations) Annotations {
	for k, v := range aa {
		(*a)[k] = v
	}
	return *a
}

func (a Annotations) AsErrors() []error {
	arr := make([]error, 0, len(a))
	for _, err := range a {
		arr = append(arr, err)
	}
	return arr
}

func (a Annotations) AsStrings(query string, maxAnnos int) []string {
	arr := make([]string, 0, len(a))
	for _, err := range a {
		if maxAnnos > 0 && len(arr) > maxAnnos {
			break
		}
		var anErr annoErr
		if errors.As(err, &anErr) {
			anErr.Query = query
			err = anErr
		}
		arr = append(arr, err.Error())
	}
	if maxAnnos > 0 && len(a) > maxAnnos {
		arr = append(arr, fmt.Sprintf("%d more annotations omitted", len(a) - maxAnnos))
	}
	return arr
}

type annoErr struct {
	Err           error
	Query         string
	PositionRange PositionRange
}

func (e annoErr) Error() string {
	if e.Query == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s (%s)", e.Err, e.PositionRange.StartPosInput(e.Query, 0))
}

func (e annoErr) Unwrap() error {
	return e.Err
}
