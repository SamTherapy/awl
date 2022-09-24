// SPDX-License-Identifier: BSD-3-Clause

package util

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

type errReverseDNS struct {
	addr string
}

func (errDNS *errReverseDNS) Error() string {
	return fmt.Sprintf("reverseDNS: invalid value %s given", errDNS.addr)
}

// ReverseDNS is given an IP or phone number and returns a canonical string to be queried.
func ReverseDNS(address string, querInt uint16) (string, error) {
	query := dns.TypeToString[querInt]
	if query == "PTR" {
		str, err := dns.ReverseAddr(address)
		if err != nil {
			return "", fmt.Errorf("PTR reverse: %w", err)
		}

		return str, nil
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

	return "", &errReverseDNS{address}
}

// Reverse a string, return the string in reverse.
func reverse(s string) string {
	rns := []rune(s)
	for i, j := 0, len(rns)-1; i < j; i, j = i+1, j-1 {
		rns[i], rns[j] = rns[j], rns[i]
	}

	return string(rns)
}
