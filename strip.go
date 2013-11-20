package sanitizer

func StripTags(s string) (r string) {
        _, r = Sanitize(s, nil)
        return
}
