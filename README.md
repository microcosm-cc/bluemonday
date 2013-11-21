Bluemonday
==========

Bluemonday is a HTML sanitizer implemented in Go. You should be able to safely feed it user generated content and it will give back HTML that is scrubbed using a whitelist of approved tags and attributes. It is fast, configurable and is safe for all utf-8 (unicode) input.

The primary purpose of bluemonday is to protect sites against [XSS](http://en.wikipedia.org/wiki/Cross-site_scripting) whenever user generated content is used.

You should run bluemonday **after** any other processing. So if you use [blackfriday](https://github.com/russross/blackfriday) first then bluemonday should always be the last bit of processing you do. This ensures that no insecurities are introduced later in your process.

Bluemonday is inspired by [OWASP Java HTML Sanitizer](https://code.google.com/p/owasp-java-html-sanitizer/) and [HTML Purifier](http://htmlpurifier.org/). Initial groundwork was laid by [Matt Jibson](https://github.com/mjibson) in his [Google Reader Clone](https://github.com/mjibson/goread).

TODO
====

The TODO list is extensive and bluemonday should **NOT** be used in production code at this point, in fact given that the API will change dramatically you probably shouldn't use it all just now.

1. Add the ability to describe policies for the sanitization process
1. Add the ability to sanitize based on a policy
1. Implement and pass the equivalent of the [OWASP Java HTML Sanitizer tests](https://code.google.com/p/owasp-java-html-sanitizer/source/browse/trunk/src/tests/org/owasp/html/)
1. Open the code to adversarial peer review similar to the [Attack Review Ground Rules](https://code.google.com/p/owasp-java-html-sanitizer/wiki/AttackReviewGroundRules)
1. Bundle a set of default policies that represent safe defaults for one or two key scenarios (content in forum/blog comments, content in RSS)
1. Raise funds and pay for an external security review
