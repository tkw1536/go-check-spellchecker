//spellchecker:words spellchecker
package spellchecker

//spellchecker:words regexp strings golang tools analysis
import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// isDisabled checks if the spellchecker has been disabled for the given file
func isDisabled(file *ast.File) bool {
	for _, group := range file.Comments {
		for _, comment := range group.List {
			// ignore multi-line comments
			if !strings.HasPrefix(comment.Text, "//") {
				continue
			}

			// parse the directive
			text := comment.Text[len("//"):]
			_, directive, _, ok := ParseSpellComment(text)
			if !ok {
				continue
			}
			if strings.EqualFold(directive, "disable") {
				return true
			}
		}
	}
	return false
}

var doNotEdit = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

// isDoNotEdit checks if the given file has a 'DO NOT EDIT' comment at the top.
func isDoNotEdit(file *ast.File) bool {
	if len(file.Comments) == 0 {
		return false
	}

	lst := file.Comments[0].List
	if len(lst) == 0 {
		return false
	}
	return doNotEdit.MatchString(lst[0].Text)
}

func parseWordComment(comment *ast.Comment) ([]string, bool) {
	// ignore multi-line comments
	if !strings.HasPrefix(comment.Text, "//") {
		return nil, false
	}

	// parse the directive
	text := comment.Text[len("//"):]
	_, directive, value, ok := ParseSpellComment(text)
	if !ok {
		return nil, false
	}

	// ignore everything that isn't a word
	if !strings.EqualFold(directive, "words") {
		return nil, false
	}

	// split the words
	return SplitWords(value), true
}

// edits for specific comments

func removeComment(pass *analysis.Pass, comment *ast.Comment, message string, fix string) {
	pass.Report(analysis.Diagnostic{
		Pos:     comment.Pos(),
		Message: message,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fix,
				TextEdits: []analysis.TextEdit{
					{
						Pos:     comment.Pos(),
						End:     comment.End(),
						NewText: nil,
					},
				},
			},
		},
	})
}

// wantComment ensure that comment has the given text
func wantComment(text string, pass *analysis.Pass, comment *ast.Comment, message string, fix string) bool {
	want := "//" + text
	if comment.Text == want {
		return true
	}

	pass.Report(analysis.Diagnostic{
		Pos:     comment.Pos(),
		Message: message,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: fix,
				TextEdits: []analysis.TextEdit{
					{
						Pos:     comment.Pos(),
						End:     comment.End(),
						NewText: []byte(want),
					},
				},
			},
		},
	})
	return false
}
