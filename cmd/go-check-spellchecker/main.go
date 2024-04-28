package main

// spellchecker:words github package spellcheck golang tools analysis singlechecker
import (
	go_package_spellcheck "github.com/tkw1536/go-package-spellcheck"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(go_package_spellcheck.SpellcheckerPackageComments)
}
