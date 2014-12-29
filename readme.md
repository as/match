## Match
Match provides functions for completing strings based on comparison functions. A practical use for this package is completing command line arguments from substrings. 

## Usage
To start matching we need a haystack and needle. A haystack is a search area. A needle is a search key. In this package, both of these things are strings. For most of the examples below, we use this haystack:
```
	products := []string{
		"Apple",
		"Eggplant",
		"Pear",
		"Peach",
	}
```
For most matching situations, use match.Best()

```
	needles := []string{
		"App",
		"eg",
		"Peac",
	}

	// Simplified example: See next part for error checking
	matches, _ := match.Best(products, needles...)

	for _, v := range needles {
		fmt.Println(v, " -> ", matches[v])
	}
	// Prints:
	// App  ->  Apple
	// eg  ->  Eggplant
	// Peac  ->  Peach

```

## Best Practices - Error checking
An error occurs if the needle isn't found, or an ambiguous needle is used

```
	needles := []string {
		"Egg",
		"Pea", // Ambiguous: matches Pear and Peach
	}

	matches, err := match.Best(products, needles...)
	if err != nil {
		// Expecting error
		fmt.Println(err)
		// Prints:
		// Pea matches 2 fields
	}
```

Pea matches 2 fields. So How do we enumerate them?
```
	// Match has its own Error type, and it contains diagnostic
	// information. Use a type assertion to access it.

	if e, ok := err.(*match.Error); ok && e.MultiMatch() {
		for _, v := range e.Matches {
			fmt.Println("\t", v)
		}
		// Prints:
		// Pear
		// Peach
	}

```

## Best Behavior
An exact match removes ambiguity. The code below will always match "Apple" unambiguously.
```
	haystack := []string{ "Rock", "Apple", "Apples" }
	needle = "Apple"
	matches, _ := matches.Best(haystack, needle)
```

## Primitive Functions
The basic functions used by Best are exposed in the match package. Use them for raw matching functionality. An example:
```
	needle := "tokyo"
	haystack := []string{ "new york", "murray hill", "kyoto", "bejing", "oykot" }
	
	matches := match.Needle(haystack, cmpAna, needle)

	fmt.Printf("There are %d anagram matches for %s:\n", len(matches), needle)
	for _, v := range matches {
		fmt.Println(v)
	}

	// ... boilerplate omitted

	// cmpAna returns true if a and b are anagrams
	func cmpAna(a, b string) bool {
		return sum(a) == sum(b)
	}

```

## Full example
See example/example.go