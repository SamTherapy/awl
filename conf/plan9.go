// SPDX-License-Identifier: BSD-3-Clause

package conf

import (
	"errors"
	"strings"

	"github.com/miekg/dns"
)

// Plan 9 stores its network data in /net/ndb, which seems to be formatted a specific way
// Yoink it and use it.
//
// See ndb(7).
func GetPlan9Config(str string) (*dns.ClientConfig, error) {
	str = strings.ReplaceAll(str, "\n", "")
	spl := strings.FieldsFunc(str, splitChars)
	var servers []string
	for _, option := range spl {
		if strings.HasPrefix(option, "dns=") {
			servers = append(servers, strings.TrimPrefix(option, "dns="))
		}
	}
	if len(servers) == 0 {
		return nil, errPlan9
	}

	// TODO: read more about how customizable Plan 9 is
	return &dns.ClientConfig{
		Servers: servers,
		Search:  []string{},
		Port:    "53",
	}, nil
}

// Split the string at either space or tabs.
func splitChars(r rune) bool {
	return r == ' ' || r == '\t'
}

var errPlan9 = errors.New("plan9Config: no DNS servers found")
