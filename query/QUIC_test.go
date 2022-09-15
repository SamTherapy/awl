// SPDX-License-Identifier: BSD-3-Clause

package query_test

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"gotest.tools/v3/assert"
)

func TestQuic(t *testing.T) {
	t.Parallel()

	opts := util.Options{
		QUIC:    true,
		Logger:  util.InitLogger(0),
		Request: util.Request{Server: "dns.adguard.com", Port: 853},
	}
	testCase := util.Request{Server: "dns.//./,,adguard.com", Type: dns.TypeA, Name: "git.froth.zone"}
	testCase2 := util.Request{Server: "dns.adguard.com", Type: dns.TypeA, Name: "git.froth.zone"}

	var testCases []util.Request

	testCases = append(testCases, testCase)
	testCases = append(testCases, testCase2)

	for i := range testCases {
		switch i {
		case 0:
			resolver, err := query.LoadResolver(opts)
			assert.NilError(t, err)
			// if the domain is not canonical, make it canonical
			if !strings.HasSuffix(testCase.Name, ".") {
				testCases[i].Name = fmt.Sprintf("%s.", testCases[i].Name)
			}

			msg := new(dns.Msg)
			msg.SetQuestion(testCase.Name, testCase.Type)
			// msg = msg.SetQuestion(testCase.Name, testCase.Type)
			res, err := resolver.LookUp(msg)

			assert.ErrorContains(t, err, "fully qualified")
			assert.Equal(t, res, util.Response{})
		case 1:
			resolver, err := query.LoadResolver(opts)
			assert.NilError(t, err)

			testCase2.Server = net.JoinHostPort(testCase2.Server, strconv.Itoa(opts.Request.Port))

			// if the domain is not canonical, make it canonical
			if !strings.HasSuffix(testCase2.Name, ".") {
				testCase2.Name = fmt.Sprintf("%s.", testCase2.Name)
			}

			msg := new(dns.Msg)
			msg.SetQuestion(testCase2.Name, testCase2.Type)

			res, err := resolver.LookUp(msg)

			assert.NilError(t, err)
			assert.Assert(t, res != util.Response{})
		}
	}
}

func TestInvalidQuic(t *testing.T) {
	t.Parallel()

	opts := util.Options{
		QUIC:    true,
		Logger:  util.InitLogger(0),
		Request: util.Request{Server: "example.com", Port: 853, Type: dns.TypeA, Name: "git.froth.zone", Timeout: 10 * time.Millisecond},
	}
	resolver, err := query.LoadResolver(opts)
	assert.NilError(t, err)

	msg := new(dns.Msg)
	msg.SetQuestion(opts.Request.Name, opts.Request.Type)
	res, err := resolver.LookUp(msg)
	assert.ErrorContains(t, err, "timeout")
	assert.Equal(t, res, util.Response{})
}
