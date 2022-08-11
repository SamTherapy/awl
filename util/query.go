// SPDX-License-Identifier: BSD-3-Clause

package util

import (
	"time"

	"github.com/miekg/dns"
)

// Response is the DNS response.
type Response struct {
	DNS *dns.Msg      // The full DNS response
	RTT time.Duration `json:"rtt"` // The time it took to make the DNS query
}

// Request is a structure for a DNS query.
type Request struct {
	Server  string `json:"server"`
	Name    string `json:"name"`
	Timeout time.Duration
	Retries int
	Type    uint16 `json:"request"`
	Class   uint16 `json:"class"`
}
