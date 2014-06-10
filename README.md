Bluemonday
==========

Bluemonday is a HTML sanitizer implemented in Go.

You should be able to safely feed it user generated content (UTF-8 strings and HTML) and it will give back HTML that has been sanitised using a whitelist of approved HTML elements and attributes. It is fast enough and highly configurable.

The primary purpose of bluemonday is to protect sites against [XSS](http://en.wikipedia.org/wiki/Cross-site_scripting) whenever user generated content is used.

You should run bluemonday **after** any other processing. So if you use [blackfriday](https://github.com/russross/blackfriday) first then bluemonday should always be the last bit of processing you do. This ensures that no insecurities are introduced later in your process.

Bluemonday is heavily inspired by both the [OWASP Java HTML Sanitizer](https://code.google.com/p/owasp-java-html-sanitizer/) and the [HTML Purifier](http://htmlpurifier.org/).

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
	p := bluemonday.MarkdownUGCPolicy()
	html := p.Sanitize(
		`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	)

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

	html := p.Sanitize(
		`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
	)

	// Should print:
	// <a href="http://www.google.com">Google</a>
	fmt.Println(html)
}
````

We ship two default policies, one is `bluemonday.StrictPolicy()` and can be thought of as equivalent to stripping all HTML elements and their attributes as it has nothing on it's whitelist.

The other is `bluemonday.MarkdownUGCPolicy()` and allows a broad selection of HTML elements and attributes that are safe for user generated content. Note that this policy does *not* whitelist iframes, object, embed, styles, script, etc.

TODO
====

The TODO list is extensive and bluemonday should **NOT** be used in production code at this point, in fact given that the API will change dramatically you probably shouldn't use it all just now.

1. ~~Add the ability to describe policies for the sanitization process~~
1. ~~Bundle a set of default policies that represent safe defaults~~
1. Add the ability to sanitize based on a policy
1. Implement and pass the equivalent of the [OWASP Java HTML Sanitizer tests](https://code.google.com/p/owasp-java-html-sanitizer/source/browse/trunk/src/tests/org/owasp/html/)
1. Open the code to adversarial peer review similar to the [Attack Review Ground Rules](https://code.google.com/p/owasp-java-html-sanitizer/wiki/AttackReviewGroundRules)
1. Raise funds and pay for an external security review
