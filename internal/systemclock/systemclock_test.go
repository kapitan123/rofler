package systemclock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSystemClock_CurrentDate(t *testing.T) {
	t.Run("Should return current date with no time", func(t *testing.T) {
		c := New()
		date := c.CurrentDate()
		now := time.Now() // the chance of false negative is really low

		assert.Equal(t, 0, date.Hour(), "hour should be 0")
		assert.Equal(t, 0, date.Minute(), "minute should be 0")
		assert.Equal(t, 0, date.Second(), "second should be 0")
		assert.Equal(t, 0, date.Nanosecond(), "nanosecond should be 0")

		assert.Equal(t, now.Year(), date.Year())
		assert.Equal(t, now.Month(), date.Month())
		assert.Equal(t, now.Day(), date.Day())
	})
}
