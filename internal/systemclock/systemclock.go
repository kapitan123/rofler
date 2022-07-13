package systemclock

import (
	"time"
)

type SystemClock struct {
}

func New() *SystemClock {
	return &SystemClock{}
}

func (sc *SystemClock) Now() time.Time {
	return time.Now()
}

func (sc *SystemClock) CurrentDate() time.Time {
	t := sc.Now()
	d := (24 * time.Hour)
	return t.Truncate(d)
}
