// SPDX-License-Identifier: BSD-3-Clause

package resolvers_test

import (
	"errors"
	"testing"

	"dns.froth.zone/awl/pkg/query"
	"dns.froth.zone/awl/pkg/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestHTTPS(t *testing.T) {
	t.Parallel()

	//nolint:govet // I could not be assed to refactor this, and it is only for tests
	tests := []struct {
		name string
		opts *util.Options
	}{
		{
			"Good",
			&util.Options{
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
			&util.Options{
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
			&util.Options{
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
			&util.Options{
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

			var (
				res util.Response
				err error
			)
			for i := 0; i <= test.opts.Request.Retries; i++ {
				res, err = query.CreateQuery(test.opts)
				if err == nil || errors.Is(err, &util.ErrHTTPStatus{}) {
					break
				}
			}

			if err == nil {
				assert.NilError(t, err)
				assert.Assert(t, res != util.Response{})
			} else {
				if errors.Is(err, &util.ErrHTTPStatus{}) {
					assert.ErrorContains(t, err, "404")
				}
				assert.Equal(t, res, util.Response{})
			}
		})
	}
}
