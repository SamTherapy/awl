// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"testing"

	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestCreateQ(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		opts util.Options
	}{
		{
			"1",
			util.Options{
				Logger: util.InitLogger(0),
				HeaderFlags: util.HeaderFlags{
					Z: true,
				},
				YAML: true,
				Request: util.Request{
					Server:  "8.8.4.4",
					Port:    53,
					Type:    dns.TypeA,
					Name:    "example.com.",
					Retries: 3,
				},
				Display: util.Display{
					Comments:   true,
					Question:   true,
					Opt:        true,
					Answer:     true,
					Authority:  true,
					Additional: true,
					Statistics: true,
					ShowQuery:  true,
				},
				EDNS: util.EDNS{
					ZFlag:      1,
					BufSize:    1500,
					EnableEDNS: true,
					Cookie:     true,
					DNSSEC:     true,
					Expire:     true,
					KeepOpen:   true,
					Nsid:       true,
					Padding:    true,
					Version:    0,
				},
			},
		},
		{
			"2",
			util.Options{
				Logger: util.InitLogger(0),
				HeaderFlags: util.HeaderFlags{
					Z: true,
				},
				XML: true,

				Request: util.Request{
					Server:  "8.8.4.4",
					Port:    53,
					Type:    dns.TypeA,
					Name:    "example.com.",
					Retries: 3,
				},
				Display: util.Display{
					Comments:       true,
					Question:       true,
					Opt:            true,
					Answer:         true,
					Authority:      true,
					Additional:     true,
					Statistics:     true,
					UcodeTranslate: true,
					ShowQuery:      true,
				},
			},
		},
		{
			"3",
			util.Options{
				Logger: util.InitLogger(0),
				JSON:   true,
				QUIC:   true,

				Request: util.Request{
					Server:  "dns.adguard.com",
					Port:    853,
					Type:    dns.TypeA,
					Name:    "example.com.",
					Retries: 3,
				},
				Display: util.Display{
					Comments:   true,
					Question:   true,
					Opt:        true,
					Answer:     true,
					Authority:  true,
					Additional: true,
					Statistics: true,
					ShowQuery:  true,
				},
				EDNS: util.EDNS{
					EnableEDNS: true,
					DNSSEC:     true,
					Cookie:     true,
					Expire:     true,
					Nsid:       true,
				},
			},
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			res, err := query.CreateQuery(test.opts)
			assert.NilError(t, err)
			assert.Assert(t, res != util.Response{})

			str, err := query.PrintSpecial(res, test.opts)

			assert.NilError(t, err)
			assert.Assert(t, str != "")

			str, err = query.ToString(res, test.opts)
			assert.NilError(t, err)
			assert.Assert(t, str != "")
		})
	}
}
