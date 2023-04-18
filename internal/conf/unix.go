// SPDX-License-Identifier: BSD-3-Clause
//go:build unix || (!windows && !plan9 && !js && !zos)

// FIXME: Can remove the or on the preprocessor when Go 1.18 becomes obsolete

package conf

import (
	"fmt"

	"github.com/miekg/dns"
)

// GetDNSConfig gets the DNS configuration, either from /etc/resolv.conf or somewhere else.
func GetDNSConfig() (*dns.ClientConfig, error) {
	conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("unix config: %w", err)
	}

	return conf, nil
}
