bluemonday
==========

bluemonday is a HTML sanitizer implemented in Go.

You should be able to safely feed it user generated content (UTF-8 strings and HTML) and it will give back HTML that has been sanitised using a whitelist of approved HTML elements and attributes. It is fast and highly configurable.

The primary purpose of bluemonday is to protect sites against [XSS](http://en.wikipedia.org/wiki/Cross-site_scripting) and other malicious content that a user interface may deliver.

Note: It is not the job of bluemonday to fix your bad HTML, it is merely the job of bluemonday to prevent malicious HTML getting through. If you have mismatched HTML elements, or non-conforming nesting of elements, those will remain. But if you have well-structured HTML bluemonday will not break it.

You should **always** run bluemonday **after** any other processing. So if you use [blackfriday](https://github.com/russross/blackfriday) or [Pandoc](http://johnmacfarlane.net/pandoc/) then bluemonday should be run after these steps. This ensures that no insecure HTML is introduced later in your process.

Bluemonday is heavily inspired by both the [OWASP Java HTML Sanitizer](https://code.google.com/p/owasp-java-html-sanitizer/) and the [HTML Purifier](http://htmlpurifier.org/).

Is it production ready yet?
===========================

*Maybe*

We are currently passing our tests (including AntiSamy tests). Please see [Issues](https://github.com/microcosm-cc/bluemonday/issues) for any current issues.

Usage
=====

Install in your $(GOPATH) using `go get -u github.com/microcosm-cc/bluemonday`

Then call it:
````
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
````

You can build your own policies:
````
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
````

We ship two default policies, one is `bluemonday.StrictPolicy()` and can be thought of as equivalent to stripping all HTML elements and their attributes as it has nothing on it's whitelist.

The other is `bluemonday.UGCPolicy()` and allows a broad selection of HTML elements and attributes that are safe for user generated content. Note that this policy does *not* whitelist iframes, object, embed, styles, script, etc.

Limitations
===========

In this early release we are focusing on sanitizing HTML elements and attributes only. We are not yet including any tools to help whitelist and sanitize CSS. Which means that unless you wish to do the heavy lifting in a single regular expression, **you should probably not allow the "style" attribute anywhere**.

TODO
====

1. Support p.RequireNoFollowOnLinks() as we're ignoring it right now
1. Add support for parsing of URLs and URL protocols more intelligently than forcing the developer to write a regexp
1. Add support for allowing the list of HTML elements that are permitted to be empty to be configured

Development
===========

We use a Makefile as there's nowt wrong with `make`.

`make` will build, test and install the library.

`make clean` will remove the library from the `${GOPATH}/pkg` directory tree

`make test` will run the tests with a coverage report

`make cover` will run the tests and open a browser window with the coverage report

`make lint` will run golint (install via `go get github.com/golang/lint/golint`)

Long term goals
===============

1. Open the code to adversarial peer review similar to the [Attack Review Ground Rules](https://code.google.com/p/owasp-java-html-sanitizer/wiki/AttackReviewGroundRules)
1. Raise funds and pay for an external security review
