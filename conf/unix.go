// SPDX-License-Identifier: BSD-3-Clause
//go:build !windows

package conf

import (
	"fmt"
	"os"
	"runtime"

	"github.com/miekg/dns"
)

// GetDNSConfig gets the DNS configuration, either from /etc/resolv.conf or somewhere else.
func GetDNSConfig() (*dns.ClientConfig, error) {
	if runtime.GOOS == "plan9" {
		dat, err := os.ReadFile("/net/ndb")
		if err != nil {
			return nil, fmt.Errorf("plan9 ndb: %w", err)
		}

		return GetPlan9Config(string(dat))
	} else {
		conf, err := dns.ClientConfigFromFile("/etc/resolv.conf")
		if err != nil {
			return nil, fmt.Errorf("unix config: %w", err)
		}

		return conf, nil
	}
}
