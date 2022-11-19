// SPDX-License-Identifier: BSD-3-Clause
//go:build !gccgo

package resolvers_test

import (
	"testing"
	"time"

	"git.froth.zone/sam/awl/pkg/query"
	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestQuic(t *testing.T) {
	t.Parallel()

	//nolint:govet // I could not be assed to refactor this, and it is only for tests
	tests := []struct {
		name string
		opts *util.Options
	}{
		{
			"Valid",
			&util.Options{
				QUIC:   true,
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server:  "dns.adguard.com",
					Type:    dns.TypeNS,
					Port:    853,
					Timeout: 750 * time.Millisecond,
					Retries: 3,
				},
			},
		},
		{
			"Bad domain",
			&util.Options{
				QUIC:   true,
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server:  "dns.//./,,adguard\a.com",
					Port:    853,
					Type:    dns.TypeA,
					Name:    "git.froth.zone",
					Timeout: 100 * time.Millisecond,
					Retries: 0,
				},
			},
		},
		{
			"Not canonical",
			&util.Options{
				QUIC:   true,
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server:  "dns.adguard.com",
					Port:    853,
					Type:    dns.TypeA,
					Name:    "git.froth.zone",
					Timeout: 100 * time.Millisecond,
					Retries: 0,
				},
			},
		},
		{
			"Invalid query domain",
			&util.Options{
				QUIC:   true,
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server:  "example.com",
					Port:    853,
					Type:    dns.TypeA,
					Name:    "git.froth.zone",
					Timeout: 10 * time.Millisecond,
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
				if err == nil {
					break
				}
			}

			if err == nil {
				assert.NilError(t, err)
				assert.Assert(t, res != util.Response{})
			} else {
				assert.Assert(t, res == util.Response{})
			}
		})
	}
}
