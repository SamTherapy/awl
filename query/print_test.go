// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"testing"

	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestRealPrint(t *testing.T) {
	t.Parallel()

	opts := []util.Options{
		{
			Logger:    util.InitLogger(0),
			Port:      53,
			TCP:       true,
			ShowQuery: true,
			RD:        true,
			ShowTTL:   true,
			HumanTTL:  true,
			JSON:      true,
			Display: util.Displays{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: true,
			},
			Request: util.Request{
				Server: "a.gtld-servers.net",
				Type:   dns.StringToType["NS"],
				Class:  1,
				Name:   "google.com.",
			},
			EDNS: util.EDNS{
				EnableEDNS: false,
			},
		},
		{
			Logger:    util.InitLogger(0),
			Port:      53,
			TCP:       true,
			ShowQuery: true,
			RD:        true,
			Verbosity: 0,
			ShowTTL:   true,
			Short:     true,
			Identify:  true,
			YAML:      false,
			Display: util.Displays{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: true,
			},
			Request: util.Request{
				Server:  "ns1.google.com",
				Type:    dns.StringToType["NS"],
				Class:   1,
				Name:    "google.com.",
				Timeout: 0,
				Retries: 0,
			},
			EDNS: util.EDNS{
				EnableEDNS: false,
			},
		},
		{
			Logger:    util.InitLogger(0),
			Port:      53,
			HTTPS:     true,
			ShowQuery: true,
			RD:        true,
			ShowTTL:   true,
			HumanTTL:  true,
			Identify:  true,
			XML:       true,
			Display: util.Displays{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: false,
			},
			Request: util.Request{
				Server:  "https://dns.froth.zone/dns-query",
				Type:    dns.StringToType["NS"],
				Class:   1,
				Name:    "freecumextremist.com.",
				Timeout: 0,
				Retries: 0,
			},
			EDNS: util.EDNS{
				EnableEDNS: false,
				DNSSEC:     true,
			},
		},
		{
			Logger:    util.InitLogger(0),
			Port:      853,
			TLS:       true,
			ShowQuery: true,
			RD:        true,
			Verbosity: 0,
			ShowTTL:   false,
			Display: util.Displays{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: true,
			},
			Request: util.Request{
				Server: "dns.google",
				Type:   dns.StringToType["NS"],
				Class:  1,
				Name:   "freecumextremist.com.",
			},
		},
		{
			Logger:    util.InitLogger(0),
			Port:      53,
			TCP:       true,
			ShowQuery: true,
			AA:        true,
			RD:        true,
			Verbosity: 0,
			ShowTTL:   true,
			YAML:      true,
			Display: util.Displays{
				Comments:       true,
				Question:       true,
				Answer:         true,
				Authority:      true,
				Additional:     true,
				Statistics:     true,
				UcodeTranslate: false,
			},
			Request: util.Request{
				Server:  "rin.froth.zone",
				Type:    dns.StringToType["A"],
				Class:   1,
				Name:    "froth.zone.",
				Timeout: 0,
				Retries: 0,
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
				str, err = query.PrintSpecial(resp.DNS, test)
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

	_, err := query.PrintSpecial(new(dns.Msg), util.Options{})
	assert.ErrorContains(t, err, "never happen")
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	str, err := query.ToString(util.Response{}, util.Options{})

	assert.Error(t, err, "no message")
	assert.Assert(t, str == "<nil> MsgHdr")
}
