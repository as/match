package main

import "github.com/as/match"
import "fmt"

/*
	go run example.go
*/

func main() {
	BestEx()   // Example of match.Best()
	NeedleEx() // Example of match.Needle()
}

func BestEx() {
	Products := []string{
		"Apple",
		"Eggplant",
		"Pear",
		"Peach",
	}

	GoodInput := []string{
		"App",
		"eg",
		"Peac",
	}

	matches, _ := match.Best(Products, GoodInput...)
	for _, v := range GoodInput {
		fmt.Println("\t", v, " -> ", matches[v])
	}
	// Prints:
	// App  ->  Apple
	// eg  ->  Eggplant
	// Peac  ->  Peach

	AmbigInput := []string {
		"Egg",
		"Pea", // Ambiguous: matches 2 Products
	}

	matches, err := match.Best(Products, AmbigInput...)
	if err != nil {
		// Expecting error
		fmt.Println(err)
		// Prints:
		// Pea matches 2 fields
	}

	// How do I print what fields were matched? 
	if e, ok := err.(*match.Error); ok && e.MultiMatch() {
		for _, v := range e.Matches {
			fmt.Println("\t", v)
		}
		// Prints:
		// Pear
		// Peach
	}
}

func NeedleEx() {
	needle := "tokyo"
	haystack := []string{ "new york", "murray hill", "kyoto", "bejing", "oykot" }
	
	matches := match.Needle(haystack, cmpAna, needle)

	fmt.Printf("There are %d anagram matches for %s:\n", len(matches), needle)
	for _, v := range matches {
		fmt.Println("\t", v)
	}

	matches = match.Needle(haystack, cmpRev, needle)
	fmt.Printf("There are %d reverse matches for %s:\n", len(matches), needle)
	for _, v := range matches {
		fmt.Println("\t", v)
	}
}

// cmpRev returns true if the reverse of string a
// is equal to string b
func cmpRev(a, b string) bool {
	return rev(a) == b
}

// cmpAna returns true if a and b are anagrams
func cmpAna(a, b string) bool {
	return sum(a) == sum(b)
}

func rev(str string) string {
	s := []byte(str)
	b := 0       // begin
	e := len(s)  // end + 1

	for b = range s {
		e--
		if e <= b {
			break
		}
		s[b], s[e] = s[e], s[b]
	}

	return string(s)
}

func sum(str string) int {
	s := []byte(str)
	total := 0

	for _, v := range s {
		total += int(v)
	}

	return total
}