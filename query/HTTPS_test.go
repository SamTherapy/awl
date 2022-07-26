// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"fmt"
	"strings"
	"testing"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"
	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"

	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestResolveHTTPS(t *testing.T) {
	t.Parallel()
	var err error
	opts := cli.Options{
		HTTPS:  true,
		Logger: util.InitLogger(0),
		Request: helpers.Request{
			Server: "https://dns9.quad9.net/dns-query",
			Type:   dns.TypeA,
			Name:   "git.froth.zone.",
		},
	}
	// testCase := helpers.Request{Server: "https://dns9.quad9.net/dns-query", Type: dns.TypeA, Name: "git.froth.zone."}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)

	msg := new(dns.Msg)
	msg.SetQuestion(opts.Request.Name, opts.Request.Type)
	// msg = msg.SetQuestion(testCase.Name, testCase.Type)
	res, err := resolver.LookUp(msg)
	assert.NilError(t, err)
	assert.Assert(t, res != helpers.Response{})
}

func Test2ResolveHTTPS(t *testing.T) {
	t.Parallel()
	opts := cli.Options{
		HTTPS:  true,
		Logger: util.InitLogger(0),
	}
	var err error
	testCase := helpers.Request{Server: "dns9.quad9.net/dns-query", Type: dns.TypeA, Name: "git.froth.zone"}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)
	msg := new(dns.Msg)
	msg.SetQuestion(testCase.Name, testCase.Type)
	// msg = msg.SetQuestion(testCase.Name, testCase.Type)
	res, err := resolver.LookUp(msg)
	assert.ErrorContains(t, err, "fully qualified")
	assert.Equal(t, res, helpers.Response{})
}

func Test3ResolveHTTPS(t *testing.T) {
	t.Parallel()
	opts := cli.Options{
		HTTPS:  true,
		Logger: util.InitLogger(0),
	}
	var err error
	testCase := helpers.Request{Server: "dns9..quad9.net/dns-query", Type: dns.TypeA, Name: "git.froth.zone."}
	// if the domain is not canonical, make it canonical
	if !strings.HasSuffix(testCase.Name, ".") {
		testCase.Name = fmt.Sprintf("%s.", testCase.Name)
	}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)
	msg := new(dns.Msg)
	msg.SetQuestion(testCase.Name, testCase.Type)
	// msg = msg.SetQuestion(testCase.Name, testCase.Type)
	res, err := resolver.LookUp(msg)
	assert.ErrorContains(t, err, "request error")
	assert.Equal(t, res, helpers.Response{})
}

func Test404ResolveHTTPS(t *testing.T) {
	t.Parallel()
	var err error
	opts := cli.Options{
		HTTPS:  true,
		Logger: util.InitLogger(0),
		Request: helpers.Request{
			Server: "https://dns9.quad9.net/dns",
			Type:   dns.TypeA,
			Name:   "git.froth.zone.",
		},
	}
	// testCase := helpers.Request{Server: "https://dns9.quad9.net/dns-query", Type: dns.TypeA, Name: "git.froth.zone."}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)

	msg := new(dns.Msg)
	msg.SetQuestion(opts.Request.Name, opts.Request.Type)
	// msg = msg.SetQuestion(testCase.Name, testCase.Type)
	res, err := resolver.LookUp(msg)
	assert.ErrorContains(t, err, "404")
	assert.Equal(t, res, helpers.Response{})
}
