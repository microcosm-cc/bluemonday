package bluemonday

import "fmt"

// ErrEmptyInput is returned when Sanitize() is called with an empty string or
// a string only composed of whitespace
var ErrEmptyInput = fmt.Errorf("Input is empty")

// ErrNotImplemented is returned when the dependent lib
// code.google.com/p/go.net/html has been updated to support token types not
// presently supported (by the lib or us).
//
// We do not anticipate ever returning this, if anyone ever sees this error in
// the wild please raise a github issue immediately so that we can update our
// code: https://github.com/microcosm-cc/bluemonday/issues
var ErrNotImplemented = fmt.Errorf("Not implemented")
