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

	in := []util.Options{
		{
			Logger:    util.InitLogger(0),
			Port:      53,
			Z:         true,
			ShowQuery: true,
			YAML:      true,

			Request: util.Request{
				Server: "8.8.4.4",
				Type:   dns.TypeA,
				Name:   "example.com.",
			},
			Display: util.Displays{
				Comments:   true,
				Question:   true,
				Opt:        true,
				Answer:     true,
				Authority:  true,
				Additional: true,
				Statistics: true,
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
		{
			Logger:    util.InitLogger(0),
			Port:      53,
			Z:         true,
			ShowQuery: true,
			XML:       true,

			Request: util.Request{
				Server: "8.8.4.4",
				Type:   dns.TypeA,
				Name:   "example.com.",
			},
			Display: util.Displays{
				Comments:       true,
				Question:       true,
				Opt:            true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: true,
			},
		},
		{
			Logger: util.InitLogger(0),
			Port:   853,
			// Z:         true,
			ShowQuery: true,
			JSON:      true,
			QUIC:      true,

			Request: util.Request{
				Server: "dns.adguard.com",
				Type:   dns.TypeA,
				Name:   "example.com.",
			},
			Display: util.Displays{
				Comments:   true,
				Question:   true,
				Opt:        true,
				Answer:     true,
				Authority:  true,
				Additional: true,
				Statistics: true,
			},
			EDNS: util.EDNS{
				EnableEDNS: true,
				DNSSEC:     true,
				Cookie:     true,
				Expire:     true,
				Nsid:       true,
			},
		},
	}

	for _, opt := range in {
		opt := opt

		t.Run("", func(t *testing.T) {
			t.Parallel()
			res, err := query.CreateQuery(opt)
			assert.NilError(t, err)
			assert.Assert(t, res != util.Response{})

			str, err := query.PrintSpecial(res.DNS, opt)

			assert.NilError(t, err)
			assert.Assert(t, str != "")

			str, err = query.ToString(res, opt)
			assert.NilError(t, err)
			assert.Assert(t, str != "")
		})
	}
}
