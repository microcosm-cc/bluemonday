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
	"strings"
)

// Policy encapsulates the whitelist of HTML elements and attributes that will
// be applied to the sanitised HTML.
type Policy struct {

	// Allows the <!DOCTYPE > tag to exist in the sanitized document
	allowDocType bool

	// When true, add rel="nofollow" to HTML anchors
	requireNoFollow bool

	// When true, URLs must be parseable by "net/url" url.Parse()
	requireParseableURLs bool

	// When true, u, _ := url.Parse("url"); !u.IsAbs() is permitted
	allowRelativeURLs bool

	// map[urlScheme]bool
	urlSchemes map[string]bool

	// map[htmlElementName]map[htmlAttributeName]attrPolicy
	elsAndAttrs map[string]map[string]attrPolicy

	// map[htmlAttributeName]attrPolicy
	globalAttrs map[string]attrPolicy

	// map[htmlElementName]bool
	elsWithoutAttrs map[string]bool
}

type attrPolicy struct {

	// optional pattern to match, when not nil the regexp needs to match
	// otherwise the attribute is removed
	regexp *regexp.Regexp
}

type attrPolicyBuilder struct {
	p *Policy

	attrNames []string
	regexp    *regexp.Regexp
}

// NewPolicy returns a blank policy with nothing whitelisted or permitted. This
// is the recommended way to start building a policy and you should now use
// AllowAttrs() and/or AllowElements() to construct the whitelist of HTML
// elements and attributes.
func NewPolicy() *Policy {

	p := Policy{}

	p.urlSchemes = make(map[string]bool)
	p.elsAndAttrs = make(map[string]map[string]attrPolicy)
	p.globalAttrs = make(map[string]attrPolicy)
	p.elsWithoutAttrs = make(map[string]bool)

	p.addDefaultElsWithoutAttrs()

	return &p
}

// AllowAttrs takes a range of HTML attribute names and returns an
// attribute policy builder that allows you to specify the pattern and scope of
// the whitelisted attribute.
//
// Examples:
//   AllowAttrs("title").Globally()
//   AllowAttrs("abbr").OnElements("td", "th")
//   AllowAttrs("colspan", "rowspan").Matching(
//           regexp.MustCompile("[0-9]+"),
//       ).OnElements("td", "th")
//
// The attribute policy is only added to the core policy when either Globally()
// or OnElements(...) are called.
func (p *Policy) AllowAttrs(attrNames ...string) *attrPolicyBuilder {

	abp := attrPolicyBuilder{p: p}

	for _, attrName := range attrNames {
		abp.attrNames = append(abp.attrNames, strings.ToLower(attrName))
	}

	return &abp
}

// Matching allows a regular expression to be applied to a nascent attribute
// policy, and returns the attribute policy. Calling this more than once will
// replace the existing regexp.
func (abp *attrPolicyBuilder) Matching(regex *regexp.Regexp) *attrPolicyBuilder {

	abp.regexp = regex

	return abp
}

// OnElements will bind an attribute policy to a given range of HTML elements
// and return the updated policy
func (abp *attrPolicyBuilder) OnElements(elements ...string) *Policy {

	for _, element := range elements {
		element = strings.ToLower(element)

		for _, attr := range abp.attrNames {

			if _, ok := abp.p.elsAndAttrs[element]; !ok {
				abp.p.elsAndAttrs[element] = make(map[string]attrPolicy)
			}

			ap := attrPolicy{}
			if abp.regexp != nil {
				ap.regexp = abp.regexp
			}

			abp.p.elsAndAttrs[element][attr] = ap
		}
	}

	return abp.p
}

// Globally will bind an attribute policy to all HTML elements and return the
// updated policy
func (abp *attrPolicyBuilder) Globally() *Policy {

	for _, attr := range abp.attrNames {
		if _, ok := abp.p.globalAttrs[attr]; !ok {
			abp.p.globalAttrs[attr] = attrPolicy{}
		}

		ap := attrPolicy{}
		if abp.regexp != nil {
			ap.regexp = abp.regexp
		}

		abp.p.globalAttrs[attr] = ap
	}

	return abp.p
}

// AllowElements will append HTML elements to the whitelist without applying an
// attribute policy to those elements (the elements are permitted
// sans-attributes)
func (p *Policy) AllowElements(names ...string) *Policy {

	for _, element := range names {
		element = strings.ToLower(element)

		if _, ok := p.elsAndAttrs[element]; !ok {
			p.elsAndAttrs[element] = make(map[string]attrPolicy)
		}
	}

	return p
}

// RequireNoFollowOnLinks will result in all <a> tags having a rel="nofollow"
// added to them if one does not already exist
func (p *Policy) RequireNoFollowOnLinks(require bool) *Policy {
	p.requireNoFollow = require

	return p
}

// RequireParseableURLs will result in all URLs requiring that they be parseable
// by "net/url" url.Parse()
// This applies to:
// - a.href
// - area.href
// - blockquote.cite
// - img.src
// - link.href
// - script.src
func (p *Policy) RequireParseableURLs(require bool) *Policy {

	p.requireParseableURLs = require

	return p
}

// AllowRelativeURLs enables RequireParseableURLs and then permits URLs that
// are parseable, have no schema information and url.IsAbs() returns false
// This permits local URLs
func (p *Policy) AllowRelativeURLs(require bool) *Policy {

	p.RequireParseableURLs(true)
	p.allowRelativeURLs = require

	return p
}

// AllowURLSchemes will append URL schems to the whitelist
// Example: p.AllowURLSchemes("mailto", "http", "https")
func (p *Policy) AllowURLSchemes(schemes ...string) *Policy {

	p.RequireParseableURLs(true)

	for _, scheme := range schemes {
		scheme = strings.ToLower(scheme)

		if _, ok := p.urlSchemes[scheme]; !ok {
			p.urlSchemes[scheme] = true
		}
	}

	return p
}

// AllowDocType states whether the HTML sanitised by the sanitizer is allowed to
// contain the HTML DocType tag: <!DOCTYPE HTML> or one of it's variants.
//
// The HTML spec only permits one doctype per document, and as you know how you
// are using the output of this, you know best as to whether we should ignore it
// (default) or not.
//
// If you are sanitizing a HTML fragment the default (false) is fine.
func (p *Policy) AllowDocType(allow bool) *Policy {

	p.allowDocType = allow

	return p
}

// addDefaultElsWithoutAttrs adds the HTML elements that we know are valid
// without any attributes to an internal map.
// i.e. we know that <table> is valid, but <bdo> isn't valid as the "dir" attr
// is mandatory
func (p *Policy) addDefaultElsWithoutAttrs() {

	p.elsWithoutAttrs["abbr"] = true
	p.elsWithoutAttrs["acronym"] = true
	p.elsWithoutAttrs["article"] = true
	p.elsWithoutAttrs["aside"] = true
	p.elsWithoutAttrs["audio"] = true
	p.elsWithoutAttrs["b"] = true
	p.elsWithoutAttrs["bdi"] = true
	p.elsWithoutAttrs["blockquote"] = true
	p.elsWithoutAttrs["body"] = true
	p.elsWithoutAttrs["br"] = true
	p.elsWithoutAttrs["button"] = true
	p.elsWithoutAttrs["canvas"] = true
	p.elsWithoutAttrs["caption"] = true
	p.elsWithoutAttrs["cite"] = true
	p.elsWithoutAttrs["code"] = true
	p.elsWithoutAttrs["col"] = true
	p.elsWithoutAttrs["colgroup"] = true
	p.elsWithoutAttrs["datalist"] = true
	p.elsWithoutAttrs["dd"] = true
	p.elsWithoutAttrs["del"] = true
	p.elsWithoutAttrs["details"] = true
	p.elsWithoutAttrs["dfn"] = true
	p.elsWithoutAttrs["div"] = true
	p.elsWithoutAttrs["dl"] = true
	p.elsWithoutAttrs["dt"] = true
	p.elsWithoutAttrs["em"] = true
	p.elsWithoutAttrs["fieldset"] = true
	p.elsWithoutAttrs["figcaption"] = true
	p.elsWithoutAttrs["figure"] = true
	p.elsWithoutAttrs["footer"] = true
	p.elsWithoutAttrs["h1"] = true
	p.elsWithoutAttrs["h2"] = true
	p.elsWithoutAttrs["h3"] = true
	p.elsWithoutAttrs["h4"] = true
	p.elsWithoutAttrs["h5"] = true
	p.elsWithoutAttrs["h6"] = true
	p.elsWithoutAttrs["head"] = true
	p.elsWithoutAttrs["header"] = true
	p.elsWithoutAttrs["hgroup"] = true
	p.elsWithoutAttrs["hr"] = true
	p.elsWithoutAttrs["html"] = true
	p.elsWithoutAttrs["i"] = true
	p.elsWithoutAttrs["ins"] = true
	p.elsWithoutAttrs["kbd"] = true
	p.elsWithoutAttrs["li"] = true
	p.elsWithoutAttrs["mark"] = true
	p.elsWithoutAttrs["nav"] = true
	p.elsWithoutAttrs["ol"] = true
	p.elsWithoutAttrs["optgroup"] = true
	p.elsWithoutAttrs["option"] = true
	p.elsWithoutAttrs["p"] = true
	p.elsWithoutAttrs["pre"] = true
	p.elsWithoutAttrs["q"] = true
	p.elsWithoutAttrs["rp"] = true
	p.elsWithoutAttrs["rt"] = true
	p.elsWithoutAttrs["ruby"] = true
	p.elsWithoutAttrs["s"] = true
	p.elsWithoutAttrs["samp"] = true
	p.elsWithoutAttrs["section"] = true
	p.elsWithoutAttrs["select"] = true
	p.elsWithoutAttrs["small"] = true
	p.elsWithoutAttrs["span"] = true
	p.elsWithoutAttrs["strike"] = true
	p.elsWithoutAttrs["strong"] = true
	p.elsWithoutAttrs["style"] = true
	p.elsWithoutAttrs["sub"] = true
	p.elsWithoutAttrs["summary"] = true
	p.elsWithoutAttrs["sup"] = true
	p.elsWithoutAttrs["svg"] = true
	p.elsWithoutAttrs["table"] = true
	p.elsWithoutAttrs["tbody"] = true
	p.elsWithoutAttrs["td"] = true
	p.elsWithoutAttrs["textarea"] = true
	p.elsWithoutAttrs["tfoot"] = true
	p.elsWithoutAttrs["th"] = true
	p.elsWithoutAttrs["thead"] = true
	p.elsWithoutAttrs["time"] = true
	p.elsWithoutAttrs["tr"] = true
	p.elsWithoutAttrs["tt"] = true
	p.elsWithoutAttrs["u"] = true
	p.elsWithoutAttrs["ul"] = true
	p.elsWithoutAttrs["var"] = true
	p.elsWithoutAttrs["video"] = true
	p.elsWithoutAttrs["wbr"] = true

}
