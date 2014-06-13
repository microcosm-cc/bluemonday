# bluemonday [![Build Status](https://travis-ci.org/microcosm-cc/bluemonday.svg?branch=master)](https://travis-ci.org/microcosm-cc/bluemonday) [![GoDoc](https://godoc.org/github.com/microcosm-cc/bluemonday?status.png)](https://godoc.org/github.com/microcosm-cc/bluemonday)

bluemonday is a HTML sanitizer implemented in Go.

Feed it user generated content (UTF-8 strings and HTML) and it will give back HTML that has been sanitised using a whitelist of approved HTML elements and attributes. It is fast and highly configurable.

The default `bluemonday.UGCPolicy().Sanitize()` turns this:

```html
Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE><A CLASS=XSS></A>World
```

Into the more harmless:

```html
Hello World
```

And it turns this:

```html
<a href="javascript:alert('XSS1')" onmouseover="alert('XSS2')">XSS<a>
```

Into this:

```html
XSS
```

Whilst still allowing this:
```html
<a href="http://www.google.com/">
  <img src="https://ssl.gstatic.com/accounts/ui/logo_2x.png"/>
</a>
```

To pass through mostly unaltered (it gained a rel="nofollow"):

```html
<a href="http://www.google.com/" rel="nofollow">
  <img src="https://ssl.gstatic.com/accounts/ui/logo_2x.png"/>
</a>
```

The primary purpose of bluemonday is to take potentially unsafe user generated content (from things like Markdown, HTML WYSIWYG tools, etc) and make it safe for you to put on your website.

It protects sites against [XSS](http://en.wikipedia.org/wiki/Cross-site_scripting) and other malicious content that a user interface may deliver. There are many [vectors for an XSS attack](https://www.owasp.org/index.php/XSS_Filter_Evasion_Cheat_Sheet) and the safest thing to do is to sanitize user input against a known safe list of HTML elements and attributes.

You should **always** run bluemonday **after** any other processing.

If you use [blackfriday](https://github.com/russross/blackfriday) or [Pandoc](http://johnmacfarlane.net/pandoc/) then bluemonday should be run after these steps. This ensures that no insecure HTML is introduced later in your process.

bluemonday is heavily inspired by both the [OWASP Java HTML Sanitizer](https://code.google.com/p/owasp-java-html-sanitizer/) and the [HTML Purifier](http://htmlpurifier.org/).

## Is it production ready yet?

*Maybe*

We are currently passing our tests (including AntiSamy tests). Please see [Issues](https://github.com/microcosm-cc/bluemonday/issues) for any current issues.

Our tests may not be complete, we invite pull requests.

## Usage

Install in your `${GOPATH}` using `go get -u github.com/microcosm-cc/bluemonday`

Then call it:
```go
package main

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"
)

func main() {
	p := bluemonday.UGCPolicy()
	html := p.Sanitize(
		`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	)

	// Should print:
	// <a href="http://www.google.com" rel="nofollow">Google</a>
	fmt.Println(html)
}
```

You can build your own policies:
```go
package main

import (
	"fmt"

	"github.com/microcosm-cc/bluemonday"
)

func main() {
	p := bluemonday.NewPolicy()

	// We only allow <p> and <a href="">
	p.AllowAttrs("href").OnElements("a")
	p.AllowElements("p")

	html := p.Sanitize(
		`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	)

	// Should print:
	// <a href="http://www.google.com">Google</a>
	fmt.Println(html)
}
```

We ship two default policies, one is `bluemonday.StrictPolicy()` and can be thought of as equivalent to stripping all HTML elements and their attributes as it has nothing on it's whitelist.

The other is `bluemonday.UGCPolicy()` and allows a broad selection of HTML elements and attributes that are safe for user generated content. Note that this policy does *not* whitelist iframes, object, embed, styles, script, etc.

## Policy Building

The essence of building a policy is to determine which HTML elements and attributes are considered safe for your scenario. OWASP provide an [XSS prevention cheat sheet](https://www.owasp.org/index.php/XSS_(Cross_Site_Scripting)_Prevention_Cheat_Sheet) to help explain the risks, but essentially:

1. Avoid anything other than plain HTML elements
1. Avoid `script`, `style`, `iframe`, `object`, `embed`, `base` elements
1. Avoid anything other than plain HTML elements with simple values that you can match to a regexp

To create a new policy:

```go
p := bluemonday.NewPolicy()
```

To add elements to a policy either add just the elements:

```go
p.AllowElements("b", "strong")
```

Or add elements as a virtue of adding an attribute:

```go
p.AllowAttrs("nowrap").OnElements("td", "th")
```

Attributes can either be added to all elements:

```go
p.AllowAttrs("dir").Matching(regexp.MustCompile("(?i)rtl|ltr")).Globally()
```

Or attributes can be added to specific elements:

```go
p.AllowAttrs("value").OnElements("li")
```

It is **always** recommended that an attribute be made to match a pattern. XSS in HTML attributes is very easy otherwise:

```go
// \p{L} matches unicode letters, \p{N} matches unicode numbers
p.AllowAttrs("title").Matching(regexp.MustCompile(`[\p{L}\p{N}\s\-_',:\[\]!\./\\\(\)&]*`)).Globally()
```

You can stop at any time and call .Sanitize():

```go
// string htmlIn passed in from a HTTP POST
htmlOut := p.Sanitize(htmlIn)
```

And you can take any existing policy and extend it:

```go
p := bluemonday.UGCPolicy()
p.AllowElements("fieldset", "select", "option")
```

### Links

Links are complex beasts and one of the biggest attack vectors for malicious content.

It is possible to do this:

```go
p.AllowAttrs("href").Matching(regexp.MustCompile(`(?i)mailto|https?`)).OnElements("a")
```

But that may not help you as the regexp is insufficient in this case to have prevented a malformed value doing something unexpected.

We provide some additional global options for working with links.

This will ensure that URLs are not considered invalid by Go's `net/url` package.

```go
p.RequireParseableURLs(true)
```

If you have enabled parseable URLs then the following option will allow relative URLs. By default this is disabled and will prevent all local and schema relative URLs (i.e. `href="//www.google.com"` is schema relative).

```go
p.AllowRelativeURLs(true)
```

If you have enabled parseable URLs then you can whitelist the schemas that are permitted. Bear in mind that allowing relative URLs in the above option allows for blank schemas.

```go
p.AllowURLSchemes("mailto", "http", "https")
```

Regardless of whether you have enabled parseable URLs, you can force all URLs to have a rel="nofollow" attribute. This will be added if it does not exist.

```go
// This applies to "a" "area" "link" elements that have a "href" attribute
p.RequireNoFollowOnLinks(true)
```

We provide a convenience method that applies all of the above, but you will still need to whitelist the linkable elements:

```go
p.AllowStandardURLs()
p.AllowAttrs("cite").OnElements("blockquote")
p.AllowAttrs("href").OnElements("a", "area")
p.AllowAttrs("src").OnElements("img")
```

### Policy Building Helpers

If you've got this far and you're bored already, we also bundle some helpers:

```go
p.AllowStandardAttributes()
p.AllowImages()
p.AllowLists()
p.AllowTables()
```

### Invalid Instructions

The following are invalid:

```go
// This does not say where the attributes are allowed, you need to add
// .Globally() or .OnElements(...)
// This will be ignored without error.
p.AllowAttrs("value")

// This does not say where the attributes are allowed, you need to add
// .Globally() or .OnElements(...)
// This will be ignored without error.
p.AllowAttrs(
	"type",
).Matching(
	regexp.MustCompile("(?i)circle|disc|square|a|A|i|I|1"),
)
```

Both examples exhibit the same issues, they declared attributes but didn't then specify whether they are whitelisted globally or only on specific elements (and which elements).

## Limitations

In this early release we are focusing on sanitizing HTML elements and attributes only. We are not yet including any tools to help whitelist and sanitize CSS. Which means that unless you wish to do the heavy lifting in a single regular expression (inadvisable), **you should not allow the "style" attribute anywhere**.

It is not the job of bluemonday to fix your bad HTML, it is merely the job of bluemonday to prevent malicious HTML getting through. If you have mismatched HTML elements, or non-conforming nesting of elements, those will remain. But if you have well-structured HTML bluemonday will not break it.

## TODO

1. Add support for allowing the list of HTML elements that are permitted to be empty to be configured
1. Add support for CSS sanitisation to allow some CSS properties based on a whitelist
1. Investigate whether devs want to blacklist elements and attributes. This would allow devs to take an existing policy (such as the `bluemonday.UGCPolicy()` ) that are 90% of what they're looking for and to remove the few things they don't want to make it 100% what they want

## Development

If you have cloned this repo you will probably need the dependency:

`go get code.google.com/p/go.net/html`

Gophers can use their familiar tools:

`go build`

`go test`

I personally use a Makefile as it spares typing the same args over and over.

`make` will build (for 64-bit linux), test and install the library.

`make clean` will remove the library from a *single* `${GOPATH}/pkg` directory tree

`make test` will run the tests with a coverage report

`make cover` will run the tests and *open a browser window* with the coverage report

`make lint` will run golint (install via `go get github.com/golang/lint/golint`)

## Long term goals

1. Open the code to adversarial peer review similar to the [Attack Review Ground Rules](https://code.google.com/p/owasp-java-html-sanitizer/wiki/AttackReviewGroundRules)
1. Raise funds and pay for an external security review
