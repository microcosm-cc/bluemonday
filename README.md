bluemonday [![Build Status](https://travis-ci.org/microcosm-cc/bluemonday.svg?branch=master)](https://travis-ci.org/microcosm-cc/bluemonday)
==========

bluemonday is a HTML sanitizer implemented in Go.

Feed it user generated content (UTF-8 strings and HTML) and it will give back HTML that has been sanitised using a whitelist of approved HTML elements and attributes. It is fast and highly configurable.

The default `bluemonday.UGCPolicy().Sanitize()` Turns this:

```html
Hello <STYLE>.XSS{background-image:url("javascript:alert('XSS')");}</STYLE><A CLASS=XSS></A>World
```

Into the more harmless:

```html
Hello World
```

The primary purpose of bluemonday is to protect sites against [XSS](http://en.wikipedia.org/wiki/Cross-site_scripting) and other malicious content that a user interface may deliver. There are many [vectors for an XSS attack](https://www.owasp.org/index.php/XSS_Filter_Evasion_Cheat_Sheet) and the safest thing for someone accepting user generated content is to sanitize user input against a safe list of HTML elements and attributes.

You should **always** run bluemonday **after** any other processing. So if you use [blackfriday](https://github.com/russross/blackfriday) or [Pandoc](http://johnmacfarlane.net/pandoc/) then bluemonday should be run after these steps. This ensures that no insecure HTML is introduced later in your process.

Bluemonday is heavily inspired by both the [OWASP Java HTML Sanitizer](https://code.google.com/p/owasp-java-html-sanitizer/) and the [HTML Purifier](http://htmlpurifier.org/).

Is it production ready yet?
===========================

*Maybe*

We are currently passing our tests (including AntiSamy tests). Please see [Issues](https://github.com/microcosm-cc/bluemonday/issues) for any current issues.

Our tests may not be complete, we invite pull requests.

Usage
=====

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
	html, err := p.Sanitize(
		`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return
	}

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

	html, err := p.Sanitize(
		`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		return
	}

	// Should print:
	// <a href="http://www.google.com">Google</a>
	fmt.Println(html)
}
```

We ship two default policies, one is `bluemonday.StrictPolicy()` and can be thought of as equivalent to stripping all HTML elements and their attributes as it has nothing on it's whitelist.

The other is `bluemonday.UGCPolicy()` and allows a broad selection of HTML elements and attributes that are safe for user generated content. Note that this policy does *not* whitelist iframes, object, embed, styles, script, etc.

Policy Building
===============

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
htmlOut, err := p.Sanitize(htmlIn)
```

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

Limitations
===========

In this early release we are focusing on sanitizing HTML elements and attributes only. We are not yet including any tools to help whitelist and sanitize CSS. Which means that unless you wish to do the heavy lifting in a single regular expression (inadvisable), **you should not allow the "style" attribute anywhere**.

It is not the job of bluemonday to fix your bad HTML, it is merely the job of bluemonday to prevent malicious HTML getting through. If you have mismatched HTML elements, or non-conforming nesting of elements, those will remain. But if you have well-structured HTML bluemonday will not break it.


TODO
====

1. Support p.RequireNoFollowOnLinks() as we're ignoring it right now
1. Add support for parsing of URLs and URL protocols more intelligently than forcing the developer to write a regexp
1. Add support for allowing the list of HTML elements that are permitted to be empty to be configured
1. Add support for CSS sanitisation to allow some CSS properties based on a whitelist

Development
===========

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

Long term goals
===============

1. Open the code to adversarial peer review similar to the [Attack Review Ground Rules](https://code.google.com/p/owasp-java-html-sanitizer/wiki/AttackReviewGroundRules)
1. Raise funds and pay for an external security review
