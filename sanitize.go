package bluemonday

import (
	"bytes"
	"fmt"
	"io"
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
		return "", fmt.Errorf("Input is empty")
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
				cleanHTML.WriteString(token.Data)
			}

		case html.CommentToken:
			// Comments are ignored by default, effectively removing them from
			// the output document.

		case html.StartTagToken:

			aps, ok := p.elsAndAttrs[token.Data]
			if !ok {
				ignore = true
				break
			}

			attrs, err := sanitizeAttrs(aps, p.globalAttrs, token.Attr)
			if err != nil {
				return "", err
			}
			token.Attr = attrs

			// Do we have any attributes?
			if len(token.Attr) == 0 {
				// Some elements make no sense without attributes, so we skip
				// those, but anything in this switch is basically permitted
				// to have no attributes.
				switch token.Data {
				case "b":
				case "div":
				case "li":
				case "ol":
				case "p":
				case "span":
				case "table":
				case "tbody":
				case "td":
				case "th":
				case "thead":
				case "tr":
				case "ul":
				default:
					skipClosingTag = true
				}
			}

			// If we're skipping the closing tag, we should skip the opening
			// one too
			if !skipClosingTag {
				cleanHTML.WriteString(token.String())
			}

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

			attrs, err := sanitizeAttrs(aps, p.globalAttrs, token.Attr)
			if err != nil {
				return "", err
			}
			token.Attr = attrs

			cleanHTML.WriteString(token.String())

		case html.TextToken:

			if !ignore {
				cleanHTML.WriteString(token.String())
			}

		default:
			return "", fmt.Errorf("Not implemented")
		}
	}
}

// sanitizeAttrs takes a set of element attribute policies and the global
// attribute policies and applies them to the []html.Attribute returning a set
// of html.Attributes that match the policies
func sanitizeAttrs(
	aps map[string]attrPolicy,
	gap map[string]attrPolicy,
	attrs []html.Attribute,
) (
	[]html.Attribute,
	error,
) {
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
		if ap, ok := gap[htmlAttr.Key]; ok {
			if ap.regexp != nil {
				if ap.regexp.MatchString(htmlAttr.Val) {
					cleanAttrs = append(cleanAttrs, htmlAttr)
				}
			} else {
				cleanAttrs = append(cleanAttrs, htmlAttr)
			}
		}
	}

	return cleanAttrs, nil
}
