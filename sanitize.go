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
// It returns a HTML string that has been sanitized by the policy or an error if
// one occurred (most likely as a consequence of malformed input)
func (p *policy) Sanitize(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", ErrEmptyInput
	}

	var cleanHTML bytes.Buffer
	tokenizer := html.NewTokenizer(strings.NewReader(s))

	ignore := false
	skipClosingTag := false
	for {
		if tokenizer.Next() == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				// End of input means end of processing
				return cleanHTML.String(), nil
			}
			return "", err
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
			return "", ErrNotImplemented
		}
	}
}

// sanitizeAttrs takes a set of element attribute policies and the global
// attribute policies and applies them to the []html.Attribute returning a set
// of html.Attributes that match the policies
func (p *policy) sanitizeAttrs(
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
			case "blockquote":
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

func (p *policy) allowNoAttrs(elementName string) bool {
	_, ok := p.elsWithoutAttrs[elementName]
	return ok
}

func (p *policy) validURL(rawurl string) (string, bool) {
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
