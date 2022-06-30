// SPDX-License-Identifier: BSD-3-Clause

package conf

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

// Plan 9 stores its network data in /net/ndb, which seems to be formatted a specific way
// Yoink it and use it.
//
// See ndb(7).
func getPlan9Config(str string) (*dns.ClientConfig, error) {
	str = strings.ReplaceAll(str, "\n", "")
	spl := strings.FieldsFunc(str, splitChars)
	var servers []string
	for _, option := range spl {
		if strings.HasPrefix(option, "dns=") {
			servers = append(servers, strings.TrimPrefix(option, "dns="))
		}
	}
	if len(servers) == 0 {
		return nil, fmt.Errorf("plan9: no DNS servers found")
	}

	// TODO: read more about how customizable Plan 9 is
	return &dns.ClientConfig{
		Servers: servers,
		Search:  []string{},
		Port:    "53",
	}, nil
}

// Split the string at either space or tabs
func splitChars(r rune) bool {
	return r == ' ' || r == '\t'
}
