package model

import (
	"math"
	"time"
)

func FromTime(t time.Time) int64 {
	return t.Unix() * 1000 + int64(t.Nanosecond()) / int64(time.Millisecond)
}

func Time(ts int64) time.Time {
	return time.Unix(ts / 1000, (ts % 1000) * int64(time.Millisecond)).UTC()
}

func FromFloatSeconds(ts float64) int64 {
	return int64(math.Round(ts * 1000))
}
