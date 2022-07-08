package sortedmap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSortedMap_Sort(t *testing.T) {
	for _, test := range []struct {
		nonSorted map[string]int
		expected PairList
	}{
	{
		nonSorted: map[string]int{
	}
	} {
	t.Run("Should sort map in descending order", func(t *testing.T) {
		sorted:= Sort(test.nonSorted)
		date := c.CurrentDate()
		now := time.Now()

		assert.Equal(t, 0, date.Minute())
	})
}
}
