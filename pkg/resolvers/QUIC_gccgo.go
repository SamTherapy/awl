// SPDX-License-Identifier: BSD-3-Clause
//go:build gccgo

// TODO: Whenever gccgo supports quic-go, delete this
package resolvers

import (
	"errors"

	"dns.froth.zone/awl/pkg/util"
	"github.com/miekg/dns"
)

// QUICResolver is for DNS-over-QUIC queries.
type QUICResolver struct {
	opts *util.Options
}

var _ Resolver = (*QUICResolver)(nil)

var errNotImplemented = errors.New("DNS-over-QUIC not supported when running gccgo!")

// LookUp cannot be used with gccgo because gccgo does not (and likely will not) support generics.
func (resolver *QUICResolver) LookUp(msg *dns.Msg) (util.Response, error) {
	return util.Response{}, errNotImplemented
}
