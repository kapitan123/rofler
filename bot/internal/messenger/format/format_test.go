package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat_AsDescendingList(t *testing.T) {
	expected := "r3: 3\nr2: 2\nr1: 1\nr0: 0\n"
	testCases := []map[string]int{
		{"r0": 0, "r2": 2, "r3": 3, "r1": 1},
	}

	for _, tc := range testCases {
		t.Run("Should build a list using the template", func(t *testing.T) {
			actual := AsDescendingList(tc, "%s: %d")
			assert.Equal(t, expected, actual)
		})
	}
}
