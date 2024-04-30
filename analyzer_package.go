//spellchecker:words spellchecker
package spellchecker

//spellchecker:words token strings github pkglib collection golang tools analysis
import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/tkw1536/pkglib/collection"
	"golang.org/x/tools/go/analysis"
)

var SpellcheckerPackageComments = &analysis.Analyzer{
	Name: "spellchecker_package_comments",
	Doc:  "Checks that each package name has exactly one 'spellchecker:words' comment containing the words in the package name",
	Run: func(pass *analysis.Pass) (interface{}, error) {
		for _, file := range pass.Files {
			// skip over files that say do not edit
			if isDoNotEdit(file) || isDisabled(file) {
				continue
			}

			// check the actual words in this file
			analyzePackageWordDirective(pass, file)
		}

		return nil, nil
	},
}

func analyzePackageWordDirective(pass *analysis.Pass, file *ast.File) {

	// collect all the comments
	comments := make([]*ast.Comment, 0)
	for _, group := range file.Comments {
		for _, comment := range group.List {
			if comment.End() > file.Package {
				continue
			}
			_, ok := parseWordComment(comment)
			if !ok {
				continue
			}
			comments = append(comments, comment)
		}
	}

	// find the words in the package name, but explicitly exclude "test"
	importWords := collection.Deduplicate(SplitWords(file.Name.Name))
	importWords = collection.KeepFunc(importWords, func(word string) bool { return len(word) >= minWordLength && !strings.EqualFold(word, "test") })

	// want no comment, but there is one
	if len(importWords) == 0 {
		for _, comment := range comments {
			removeComment(
				pass, comment,
				"'spellchecker:words' directive in header should only refer to package words (of which there are none)",
				"remove extra directive",
			)
		}

		return

	} else if len(comments) > 0 {
		wantComment(
			FormatDirective("words", strings.Join(importWords, " ")),
			pass, comments[0],
			"'spellchecker:words' directive in header doc should only contain package words",
			"update package words directive",
		)
	} else {
		want := fmt.Sprintf("//%s\n", FormatDirective("words", strings.Join(importWords, " ")))
		if len(comments) > 1 {
			want = "//\n" + want
		}

		pass.Report(analysis.Diagnostic{
			Pos:     file.Package,
			Message: "missing 'spellchecker:words' directive for package documentation",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "insert 'spellchecker:words' directive in package header",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     file.Package,
							End:     token.NoPos,
							NewText: []byte(want),
						},
					},
				},
			},
		})
	}

	// extra spellchecker:words comments should be remove
	if len(comments) > 1 {
		for _, comment := range comments[1:] {
			removeComment(
				pass, comment,
				"there should be at most one 'spellchecker:words' directive in header",
				"remove extra directive",
			)
		}
	}
}
