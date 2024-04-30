//spellchecker:words main
package main

//spellchecker:words github check spellchecker golang tools analysis multichecker
import (
	spellchecker "github.com/tkw1536/go-check-spellchecker"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		spellchecker.SpellcheckerPackageComments,
		spellchecker.SpellcheckerImportComments,
		spellchecker.SpellcheckerWords,
	)
}
