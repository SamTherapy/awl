// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"errors"
)

// Message is for overall DNS responses.
//
//nolint:govet,tagliatelle // Better looking output is worth a few bytes.
type Message struct {
	DateString  string `json:"dateString,omitempty" xml:"dateString,omitempty" yaml:"dateString,omitempty"`
	DateSeconds int64  `json:"dateSeconds,omitempty" xml:"dateSeconds,omitempty" yaml:"dateSeconds,omitempty"`
	MsgSize     int    `json:"msgLength,omitempty" xml:"msgSize,omitempty" yaml:"msgSize,omitempty"`
	ID          uint16 `json:"ID" xml:"ID" yaml:"ID" example:"12"`

	Opcode             int  `json:"opcode" xml:"opcode" yaml:"opcode" example:"QUERY"`
	Response           bool `json:"QR" xml:"QR" yaml:"QR" example:"true"`
	Authoritative      bool `json:"AA" xml:"AA" yaml:"AA" example:"false"`
	Truncated          bool `json:"TC" xml:"TC" yaml:"TC" example:"false"`
	RecursionDesired   bool `json:"RD" xml:"RD" yaml:"RD" example:"true"`
	RecursionAvailable bool `json:"RA" xml:"RA" yaml:"RA" example:"true"`
	AuthenticatedData  bool `json:"AD" xml:"AD" yaml:"AD" example:"false"`
	CheckingDisabled   bool `json:"CD" xml:"CD" yaml:"CD" example:"false"`
	Zero               bool `json:"Z" xml:"Z" yaml:"Z" example:"false"`

	QdCount int `json:"QDCOUNT" xml:"QDCOUNT" yaml:"QDCOUNT" example:"0"`
	AnCount int `json:"ANCOUNT" xml:"ANCOUNT" yaml:"ANCOUNT" example:"0"`
	NsCount int `json:"NSCOUNT" xml:"NSCOUNT" yaml:"NSCOUNT" example:"0"`
	ArCount int `json:"ARCOUNT" xml:"ARCOUNT" yaml:"ARCOUNT" example:"0"`

	Name      string `json:"QNAME,omitempty" xml:"QNAME,omitempty" yaml:"QNAME,omitempty" example:"localhost"`
	Type      uint16 `json:"QTYPE,omitempty" xml:"QTYPE,omitempty" yaml:"QTYPE,omitempty" example:"IN"`
	TypeName  string `json:"QTYPEname,omitempty" xml:"QTYPEname,omitempty" yaml:"QTYPEname,omitempty" example:"IN"`
	Class     uint16 `json:"QCLASS,omitempty" xml:"QCLASS,omitempty" yaml:"QCLASS,omitempty" example:"A"`
	ClassName string `json:"QCLASSname,omitempty" xml:"QCLASSname,omitempty" yaml:"QCLASSname,omitempty" example:"1"`

	EDNS0 EDNS0 `json:",omitempty" xml:",omitempty" yaml:",omitempty"`

	// Answer Section
	AnswerRRs        []Answer `json:"answersRRs,omitempty" xml:"answersRRs,omitempty" yaml:"answersRRs,omitempty" example:"false"`
	AuthoritativeRRs []Answer `json:"authorityRRs,omitempty" xml:"authorityRRs,omitempty" yaml:"authorityRRs,omitempty" example:"false"`
	AdditionalRRs    []Answer `json:"additionalRRs,omitempty" xml:"additionalRRs,omitempty" yaml:"additionalRRs,omitempty" example:"false"`
}

// Answer is for DNS Resource Headers.
//
//nolint:govet,tagliatelle
type Answer struct {
	Name      string `json:"NAME,omitempty" xml:"NAME,omitempty" yaml:"NAME,omitempty" example:"127.0.0.1"`
	Type      uint16 `json:"TYPE,omitempty" xml:"TYPE,omitempty" yaml:"TYPE,omitempty" example:"1"`
	TypeName  string `json:"TYPEname,omitempty" xml:"TYPEname,omitempty" yaml:"TYPEname,omitempty" example:"A"`
	Class     uint16 `json:"CLASS,omitempty" xml:"CLASS,omitempty" yaml:"CLASS,omitempty" example:"1"`
	ClassName string `json:"CLASSname,omitempty" xml:"CLASSname,omitempty" yaml:"CLASSname,omitempty" example:"IN"`
	TTL       any    `json:"TTL,omitempty" xml:"TTL,omitempty" yaml:"TTL,omitempty" example:"0ms"`
	Value     string `json:"rdata,omitempty" xml:"rdata,omitempty" yaml:"rdata,omitempty"`
	Rdlength  uint16 `json:"RDLENGTH,omitempty" xml:"RDLENGTH,omitempty" yaml:"RDLENGTH,omitempty"`
	Rdhex     string `json:"RDATAHEX,omitempty" xml:"RDATAHEX,omitempty" yaml:"RDATAHEX,omitempty"`
}

// EDNS0 is for all EDNS options.
//
// RFC: https://datatracker.ietf.org/docs/draft-peltan-edns-presentation-format/
//
//nolint:govet,tagliatelle
type EDNS0 struct {
	Flags       []string    `json:"FLAGS" xml:"FLAGS" yaml:"FLAGS"`
	Rcode       string      `json:"RCODE" xml:"RCODE" yaml:"RCODE"`
	PayloadSize uint16      `json:"UDPSIZE" xml:"UDPSIZE" yaml:"UDPSIZE"`
	LLQ         *EdnsLLQ    `json:"LLQ,omitempty" xml:"LLQ,omitempty" yaml:"LLQ,omitempty"`
	NsidHex     string      `json:"NSIDHEX,omitempty" xml:"NSIDHEX,omitempty" yaml:"NSIDHEX,omitempty"`
	Nsid        string      `json:"NSID,omitempty" xml:"NSID,omitempty" yaml:"NSID,omitempty"`
	Dau         []uint8     `json:"DAU,omitempty" xml:"DAU,omitempty" yaml:"DAU,omitempty"`
	Dhu         []uint8     `json:"DHU,omitempty" xml:"DHU,omitempty" yaml:"DHU,omitempty"`
	N3u         []uint8     `json:"N3U,omitempty" xml:"N3U,omitempty" yaml:"N3U,omitempty"`
	Subnet      *EDNSSubnet `json:"ECS,omitempty" xml:"ECS,omitempty" yaml:"ECS,omitempty"`
	Expire      uint32      `json:"EXPIRE,omitempty" xml:"EXPIRE,omitempty" yaml:"EXPIRE,omitempty"`
	Cookie      []string    `json:"COOKIE,omitempty" xml:"COOKIE,omitempty" yaml:"COOKIE,omitempty"`
	KeepAlive   uint16      `json:"KEEPALIVE,omitempty" xml:"KEEPALIVE,omitempty" yaml:"KEEPALIVE,omitempty"`
	Padding     string      `json:"PADDING,omitempty" xml:"PADDING,omitempty" yaml:"PADDING,omitempty"`
	Chain       string      `json:"CHAIN,omitempty" xml:"CHAIN,omitempty" yaml:"CHAIN,omitempty"`
	EDE         *EDNSErr    `json:"EDE,omitempty" xml:"EDE,omitempty" yaml:"EDE,omitempty"`
}

// EdnsLLQ is for Long-lived queries.
//
//nolint:tagliatelle
type EdnsLLQ struct {
	Version uint16 `json:"LLQ-VERSION" xml:"LLQ-VERSION" yaml:"LLQ-VERSION"`
	Opcode  uint16 `json:"LLQ-OPCODE" xml:"LLQ-OPCODE" yaml:"LLQ-OPCODE"`
	Error   uint16 `json:"LLQ-ERROR" xml:"LLQ-ERROR" yaml:"LLQ-ERROR"`
	ID      uint64 `json:"LLQ-ID" xml:"LLQ-ID" yaml:"LLQ-ID"`
	Lease   uint32 `json:"LLQ-LEASE" xml:"LLQ-LEASE" yaml:"LLQ-LEASE"`
}

// EDNSSubnet is for EDNS subnet options,
//
//nolint:govet,tagliatelle
type EDNSSubnet struct {
	Family uint16 `json:"FAMILY" xml:"FAMILY" yaml:"FAMILY"`
	IP     string
	Source uint8 `json:"SOURCE" xml:"SOURCE" yaml:"SOURCE"`
	Scope  uint8 `json:"SCOPE,omitempty" xml:"SCOPE,omitempty" yaml:"SCOPE,omitempty"`
}

// EDNSErr is for EDE codes
//
//nolint:govet,tagliatelle
type EDNSErr struct {
	Code    uint16 `json:"INFO-CODE" xml:"INFO-CODE" yaml:"INFO-CODE"`
	Purpose string
	Text    string `json:"EXTRA-TEXT,omitempty" xml:"EXTRA-TEXT,omitempty" yaml:"EXTRA-TEXT,omitempty"`
}

var errNoMessage = errors.New("no message")
