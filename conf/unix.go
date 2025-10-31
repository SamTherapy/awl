// SPDX-License-Identifier: BSD-3-Clause
//go:build unix || (!windows && !plan9 && !js && !zos)

// FIXME: Can remove the or on the preprocessor when Go 1.18 becomes obsolete

package conf

import (
	"fmt"

	"codeberg.org/miekg/dns/dnsconf"
)

// GetDNSConfig gets the DNS configuration, either from /etc/resolv.conf or somewhere else.
func GetDNSConfig() (*dnsconf.Config, error) {
	conf, err := dnsconf.FromFile("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("unix config: %w", err)
	}

	return conf, nil
}
