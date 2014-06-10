package bluemonday

import (
	"regexp"
)

// StrictPolicy returns an empty policy, which will effectively strip all HTML
// elements and their attributes from a document.
func StrictPolicy() *policy {
	return NewPolicy()
}

// MarkdownUGCPolicy returns a policy aimed at user generated content generated
// by Markdown. This is expected to be a fairly rich document where as much
// markup as possible should be retained. Markdown permits raw HTML so we are
// basically providing a policy to sanitise HTML5 documents safely but with the
// least intrusion on the expectations of the user.
func MarkdownUGCPolicy() *policy {
	p := NewPolicy()

	// This looks scary but isn't, it's taken from here:
	// http://www.pelagodesign.com/blog/2009/05/20/iso-8601-date-validation-that-doesnt-suck/
	// Minor changes have been made to remove PERL specific syntax that requires
	// regexp backtracking which are not supported in Go
	iso8601 := regexp.MustCompile(
		`^([\+-]?\d{4})((-?)((0[1-9]|1[0-2])` +
			`([12]\d|0[1-9]|3[01])?|W([0-4]\d|5[0-2])` +
			`(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))` +
			`([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24\:?00)([\.,]\d+[^:])?)?` +
			`([0-5]\d([\.,]\d+)?)?([zZ]|([\+-])` +
			`([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$`,
	)
	name := regexp.MustCompile(`[a-zA-Z0-9\-_\$]+`)
	number := regexp.MustCompile(`[+-]?(?:(?:[0-9]+(?:\.[0-9]*)?)|\.[0-9]+)`)
	numOrPercent := regexp.MustCompile(`[0-9]+%?`)
	paragraph := regexp.MustCompile(`(?:[\p{L}\p{N},'\.\s\-_\(\)]|&[0-9]{2};)*`)

	///////////////////////
	// Global attributes //
	///////////////////////
	p.AllowAttrs(
		"id",
	).Matching(regexp.MustCompile(`[a-zA-Z0-9\:\-_\.]+`)).Globally()

	p.AllowAttrs(
		"class",
	).Matching(regexp.MustCompile(`[a-zA-Z0-9\s,\-_]+`)).Globally()

	p.AllowAttrs(
		"dir",
	).Matching(regexp.MustCompile(`(?i)rtl|ltr`)).Globally()

	p.AllowAttrs(
		"lang",
	).Matching(regexp.MustCompile(`[a-zA-Z]{2,20}`)).Globally()

	// Seldom seen regexp: \p{L} matches unicode letters, \p{N} matches unicode
	// numbers
	p.AllowAttrs(
		"title",
	).Matching(
		regexp.MustCompile(`[\p{L}\p{N}\s\-_',:\[\]!\./\\\(\)&]*`),
	).Globally()

	/////////////////////////////////////
	// Attributes on specific elements //
	/////////////////////////////////////
	p.AllowAttrs("abbr").Matching(paragraph).OnElements("td", "th")

	p.AllowAttrs(
		"align",
	).Matching(
		regexp.MustCompile(`(?i)center|left|right|justify|char`),
	).OnElements(
		"col",
		"colgroup",
		"img",
		"tbody",
		"td",
		"tfoot",
		"th",
		"thead",
		"tr",
	)

	p.AllowAttrs("alt").Matching(paragraph).OnElements("area", "img")

	p.AllowAttrs(
		"cite",
	).Matching(paragraph).OnElements("blockquote", "del", "ins", "q")

	p.AllowAttrs("colspan", "rowspan").Matching(number).OnElements("td", "th")

	p.AllowAttrs(
		"datetime",
	).Matching(iso8601).OnElements("del", "ins", "time")

	p.AllowAttrs("headers").Matching(name).OnElements("td", "th")

	p.AllowAttrs(
		"height",
		"width",
	).Matching(numOrPercent).OnElements(
		"col",
		"colgroup",
		"img",
		"table",
		"td",
		"th",
	)

	p.AllowAttrs("href").OnElements("a").RequireNoFollowOnLinks()

	p.AllowAttrs("nowrap").OnElements("td", "th")

	p.AllowAttrs(
		"open",
	).Matching(regexp.MustCompile(`(?i)|open`)).OnElements("details")

	p.AllowAttrs(
		"scope",
	).Matching(
		regexp.MustCompile(`(?i)(?:row|col)(?:group)?`),
	).OnElements("td", "th")

	p.AllowAttrs("span").Matching(numOrPercent).OnElements("colgroup", "col")

	p.AllowAttrs("src").OnElements("img")

	p.AllowAttrs("summary").Matching(paragraph).OnElements("table")

	p.AllowAttrs(
		"type",
	).Matching(
		regexp.MustCompile(`(?i)circle|disc|square|a|A|i|I|1`),
	).OnElements("li", "ol", "ul")

	p.AllowAttrs(
		"valign",
	).Matching(
		regexp.MustCompile(`(?i)baseline|bottom|middle|top`),
	).OnElements("col", "colgroup", "tbody", "td", "tfoot", "th", "thead", "tr")

	p.AllowAttrs("value").OnElements("li")

	p.AllowAttrs(
		"value",
		"min",
		"max",
		"low",
		"high",
		"optimum",
	).Matching(number).OnElements("meter")

	p.AllowAttrs("value", "max").Matching(number).OnElements("progress")

	p.AllowElements(
		"a", "abbr", "acronym", "article", "aside", "b", "bdi", "bdo",
		"blockquote", "br", "caption", "cite", "code", "col", "colgroup", "del",
		"details", "dfn", "div", "dl", "em", "figcaption", "figure", "footer",
		"h1", "h2", "h3", "h4", "h5", "h6", "hr", "i", "img", "ins", "li",
		"mark", "meter", "ol", "p", "pre", "progress", "q", "s", "samp",
		"section", "small", "span", "strike", "strong", "sub", "summary", "sup",
		"table", "tbody", "td", "tfoot", "th", "thead", "time", "tr", "tt", "u",
		"ul", "var",
	)

	return p
}
