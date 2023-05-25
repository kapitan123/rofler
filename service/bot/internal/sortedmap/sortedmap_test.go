package sortedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedMap_Sort(t *testing.T) {
	expected := DescendingPairList{Pair{Key: "r3", Value: 3}, Pair{Key: "r2", Value: 2}, Pair{Key: "r1", Value: 1}, Pair{Key: "r0", Value: 0}}
	testCases := []map[string]int{
		{"r0": 0, "r1": 1, "r2": 2, "r3": 3},
		{"r3": 3, "r2": 2, "r1": 1, "r0": 0},
		{"r0": 0, "r2": 2, "r3": 3, "r1": 1},
	}

	for _, tc := range testCases {
		t.Run("Should sort map in descending order", func(t *testing.T) {
			actual := Sort(tc)
			assert.Equal(t, expected, actual, "Wasn't sorted correctly")
		})
	}
}
