package showstats

import "time"

type StatPoint struct {
	Value float64
	Day   time.Time
}

func (sp StatPoint) FloatDate() float64 {
	return float64(sp.Day.Day()) // this is a simple hack and will work only inside one month
}

// probably Should have a property which converts Day, to a float representation
