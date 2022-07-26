// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"errors"

	"git.froth.zone/sam/awl/internal/helpers"
	"git.froth.zone/sam/awl/logawl"
)

// CLI options structure.
type Options struct {
	Logger    *logawl.Logger // Logger
	Port      int            // DNS port
	IPv4      bool           // Force IPv4
	IPv6      bool           // Force IPv6
	DNSSEC    bool           // Enable DNSSEC
	TCP       bool           // Query with TCP
	DNSCrypt  bool           // Query over DNSCrypt
	TLS       bool           // Query over TLS
	HTTPS     bool           // Query over HTTPS
	QUIC      bool           // Query over QUIC
	Truncate  bool           // Ignore truncation
	AA        bool           // Set Authoratative Answer
	AD        bool           // Set Authenticated Data
	CD        bool           // Set CD
	QR        bool           // Set QueRy
	RD        bool           // Set Recursion Desired
	RA        bool           // Set Recursion Available
	TC        bool           // Set TC (TrunCated)
	Z         bool           // Set Z (Zero)
	Reverse   bool           // Make reverse query
	Verbosity int            // Set logawl verbosity
	// HumanTTL  bool           // Make TTL human readable
	Short bool // Short output
	JSON  bool // Outout as JSON
	XML   bool // Output as XML
	YAML  bool // Output at YAML

	Display Displays        // Display options
	Request helpers.Request // DNS reuqest
}

// What to (and not to) display
type Displays struct {
	// Comments   bool
	Question   bool // QUESTION SECTION
	Answer     bool // ANSWER SECTION
	Authority  bool // AUTHORITY SECTION
	Additional bool // ADDITIONAL SECTION
	Statistics bool // Query time, message size, etc.
}

var ErrNotError = errors.New("not an error")
