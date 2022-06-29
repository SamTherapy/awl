// SPDX-License-Identifier: BSD-3-Clause

package util

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Response struct {
	Answers Answers `json:"Response"` // These be DNS query answers
}

// The Answers struct is the basic structure of a DNS request
// to be returned to the user upon making a request
type Answers struct {
	Server  string        `json:"Server"`  // The server to make the DNS request from
	Request uint16        `json:"Request"` // The type of request
	Name    string        `json:"Name"`    // The domain name to make a DNS request for
	RTT     time.Duration `json:"RTT"`     // The time it took to make the DNS query
}

// Given an IP or phone number, return a canonical string to be queried
func ReverseDNS(address string, query string) (string, error) {
	if query == "PTR" {
		return dns.ReverseAddr(address)
	} else if query == "NAPTR" {
		// get rid of characters not needed
		replacer := strings.NewReplacer("+", "", " ", "", "-", "")
		address = replacer.Replace(address)
		// reverse the order of the string
		address = reverse(address)
		var arpa strings.Builder
		// Make it canonical
		for _, c := range address {
			fmt.Fprintf(&arpa, "%c.", c)
		}
		arpa.WriteString("e164.arpa.")
		return arpa.String(), nil
	}

	return "", errors.New("ReverseDNS: -x flag given but no IP found")
}

// Reverse a string, return the string in reverse
func reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {

		rns[i], rns[j] = rns[j], rns[i]
	}
	return string(rns)
}
