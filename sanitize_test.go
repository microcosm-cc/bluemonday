package bluemonday

import (
	"testing"
)

func TestUGCPolicy(t *testing.T) {

	type test struct {
		in       string
		expected string
	}

	tests := []test{
		// Simple formatting
		test{in: "Hello, World!", expected: "Hello, World!"},
		test{in: "Hello, <b>World</b>!", expected: "Hello, <b>World</b>!"},
		// Blocks and formatting
		test{
			in:       "<p>Hello, <b onclick=alert(1337)>World</b>!</p>",
			expected: "<p>Hello, <b>World</b>!</p>",
		},
		test{
			in:       "<p onclick=alert(1337)>Hello, <b>World</b>!</p>",
			expected: "<p>Hello, <b>World</b>!</p>",
		},
		// Inline tags featuring globals
		test{
			in: `<a href="http://foo.com/" rel="nofollow">Hello, <b>World</b></a>` +
				`<a href="https://foo.com/#!" rel="nofollow">!</a>`,
			expected: `<a href="http://foo.com/">Hello, <b>World</b></a>` +
				`<a href="https://foo.com/#!">!</a>`,
		},
		test{
			in: `Hello, <b>World</b>` +
				`<a title="!" href="https://foo.com/#!" rel="nofollow">!</a>`,
			expected: `Hello, <b>World</b>` +
				`<a title="!" href="https://foo.com/#!">!</a>`,
		},
		// Images
		test{
			in:       `<a href="javascript:alert(1337)">foo</a>`,
			expected: `foo`,
		},
	}

	p := UGCPolicy()

	for ii, test := range tests {
		out, err := p.Sanitize(test.in)
		if err != nil {
			t.Error(err)
		}
		if out != test.expected {
			t.Errorf(
				"test %d failed;\ninput   : %s\noutput  : %s\nexpected: %s",
				ii,
				test.in,
				out,
				test.expected,
			)
		}
	}
}
