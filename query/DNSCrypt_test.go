// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"testing"

	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestDNSCrypt(t *testing.T) {
	t.Parallel()

	tests := []struct {
		opt util.Options
	}{
		{
			util.Options{
				Logger:   util.InitLogger(0),
				DNSCrypt: true,
				Request: util.Request{
					Server: "sdns://AQMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20",
					Type:   dns.TypeA,
					Name:   "example.com.",
				},
			},
		},
		{
			util.Options{
				Logger:   util.InitLogger(0),
				DNSCrypt: true,
				TCP:      true,
				IPv4:     true,
				Request: util.Request{
					Server: "sdns://AQMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20",
					Type:   dns.TypeAAAA,
					Name:   "example.com.",
				},
			},
		},
		{
			util.Options{
				Logger:   util.InitLogger(0),
				DNSCrypt: true,
				TCP:      true,
				IPv4:     true,
				Request: util.Request{
					Server: "QMAAAAAAAAAETk0LjE0MC4xNC4xNDo1NDQzINErR_JS3PLCu_iZEIbq95zkSV2LFsigxDIuUso_OQhzIjIuZG5zY3J5cHQuZGVmYXVsdC5uczEuYWRndWFyZC5jb20",
					Type:   dns.TypeAAAA,
					Name:   "example.com.",
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run("", func(t *testing.T) {
			t.Parallel()

			res, err := query.CreateQuery(test.opt)
			if err == nil {
				assert.Assert(t, res != util.Response{})
			} else {
				assert.ErrorContains(t, err, "unsupported stamp")
			}
		})
	}
}
