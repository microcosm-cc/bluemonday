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
			// TODO: Need to add rel="nofollow" to this
			in: `<a href="http://example.org/" rel="nofollow">Hello, <b>World</b></a>` +
				`<a href="https://example.org/#!" rel="nofollow">!</a>`,
			expected: `<a href="http://example.org/">Hello, <b>World</b></a>` +
				`<a href="https://example.org/#!">!</a>`,
		},
		test{
			// TODO: Need to add rel="nofollow" to this
			in: `Hello, <b>World</b>` +
				`<a title="!" href="https://example.org/#!" rel="nofollow">!</a>`,
			expected: `Hello, <b>World</b>` +
				`<a title="!" href="https://example.org/#!">!</a>`,
		},
		// Images
		test{
			in:       `<a href="javascript:alert(1337)">foo</a>`,
			expected: `foo`,
		},
		test{
			in:       `<img src="http://example.org/foo.gif">`,
			expected: `<img src="http://example.org/foo.gif">`,
		},
		test{
			in:       `<img src="http://example.org/x.gif" alt="y" width=96 height=64 border=0>`,
			expected: `<img src="http://example.org/x.gif" alt="y" width="96" height="64">`,
		},
		test{
			in:       `<img src="http://example.org/x.png" alt="y" width="widgy" height=64 border=0>`,
			expected: `<img src="http://example.org/x.png" alt="y" height="64">`,
		},
		// Anchors
		// TODO: Need to add rel="nofollow" to all of these
		// test{
		// 	// TODO: Need to add support for local links
		// 	in:       `<a href="foo.html">Link text</a>`,
		// 	expected: `<a href="foo.html">Link text</a>`,
		// },
		// test{
		// 	// TODO: Need to add support for local links
		// 	in:       `<a href="foo.html" onclick="alert(1337)">Link text</a>`,
		// 	expected: `<a href="foo.html">Link text</a>`,
		// },
		test{
			in:       `<a href="http://example.org/x.html" onclick="alert(1337)">Link text</a>`,
			expected: `<a href="http://example.org/x.html">Link text</a>`,
		},
		test{
			in:       `<a href="https://example.org/x.html" onclick="alert(1337)">Link text</a>`,
			expected: `<a href="https://example.org/x.html">Link text</a>`,
		},
		test{
			in:       `<a href="HTTPS://example.org/x.html" onclick="alert(1337)">Link text</a>`,
			expected: `<a href="HTTPS://example.org/x.html">Link text</a>`,
		},
		// test{
		// 	// TODO: Need to add support for protocol links: //example.org
		// 	in:       `<a href="//example.org/x.html" onclick="alert(1337)">Link text</a>`,
		// 	expected: `<a href="//example.org/x.html">Link text</a>`,
		// },
		test{
			in:       `<a href="javascript:alert(1337).html" onclick="alert(1337)">Link text</a>`,
			expected: `Link text`,
		},
		test{
			in:       `<a name="header" id="header">Header text</a>`,
			expected: `Header text`,
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
