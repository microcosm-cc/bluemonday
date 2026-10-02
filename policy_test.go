// Copyright (c) 2014, David Kitchen <david@buro9.com>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//   list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//   this list of conditions and the following disclaimer in the documentation
//   and/or other materials provided with the distribution.
//
// * Neither the name of the organisation (Microcosm) nor the names of its
//   contributors may be used to endorse or promote products derived from
//   this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package bluemonday

import (
	"regexp"
	"testing"
)

func TestAllowElementsContent(t *testing.T) {
	policy := NewPolicy().AllowElementsContent("iframe", "script").AllowUnsafe(true)

	tests := []test{
		{
			in:       "<iframe src='http://url.com/test'>this is fallback content</iframe>",
			expected: "this is fallback content",
		},
		{
			in:       "<script>var a = 10; alert(a);</script>",
			expected: "var a = 10; alert(a);",
		},
	}

	for ii, test := range tests {
		out := policy.Sanitize(test.in)
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

func TestAllowElementsMatching(t *testing.T) {
	tests := map[string]struct {
		policyFn func(policy *Policy)
		in       string
		expected string
	}{
		"Self closing tags with regex prefix should strip any that do not match": {
			policyFn: func(policy *Policy) {
				policy.AllowElementsMatching(regexp.MustCompile(`^my-element-`))
			},
			in: `<div>
							<my-element-demo-one data-test="test" my-attr="test"/>
							<my-element-demo-two data-test="test"/>
							<not-my-element-demo-one data-test="test"/>
						</div>`,
			expected: `<div>
							<my-element-demo-one data-test="test"/>
							<my-element-demo-two data-test="test"/>
							
						</div>`,
		}, "Standard elements regex prefix should strip any that do not match": {
			policyFn: func(policy *Policy) {
				policy.AllowElementsMatching(regexp.MustCompile(`^my-element-`))
			},
			in: `<div>
							<my-element-demo-one data-test="test"></my-element-demo-one>
							<my-element-demo-two data-test="test"></my-element-demo-two>
							<not-my-element-demo-one data-test="test"></not-my-element-demo-one>
						</div>`,
			expected: `<div>
							<my-element-demo-one data-test="test"></my-element-demo-one>
							<my-element-demo-two data-test="test"></my-element-demo-two>
							
						</div>`,
		}, "Self closing tags with regex prefix and custom attr should strip any that do not match": {
			policyFn: func(policy *Policy) {
				policy.AllowElementsMatching(regexp.MustCompile(`^my-element-`))
				policy.AllowElements("not-my-element-demo-one")
			},
			in: `<div>
							<my-element-demo-one data-test="test" my-attr="test"/>
							<my-element-demo-two data-test="test"/>
							<not-my-element-demo-one data-test="test"/>
						</div>`,
			expected: `<div>
							<my-element-demo-one data-test="test"/>
							<my-element-demo-two data-test="test"/>
							<not-my-element-demo-one data-test="test"/>
						</div>`,
		},
	}

	for name, test := range tests {
		policy := NewPolicy().AllowElements("div")
		policy.AllowDataAttributes()
		if test.policyFn != nil {
			test.policyFn(policy)
		}
		out := policy.Sanitize(test.in)
		if out != test.expected {
			t.Errorf(
				"test %s failed;\ninput   : %s\noutput  : %s\nexpected: %s",
				name,
				test.in,
				out,
				test.expected,
			)
		}
	}
}

func TestAttrOnElementMatching(t *testing.T) {
	tests := map[string]struct {
		policyFn func(policy *Policy)
		in       string
		expected string
	}{
		"Self closing tags with regex prefix should strip any that do not match with custom attr": {
			policyFn: func(policy *Policy) {
				policy.AllowAttrs("my-attr").OnElementsMatching(regexp.MustCompile(`^my-element-`))
			},
			in: `<div>
							<my-element-demo-one data-test="test" my-attr="test"/>
							<my-element-demo-two data-test="test" other-attr="test"/>
							<not-my-element-demo-one data-test="test"/>
						</div>`,
			expected: `<div>
							<my-element-demo-one data-test="test" my-attr="test"/>
							<my-element-demo-two data-test="test"/>
							
						</div>`,
		}, "Standard elements regex prefix should strip any that do not match": {
			policyFn: func(policy *Policy) {
				policy.AllowAttrs("my-attr").OnElementsMatching(regexp.MustCompile(`^my-element-`))
			},
			in: `<div>
							<my-element-demo-one data-test="test" my-attr="test" other-attr="test"></my-element-demo-one>
							<my-element-demo-two data-test="test" other-attr="test"></my-element-demo-two>
							<not-my-element-demo-one data-test="test" other-attr="test"></not-my-element-demo-one>
						</div>`,
			expected: `<div>
							<my-element-demo-one data-test="test" my-attr="test"></my-element-demo-one>
							<my-element-demo-two data-test="test"></my-element-demo-two>
							
						</div>`,
		}, "Specific element rule defined should override matching rules": {
			policyFn: func(policy *Policy) {
				// specific element rule
				policy.AllowAttrs("my-other-attr").OnElements("my-element-demo-one")
				// matched rule takes lower precedence
				policy.AllowAttrs("my-attr").OnElementsMatching(regexp.MustCompile(`^my-element-`))
			},
			in: `<div>
							<my-element-demo-one data-test="test" my-attr="test" my-other-attr="test"/>
							<my-element-demo-two data-test="test" my-attr="test" my-other-attr="test"/>
							<not-my-element-demo-one data-test="test"/>
						</div>`,
			expected: `<div>
							<my-element-demo-one data-test="test" my-other-attr="test"/>
							<my-element-demo-two data-test="test" my-attr="test"/>
							
						</div>`,
		},
	}

	for name, test := range tests {
		policy := NewPolicy().AllowElements("div")
		policy.AllowDataAttributes()
		if test.policyFn != nil {
			test.policyFn(policy)
		}
		out := policy.Sanitize(test.in)
		if out != test.expected {
			t.Errorf(
				"test %s failed;\ninput   : %s\noutput  : %s\nexpected: %s",
				name,
				test.in,
				out,
				test.expected,
			)
		}
	}
}

func TestStyleOnElementMatching(t *testing.T) {
	tests := map[string]struct {
		policyFn func(policy *Policy)
		in       string
		expected string
	}{
		"Self closing tags with style policy matching prefix should strip any that do not match with custom attr": {
			policyFn: func(policy *Policy) {
				policy.AllowAttrs("style").
					OnElementsMatching(regexp.MustCompile(`^my-element-`))
				policy.AllowStyles("color", "mystyle").
					MatchingHandler(func(s string) bool {
						return true
					}).OnElementsMatching(regexp.MustCompile(`^my-element-`))
			},
			in: `<div>
							<my-element-demo-one data-test="test" style="color:#ffffff;mystyle:test;other:value"/>
							<my-element-demo-two data-test="test" other-attr="test" style="other:value"/>
							<not-my-element-demo-one data-test="test"/>
						</div>`,
			expected: `<div>
							<my-element-demo-one data-test="test" style="color: #ffffff; mystyle: test"/>
							<my-element-demo-two data-test="test"/>
							
						</div>`,
		}, "Standard elements with style policy and matching elements should strip any styles not allowed": {
			policyFn: func(policy *Policy) {
				policy.AllowAttrs("style").
					OnElementsMatching(regexp.MustCompile(`^my-element-`))
				policy.AllowStyles("color", "mystyle").
					MatchingHandler(func(s string) bool {
						return true
					}).OnElementsMatching(regexp.MustCompile(`^my-element-`))
			},
			in: `<div>
							<my-element-demo-one data-test="test" style="color:#ffffff;mystyle:test;other:value"></my-element-demo-one>
							<my-element-demo-two data-test="test" other-attr="test" style="other:value"></my-element-demo-two>
							<not-my-element-demo-one data-test="test" other-attr="test"></not-my-element-demo-one>
						</div>`,
			expected: `<div>
							<my-element-demo-one data-test="test" style="color: #ffffff; mystyle: test"></my-element-demo-one>
							<my-element-demo-two data-test="test"></my-element-demo-two>
							
						</div>`,
		}, "Specific element rule defined should override matching rules": {
			policyFn: func(policy *Policy) {
				policy.AllowAttrs("style").
					OnElements("my-element-demo-one")
				policy.AllowStyles("color", "mystyle").
					MatchingHandler(func(s string) bool {
						return true
					}).OnElements("my-element-demo-one")

				policy.AllowAttrs("style").
					OnElementsMatching(regexp.MustCompile(`^my-element-`))
				policy.AllowStyles("color", "customstyle").
					MatchingHandler(func(s string) bool {
						return true
					}).OnElementsMatching(regexp.MustCompile(`^my-element-`))
			},
			in: `<div>
							<my-element-demo-one data-test="test" style="color:#ffffff;mystyle:test;other:value"/>
							<my-element-demo-two data-test="test" style="color:#ffffff;mystyle:test;customstyle:value"/>
							<not-my-element-demo-one data-test="test"/>
						</div>`,
			expected: `<div>
							<my-element-demo-one data-test="test" style="color: #ffffff; mystyle: test"/>
							<my-element-demo-two data-test="test" style="color: #ffffff; customstyle: value"/>
							
						</div>`,
		},
	}

	for name, test := range tests {
		policy := NewPolicy().AllowElements("div")
		policy.AllowDataAttributes()
		if test.policyFn != nil {
			test.policyFn(policy)
		}
		out := policy.Sanitize(test.in)
		if out != test.expected {
			t.Errorf(
				"test %s failed;\ninput   : %s\noutput  : %s\nexpected: %s",
				name,
				test.in,
				out,
				test.expected,
			)
		}
	}
}

func TestAllowImportantInStyles(t *testing.T) {
	tests := []struct {
		name     string
		policyFn func(p *Policy)
		in       string
		expected string
	}{
		{
			name: "Default policy strips !important",
			policyFn: func(p *Policy) {
				p.AllowElements("div")
				p.AllowAttrs("style").OnElements("div")
				p.AllowStyles("font-family", "color").OnElements("div")
			},
			in:       `<div style="font-family: Arial, sans-serif !important;"></div>`,
			expected: `<div style="font-family: Arial, sans-serif"></div>`,
		},
		{
			name: "AllowImportantInStyles retains !important on allowed element styles",
			policyFn: func(p *Policy) {
				p.AllowElements("div")
				p.AllowAttrs("style").OnElements("div")
				p.AllowStyles("font-family", "color").OnElements("div")
				p.AllowImportantInStyles(true)
			},
			in:       `<div style="font-family: Arial, sans-serif !important;"></div>`,
			expected: `<div style="font-family: Arial, sans-serif !important"></div>`,
		},
		{
			name: "AllowImportantInStyles does not add !important when not present",
			policyFn: func(p *Policy) {
				p.AllowElements("div")
				p.AllowAttrs("style").OnElements("div")
				p.AllowStyles("font-family", "color").OnElements("div")
				p.AllowImportantInStyles(true)
			},
			in:       `<div style="font-family: Arial, sans-serif;"></div>`,
			expected: `<div style="font-family: Arial, sans-serif"></div>`,
		},
		{
			name: "AllowImportantInStyles mixed declarations",
			policyFn: func(p *Policy) {
				p.AllowElements("div")
				p.AllowAttrs("style").OnElements("div")
				p.AllowStyles("font-family", "color").OnElements("div")
				p.AllowImportantInStyles(true)
			},
			in:       `<div style="font-family: Arial, sans-serif !important; color: red;"></div>`,
			expected: `<div style="font-family: Arial, sans-serif !important; color: red"></div>`,
		},
		{
			name: "AllowImportantInStyles both declarations important",
			policyFn: func(p *Policy) {
				p.AllowElements("div")
				p.AllowAttrs("style").OnElements("div")
				p.AllowStyles("font-family", "color").OnElements("div")
				p.AllowImportantInStyles(true)
			},
			in:       `<div style="font-family: Arial, sans-serif !important; color: red !important;"></div>`,
			expected: `<div style="font-family: Arial, sans-serif !important; color: red !important"></div>`,
		},
		{
			name: "AllowImportantInStyles globally allowed styles",
			policyFn: func(p *Policy) {
				p.AllowElements("span")
				p.AllowAttrs("style").Globally()
				p.AllowStyles("color").Globally()
				p.AllowImportantInStyles(true)
			},
			in:       `<span style="color: blue !important;">test</span>`,
			expected: `<span style="color: blue !important">test</span>`,
		},
		{
			name: "AllowImportantInStyles strips disallowed properties even with !important",
			policyFn: func(p *Policy) {
				p.AllowElements("div")
				p.AllowAttrs("style").OnElements("div")
				p.AllowStyles("color").OnElements("div")
				p.AllowImportantInStyles(true)
			},
			in:       `<div style="evil: bad !important; color: red !important;"></div>`,
			expected: `<div style="color: red !important"></div>`,
		},
		{
			name: "AllowImportantInStyles strips invalid property value with !important",
			policyFn: func(p *Policy) {
				p.AllowElements("div")
				p.AllowAttrs("style").OnElements("div")
				p.AllowStyles("color").OnElements("div")
				p.AllowImportantInStyles(true)
			},
			in:       `<div style="color: invalid-color-123 !important;"></div>`,
			expected: `<div></div>`,
		},
		{
			name: "AllowImportantInStyles disabled after enabled strips !important",
			policyFn: func(p *Policy) {
				p.AllowElements("div")
				p.AllowAttrs("style").OnElements("div")
				p.AllowStyles("font-family").OnElements("div")
				p.AllowImportantInStyles(true)
				p.AllowImportantInStyles(false)
			},
			in:       `<div style="font-family: Arial, sans-serif !important;"></div>`,
			expected: `<div style="font-family: Arial, sans-serif"></div>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPolicy()
			tt.policyFn(p)
			out := p.Sanitize(tt.in)
			if out != tt.expected {
				t.Errorf("\ninput   : %s\noutput  : %s\nexpected: %s", tt.in, out, tt.expected)
			}
		})
	}
}
