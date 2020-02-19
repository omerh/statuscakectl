package helpers

import (
	"fmt"
	"net/url"
	"strings"
)

// DomainValidation validate domain is inserted ok, due to url.Parse parses everything
func DomainValidation(domain string) bool {
	if len(domain) == 0 {
		fmt.Println("Error: Domain is required")
		return false
	}
	if !strings.Contains(domain, ".") {
		fmt.Println("Error: Domain must be in the form of domain.com")
		return false
	}
	_, err := url.Parse(domain)
	if err != nil {
		fmt.Println("Error: failed to parse domain, must be in the form of domain.com")
		return false
	}

	return true
}
