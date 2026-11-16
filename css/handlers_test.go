package css

import (
	"strings"
	"testing"
)

// TestFontHandlerWhitespace covers the repeated whitespace from #225.
func TestFontHandlerWhitespace(t *testing.T) {
	font := "normal\n" + strings.Repeat(" ", 12) + "16px\n  Roboto"
	if !FontHandler(font) {
		t.Fatalf("FontHandler(%q) = false, want true", font)
	}
}
