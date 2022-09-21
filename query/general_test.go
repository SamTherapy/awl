// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"os"
	"testing"
	"time"

	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestResolve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts util.Options
	}{
		{
			"UDP",
			util.Options{
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server:  "8.8.4.4",
					Port:    53,
					Type:    dns.TypeAAAA,
					Name:    "example.com.",
					Retries: 3,
				},
			},
		},
		{
			"UDP (Bad Cookie)",
			util.Options{
				Logger:    util.InitLogger(0),
				BadCookie: false,
				Request: util.Request{
					Server:  "b.root-servers.net",
					Port:    53,
					Type:    dns.TypeNS,
					Name:    "example.com.",
					Retries: 3,
				},
				EDNS: util.EDNS{
					EnableEDNS: true,
					Cookie:     true,
				},
			},
		},
		{
			"UDP (Truncated)",
			util.Options{
				Logger: util.InitLogger(0),
				IPv4:   true,
				Request: util.Request{
					Server:  "madns.binarystar.systems",
					Port:    5301,
					Type:    dns.TypeTXT,
					Name:    "limit.txt.example.",
					Retries: 3,
				},
			},
		},
		{
			"TCP",
			util.Options{
				Logger: util.InitLogger(0),
				TCP:    true,

				Request: util.Request{
					Server:  "8.8.4.4",
					Port:    53,
					Type:    dns.TypeA,
					Name:    "example.com.",
					Retries: 3,
				},
			},
		},
		{
			"TLS",
			util.Options{
				Logger: util.InitLogger(0),
				TLS:    true,
				Request: util.Request{
					Server:  "dns.google",
					Port:    853,
					Type:    dns.TypeAAAA,
					Name:    "example.com.",
					Retries: 3,
				},
			},
		},
		{
			"Timeout",
			util.Options{
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server:  "8.8.4.1",
					Port:    1,
					Type:    dns.TypeA,
					Name:    "example.com.",
					Timeout: time.Millisecond * 100,
					Retries: 0,
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			res, err := query.CreateQuery(test.opts)
			if err == nil {
				assert.NilError(t, err)
				assert.Assert(t, res != util.Response{})
			} else {
				assert.ErrorIs(t, err, os.ErrDeadlineExceeded)
			}
		})
	}
}
