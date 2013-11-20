Bluemonday
==========

Bluemonday is a HTML sanitizer implemented in Go. You can safely feed it user generated content and it will give back HTML that is scrubbed using a whitelist of approved tags and attributes. It is fast, configurable and is safe for all utf-8 (unicode) input.

You should run bluemonday **after** any other processing. So if you use [blackfriday](https://github.com/russross/blackfriday) first then bluemonday should always be the last bit of processing you do. This ensures that no insecurities are introduced later in your process.
