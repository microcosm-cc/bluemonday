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
	"bytes"
	"io"
	"net/url"
	"strings"

	"code.google.com/p/go.net/html"
)

// Sanitize takes a string that contains a HTML fragment or document and applies
// the given policy whitelist.
// It returns a HTML string that has been sanitized by the policy or an empty
// string if an error has occurred (most likely as a consequence of extremely
// malformed input)
func (p *Policy) Sanitize(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	// It is possible that the developer has created the policy via:
	//   p := bluemonday.Policy{}
	// rather than:
	//   p := bluemonday.NewPolicy()
	// If this is the case, and if they haven't yet triggered an action that
	// would initiliaze the maps, then we need to do that.
	p.init()

	var cleanHTML bytes.Buffer
	tokenizer := html.NewTokenizer(strings.NewReader(s))

	ignore := false
	skipClosingTag := false
	for {
		if tokenizer.Next() == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				// End of input means end of processing
				return cleanHTML.String()
			}

			// Raw tokenizer error
			return ""
		}

		token := tokenizer.Token()
		switch token.Type {
		case html.DoctypeToken:

			if p.allowDocType {
				cleanHTML.WriteString(token.String())
			}

		case html.CommentToken:

			// Comments are ignored by default

		case html.StartTagToken:

			aps, ok := p.elsAndAttrs[token.Data]
			if !ok {
				ignore = true
				break
			}

			if len(token.Attr) != 0 {
				token.Attr = p.sanitizeAttrs(token.Data, token.Attr, aps)
			}

			if len(token.Attr) == 0 {
				if !p.allowNoAttrs(token.Data) {
					skipClosingTag = true
					break
				}
			}

			cleanHTML.WriteString(token.String())

		case html.EndTagToken:

			if skipClosingTag {
				skipClosingTag = false
				break
			}

			_, ok := p.elsAndAttrs[token.Data]
			if !ok {
				ignore = false
				break
			}

			cleanHTML.WriteString(token.String())

		case html.SelfClosingTagToken:

			aps, ok := p.elsAndAttrs[token.Data]
			if !ok {
				break
			}

			if len(token.Attr) != 0 {
				token.Attr = p.sanitizeAttrs(token.Data, token.Attr, aps)
			}

			if len(token.Attr) == 0 && !p.allowNoAttrs(token.Data) {
				break
			}

			cleanHTML.WriteString(token.String())

		case html.TextToken:

			if !ignore {
				cleanHTML.WriteString(token.String())
			}

		default:
			// A token that didn't exist in the html package when we wrote this
			return ""
		}
	}
}

// sanitizeAttrs takes a set of element attribute policies and the global
// attribute policies and applies them to the []html.Attribute returning a set
// of html.Attributes that match the policies
func (p *Policy) sanitizeAttrs(
	elementName string,
	attrs []html.Attribute,
	aps map[string]attrPolicy,
) []html.Attribute {

	if len(attrs) == 0 {
		return attrs
	}

	cleanAttrs := []html.Attribute{}

	for _, htmlAttr := range attrs {
		// Is there an element specific attribute policy that applies?
		if ap, ok := aps[htmlAttr.Key]; ok {
			if ap.regexp != nil {
				if ap.regexp.MatchString(htmlAttr.Val) {
					cleanAttrs = append(cleanAttrs, htmlAttr)
					continue
				}
			} else {
				cleanAttrs = append(cleanAttrs, htmlAttr)
				continue
			}
		}

		// Is there a global attribute policy that applies?
		if ap, ok := p.globalAttrs[htmlAttr.Key]; ok {
			if ap.regexp != nil {
				if ap.regexp.MatchString(htmlAttr.Val) {
					cleanAttrs = append(cleanAttrs, htmlAttr)
				}
			} else {
				cleanAttrs = append(cleanAttrs, htmlAttr)
			}
		}
	}

	if linkable(elementName) && p.requireParseableURLs {
		// Ensure URLs are parseable:
		// - a.href
		// - area.href
		// - link.href
		// - blockquote.cite
		// - q.cite
		// - img.src
		// - script.src
		tmpAttrs := []html.Attribute{}
		for _, htmlAttr := range cleanAttrs {
			switch elementName {
			case "a", "area", "link":
				if htmlAttr.Key == "href" {
					if u, ok := p.validURL(htmlAttr.Val); ok {
						htmlAttr.Val = u
						tmpAttrs = append(tmpAttrs, htmlAttr)
					}
					break
				}
				tmpAttrs = append(tmpAttrs, htmlAttr)
			case "blockquote", "q":
				if htmlAttr.Key == "cite" {
					if u, ok := p.validURL(htmlAttr.Val); ok {
						htmlAttr.Val = u
						tmpAttrs = append(tmpAttrs, htmlAttr)
					}
					break
				}
				tmpAttrs = append(tmpAttrs, htmlAttr)
			case "img", "script":
				if htmlAttr.Key == "src" {
					if u, ok := p.validURL(htmlAttr.Val); ok {
						htmlAttr.Val = u
						tmpAttrs = append(tmpAttrs, htmlAttr)
					}
					break
				}
				tmpAttrs = append(tmpAttrs, htmlAttr)
			default:
				tmpAttrs = append(tmpAttrs, htmlAttr)
			}
		}
		cleanAttrs = tmpAttrs
	}

	if linkable(elementName) && p.requireNoFollow && len(cleanAttrs) > 0 {
		// Add rel="nofollow" if a "href" exists
		switch elementName {
		case "a", "area", "link":
			var hrefFound bool
			for _, htmlAttr := range cleanAttrs {
				if htmlAttr.Key == "href" {
					hrefFound = true
					continue
				}
			}

			if hrefFound {
				tmpAttrs := []html.Attribute{}
				var relFound bool
				var noFollowFound bool
				for _, htmlAttr := range cleanAttrs {
					if htmlAttr.Key == "rel" {
						relFound = true
						if strings.Contains(htmlAttr.Val, "nofollow") {
							noFollowFound = true
							continue
						}

						htmlAttr.Val += " nofollow"
						tmpAttrs = append(tmpAttrs, htmlAttr)
					} else {

						tmpAttrs = append(tmpAttrs, htmlAttr)
					}
				}
				if noFollowFound {
					break
				}
				if relFound {
					cleanAttrs = tmpAttrs
					break
				}

				rel := html.Attribute{}
				rel.Key = "rel"
				rel.Val = "nofollow"
				cleanAttrs = append(cleanAttrs, rel)
			}
		default:
		}
	}

	return cleanAttrs
}

func (p *Policy) allowNoAttrs(elementName string) bool {
	_, ok := p.elsWithoutAttrs[elementName]
	return ok
}

func (p *Policy) validURL(rawurl string) (string, bool) {
	if p.requireParseableURLs {
		// URLs do not contain whitespace
		if strings.Contains(rawurl, " ") ||
			strings.Contains(rawurl, "\t") ||
			strings.Contains(rawurl, "\n") {
			return "", false
		}

		u, err := url.Parse(rawurl)
		if err != nil {
			return "", false
		}

		if u.Scheme != "" {
			_, ok := p.urlSchemes[u.Scheme]
			if ok {
				return u.String(), true
			}

			return "", false
		}

		if p.allowRelativeURLs {
			if u.String() != "" {
				return u.String(), true
			}
		}

		return "", false
	}

	return rawurl, true
}

func linkable(elementName string) bool {
	switch elementName {
	case "a", "area", "blockquote", "img", "link", "script":
		return true
	default:
		return false
	}
}
