// SPDX-License-Identifier: BSD-3-Clause

package util

import (
	"time"

	"codeberg.org/miekg/dns"
)

// Response is the DNS response.
type Response struct {
	// The full DNS response
	DNS *dns.Msg `json:"response"`
	// The time it took to make the DNS query
	RTT time.Duration `json:"rtt" example:"2000000000"`
}

// Request is a structure for a DNS query.
type Request struct {
	// Server to query
	Server string `json:"server" example:"1.0.0.1"`
	// Domain to query
	Name string `json:"name" example:"example.com"`
	// Duration to wait until marking request as failed
	Timeout time.Duration `json:"timeout" example:"2000000000"`
	// Port to make DNS request on
	Port int `json:"port" example:"53"`
	// Number of failures to make before giving up
	Retries int `json:"retries" example:"2"`
	// Request type, eg. A, AAAA, NAPTR
	Type uint16 `json:"type" example:"1"`
	// Request class, eg. IN
	Class uint16 `json:"class" example:"1"`
}
