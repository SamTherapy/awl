// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"testing"

	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestResolveHTTPS(t *testing.T) {
	t.Parallel()

	var err error

	opts := util.Options{
		HTTPS:  true,
		Logger: util.InitLogger(0),
		Request: util.Request{
			Server: "https://dns9.quad9.net/dns-query",
			Type:   dns.TypeA,
			Name:   "git.froth.zone.",
		},
	}
	// testCase := util.Request{Server: "https://dns9.quad9.net/dns-query", Type: dns.TypeA, Name: "git.froth.zone."}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)

	msg := new(dns.Msg)
	msg.SetQuestion(opts.Request.Name, opts.Request.Type)
	// msg = msg.SetQuestion(testCase.Name, testCase.Type)
	res, err := resolver.LookUp(msg)
	assert.NilError(t, err)
	assert.Assert(t, res != util.Response{})
}

func Test2ResolveHTTPS(t *testing.T) {
	t.Parallel()

	opts := util.Options{
		HTTPS:   true,
		Logger:  util.InitLogger(0),
		Request: util.Request{Server: "dns9.quad9.net/dns-query", Type: dns.TypeA, Name: "git.froth.zone."},
	}

	var err error

	testCase := util.Request{Type: dns.TypeA, Name: "git.froth.zone"}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)

	msg := new(dns.Msg)
	msg.SetQuestion(testCase.Name, testCase.Type)
	// msg = msg.SetQuestion(testCase.Name, testCase.Type)
	res, err := resolver.LookUp(msg)
	assert.ErrorContains(t, err, "fully qualified")
	assert.Equal(t, res, util.Response{})
}

func Test3ResolveHTTPS(t *testing.T) {
	t.Parallel()

	opts := util.Options{
		HTTPS:   true,
		Logger:  util.InitLogger(0),
		Request: util.Request{Server: "dns9..quad9.net/dns-query", Type: dns.TypeA, Name: "git.froth.zone."},
	}

	var err error

	// testCase :=
	// if the domain is not canonical, make it canonical
	// if !strings.HasSuffix(testCase.Name, ".") {
	// 	testCase.Name = fmt.Sprintf("%s.", testCase.Name)
	// }

	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)

	msg := new(dns.Msg)
	// msg.SetQuestion(testCase.Name, testCase.Type)
	// msg = msg.SetQuestion(testCase.Name, testCase.Type)
	res, err := resolver.LookUp(msg)
	assert.ErrorContains(t, err, "doh: HTTP request")
	assert.Equal(t, res, util.Response{})
}

func Test404ResolveHTTPS(t *testing.T) {
	t.Parallel()

	var err error

	opts := util.Options{
		HTTPS:  true,
		Logger: util.InitLogger(0),
		Request: util.Request{
			Server: "https://dns9.quad9.net/dns",
			Type:   dns.TypeA,
			Name:   "git.froth.zone.",
		},
	}
	// testCase := util.Request{Server: "https://dns9.quad9.net/dns-query", Type: dns.TypeA, Name: "git.froth.zone."}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)

	msg := new(dns.Msg)
	msg.SetQuestion(opts.Request.Name, opts.Request.Type)
	// msg = msg.SetQuestion(testCase.Name, testCase.Type)
	res, err := resolver.LookUp(msg)
	assert.ErrorContains(t, err, "404")
	assert.Equal(t, res, util.Response{})
}
