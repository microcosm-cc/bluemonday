package bluemonday

import (
	"regexp"
)

// policy encapsulates the whitelist of HTML elements and attributes that will
// be applied to the sanitised HTML.
type policy struct {
	elsAndAttrs map[string][]attrPolicy

	globalAttrs map[string]attrPolicy

	requireNoFollow bool
}

type attrPolicy struct {
	name   string
	regexp *regexp.Regexp
}

type attrPolicyBuilder struct {
	p *policy

	attrNames []string
	regexp    *regexp.Regexp
}

func NewPolicy() *policy {
	p := policy{}
	p.elsAndAttrs = make(map[string][]attrPolicy)
	p.globalAttrs = make(map[string]attrPolicy)

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
func (p *policy) AllowAttrs(attrNames ...string) *attrPolicyBuilder {

	abp := attrPolicyBuilder{p: p}

	for _, attrName := range attrNames {
		abp.attrNames = append(abp.attrNames, attrName)
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
func (abp *attrPolicyBuilder) OnElements(elements ...string) *policy {

	for _, element := range elements {
		for _, attr := range abp.attrNames {
			if _, ok := abp.p.elsAndAttrs[element]; !ok {
				abp.p.elsAndAttrs[element] = []attrPolicy{}
			}

			ap := attrPolicy{name: attr}
			if abp.regexp != nil {
				ap.regexp = abp.regexp
			}

			abp.p.elsAndAttrs[element] = append(
				abp.p.elsAndAttrs[element],
				ap,
			)
		}
	}

	return abp.p
}

// Globally will bind an attribute policy to all HTML elements and return the
// updated policy
func (abp *attrPolicyBuilder) Globally() *policy {

	for _, attr := range abp.attrNames {
		if _, ok := abp.p.globalAttrs[attr]; !ok {
			abp.p.globalAttrs[attr] = attrPolicy{}
		}

		ap := attrPolicy{name: attr}
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
func (p *policy) AllowElements(names ...string) *policy {

	for _, element := range names {
		if _, ok := p.elsAndAttrs[element]; !ok {
			p.elsAndAttrs[element] = []attrPolicy{}
		}
	}

	return p
}

// RequireNoFollowOnLinks will result in all <a> tags having a rel="nofollow"
// added to them if one does not already exist
func (p *policy) RequireNoFollowOnLinks() *policy {
	p.requireNoFollow = true

	return p
}
