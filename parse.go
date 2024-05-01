//spellchecker:words spellchecker
package spellchecker

//spellchecker:words strings unicode
import (
	"fmt"
	"strings"
	"unicode"
)

const correctKeyword = "spellchecker"

var keywords = []string{correctKeyword, "cSpell", "spell-checker"}

// CommentText represents the text belonging to a parsed spellchecker directive.
// A comment looks like:
//
//	keyword:directive Value
//
// A keyword must be one of 'spellchecker', 'cSpell' or 'spell-checker'.
// Keywords and directives are parsed under case-folding.
type CommentText struct {
	Keyword   string
	Directive string
	Value     string
}

// IsDirective checks if this CommentText represents the given directive.
// Directives are compared under case folding.
func (ct CommentText) IsDirective(directive string) bool {
	return strings.EqualFold(ct.Directive, directive)
}

// CommentText returns a normalized copy of comment.
//
// The picks the default keyword, and lowercases the directive.
func (ct CommentText) Normalize() CommentText {
	ct.Keyword = correctKeyword
	ct.Directive = strings.ToLower(ct.Directive)
	return ct
}

// String formats the contents of this CommentText, that is it brings it back into the form:
//
// Keyword:Directive Value
func (ct CommentText) String() string {
	if ct.Value == "" {
		return fmt.Sprintf("%s:%s", ct.Keyword, ct.Directive)
	}
	return fmt.Sprintf("%s:%s %s", ct.Keyword, ct.Directive, ct.Value)
}

// Text is like string, but includes the comment character ("//").
func (ct CommentText) Text() string {
	if ct.Value == "" {
		return fmt.Sprintf("//%s:%s", ct.Keyword, ct.Directive)
	}
	return fmt.Sprintf("//%s:%s %s", ct.Keyword, ct.Directive, ct.Value)
}

// Parse parses the given text into a comment.
// If the comment does not represent a text, returns false.
func (ct *CommentText) Parse(text string) bool {
	// if there is a newline, it can't be a directive.
	if strings.ContainsRune(text, '\n') {
		return false
	}

	// split by the ':' to check if we have a keyword
	before, after, found := strings.Cut(text, ":")
	if !found {
		return false
	}

	var ok bool
	var keyword, directive, value string

	// check the before and attempt to match a keyword
	// (under case-folding)
	before = strings.TrimSpace(before)
	for _, word := range keywords {
		if strings.EqualFold(before, word) {
			keyword = before
			ok = true
			break
		}
	}

	if !ok {
		return false
	}

	// use the directive
	directive = strings.TrimSpace(after)

	// if there is a space, split off the value of the directive
	space := strings.IndexFunc(directive, unicode.IsSpace)
	if space > 0 {
		value = strings.TrimRightFunc(directive[space+1:], unicode.IsSpace)
		directive = directive[:space]
	}

	// check that there was a directive
	if directive == "" {
		return false
	}

	// store the parsed values
	ct.Directive = directive
	ct.Keyword = keyword
	ct.Value = value
	return true
}

// ParseSpellComment parses text belonging to a spellchecker comment.
// It is expected to be of the form `keyword:directive args`.
// Each part may contain spaces at the edges and is matched under case folding.
func ParseSpellComment(text string) (keyword, directive, value string, ok bool) {
	var ct CommentText
	if !ct.Parse(text) {
		return "", "", "", false
	}
	return ct.Keyword, ct.Directive, ct.Value, true
}

// FormatDirective formats a directive into a string
func FormatDirective(directive, value string) string {
	comment := CommentText{Keyword: correctKeyword, Directive: directive, Value: value}
	return comment.Normalize().String()
}
