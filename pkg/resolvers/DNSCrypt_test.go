// SPDX-License-Identifier: BSD-3-Clause

package resolvers_test

import (
	"errors"
	"testing"

	"git.froth.zone/sam/awl/pkg/query"
	"git.froth.zone/sam/awl/pkg/util"
	"github.com/ameshkov/dnscrypt/v2"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestDNSCrypt(t *testing.T) {
	t.Parallel()

	//nolint:govet // I could not be assed to refactor this, and it is only for tests
	tests := []struct {
		name string
		opts *util.Options
	}{
		{
			"Valid",
			&util.Options{
				Logger:   util.InitLogger(0),
				DNSCrypt: true,
				Request: util.Request{
					Server:  "sdns://AQMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20",
					Type:    dns.TypeA,
					Name:    "example.com.",
					Retries: 3,
				},
			},
		},
		{
			"Valid (TCP)",
			&util.Options{
				Logger:   util.InitLogger(0),
				DNSCrypt: true,
				TCP:      true,
				IPv4:     true,
				Request: util.Request{
					Server:  "sdns://AQMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20",
					Type:    dns.TypeAAAA,
					Name:    "example.com.",
					Retries: 3,
				},
			},
		},
		{
			"Invalid",
			&util.Options{
				Logger:   util.InitLogger(0),
				DNSCrypt: true,
				TCP:      true,
				IPv4:     true,
				Request: util.Request{
					Server:  "QMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20",
					Type:    dns.TypeAAAA,
					Name:    "example.com.",
					Retries: 0,
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			var (
				res util.Response
				err error
			)
			for i := 0; i <= test.opts.Request.Retries; i++ {
				res, err = query.CreateQuery(test.opts)
				if err == nil || errors.Is(err, dnscrypt.ErrInvalidDNSStamp) {
					break
				}
			}

			if err == nil {
				assert.Assert(t, res != util.Response{})
			} else {
				assert.ErrorContains(t, err, "unsupported stamp")
			}
		})
	}
}
