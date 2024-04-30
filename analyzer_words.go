//spellchecker:words spellchecker
package spellchecker

//spellchecker:words strings golang tools analysis
import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// analyzeWordsDirectives processes all words directives for the given file
func analyzeWordsDirectives(pass *analysis.Pass, file *ast.File) {
	for _, group := range file.Comments {
		for _, comment := range group.List {
			words, ok := parseWordComment(comment)
			if !ok {
				continue
			}

			// complain if there are no words, we should remove it
			if len(words) == 0 {
				removeComment(
					pass, comment,
					"empty words directive",
					"remove comment",
				)
				continue
			}

			// ensure that the directive is spelled correctly
			wantComment(
				FormatDirective("words", strings.Join(words, " ")),
				pass, comment,
				"improperly formatted 'words' directive",
				"reformat 'words' directive",
			)
		}
	}
}
