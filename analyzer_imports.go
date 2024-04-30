//spellchecker:words spellchecker
package spellchecker

//spellchecker:words token strings golang tools analysis
import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
)

// analyzeImportWordDirective analyzes all import GenDecls for imports.
func analyzeImportWordDirective(pass *analysis.Pass, file *ast.File) {
	for _, decl := range file.Decls {
		// ensure that we have a generic declaration
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.IMPORT {
			continue
		}

		// we must have import specs!
		specs := make([]*ast.ImportSpec, len(gen.Specs))
		for idx, spec := range gen.Specs {
			specs[idx] = spec.(*ast.ImportSpec)
		}

		// deal with the spec
		doImportSpecWords(pass, gen, specs)
	}

}

// doImportSpecWords handles words for the provided import declaration
func doImportSpecWords(pass *analysis.Pass, decl *ast.GenDecl, specs []*ast.ImportSpec) {
	// if there are no specs, we don't need to do anything
	if len(specs) == 0 {
		return
	}

	// get the comments in the documentation of the import group
	var comments []*ast.Comment
	if decl.Doc != nil {
		comments = make([]*ast.Comment, 0, len(decl.Doc.List))
		for _, comment := range decl.Doc.List {
			_, ok := parseWordComment(comment)
			if !ok {
				continue
			}
			comments = append(comments, comment)
		}
	}

	// find the comment we want
	importWords := makeImportWords(specs)

	// want no comment, but there is one
	if len(importWords) == 0 {
		for _, comment := range comments {
			removeComment(
				pass, comment,
				"'spellchecker:words' directive in import doc should only refer to import words",
				"remove extra directive",
			)
		}

		return

	} else if len(comments) > 0 {
		wantComment(
			FormatDirective("words", strings.Join(importWords, " ")),
			pass, comments[0],
			"'spellchecker:words' directive in import doc should only contain import words",
			"update import words directive",
		)
	} else {
		want := fmt.Sprintf("//%s\n", FormatDirective("words", strings.Join(importWords, " ")))

		// if the documentation has a preceding '//' comment, then insert a newline
		if decl.Doc != nil && len(decl.Doc.List) > 0 && strings.HasPrefix(decl.Doc.List[0].Text, "//") {
			want = "//\n" + want
		}

		pass.Report(analysis.Diagnostic{
			Pos:     decl.Pos(),
			Message: "missing 'spellchecker:words' directive in import doc",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "insert 'spellchecker:words' directive in import doc",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     decl.Pos(),
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
				"there should be at most one 'spellchecker:words' directive in import doc",
				"remove extra directive",
			)
		}
	}
}

// TODO: make this configurable
const minWordLength = 4

func makeImportWords(imports []*ast.ImportSpec) []string {
	// guess the number of words for all the imports
	sizeGuess := 5 * len(imports)

	// record the words we need for the import
	importWords := make([]string, 0, sizeGuess)
	hadImportWords := make(map[string]struct{}, sizeGuess) // for de-duping

	// a function to add some text to the known import words
	add := func(text string) {
		for _, word := range SplitWords(text) {
			if len(word) < minWordLength {
				continue
			}
			if _, ok := hadImportWords[word]; ok {
				continue
			}
			hadImportWords[word] = struct{}{}
			importWords = append(importWords, word)
		}
	}

	// process all of the imports automatically
	for _, pkg := range imports {
		add(pkg.Path.Value)
		if pkg.Name == nil {
			continue
		}
		add(pkg.Name.Name)
	}

	return importWords
}
