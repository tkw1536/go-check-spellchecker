package go_package_spellcheck

// spellchecker:words regexp strings golang tools analysis
import (
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var SpellcheckerPackageComments = &analysis.Analyzer{
	Name: "spellcomments",
	Doc:  "reports incorrect or missing comments for package names",
	Run: func(pass *analysis.Pass) (any, error) {
		for _, file := range pass.Files {
			// skip over files that say 'DO NOT EDIT'
			if isDoNotEdit(file) {
				continue
			}

			// if someone turned off the spellchecker for the file
			// don't even bother doing any more checking
			if isDisabled(pass, file) {
				continue
			}

			analyzeImportWordDirective(pass, file)
			analyzeWordsDirectives(pass, file)
		}

		return nil, nil
	},
}

// isDisabled checks if the spellchecker has been disabled for the given file
func isDisabled(_ *analysis.Pass, file *ast.File) bool {
	for _, group := range file.Comments {
		for _, comment := range group.List {
			// ignore multi-line comments
			if !strings.HasPrefix(comment.Text, "//") {
				continue
			}

			// parse the directive
			text := comment.Text[len("//"):]
			_, directive, _, ok := ParseSpellDirective(text)
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
	_, directive, value, ok := ParseSpellDirective(text)
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
	want := "// " + text
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
