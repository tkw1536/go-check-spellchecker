package go_package_spellcheck

// spellchecker:words package spellcheck

// spellchecker:words token strings github pkglib collection golang tools analysis
import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/tkw1536/pkglib/collection"
	"golang.org/x/tools/go/analysis"
)

func getHeaderRange(file *ast.File) (from, to token.Pos) {
	from = file.FileStart
	to = file.FileEnd

	if file.Name != nil {
		from = file.Name.End()
	}

	if len(file.Decls) == 0 {
		return
	}

	decl := file.Decls[0]

	if gd, ok := decl.(*ast.GenDecl); ok && gd.Doc != nil {
		to = gd.Doc.Pos()
		return
	}

	if fd, ok := decl.(*ast.FuncDecl); ok && fd.Doc != nil {
		to = fd.Doc.Pos()
		return
	}

	to = decl.Pos()
	return
}

func analyzePackageWordDirective(pass *analysis.Pass, file *ast.File) {
	// first possible location of the comment
	commentFrom, commentUntil := getHeaderRange(file)

	// do we want to automatically delete comments?
	doDelete := len(file.Imports) != 0

	// collect all the comments
	comments := make([]*ast.Comment, 0)
	for _, group := range file.Comments {
		for _, comment := range group.List {
			start, end := comment.Pos(), comment.End()
			if start <= commentFrom || start >= commentUntil || end <= commentFrom || end >= commentUntil {
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
		if doDelete {
			for _, comment := range comments {
				removeComment(
					pass, comment,
					"'spellchecker:words' directive in header should only refer to package words (of which there are none)",
					"remove extra directive",
				)
			}
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
		want := fmt.Sprintf("\n// %s\n\n", FormatDirective("words", strings.Join(importWords, " ")))

		pass.Report(analysis.Diagnostic{
			Pos:     commentFrom,
			End:     commentUntil,
			Message: "missing 'spellchecker:words' directive in header",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "insert 'spellchecker:words' directive in header",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     commentUntil,
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
		if doDelete {
			for _, comment := range comments[1:] {
				removeComment(
					pass, comment,
					"there should be at most one 'spellchecker:words' directive in header",
					"remove extra directive",
				)
			}
		}
	}
}
