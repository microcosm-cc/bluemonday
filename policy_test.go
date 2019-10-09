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
	policy := NewPolicy().AllowElementsContent("iframe", "script")

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

func TestElementsMatching(t *testing.T) {
	tests := map[string]struct{
		regexs []*regexp.Regexp
		in string
		expected string
	}{
		"Self closing tags with regex prefix should strip any that do not match":{
			regexs: []*regexp.Regexp{
				regexp.MustCompile(`^my-element-`),
			},
			in:`<div><my-element-demo-one /><my-element-demo-two /><not-my-element-demo-one /></div>`,
			  expected:`<div><my-element-demo-one/><my-element-demo-two/></div>`,
		},"Standard elements regex prefix should strip any that do not match":{
			regexs: []*regexp.Regexp{
				regexp.MustCompile(`^my-element-`),
			},
			in:`<div><my-element-demo-one data-test='test'></my-element-demo-one ><my-element-demo-two></my-element-demo-two><not-my-element-demo-one></not-my-element-demo-one></div>`,
			expected:`<div><my-element-demo-one></my-element-demo-one ><my-element-demo-two></my-element-demo-two></div>`,
		},
	}

	for name, test := range tests {
		policy := NewPolicy().AllowElements("div")
		policy.AllowNoAttrs().OnElementsMatching(test.regexs[0])
		policy.AllowDataAttributes()
		for _, regex := range  test.regexs{
			policy.AllowElementsMatching(regex)
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

