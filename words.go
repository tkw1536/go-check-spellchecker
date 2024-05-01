//spellchecker:words spellchecker
package spellchecker

//spellchecker:words unicode
import "unicode"

// SplitWords splits text into words.
//
// A word is a sequence of runes within the text.
// Each word may consist of upper and lowercase letters.
//
// Uppercase letters may only appear contiguously at the beginning of the word.
// "HELLOworld" is one word, whereas "HelloWorld" is two words "Hello" and "World".
//
// To count the number of words in text, use CountWords instead.
func SplitWords(text string) []string {
	// NOTE: Keep this in sync with CountWords.
	words := make([]string, 0, CountWords(text))

	lastStart := -1       // index where the last word started
	lastWasUpper := false // was the last letter of a word upper case?
	for index, char := range text {
		// not inside a word
		if lastStart == -1 {
			if unicode.IsLetter(char) { // letter starts a new word
				lastWasUpper = unicode.IsUpper(char)
				lastStart = index
			}
			continue
		}

		// word can only continue if we had a letter
		isLetter := unicode.IsLetter(char)
		if isLetter {
			// contiguous upper-case at the beginning of a word
			if lastWasUpper && unicode.IsUpper(char) {
				continue
			}

			// switched to lower-case
			if unicode.IsLower(char) {
				lastWasUpper = false
				continue
			}
			lastWasUpper = true
		}

		// word has ended => add it to the seen ones
		words = append(words, text[lastStart:index])

		// if we saw a letter, we have started a new word
		// start a new word if we saw a letter
		if isLetter {
			lastStart = index
		} else {
			lastStart = -1
		}
	}
	// finish closing the last word
	if lastStart != -1 {
		words = append(words, text[lastStart:])
	}

	// and return the words
	return words
}

// CountWords counts the number of words in the given text.
// It is an efficient version of len(SplitWords(text))
func CountWords(text string) int {
	// NOTE: Keep this in sync with SplitWords.
	words := 0

	insideWord := false   // are we currently inside a word.
	lastWasUpper := false // was the last letter of a word upper case?
	for _, char := range text {
		// not inside a word
		if !insideWord {
			if unicode.IsLetter(char) { // letter starts a new word
				lastWasUpper = unicode.IsUpper(char)
				insideWord = true
			}
			continue
		}

		// word can only continue if we had a letter
		isLetter := unicode.IsLetter(char)
		if isLetter {
			// contiguous upper-case at the beginning of a word
			if lastWasUpper && unicode.IsUpper(char) {
				continue
			}

			// switched to lower-case
			if unicode.IsLower(char) {
				lastWasUpper = false
				continue
			}
			lastWasUpper = true
		}

		// word has ended => start a new one if we had a letter
		words++
		insideWord = isLetter
	}

	// last word was not closed
	if insideWord {
		words++
	}
	return words
}
