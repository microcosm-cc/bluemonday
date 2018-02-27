package bluemonday

import (
	"testing"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type testRemoverReader struct {
	source  TokenReader
	tagAtom atom.Atom
}

func (r *testRemoverReader) Token() (*html.Token, error) {
	t, err := r.source.Token()
	if err != nil {
		return t, err
	}
	if (t.Type == html.StartTagToken || t.Type == html.EndTagToken) && t.DataAtom == r.tagAtom {
		// Skip bold, return next token
		return r.source.Token()
	}
	return t, nil
}

func (r *testRemoverReader) Source(s TokenReader) {
	r.source = s
}

func TestTokenReader(t *testing.T) {
	p := UGCPolicy()

	input := "<p><b>A bold statement.</b></p>"
	want := "<p><b>A bold statement.</b></p>"
	got := p.Sanitize(input)
	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}

	removeBold := &testRemoverReader{tagAtom: atom.B}
	input = "<p><b>A bold statement.</b></p>"
	want = "<p>A bold statement.</p>"
	got = p.Sanitize(input, removeBold)
	if got != want {
		t.Errorf("got: %q, want: %q", got, want)
	}
}
