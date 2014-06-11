package bluemonday

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"code.google.com/p/go.net/html"
)

// Sanitize takes a string that contains HTML elements and applies the policy
// whitelist to return a string that has been sanitized by the policy.
func (p *policy) Sanitize(s string) (string, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return "", fmt.Errorf("Input is empty")
	}

	var cleanHTML bytes.Buffer
	tokenizer := html.NewTokenizer(strings.NewReader(s))

	depth := 0
	ignore := false
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
			depth++

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

			cleanHTML.WriteString(token.String())

		case html.EndTagToken:
			depth--

			aps, ok := p.elsAndAttrs[token.Data]
			if !ok {
				ignore = false
				break
			}

			attrs, err := sanitizeAttrs(aps, p.globalAttrs, token.Attr)
			if err != nil {
				return "", err
			}
			token.Attr = attrs

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
		if ap, ok := aps[htmlAttr.Key]; ok {
			if ap.regexp != nil {
				if ap.regexp.MatchString(htmlAttr.Val) {
					cleanAttrs = append(cleanAttrs, htmlAttr)
				}
			} else {
				cleanAttrs = append(cleanAttrs, htmlAttr)
			}
			continue
		}

		if ap, ok := aps[htmlAttr.Key]; ok {
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
