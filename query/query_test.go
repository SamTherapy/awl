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
	opts := cli.Options{
		Logger: util.InitLogger(0),
		Port:   53,
		QR:     false,
		Z:      true,
		RD:     false,

		Request: helpers.Request{
			Server: "8.8.4.4",
			Type:   dns.TypeA,
			Name:   "example.com.",
		},
		EDNS: cli.EDNS{
			EnableEDNS: true,
			DNSSEC:     true,
			Cookie:     true,
			Expire:     true,
			KeepOpen:   true,
			Nsid:       true,
		},
	}
	res, err := query.CreateQuery(opts)
	assert.NilError(t, err)
	assert.Assert(t, res != helpers.Response{})
}

func TestCreateQr(t *testing.T) {
	opts := cli.Options{
		Logger:    util.InitLogger(0),
		Port:      53,
		QR:        false,
		Z:         true,
		RD:        false,
		ShowQuery: true,

		Request: helpers.Request{
			Server: "8.8.4.4",
			Type:   dns.TypeA,
			Name:   "example.com.",
		},
		EDNS: cli.EDNS{
			EnableEDNS: false,
			DNSSEC:     true,
			Cookie:     false,
			Expire:     false,
			KeepOpen:   false,
			Nsid:       false,
		},
		Display: cli.Displays{
			Comments:       true,
			Question:       true,
			Answer:         true,
			Authority:      true,
			Additional:     true,
			Statistics:     true,
			UcodeTranslate: true,
		},
	}
	res, err := query.CreateQuery(opts)
	assert.NilError(t, err)
	assert.Assert(t, res != helpers.Response{})
}

func TestCreateQr2(t *testing.T) {
	opts := cli.Options{
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
		EDNS: cli.EDNS{
			EnableEDNS: false,
			DNSSEC:     false,
			Cookie:     false,
			Expire:     false,
			KeepOpen:   false,
			Nsid:       false,
		},
	}
	res, err := query.CreateQuery(opts)
	assert.NilError(t, err)
	assert.Assert(t, res != helpers.Response{})
}
