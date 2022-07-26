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

func TestResolve(t *testing.T) {
	opts := cli.Options{
		Logger: util.InitLogger(0),
		Port:   53,
		Request: helpers.Request{
			Server: "8.8.4.4",
			Type:   dns.TypeA,
			Name:   "example.com.",
		},
	}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)
	msg := new(dns.Msg)
	msg.SetQuestion(opts.Request.Name, opts.Request.Type)
	res, err := resolver.LookUp(msg)
	assert.NilError(t, err)
	assert.Assert(t, res != helpers.Response{})
}

func TestTruncate(t *testing.T) {
	opts := cli.Options{
		Logger: util.InitLogger(0),
		IPv4:   true,
		Port:   5301,
		Request: helpers.Request{
			Server: "madns.binarystar.systems",
			Type:   dns.TypeTXT,
			Name:   "limit.txt.example.",
		},
	}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)
	msg := new(dns.Msg)
	msg.SetQuestion(opts.Request.Name, opts.Request.Type)
	res, err := resolver.LookUp(msg)
	assert.NilError(t, err)
	assert.Assert(t, res != helpers.Response{})
}

func TestResolveAgain(t *testing.T) {
	tests := []struct {
		opt cli.Options
	}{
		{
			cli.Options{
				Logger: util.InitLogger(0),
				TCP:    true,
				Port:   53,
				Request: helpers.Request{
					Server: "8.8.4.4",
					Type:   dns.TypeA,
					Name:   "example.com.",
				},
			},
		},
		{
			cli.Options{
				Logger: util.InitLogger(0),
				Port:   53,
				Request: helpers.Request{
					Server: "8.8.4.4",
					Type:   dns.TypeAAAA,
					Name:   "example.com.",
				},
			},
		},
	}
	for _, test := range tests {
		res, err := query.CreateQuery(test.opt)
		assert.NilError(t, err)
		assert.Assert(t, res != helpers.Response{})
	}

}
