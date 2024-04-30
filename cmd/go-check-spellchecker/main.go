//spellchecker:words main
package main

//spellchecker:words github check spellchecker golang tools analysis singlechecker
import (
	spellchecker "github.com/tkw1536/go-check-spellchecker"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(spellchecker.SpellcheckerPackageComments)
}
