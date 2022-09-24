// SPDX-License-Identifier: BSD-3-Clause

package resolvers_test

import (
	"errors"
	"testing"

	"git.froth.zone/sam/awl/pkg/resolvers"
	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestHTTPS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts util.Options
	}{
		{
			"Good",
			util.Options{
				HTTPS:  true,
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server:  "https://dns9.quad9.net/dns-query",
					Type:    dns.TypeA,
					Name:    "git.froth.zone.",
					Retries: 3,
				},
			},
		},
		{
			"404",
			util.Options{
				HTTPS:  true,
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server: "https://dns9.quad9.net/dns",
					Type:   dns.TypeA,
					Name:   "git.froth.zone.",
				},
			},
		},
		{
			"Bad request domain",
			util.Options{
				HTTPS:  true,
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server: "dns9.quad9.net/dns-query",
					Type:   dns.TypeA,
					Name:   "git.froth.zone",
				},
			},
		},
		{
			"Bad server domain",
			util.Options{
				HTTPS:  true,
				Logger: util.InitLogger(0),
				Request: util.Request{
					Server: "dns9..quad9.net/dns-query",
					Type:   dns.TypeA,
					Name:   "git.froth.zone.",
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
			// msg = msg.SetQuestion(testCase.Name, testCase.Type)
			res, err := resolver.LookUp(msg)

			if err == nil {
				assert.NilError(t, err)
				assert.Assert(t, res != util.Response{})
			} else {
				if errors.Is(err, &resolvers.ErrHTTPStatus{}) {
					assert.ErrorContains(t, err, "404")
				}
				assert.Equal(t, res, util.Response{})
			}
		})
	}
}
