// SPDX-License-Identifier: BSD-3-Clause
//go:build !windows
// +build !windows

package conf

import (
	"os"
	"runtime"

	"github.com/miekg/dns"
)

// Get the DNS configuration, either from /etc/resolv.conf or somewhere else
func GetDNSConfig() (*dns.ClientConfig, error) {
	if runtime.GOOS == "plan9" {
		dat, err := os.ReadFile("/net/ndb")
		if err != nil {
			return nil, err
		}
		return getPlan9Config(string(dat))
	} else {
		return dns.ClientConfigFromFile("/etc/resolv.conf")
	}
}
