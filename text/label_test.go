package text

import (
	"reflect"
	"sort"
	"testing"
)

func TestUnescapeLabelValues(t *testing.T) {
	type testCase struct {
		values []string
	}

	testCases := []testCase{
		{
			values: []string{
				"print1",
				"print1$",
				"$print1",
				"$print_1$Z$",
				"$print_中国1$Z$",
			},
		},
	}

	for _, tc := range testCases {
		sort.Strings(tc.values)
		escaped := EscapeLabelValues(tc.values)
		unescaped, err := UnescapeLabelValues(escaped)
		if err != nil || !reflect.DeepEqual(tc.values, unescaped) {
			t.Fatalf("origin %+v, escaped %+v, unescaped %+v, err %v", tc.values, escaped, unescaped, err)
		}
	}
}
