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
		DNSSEC: true,
		Request: helpers.Request{
			Server: "8.8.4.4",
			Type:   dns.TypeA,
			Name:   "example.com.",
		},
	}
	res, err := query.CreateQuery(opts)
	assert.NilError(t, err)
	assert.Assert(t, res != helpers.Response{})
}
