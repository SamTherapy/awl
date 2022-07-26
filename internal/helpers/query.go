// SPDX-License-Identifier: BSD-3-Clause

package helpers

import (
	"time"

	"github.com/miekg/dns"
)

// The DNS response.
type Response struct {
	DNS *dns.Msg      // The full DNS response
	RTT time.Duration `json:"rtt"` // The time it took to make the DNS query
}

// A structure for a DNS query.
type Request struct {
	Server  string        `json:"server"`  // The server to make the DNS request from
	Type    uint16        `json:"request"` // The type of request
	Class   uint16        `json:"class"`   // DNS Class
	Name    string        `json:"name"`    // The domain name to make a DNS request for
	Timeout time.Duration // The maximum timeout
	Retries int           // Number of queries to retry
}
