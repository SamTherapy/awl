// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"testing"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"
	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestCreateQ(t *testing.T) {
	t.Parallel()
	in := []cli.Options{
		{
			Logger:    util.InitLogger(0),
			Port:      53,
			QR:        false,
			Z:         true,
			RD:        false,
			ShowQuery: true,
			YAML:      true,

			Request: helpers.Request{
				Server: "8.8.4.4",
				Type:   dns.TypeA,
				Name:   "example.com.",
			},
			Display: cli.Displays{
				Comments:       true,
				Question:       true,
				Opt:            true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: false,
			},
			EDNS: cli.EDNS{
				EnableEDNS: true,
				DNSSEC:     true,
				Cookie:     true,
				Expire:     true,
				KeepOpen:   true,
				Nsid:       true,
			},
		},
		{
			Logger:    util.InitLogger(0),
			Port:      53,
			QR:        false,
			Z:         true,
			RD:        false,
			ShowQuery: true,
			XML:       true,

			Request: helpers.Request{
				Server: "8.8.4.4",
				Type:   dns.TypeA,
				Name:   "example.com.",
			},
			Display: cli.Displays{
				Comments:       true,
				Question:       true,
				Opt:            true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: true,
			},
			EDNS: cli.EDNS{
				EnableEDNS: false,
				DNSSEC:     false,
				Cookie:     false,
				Expire:     false,
				KeepOpen:   false,
				Nsid:       false,
			},
		},
	}
	for _, opt := range in {
		opt := opt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			res, err := query.CreateQuery(opt)
			assert.NilError(t, err)
			assert.Assert(t, res != helpers.Response{})
			str, err := query.PrintSpecial(res.DNS, opt)
			assert.NilError(t, err)
			assert.Assert(t, str != "")
			str = query.ToString(res, opt)
			assert.Assert(t, str != "")
		})
	}
}
