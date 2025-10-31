// SPDX-License-Identifier: BSD-3-Clause
//go:build js

package conf

import (
	"errors"

	"codeberg.org/miekg/dns"
)

// GetDNSConfig doesn't do anything, because it is impossible (and bad security)
// if it could, as that is the definition of a DNS leak.
func GetDNSConfig() (*dns.ClientConfig, error) {
	return nil, errNotImplemented
}

var errNotImplemented = errors.New("not implemented")
