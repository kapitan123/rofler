package showStats

import "time"

type StatPoint struct {
	Value float64
	Day   time.Time
}

// probably Should have a property which converts Day, to a float representation
