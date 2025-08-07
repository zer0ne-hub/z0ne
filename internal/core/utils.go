// performs utilities tasks like parsing target input
package core

import (
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type TargetType int

const (
	UNKNOWN TargetType = iota
	IP
	DOMAIN
	URL
	FILE
)

func detectTargetType(input string) TargetType {
	input = strings.TrimSpace(input)
	if _, err := os.Stat(input); err == nil || filepath.IsAbs(input) {
		return FILE
	}
	if net.ParseIP(input) != nil {
		return IP
	}
	domainRegex := regexp.MustCompile(`^(?:[a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	if domainRegex.MatchString(input) {
		return DOMAIN
	}
	if u, err := url.Parse(input); err == nil && u.Scheme != "" && u.Host != "" {
		return URL
	}
	return UNKNOWN
}
