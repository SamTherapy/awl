package query

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"

	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func TestResolveHTTPS(t *testing.T) {
	var err error
	opts := Options{
		HTTPS:  true,
		Logger: util.InitLogger(false),
	}
	testCase := Answers{Server: "dns9.quad9.net/dns-query", Request: dns.TypeA, Name: "git.froth.zone"}
	resolver, err := LoadResolver(testCase.Server, opts)

	if !strings.HasPrefix(testCase.Server, "https://") {
		testCase.Server = "https://" + testCase.Server
	}
	// if the domain is not canonical, make it canonical
	if !strings.HasSuffix(testCase.Name, ".") {
		testCase.Name = fmt.Sprintf("%s.", testCase.Name)
	}

	msg := new(dns.Msg)
	msg.SetQuestion(testCase.Name, testCase.Request)
	msg = msg.SetQuestion(testCase.Name, testCase.Request)
	var in *dns.Msg
	in, testCase.RTT, err = resolver.LookUp(msg)
	assert.Nil(t, err)
	assert.NotNil(t, in)

}

func Test2ResolveHTTPS(t *testing.T) {
	opts := Options{
		HTTPS:  true,
		Logger: util.InitLogger(false),
	}
	var err error
	testCase := Answers{Server: "dns9.quad9.net/dns-query", Request: dns.TypeA, Name: "git.froth.zone"}
	resolver, err := LoadResolver(testCase.Server, opts)
	msg := new(dns.Msg)
	msg.SetQuestion(testCase.Name, testCase.Request)
	msg = msg.SetQuestion(testCase.Name, testCase.Request)
	var in *dns.Msg
	in, testCase.RTT, err = resolver.LookUp(msg)
	assert.NotNil(t, err)
	assert.Nil(t, in)

}
func Test3ResolveHTTPS(t *testing.T) {
	opts := Options{
		HTTPS:  true,
		Logger: util.InitLogger(false),
	}
	var err error
	testCase := Answers{Server: "dns9..quad9.net/dns-query", Request: dns.TypeA, Name: "git.froth.zone."}
	if !strings.HasPrefix(testCase.Server, "https://") {
		testCase.Server = "https://" + testCase.Server
	}
	// if the domain is not canonical, make it canonical
	if !strings.HasSuffix(testCase.Name, ".") {
		testCase.Name = fmt.Sprintf("%s.", testCase.Name)
	}
	resolver, err := LoadResolver(testCase.Server, opts)
	msg := new(dns.Msg)
	msg.SetQuestion(testCase.Name, testCase.Request)
	msg = msg.SetQuestion(testCase.Name, testCase.Request)
	var in *dns.Msg
	in, testCase.RTT, err = resolver.LookUp(msg)
	assert.NotNil(t, err)
	assert.Nil(t, in)

}

func TestQuic(t *testing.T) {
	opts := Options{
		QUIC:    true,
		Logger:  util.InitLogger(false),
		Port:    853,
		Answers: Answers{Server: "dns.adguard.com"},
	}
	testCase := Answers{Server: "dns.//./,,adguard.com", Request: dns.TypeA, Name: "git.froth.zone"}
	testCase2 := Answers{Server: "dns.adguard.com", Request: dns.TypeA, Name: "git.froth.zone"}
	var testCases []Answers
	testCases = append(testCases, testCase)
	testCases = append(testCases, testCase2)
	for i := range testCases {
		switch i {
		case 0:
			resolver, err := LoadResolver(testCases[i].Server, opts)
			// if the domain is not canonical, make it canonical
			if !strings.HasSuffix(testCase.Name, ".") {
				testCases[i].Name = fmt.Sprintf("%s.", testCases[i].Name)
			}
			msg := new(dns.Msg)
			msg.SetQuestion(testCase.Name, testCase.Request)
			msg = msg.SetQuestion(testCase.Name, testCase.Request)
			var in *dns.Msg
			in, testCase.RTT, err = resolver.LookUp(msg)
			assert.NotNil(t, err)
			assert.Nil(t, in)
		case 1:
			resolver, err := LoadResolver(testCase2.Server, opts)
			testCase2.Server = net.JoinHostPort(testCase2.Server, strconv.Itoa(opts.Port))
			// if the domain is not canonical, make it canonical
			if !strings.HasSuffix(testCase2.Name, ".") {
				testCase2.Name = fmt.Sprintf("%s.", testCase2.Name)
			}
			msg := new(dns.Msg)
			msg.SetQuestion(testCase2.Name, testCase2.Request)
			msg = msg.SetQuestion(testCase2.Name, testCase2.Request)
			var in *dns.Msg
			in, testCase.RTT, err = resolver.LookUp(msg)
			assert.Nil(t, err)
			assert.NotNil(t, in)
		}

	}

}
