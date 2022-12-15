// SPDX-License-Identifier: BSD-3-Clause

package query

import (
	"errors"
)

// Message is for overall DNS responses.
//
//nolint:govet // Better looking output is worth a few bytes.
type Message struct {
	// Header section
	Header `json:"header,omitempty" xml:"header,omitempty" yaml:"header,omitempty"`
	// Opt Pseudosection
	Opt []Opts `json:"opt,omitempty" xml:"opt,omitempty" yaml:"opt,omitempty"`
	// Question Section
	Question []Question `json:"question,omitempty" xml:"question,omitempty" yaml:"question,omitempty"`
	// Answer Section
	Answer []Answer `json:"answer,omitempty" xml:"answer,omitempty" yaml:"answer,omitempty"`
	// Authority Section
	Authority []Answer `json:"authority,omitempty" xml:"authority,omitempty" yaml:"authority,omitempty"`
	// Additional Section
	Additional []Answer `json:"additional,omitempty" xml:"additional,omitempty" yaml:"additional,omitempty"`
	// Statistics :)
	Statistics `json:"statistics,omitempty" xml:"statistics,omitempty" yaml:"statistics,omitempty"`
}

// Header is the header.
type Header struct {
	Opcode             string `json:"opcode," xml:"opcode," yaml:"opcode" example:"QUERY"`
	Status             string `json:"status," xml:"status," yaml:"status" example:"NOERR"`
	ID                 uint16 `json:"id," xml:"id," yaml:"id" example:"12"`
	Response           bool   `json:"response," xml:"response," yaml:"response" example:"true"`
	Authoritative      bool   `json:"authoritative," xml:"authoritative," yaml:"authoritative" example:"false"`
	Truncated          bool   `json:"truncated," xml:"truncated," yaml:"truncated" example:"false"`
	RecursionDesired   bool   `json:"recursionDesired," xml:"recursionDesired," yaml:"recursionDesired" example:"true"`
	RecursionAvailable bool   `json:"recursionAvailable," xml:"recursionAvailable," yaml:"recursionAvailable" example:"true"`
	Zero               bool   `json:"zero," xml:"zero," yaml:"zero" example:"false"`
	AuthenticatedData  bool   `json:"authenticatedData," xml:"authenticatedData," yaml:"authenticatedData" example:"false"`
	CheckingDisabled   bool   `json:"checkingDisabled," xml:"checkingDisabled," yaml:"checkingDisabled" example:"false"`
}

// Question is a DNS Query.
type Question struct {
	Name  string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty" example:"localhost"`
	Class string `json:"class,omitempty" xml:"class,omitempty" yaml:"class,omitempty" example:"A"`
	Type  string `json:"type,omitempty" xml:"type,omitempty" yaml:"type,omitempty" example:"IN"`
}

// RRHeader is for DNS Resource Headers.
type RRHeader struct {
	Name     string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty" example:"127.0.0.1"`
	TTL      any    `json:"ttl,omitempty" xml:"ttl,omitempty" yaml:"ttl,omitempty" example:"0ms"`
	Class    string `json:"class,omitempty" xml:"class,omitempty" yaml:"class,omitempty" example:"A"`
	Type     string `json:"type,omitempty" xml:"type,omitempty" yaml:"type,omitempty" example:"IN"`
	Rdlength uint16 `json:"-" xml:"-" yaml:"-"`
}

// Opts is for the OPT pseudosection, nearly exclusively for EDNS.
type Opts struct {
	Name  string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty"`
	Value string `json:"value" xml:"value" yaml:"value"`
}

// Answer is for a DNS Response.
type Answer struct {
	Value    string `json:"response,omitempty" xml:"response,omitempty" yaml:"response,omitempty"`
	RRHeader `json:"header,omitempty" xml:"header,omitempty" yaml:"header,omitempty"`
}

// Statistics is the little bit at the bottom :).
type Statistics struct {
	RTT     string `json:"queryTime,omitempty" xml:"queryTime,omitempty" yaml:"queryTime,omitempty"`
	Server  string `json:"server,omitempty" xml:"server,omitempty" yaml:"server,omitempty"`
	When    string `json:"when,omitempty" xml:"when,omitempty" yaml:"when,omitempty"`
	MsgSize int    `json:"msgSize,omitempty" xml:"msgSize,omitempty" yaml:"msgSize,omitempty"`
}

var errNoMessage = errors.New("no message")
