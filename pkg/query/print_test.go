// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"testing"

	"git.froth.zone/sam/awl/pkg/query"
	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestRealPrint(t *testing.T) {
	t.Parallel()

	opts := []util.Options{
		{
			Logger: util.InitLogger(0),

			TCP: true,

			HeaderFlags: util.HeaderFlags{
				RD: true,
			},

			JSON: true,
			Display: util.Display{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: true,
				TTL:            true,
				HumanTTL:       true,
				ShowQuery:      true,
			},
			Request: util.Request{
				Server:  "a.gtld-servers.net",
				Port:    53,
				Type:    dns.StringToType["NS"],
				Class:   1,
				Name:    "google.com.",
				Retries: 3,
			},
			EDNS: util.EDNS{
				EnableEDNS: false,
			},
		},
		{
			Logger: util.InitLogger(0),

			TCP: true,
			HeaderFlags: util.HeaderFlags{
				RD: true,
			},
			Verbosity: 0,

			Short:    true,
			Identify: true,
			YAML:     false,
			Display: util.Display{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: true,
				TTL:            true,
				ShowQuery:      true,
			},
			Request: util.Request{
				Server:  "ns1.google.com",
				Port:    53,
				Type:    dns.StringToType["NS"],
				Class:   1,
				Name:    "google.com.",
				Retries: 3,
			},
			EDNS: util.EDNS{
				EnableEDNS: false,
			},
		},
		{
			Logger: util.InitLogger(0),
			HTTPS:  true,
			HeaderFlags: util.HeaderFlags{
				RD: true,
			},
			Identify: true,
			XML:      true,
			Display: util.Display{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: true,
				TTL:            true,
				HumanTTL:       true,
				ShowQuery:      true,
			},
			Request: util.Request{
				Server:  "https://dns.froth.zone/dns-query",
				Port:    443,
				Type:    dns.StringToType["NS"],
				Class:   1,
				Name:    "freecumextremist.com.",
				Retries: 3,
			},
			EDNS: util.EDNS{
				EnableEDNS: false,
				DNSSEC:     true,
			},
		},
		{
			Logger: util.InitLogger(0),
			TLS:    true,
			HeaderFlags: util.HeaderFlags{
				RD: true,
			},
			Verbosity: 0,
			Display: util.Display{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: true,
				TTL:            false,
				ShowQuery:      true,
			},
			Request: util.Request{
				Server:  "dns.google",
				Port:    853,
				Type:    dns.StringToType["NS"],
				Class:   1,
				Name:    "freecumextremist.com.",
				Retries: 3,
			},
		},
		{
			Logger: util.InitLogger(0),
			TCP:    true,

			HeaderFlags: util.HeaderFlags{
				AA: true,
				RD: true,
			},
			Verbosity: 0,

			YAML: true,
			Display: util.Display{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: false,
				TTL:            true,
				ShowQuery:      true,
			},
			Request: util.Request{
				Server:  "rin.froth.zone",
				Port:    53,
				Type:    dns.StringToType["A"],
				Class:   1,
				Name:    "froth.zone.",
				Retries: 3,
			},
			EDNS: util.EDNS{
				EnableEDNS: true,
				Cookie:     true,
				Padding:    true,
			},
		},
	}

	for _, test := range opts {
		test := test

		t.Run("", func(t *testing.T) {
			t.Parallel()
			resp, err := query.CreateQuery(test)
			assert.NilError(t, err)

			if test.JSON || test.XML || test.YAML {
				str := ""
				str, err = query.PrintSpecial(resp, test)
				assert.NilError(t, err)
				assert.Assert(t, str != "")
			}
			str, err := query.ToString(resp, test)
			assert.NilError(t, err)
			assert.Assert(t, str != "")
		})
	}
}

func TestBadFormat(t *testing.T) {
	t.Parallel()

	_, err := query.PrintSpecial(util.Response{DNS: new(dns.Msg)}, util.Options{})
	assert.ErrorContains(t, err, "never happen")
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	str, err := query.ToString(util.Response{}, util.Options{})

	assert.Error(t, err, "no message")
	assert.Assert(t, str == "<nil> MsgHdr")
}
