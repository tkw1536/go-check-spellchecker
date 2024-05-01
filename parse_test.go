//spellchecker:words spellchecker
package spellchecker

//spellchecker:words testing
import (
	"testing"
)

func TestParseSpellDirective(t *testing.T) {
	tests := []struct {
		text          string
		wantKeyword   string
		wantDirective string
		wantValue     string
		wantOk        bool
	}{
		// positive matches
		{text: "spellchecker:words argument", wantKeyword: "spellchecker", wantDirective: "words", wantValue: "argument", wantOk: true},
		{text: "SpEllchECker:words argument", wantKeyword: "SpEllchECker", wantDirective: "words", wantValue: "argument", wantOk: true},
		{text: "cspell : words argument", wantKeyword: "cspell", wantDirective: "words", wantValue: "argument", wantOk: true},
		{text: "cspell : words argument         ", wantKeyword: "cspell", wantDirective: "words", wantValue: "argument", wantOk: true},
		{text: "cspell : words               ", wantKeyword: "cspell", wantDirective: "words", wantValue: "", wantOk: true},
		{text: "SpEllchECker:words", wantKeyword: "SpEllchECker", wantDirective: "words", wantValue: "", wantOk: true},
		{text: "SpEllchECker:words", wantKeyword: "SpEllchECker", wantDirective: "words", wantValue: "", wantOk: true},

		// negative matches
		{text: "some-other-keyword:word argument", wantKeyword: "", wantDirective: "", wantValue: "", wantOk: false},
		{text: "spellchecker:", wantKeyword: "", wantDirective: "", wantValue: "", wantOk: false},
		{text: "spellchecker:              ", wantKeyword: "", wantDirective: "", wantValue: "", wantOk: false},
	}
	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			gotKeyword, gotDirective, gotValue, gotOk := ParseSpellComment(tt.text)
			if gotKeyword != tt.wantKeyword {
				t.Errorf("ParseSpellDirective() gotKeyword = %v, want %v", gotKeyword, tt.wantKeyword)
			}
			if gotDirective != tt.wantDirective {
				t.Errorf("ParseSpellDirective() gotDirective = %v, want %v", gotDirective, tt.wantDirective)
			}
			if gotValue != tt.wantValue {
				t.Errorf("ParseSpellDirective() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOk != tt.wantOk {
				t.Errorf("ParseSpellDirective() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
