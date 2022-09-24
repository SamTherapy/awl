// SPDX-License-Identifier: BSD-3-Clause

package resolvers_test

import (
	"testing"
	"time"

	"git.froth.zone/sam/awl/pkg/resolvers"
	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestQuic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts util.Options
	}{
		{
			"Valid",
			util.Options{
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
			util.Options{
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
			util.Options{
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
			util.Options{
				QUIC:   true,
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server:  "example.com",
					Port:    853,
					Type:    dns.TypeA,
					Name:    "git.froth.zone",
					Timeout: 10 * time.Millisecond,
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			resolver, err := resolvers.LoadResolver(test.opts)
			assert.NilError(t, err)

			msg := new(dns.Msg)
			msg.SetQuestion(test.opts.Request.Name, test.opts.Request.Type)

			res, err := resolver.LookUp(msg)

			if err == nil {
				assert.NilError(t, err)
				assert.Assert(t, res != util.Response{})
			} else {
				assert.Assert(t, res == util.Response{})
			}
		})
	}
}
