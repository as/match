package match

import (
	"testing"
)

type BestParams struct {
	Haystack, Needles []string
	Expected          map[string]string
	Error             bool
}

var Empty = []string{}

var EmptyMap = map[string]string{}

var Fruit = []string{
	"Apple",
	"Apples",
	"Bananna",
	"Orange",
	"Pear",
	"Peach",
	"Pineapple",
	"Tomato",
	"Strawberry",
}

// Below are needles

var ShortValid = []string{
	"Apple",
	"Ban",
	"Oran",
	"Pear",
	"Peac",
	"Pinea",
	"Tom",
	"St",
}

var Ambiguous = []string{
	"App",
	"Pea",
}

var NotExistF = []string{
	"Alligator", // here
	"Apple",
	"Ban",
	"Oran",
	"Pear",
	"Peac",
	"Pinea",
	"Too",
	"St",
}

var NotExistM = []string{
	"Apple",
	"Ban",
	"Oran",
	"Pear",
	"Castle", // here
	"Peac",
	"Pinea",
	"Too",
	"St",
}

var NotExistE = []string{
	"Apple",
	"Ban",
	"Oran",
	"Pear",
	"Peac",
	"Pinea",
	"Too",
	"St",
	"Porpoise", // here
}

func TestBest(t *testing.T) {
	TestCases := []BestParams{
		BestParams{
			Haystack: Empty,
			Needles:  Empty,
			Expected: EmptyMap,
			Error:    false,
		},
		BestParams{
			Haystack: Fruit,
			Needles:  Empty,
			Expected: EmptyMap,
			Error:    false,
		},
		BestParams{
			Haystack: Empty,
			Needles:  Fruit,
			Expected: nil,
			Error:    true,
		},
		BestParams{
			Haystack: Fruit,
			Needles:  Fruit,
			Expected: map[string]string{
				"Apple":      "Apple",
				"Apples":     "Apples",
				"Bananna":    "Bananna",
				"Orange":     "Orange",
				"Pear":       "Pear",
				"Peach":      "Peach",
				"Pineapple":  "Pineapple",
				"Tomato":     "Tomato",
				"Strawberry": "Strawberry",
			},
			Error: false,
		},
		BestParams{
			Haystack: Fruit,
			Needles:  ShortValid,
			Expected: map[string]string{
				"Apple": "Apple",
				"Ban":   "Bananna",
				"Oran":  "Orange",
				"Pear":  "Pear",
				"Peac":  "Peach",
				"Pinea": "Pineapple",
				"Tom":   "Tomato",
				"St":    "Strawberry",
			},
			Error: false,
		},
		BestParams{
			Haystack: Fruit,
			Needles:  Ambiguous,
			Expected: nil,
			Error:    true,
		},
		BestParams{
			Haystack: Fruit,
			Needles:  NotExistF,
			Expected: nil,
			Error:    true,
		},
		BestParams{
			Haystack: Fruit,
			Needles:  NotExistM,
			Expected: nil,
			Error:    true,
		},
		BestParams{
			Haystack: Fruit,
			Needles:  NotExistE,
			Expected: nil,
			Error:    true,
		},
	}

	for i, v := range TestCases {
		actual, err := Best(v.Haystack, v.Needles...)
		if err != nil {
			if !v.Error {
				t.Logf("TestCase %03d unexpected error: %s", i, err)
				if e, ok := err.(*Error); ok {
					t.Log("Needle = ", e.Needle)
					t.Log(e.Matches.Slice())
				}
				t.Fail()
				return
			}

			continue
		}
		// No error from Best

		if v.Expected == nil && actual != nil {
			t.Logf("TestCase %03d: Actual != Expected\nActual  : %#v\nExpected: %#v\n", i, actual, v.Expected)
			t.Fail()
			return
		}

		for _, w := range v.Needles {
			expect := v.Expected[w]

			if actual[w].String() != expect {
				t.Logf("TestCase %03d: Actual != Expected\nActual  : %s\nExpected: %s\n", i, actual[w], expect)
				t.Fail()
				return
			}
		}

	}
}

func TestError(t *testing.T) {
	haystack := []string{
		"AB",
		"ABC",
		"ABCD",
		"ABCDE",
	}

	needle := "A"

	_, err := Best(haystack, needle)
	if err == nil {
		t.Logf("TestError: expected error")
		t.Fail()
	}

	if err, ok := err.(*Error); ok {
		if !err.MultiMatch() {
			t.Logf("TestError: expected multimatch error")
			t.Fail()
		}

		if len(err.Matches) != len(haystack) {
			t.Logf("TestError: expected #matched to be same size as haystack")
			t.Fail()
		}

		for i, v := range err.Matches.Slice() {
			if v != haystack[i] {
				t.Logf("TestError: %s != %s")
				t.Fail()
			}
		}
	} else {
		t.Logf("TestError: type assertion fail")
		t.Fail()
	}

}
