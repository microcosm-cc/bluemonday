package bluemonday

import (
	"regexp"
)

// StrictPolicy returns an empty policy, which will effectively strip all HTML
// elements and their attributes from a document.
func StrictPolicy() *policy {
	return NewPolicy()
}

// UGCPolicy returns a policy aimed at user generated content that is a result
// of HTML WYSIWYG tools and Markdown conversions.
//
// This is expected to be a fairly rich document where as much markup as
// possible should be retained. Markdown permits raw HTML so we are basically
// providing a policy to sanitise HTML5 documents safely but with the
// least intrusion on the formatting expectations of the user.
func UGCPolicy() *policy {
	p := NewPolicy()

	align := regexp.MustCompile(`(?i)center|left|right|justify|char`)
	valign := regexp.MustCompile(`(?i)baseline|bottom|middle|top`)
	direction := regexp.MustCompile(`(?i)auto|rtl|ltr`)
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
	listType := regexp.MustCompile(`(?i)circle|disc|square|a|A|i|I|1`)
	name := regexp.MustCompile(`[a-zA-Z0-9\-_\$]+`)
	namesAndSpaces := regexp.MustCompile(`[a-zA-Z0-9\-_\$]+`)
	number := regexp.MustCompile(`[+-]?(?:(?:[0-9]+(?:\.[0-9]*)?)|\.[0-9]+)`)
	numOrPercent := regexp.MustCompile(`[0-9]+%?`)
	paragraph := regexp.MustCompile(`(?:[\p{L}\p{N},'\.\s\-_\(\)]|&[0-9]{2};)*`)
	standardURLs := regexp.MustCompile(`(?i)^https?|mailto`)

	///////////////////////
	// Global attributes //
	///////////////////////

	// "class" is not permitted as we are not allowing users to style their own
	// content

	// "dir" "lang" are permitted as both language attributes affect charsets
	// and direction of text.
	p.AllowAttrs("dir").Matching(direction).Globally()
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

	////////////////////////////////
	// Declarations and structure //
	////////////////////////////////

	// "xml" "xslt" "DOCTYPE" "html" "head" are not permitted as we are
	// expecting user generated content to be a fragment of HTML and not a full
	// document.

	//////////////////////////
	// Sectioning root tags //
	//////////////////////////

	// "article" and "aside" are permitted and takes no attributes
	p.AllowElements("article", "aside")

	// "body" is not permitted as we are expecting user generated content to be a fragment
	// of HTML and not a full document.

	// "details" is permitted, including the "open" attribute which can either
	// be blank or the value "open".
	p.AllowAttrs(
		"open",
	).Matching(regexp.MustCompile(`(?i)|open`)).OnElements("details")

	// "fieldset" is not permitted as we are not allowing forms to be created.

	// "figure" is permitted and takes no attributes
	p.AllowElements("figure")

	// "nav" is not permitted as it is assumed that the site (and not the user)
	// has defined navigation elements

	// "section" is permitted and takes no attributes
	p.AllowElements("section")

	// "summary" is permitted and takes no attributes
	p.AllowElements("summary")

	//////////////////////////
	// Headings and footers //
	//////////////////////////

	// "footer" is not permitted as we expect user content to be a fragment and
	// not structural to this extent

	// "h1" through "h6" are permitted and take no attributes
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")

	// "header" is not permitted as we expect user content to be a fragment and
	// not structural to this extent

	// "hgroup" is permitted and takes no attributes
	p.AllowElements("hgroup")

	/////////////////////////////////////
	// Content grouping and separating //
	/////////////////////////////////////

	// "blockquote" is permitted, including the "cite" attribute which must be
	// a standard URL.
	p.AllowAttrs(
		"cite",
	).Matching(standardURLs).OnElements("blockquote")

	// "br" "div" "hr" "p" "span" "wbr" are permitted and takes no attributes
	p.AllowElements("br", "div", "hr", "p", "span", "wbr")

	///////////
	// Links //
	///////////

	// For all anchors we will add rel="nofollow" if it does not already exist
	// This applies to "a" "area" "link"
	p.RequireNoFollowOnLinks(true)

	// "a" is permitted
	p.AllowAttrs("href").Matching(standardURLs).OnElements("a")

	// "area" is permitted along with the attributes that map image maps work
	p.AllowAttrs("alt").Matching(paragraph).OnElements("area")
	p.AllowAttrs("coords").Matching(
		regexp.MustCompile(`([0-9]+,){2}(,[0-9]+)*`),
	).OnElements("area")
	p.AllowAttrs("href").Matching(standardURLs).OnElements("area")
	p.AllowAttrs("rel").Matching(namesAndSpaces).OnElements("area")
	p.AllowAttrs("shape").Matching(
		regexp.MustCompile(`(?i)|default|circle|rect|poly`),
	).OnElements("area")

	// "link" is not permitted

	/////////////////////
	// Phrase elements //
	/////////////////////

	// The following are all inline phrasing elements
	p.AllowElements("abbr", "acronym", "cite", "code", "dfn", "em",
		"figcaption", "mark", "s", "samp", "strong", "sub", "sup", "var")

	// "q" is permitted
	p.AllowAttrs("cite").Matching(paragraph).OnElements("q")

	// "time" is permitted
	p.AllowAttrs("datetime").Matching(iso8601).OnElements("time")

	////////////////////
	// Style elements //
	////////////////////

	// block and inline elements that impart no semantic meaning but style the
	// document
	p.AllowElements("b", "i", "pre", "small", "strike", "tt", "u")

	// "style" is not permitted as we are not yet sanitising CSS and it is an
	// XSS attack vector

	//////////////////////
	// HTML5 Formatting //
	//////////////////////

	// "bdi" "bdo" are permitted
	p.AllowAttrs("dir").Matching(direction).OnElements("bdi", "bdo")

	// "rp" "rt" "ruby" are permitted
	p.AllowElements("rp", "rt", "ruby")

	///////////////////////////
	// HTML5 Change tracking //
	///////////////////////////

	// "del" "ins" are permitted
	p.AllowAttrs("cite").Matching(paragraph).OnElements("del", "ins")
	p.AllowAttrs("datetime").Matching(iso8601).OnElements("del", "ins")

	///////////
	// Lists //
	///////////

	// "ol" "ul" are permitted
	p.AllowAttrs("type").Matching(listType).OnElements("ol", "ul")

	// "li" is permitted
	p.AllowAttrs("type").Matching(listType).OnElements("li")
	p.AllowAttrs("value").Matching(regexp.MustCompile(`[0-9]+`)).OnElements("li")

	// "dl" "dt" "dd" are permitted
	p.AllowElements("dl", "dt", "dd")

	////////////
	// Tables //
	////////////

	// "table" is permitted
	p.AllowAttrs("height", "width").Matching(numOrPercent).OnElements("table")
	p.AllowAttrs("summary").Matching(paragraph).OnElements("table")

	// "caption" is permitted
	p.AllowElements("caption")

	// "col" "colgroup" are permitted
	p.AllowAttrs("align").Matching(align).OnElements("col", "colgroup")
	p.AllowAttrs("height", "width").Matching(
		numOrPercent,
	).OnElements("col", "colgroup")
	p.AllowAttrs("span").Matching(numOrPercent).OnElements("colgroup", "col")
	p.AllowAttrs("valign").Matching(valign).OnElements("col", "colgroup")

	// "thead" "tr" are permitted
	p.AllowAttrs("align").Matching(align).OnElements("thead", "tr")
	p.AllowAttrs("valign").Matching(valign).OnElements("thead", "tr")

	// "td" "th" are permitted
	p.AllowAttrs("abbr").Matching(paragraph).OnElements("td", "th")
	p.AllowAttrs("align").Matching(align).OnElements("td", "th")
	p.AllowAttrs("colspan", "rowspan").Matching(number).OnElements("td", "th")
	p.AllowAttrs("headers").Matching(name).OnElements("td", "th")
	p.AllowAttrs("height", "width").Matching(
		numOrPercent,
	).OnElements("td", "th")
	p.AllowAttrs(
		"scope",
	).Matching(
		regexp.MustCompile(`(?i)(?:row|col)(?:group)?`),
	).OnElements("td", "th")
	p.AllowAttrs("valign").Matching(valign).OnElements("td", "th")
	p.AllowAttrs("nowrap").OnElements("td", "th")

	// "tbody" "tfoot"
	p.AllowAttrs("align").Matching(align).OnElements("tbody", "tfoot")
	p.AllowAttrs("valign").Matching(valign).OnElements("tbody", "tfoot")

	///////////
	// Forms //
	///////////

	// By and large, forms are not permitted. However there are some form
	// elements that can be used to present data, and we do permit those
	//
	// "button" "fieldset" "input" "keygen" "label" "output" "select" "datalist"
	// "textarea" "optgroup" "option" are all not permitted

	// "meter" is permitted
	p.AllowAttrs(
		"value",
		"min",
		"max",
		"low",
		"high",
		"optimum",
	).Matching(number).OnElements("meter")

	// "progress" is permitted
	p.AllowAttrs("value", "max").Matching(number).OnElements("progress")

	//////////////////////
	// Embedded content //
	//////////////////////

	// Vast majority not permitted
	// "audio" "canvas" "embed" "iframe" "object" "param" "source" "svg" "track"
	// "video" are all not permitted

	// "img" is permitted
	p.AllowAttrs("align").Matching(align).OnElements("img")
	p.AllowAttrs("alt").Matching(paragraph).OnElements("img")
	p.AllowAttrs("height", "width").Matching(numOrPercent).OnElements("img")
	p.AllowAttrs("src").Matching(standardURLs).OnElements("img")

	return p
}
