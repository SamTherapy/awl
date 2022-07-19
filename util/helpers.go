// SPDX-License-Identifier: BSD-3-Clause

package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

// Given an IP or phone number, return a canonical string to be queried
func ReverseDNS(address string, querInt uint16) (string, error) {
	query := dns.TypeToString[querInt]
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
