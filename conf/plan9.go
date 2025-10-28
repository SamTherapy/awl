// SPDX-License-Identifier: BSD-3-Clause
//go:build plan9

package conf

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/miekg/dns"
)

// GetDNSConfig gets DNS information from Plan 9, because it's different from UNIX and Windows.
//
// Plan 9 stores its network data in /net/ndb, which seems to be formatted a specific way
// Yoink it and use it.
//
// See ndb(7).
func GetDNSConfig() (*dns.ClientConfig, error) {
	dat, err := os.ReadFile("/net/ndb")
	if err != nil {
		return nil, fmt.Errorf("read ndb: %w", err)
	}

	str := string(dat)

	// str = strings.ReplaceAll(str, "\n", "")
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
	return r == ' ' || r == '\t' || r == '\n'
}

var errPlan9 = errors.New("plan9Config: no DNS servers found")
