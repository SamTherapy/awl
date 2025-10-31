// SPDX-License-Identifier: BSD-3-Clause
//go:build windows

package conf

import (
	"fmt"
	"strings"
	"unsafe"

	"codeberg.org/miekg/dns/dnsconf"
	"golang.org/x/sys/windows"
)

/*
"Stolen" from
https://gist.github.com/moloch--/9fb1c8497b09b45c840fe93dd23b1e98
*/

// GetDNSConfig (Windows version) returns all DNS server addresses using windows fuckery.
//
// Here be dragons.
func GetDNSConfig() (*dnsconf.Config, error) {
	length := uint32(100000)
	byt := make([]byte, length)

	// Windows is an utter fucking trash fire of an operating system.
	//nolint:gosec // This is necessary unless we want to drop 1.18
	if err := windows.GetAdaptersAddresses(windows.AF_UNSPEC, windows.GAA_FLAG_INCLUDE_PREFIX, 0, (*windows.IpAdapterAddresses)(unsafe.Pointer(&byt[0])), &length); err != nil {
		return nil, fmt.Errorf("config, windows: %w", err)
	}

	var addresses []*windows.IpAdapterAddresses
	//nolint:gosec // This is necessary unless we want to drop 1.18
	for addr := (*windows.IpAdapterAddresses)(unsafe.Pointer(&byt[0])); addr != nil; addr = addr.Next {
		addresses = append(addresses, addr)
	}

	resolvers := map[string]bool{}

	for _, addr := range addresses {
		for next := addr.FirstUnicastAddress; next != nil; next = next.Next {
			if addr.OperStatus != windows.IfOperStatusUp {
				continue
			}

			if next.Address.IP() != nil {
				for dnsServer := addr.FirstDnsServerAddress; dnsServer != nil; dnsServer = dnsServer.Next {
					ip := dnsServer.Address.IP()

					if ip.IsMulticast() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() || ip.IsUnspecified() {
						continue
					}

					if ip.To16() != nil && strings.HasPrefix(ip.To16().String(), "fec0:") {
						continue
					}

					resolvers[ip.String()] = true
				}

				break
			}
		}
	}

	// Take unique values only
	servers := []string{}
	for server := range resolvers {
		servers = append(servers, server)
	}

	// TODO: Make configurable, based on defaults in https://codeberg.org/miekg/dns/blob/master/clientconfig.go
	return &dnsconf.Config{
		Servers:  servers,
		Search:   []string{},
		Port:     "53",
		Ndots:    1,
		Timeout:  5,
		Attempts: 1,
	}, nil
}
