package match

// Match provides functions for completing strings based on comparison
// functions. A practical use for this package is completing command
// line arguments from substrings. See example/example.go for details.

import (
	"fmt"
	"strings"
)

// Terminology: Needles & Haystacks
// The terms "haystack" and "needle" are used to describe
// search area and search key respectively.

// Match is a record storing the output of a match function.
// The index reffers to the index in the haystack where
// the needle was found. The data points to what was found
// there.
type Match struct {
	Index int
	data  *string
}

// Matches is a pluralization of Match. Matches implements
// methods for aggregate operation.
type Matches []Match

// Error holds a package-specific error state along with
// an error message. See example/example.go for idiomatic
// usage.
type Error struct {
	msg     string
	Needle  string
	Matches Matches
}

// String returns the match data
func (m Match) String() string {
	if m.data == nil {
		return ""
	}

	return *m.data
}

// Exists checks for the existance of a needle using the cmp
// function and returns a pointer to the Match if its found.
func (m Matches) Exists(needle string, cmp func(string, string) bool) (bool, *Match) {
	for i, v := range m {
		if cmp(v.String(), needle) {
			return true, &m[i]
		}
	}
	return false, nil
}

// Slice returns an ordered string slice of matches
func (m Matches) Slice() []string {
	s := make([]string, len(m))
	for i, v := range m {
		s[i] = v.String()
	}
	return s
}

// Swap swaps two matches by their indices
func (m Matches) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Len returns the number of matches
func (m Matches) Len() int {
	return len(m)
}

// Needle finds a needle in a haystack. Cmp should be a function that
// compares the needle to the haystack.
func Needle(hays []string, cmp func(string, string) bool, needle string) (ma Matches) {
	for i, v := range hays {
		if cmp(v, needle) {
			x := Match{i, &hays[i]} // create new match
			ma = append(ma, x)      // append it to our matches
		}
	}
	return ma
}

// Needles is the plural of Needle, finding multiple needles in the
// haystack instead of one. Because Needles's last parameter is variadic,
// it returns a matrix of values in the form of a Matches slice
func Needles(hays []string, cmp func(string, string) bool, needles ...string) []Matches {
	ma := make([]Matches, 0, len(needles))

	for _, n := range needles {
		m := Needle(hays, cmp, n)
		ma = append(ma, m)
	}

	return ma
}

// NeedlesMap is like Needles, except it returns a map containing the needle name
// as the key, and Matches as the value
func NeedlesMap(hays []string, cmp func(string, string) bool, needles ...string) map[string]Matches {
	mslice := Needles(hays, cmp, needles...)
	mmap := make(map[string]Matches)

	for i, n := range needles {
		mmap[n] = mslice[i]
	}

	return mmap
}

// filter filters a map of Matches by looking at every Matches value
// in the map and moving any exact matches to the front of each Matches
// slice. If a Matches value contains more than one match, but niether
// are exact, filter returns an error.
func filter(matches map[string]Matches) error {
	for k, v := range matches {
		if v.Len() == 1 {
			continue
		}

		if v.Len() == 0 {
			return Errorf("%s matches no fields", k)
		}

		exact, eptr := v.Exists(k, Cmp)
		if !exact {
			err := Errorf("%s matches %d fields", k, v.Len())
			err.Needle = k
			err.Matches = v
			return err
		}
		v.Swap(0, eptr.Index)
	}

	return nil
}

// Best maps every needle to its best match in the haystack. The Best
// algorithm favors exact matches over lazy matches, and returns
// an error if a needle matches two or more haystack elements without
// a single exact match.
//
// By default, Best uses CmpPrefix as its comparison function for lazy
// matching. Use BestFunc to provide your own comparison function.
func Best(hays []string, needles ...string) (map[string]Match, error) {
	return BestFunc(hays, CmpPrefix, needles...)
}

// BestFunc is like Best, except it provides a third parameter: a function
// comparing two strings.
//
// Note: BestFunc still favors exact matches over lazy matches
func BestFunc(hays []string, lazyCmp func(string, string) bool, needles ...string) (map[string]Match, error) {
	var (
		mm map[string]Matches // Multiple matches per key
		sm map[string]Match   // Single match per key
	)

	// Generate matches
	mm = NeedlesMap(hays, lazyCmp, needles...)
	err := filter(mm)
	if err != nil {
		return nil, err
	}

	sm = make(map[string]Match)
	// Reduce multiple matches to single matches
	for k := range mm {
		sm[k] = mm[k][0]
	}

	return sm, nil
}

// Errorf formats according to a format specifier and
// returns the string as an Error
func Errorf(format string, a ...interface{}) *Error {
	e := new(Error)
	e.msg = fmt.Sprintf(format, a...)
	return e
}

// MultiMatch returns true if the error occured due to
// multiple matches being found
func (e *Error) MultiMatch() bool {
	return len(e.Matches) > 1
}

// Error returns an error message
func (e *Error) Error() string {
	return e.msg
}

// Cmp returns true if 'a' and 'b' are equal
func Cmp(a, b string) bool {
	return a == b
}

// CmpLower returns true if 'a' and 'b' are equal in lower case
func CmpLower(a, b string) bool {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	return a == b
}

// CmpPrefix returns true if 'b' string  is a prefix of 'a' string
// when both strings are lowercase.
func CmpPrefix(a, b string) bool {
	a = strings.ToLower(a)
	b = strings.ToLower(b)

	if strings.Index(a, b) == 0 {
		return true
	}

	return false
}
