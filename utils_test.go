package i3config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeString(t *testing.T) {
	testCases := []struct {
		str      string
		expected string
	}{
		{
			str:      "a",
			expected: `"a"`,
		},
		{
			str:      `bar"`,
			expected: `"bar\\""`,
		},
		{
			str:      `\"`,
			expected: `"\\\\""`,
		},
		{
			str:      `\ `,
			expected: `"\\ "`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.str, func(t *testing.T) {
			result := escapeString(tc.str)
			assert.Equal(t, tc.expected, result)
		})
	}
}
