//spellchecker:words spellchecker
package spellchecker_test

//spellchecker:words reflect testing github check spellchecker
import (
	"reflect"
	"testing"

	spellchecker "github.com/tkw1536/go-check-spellchecker"
)

func TestSplitWords(t *testing.T) {
	tests := []struct {
		text string
		want []string
	}{
		{text: "hello world", want: []string{"hello", "world"}},
		{text: "HelloWorld", want: []string{"Hello", "World"}},
		{text: "HelloWORLD", want: []string{"Hello", "WORLD"}},
		{text: "HELLOworld", want: []string{"HELLOworld"}},
		{text: "hellO/world", want: []string{"hell", "O", "world"}},
		{text: "", want: []string{}},
		{text: "///", want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			if got := spellchecker.SplitWords(tt.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitWords() = %v, want %v", got, tt.want)
			}
			if got := spellchecker.CountWords(tt.text); got != len(tt.want) {
				t.Errorf("CountWords() = %v, want %v", got, len(tt.want))
			}
		})
	}
}
