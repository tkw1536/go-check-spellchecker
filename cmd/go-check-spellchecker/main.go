//spellchecker:words main
package main

//spellchecker:words check spellchecker golang tools analysis multichecker
import (
	spellchecker "go.tkw01536.de/go-check-spellchecker"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		spellchecker.SpellcheckerPackageComments,
		spellchecker.SpellcheckerImportComments,
		spellchecker.SpellcheckerWords,
	)
}
