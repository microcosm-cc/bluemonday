package bluemonday

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"io"
	"net/url"
	"strings"
)

func sanitizeLink(u *url.URL, v string) string {
	var p *url.URL
	var err error
	if u == nil {
		p, err = url.Parse(v)
		if err != nil {
			return ""
		}
	} else {
		p, err = u.Parse(v)
		if err != nil {
			return ""
		}
	}
	if !acceptableUriSchemes[p.Scheme] {
		return ""
	}

	return p.String()
}

func sanitizeStyle(v string) string {
	return v
}

func sanitizeAttributes(u *url.URL, t *html.Token) {
	var attrs []html.Attribute
	var isLink = false
	for _, a := range t.Attr {
		if a.Key == "target" {
		} else if a.Key == "style" {
			a.Val = sanitizeStyle(a.Val)
			attrs = append(attrs, a)
		} else if acceptableAttributes[a.Key] {
			if a.Key == "href" || a.Key == "src" {
				a.Val = sanitizeLink(u, a.Val)
			}
			if a.Key == "href" {
				isLink = true
			}
			attrs = append(attrs, a)
		}
	}
	if isLink {
		attrs = append(attrs, html.Attribute{
			Key: "target",
			Val: "_blank",
		})
	}
	t.Attr = attrs
}

func Sanitize(s string, u *url.URL) (string, string) {
	r := bytes.NewReader([]byte(strings.TrimSpace(s)))
	z := html.NewTokenizer(r)
	buf := &bytes.Buffer{}
	strip := &bytes.Buffer{}
	skip := 0
	if u != nil {
		u.RawQuery = ""
		u.Fragment = ""
	}
	for {
		if z.Next() == html.ErrorToken {
			if err := z.Err(); err == io.EOF {
				break
			} else {
				return s, s
			}
		}

		t := z.Token()
		if t.Type == html.StartTagToken || t.Type == html.SelfClosingTagToken {
			if !acceptableElements[t.Data] {
				if unacceptableElementsWithEndTag[t.Data] && t.Type != html.SelfClosingTagToken {
					skip += 1
				}
			} else {
				sanitizeAttributes(u, &t)
				buf.WriteString(t.String())
			}
		} else if t.Type == html.EndTagToken {
			if !acceptableElements[t.Data] {
				if unacceptableElementsWithEndTag[t.Data] {
					skip -= 1
				}
			} else {
				buf.WriteString(t.String())
			}
		} else if skip == 0 {
			buf.WriteString(t.String())
			if t.Type == html.TextToken {
				strip.WriteString(t.String())
			}
		}
	}

	return buf.String(), strip.String()
}

// Based on list from MDN's HTML5 element list
// https://developer.mozilla.org/en-US/docs/Web/Guide/HTML/HTML5/HTML5_element_list
var acceptableElements = map[string]bool{
	// Root element
	// "html": true,

	// Document metadata
	// "head":  true,
	// "title": true,
	// "base":  true,
	// "link":  true,
	// "meta":  true,
	// "style": true,

	// Scripting
	"noscript": true,
	// "script":   true,

	// Sections
	// "body":    true,
	"section": true,
	"nav":     true,
	"article": true,
	"aside":   true,
	"h1":      true,
	"h2":      true,
	"h3":      true,
	"h4":      true,
	"h5":      true,
	"h6":      true,
	"header":  true,
	"footer":  true,
	"address": true,
	"main":    true,

	// Grouping content
	"p":          true,
	"hr":         true,
	"pre":        true,
	"blockquote": true,
	"ol":         true,
	"ul":         true,
	"li":         true,
	"dl":         true,
	"dt":         true,
	"dd":         true,
	"figure":     true,
	"figcaption": true,
	"div":        true,

	// Text-level semantics
	"a":      true,
	"em":     true,
	"strong": true,
	"small":  true,
	"s":      true,
	"cite":   true,
	"q":      true,
	"dfn":    true,
	"abbr":   true,
	"data":   true,
	"time":   true,
	"code":   true,
	"var":    true,
	"samp":   true,
	"kbd":    true,
	"sub":    true,
	"sup":    true,
	"i":      true,
	"b":      true,
	"u":      true,
	"mark":   true,
	"ruby":   true,
	"rt":     true,
	"rp":     true,
	"bdi":    true,
	"bdo":    true,
	"span":   true,
	"br":     true,
	"wbr":    true,

	// Edits
	"ins": true,
	"del": true,

	// Embedded content
	"img":    true,
	"iframe": true,
	"embed":  true,
	"object": true,
	"param":  true,
	"video":  true,
	"audio":  true,
	"source": true,
	"track":  true,
	"canvas": true,
	"map":    true,
	"area":   true,
	"svg":    true,
	"math":   true,

	// Tabular data
	"table":    true,
	"caption":  true,
	"colgroup": true,
	"col":      true,
	"tbody":    true,
	"thead":    true,
	"tfoot":    true,
	"tr":       true,
	"td":       true,
	"th":       true,

	// Forms
	"form":     true,
	"fieldset": true,
	"legend":   true,
	"label":    true,
	"input":    true,
	"button":   true,
	"select":   true,
	"datalist": true,
	"optgroup": true,
	"option":   true,
	"textarea": true,
	"keygen":   true,
	"output":   true,
	"progress": true,
	"meter":    true,

	// Interactive elements
	// "details":  true,
	// "summary":  true,
	// "menuitem": true,
	// "menu":     true,
}

var unacceptableElementsWithEndTag = map[string]bool{
	"script": true,
	"applet": true,
	"style":  true,
}

// Based on list from MDN's HTML attribute reference
// https://developer.mozilla.org/en-US/docs/Web/HTML/Attributes
var acceptableAttributes = map[string]bool{
	"accept":         true,
	"accept-charset": true,
	// "accesskey":       true,
	"action":       true,
	"align":        true,
	"alt":          true,
	"async":        true,
	"autocomplete": true,
	// "autofocus":       true,
	// "autoplay":        true,
	"bgcolor":         true,
	"border":          true,
	"buffered":        true,
	"challenge":       true,
	"charset":         true,
	"checked":         true,
	"cite":            true,
	"class":           true,
	"code":            true,
	"codebase":        true,
	"color":           true,
	"cols":            true,
	"colspan":         true,
	"content":         true,
	"contenteditable": true,
	"contextmenu":     true,
	"controls":        true,
	"coords":          true,
	"data":            true,
	"data-custom":     true,
	"datetime":        true,
	"default":         true,
	"defer":           true,
	"dir":             true,
	"dirname":         true,
	"disabled":        true,
	"download":        true,
	"draggable":       true,
	"dropzone":        true,
	"enctype":         true,
	"for":             true,
	"form":            true,
	"headers":         true,
	"height":          true,
	"hidden":          true,
	"high":            true,
	"href":            true,
	"hreflang":        true,
	"http-equiv":      true,
	"icon":            true,
	"id":              true,
	"ismap":           true,
	"itemprop":        true,
	"keytype":         true,
	"kind":            true,
	"label":           true,
	"lang":            true,
	"language":        true,
	"list":            true,
	"loop":            true,
	"low":             true,
	"manifest":        true,
	"max":             true,
	"maxlength":       true,
	"media":           true,
	"method":          true,
	"min":             true,
	"multiple":        true,
	"name":            true,
	"novalidate":      true,
	"open":            true,
	"optimum":         true,
	"pattern":         true,
	"ping":            true,
	"placeholder":     true,
	"poster":          true,
	// "preload":         true,
	"pubdate":    true,
	"radiogroup": true,
	"readonly":   true,
	"rel":        true,
	"required":   true,
	"reversed":   true,
	"rows":       true,
	"rowspan":    true,
	"sandbox":    true,
	"spellcheck": true,
	"scope":      true,
	// "scoped":          true,
	// "seamless":        true,
	"selected": true,
	"shape":    true,
	"size":     true,
	"sizes":    true,
	"span":     true,
	"src":      true,
	// "srcdoc":          true,
	"srclang": true,
	"start":   true,
	// "step":            true,
	"style":   true,
	"summary": true,
	// "tabindex":        true,
	// "target":          true,
	"title": true,
	"type":  true,
	// "usemap":          true,
	"value": true,
	"width": true,
	// "wrap":            true,

	// Older HTML attributes
	// http://www.w3.org/TR/html5-diff/#obsolete-attributes
	"alink":        true,
	"background":   true,
	"cellpadding":  true,
	"cellspacing":  true,
	"char":         true,
	"clear":        true,
	"compact":      true,
	"frameborder":  true,
	"frame":        true,
	"hspace":       true,
	"marginheight": true,
	"noshade":      true,
	"nowrap":       true,
	"rules":        true,
	"scrolling":    true,
	"valign":       true,
}

// Based on list from Wikipedia's URI scheme
// http://en.wikipedia.org/wiki/URI_scheme
var acceptableUriSchemes = map[string]bool{
	"aim":      true,
	"apt":      true,
	"bitcoin":  true,
	"callto":   true,
	"cvs":      true,
	"facetime": true,
	"feed":     true,
	"ftp":      true,
	"git":      true,
	"gopher":   true,
	"gtalk":    true,
	"http":     true,
	"https":    true,
	"imap":     true,
	"irc":      true,
	"itms":     true,
	"jabber":   true,
	"magnet":   true,
	"mailto":   true,
	"mms":      true,
	"msnim":    true,
	"news":     true,
	"nntp":     true,
	"rtmp":     true,
	"rtsp":     true,
	"sftp":     true,
	"skype":    true,
	"svn":      true,
	"ymsgr":    true,
}
