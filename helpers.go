package bluemonday

import (
	"regexp"
)

// A selection of regular expressions that can be used as .Matching() rules on
// HTML attributes.
var (
	Align     = regexp.MustCompile(`(?i)center|left|right|justify|char`)
	Valign    = regexp.MustCompile(`(?i)baseline|bottom|middle|top`)
	Direction = regexp.MustCompile(`(?i)auto|rtl|ltr`)
	Integer   = regexp.MustCompile(`[0-9]+`)

	// ISO8601 looks scary but isn't, it's taken from here:
	// http://www.pelagodesign.com/blog/2009/05/20/iso-8601-date-validation-that-doesnt-suck/
	// Minor changes have been made to remove PERL specific syntax that requires
	// regexp backtracking which are not supported in Go
	ISO8601 = regexp.MustCompile(
		`^([\+-]?\d{4})((-?)((0[1-9]|1[0-2])` +
			`([12]\d|0[1-9]|3[01])?|W([0-4]\d|5[0-2])` +
			`(-?[1-7])?|(00[1-9]|0[1-9]\d|[12]\d{2}|3([0-5]\d|6[1-6])))` +
			`([T\s]((([01]\d|2[0-3])((:?)[0-5]\d)?|24\:?00)([\.,]\d+[^:])?)?` +
			`([0-5]\d([\.,]\d+)?)?([zZ]|([\+-])` +
			`([01]\d|2[0-3]):?([0-5]\d)?)?)?)?$`,
	)
	ListType        = regexp.MustCompile(`(?i)circle|disc|square|a|A|i|I|1`)
	Name            = regexp.MustCompile(`[a-zA-Z0-9\-_\$]+`)
	NamesAndSpaces  = regexp.MustCompile(`[a-zA-Z0-9\-_\$]+`)
	Number          = regexp.MustCompile(`[+-]?(?:(?:[0-9]+(?:\.[0-9]*)?)|\.[0-9]+)`)
	NumberOrPercent = regexp.MustCompile(`[0-9]+%?`)
	Paragraph       = regexp.MustCompile(`(?:[\p{L}\p{N},'\.\s\-_\(\)]|&[0-9]{2};)*`)
)

// AllowStandardURLs is a convenience function that will enable rel="nofollow"
// on "a", "area" and "link" (if you have allowed those elements) and will
// ensure that the URL values are parseable and either relative or belong to the
// "mailto", "http", or "https" schemes
func (p *policy) AllowStandardURLs() {
	// URLs must be parseable by net/url.Parse()
	p.RequireParseableURLs(true)

	// !url.IsAbs() is permitted
	p.AllowRelativeURLs(true)

	// Most common URL schemes only
	p.AllowURLSchemes("mailto", "http", "https")

	// For all anchors we will add rel="nofollow" if it does not already exist
	// This applies to "a" "area" "link"
	p.RequireNoFollowOnLinks(true)
}

// AllowStandardAttributes will enable "id", "title" and the language specific
// attributes "dir" and "lang" on all elements that are whitelisted
func (p *policy) AllowStandardAttributes() {
	// "dir" "lang" are permitted as both language attributes affect charsets
	// and direction of text.
	p.AllowAttrs("dir").Matching(Direction).Globally()
	p.AllowAttrs(
		"lang",
	).Matching(regexp.MustCompile(`[a-zA-Z]{2,20}`)).Globally()

	// "id" is permitted. This is pretty much as some HTML elements require this
	// to work well ("dfn" is an example of a "id" being value)
	// This does create a risk that JavaScript and CSS within your web page
	// might identify the wrong elements. Ensure that you select things
	// accurately
	p.AllowAttrs("id").Matching(
		regexp.MustCompile(`[a-zA-Z0-9\:\-_\.]+`),
	).Globally()

	// "title" is permitted as it improves accessibility.
	// Regexp: \p{L} matches unicode letters, \p{N} matches unicode numbers
	p.AllowAttrs(
		"title",
	).Matching(
		regexp.MustCompile(`[\p{L}\p{N}\s\-_',:\[\]!\./\\\(\)&]*`),
	).Globally()
}

// AllowImages enables the img element and some popular attributes. It will also
// ensure that URL values are parseable
func (p *policy) AllowImages() {

	// "img" is permitted
	p.AllowAttrs("align").Matching(Align).OnElements("img")
	p.AllowAttrs("alt").Matching(Paragraph).OnElements("img")
	p.AllowAttrs("height", "width").Matching(NumberOrPercent).OnElements("img")

	// Standard URLs enabled, which disables data URI images as the data scheme
	// isn't included in the policy
	p.AllowStandardURLs()
	p.AllowAttrs("src").OnElements("img")
}

// AllowLists will enabled ordered and unordered lists, as well as definition
// lists
func (p *policy) AllowLists() {
	// "ol" "ul" are permitted
	p.AllowAttrs("type").Matching(ListType).OnElements("ol", "ul")

	// "li" is permitted
	p.AllowAttrs("type").Matching(ListType).OnElements("li")
	p.AllowAttrs("value").Matching(Integer).OnElements("li")

	// "dl" "dt" "dd" are permitted
	p.AllowElements("dl", "dt", "dd")
}

// AllowTables will enable a rich set of elements and attributes to describe
// HTML tables
func (p *policy) AllowTables() {

	// "table" is permitted
	p.AllowAttrs("height", "width").Matching(NumberOrPercent).OnElements("table")
	p.AllowAttrs("summary").Matching(Paragraph).OnElements("table")

	// "caption" is permitted
	p.AllowElements("caption")

	// "col" "colgroup" are permitted
	p.AllowAttrs("align").Matching(Align).OnElements("col", "colgroup")
	p.AllowAttrs("height", "width").Matching(
		NumberOrPercent,
	).OnElements("col", "colgroup")
	p.AllowAttrs("span").Matching(NumberOrPercent).OnElements("colgroup", "col")
	p.AllowAttrs("valign").Matching(Valign).OnElements("col", "colgroup")

	// "thead" "tr" are permitted
	p.AllowAttrs("align").Matching(Align).OnElements("thead", "tr")
	p.AllowAttrs("valign").Matching(Valign).OnElements("thead", "tr")

	// "td" "th" are permitted
	p.AllowAttrs("abbr").Matching(Paragraph).OnElements("td", "th")
	p.AllowAttrs("align").Matching(Align).OnElements("td", "th")
	p.AllowAttrs("colspan", "rowspan").Matching(Number).OnElements("td", "th")
	p.AllowAttrs("headers").Matching(Name).OnElements("td", "th")
	p.AllowAttrs("height", "width").Matching(
		NumberOrPercent,
	).OnElements("td", "th")
	p.AllowAttrs(
		"scope",
	).Matching(
		regexp.MustCompile(`(?i)(?:row|col)(?:group)?`),
	).OnElements("td", "th")
	p.AllowAttrs("valign").Matching(Valign).OnElements("td", "th")
	p.AllowAttrs("nowrap").Matching(
		regexp.MustCompile(`(?i)|nowrap`),
	).OnElements("td", "th")

	// "tbody" "tfoot"
	p.AllowAttrs("align").Matching(Align).OnElements("tbody", "tfoot")
	p.AllowAttrs("valign").Matching(Valign).OnElements("tbody", "tfoot")
}
