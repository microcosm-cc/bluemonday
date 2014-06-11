package bluemonday

import (
	"fmt"
	"strings"
)

// Sanitize takes a string that contains HTML elements and applies the policy
// whitelist to return a string that has been sanitized by the policy.
func (p *policy) Sanitize(s string) (string, error) {
	if strings.TrimSpace(s) == "" {
		return "", nil
	}

	return "", fmt.Errorf("Sanitize is not yet implemented")
}
