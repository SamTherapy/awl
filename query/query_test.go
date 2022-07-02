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
	testCase := util.Answers{Server: "dns9.quad9.net/dns-query", Request: dns.TypeA, Name: "git.froth.zone"}
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
	in, testCase.RTT, err = ResolveHTTPS(msg, testCase.Server)
	assert.Nil(t, err)
	assert.NotNil(t, in)

}

func Test2ResolveHTTPS(t *testing.T) {
	var err error
	testCase := util.Answers{Server: "dns9.quad9.net/dns-query", Request: dns.TypeA, Name: "git.froth.zone"}

	msg := new(dns.Msg)
	msg.SetQuestion(testCase.Name, testCase.Request)
	msg = msg.SetQuestion(testCase.Name, testCase.Request)
	var in *dns.Msg
	in, testCase.RTT, err = ResolveHTTPS(msg, testCase.Server)
	assert.NotNil(t, err)
	assert.Nil(t, in)

}
func Test3ResolveHTTPS(t *testing.T) {
	var err error
	testCase := util.Answers{Server: "dns9..quad9.net/dns-query", Request: dns.TypeA, Name: "git.froth.zone."}
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
	in, testCase.RTT, err = ResolveHTTPS(msg, testCase.Server)
	assert.NotNil(t, err)
	assert.Nil(t, in)

}

func TestQuic(t *testing.T) {
	var err error
	testCase := util.Answers{Server: "dns.adguard.com", Request: dns.TypeA, Name: "git.froth.zone"}
	testCase2 := util.Answers{Server: "dns.adguard.com", Request: dns.TypeA, Name: "git.froth.zone"}
	var testCases []util.Answers
	testCases = append(testCases, testCase)
	testCases = append(testCases, testCase2)

	for i := range testCases {
		switch i {
		case 0:
			port := 853
			testCases[i].Server = net.JoinHostPort(testCases[i].Server, strconv.Itoa(port))
			fmt.Println(testCases[i].Server)
			// if the domain is not canonical, make it canonical
			if !strings.HasSuffix(testCase.Name, ".") {
				testCases[i].Name = fmt.Sprintf("%s.", testCases[i].Name)
			}
			msg := new(dns.Msg)
			msg.SetQuestion(testCase.Name, testCase.Request)
			msg = msg.SetQuestion(testCase.Name, testCase.Request)
			var in *dns.Msg
			in, testCase.RTT, err = ResolveQUIC(msg, testCase.Server)
			assert.NotNil(t, err)
			assert.Nil(t, in)
		case 1:
			port := 853
			testCases[i].Server = net.JoinHostPort(testCases[i].Server, strconv.Itoa(port))
			// if the domain is not canonical, make it canonical
			if !strings.HasSuffix(testCase.Name, ".") {
				testCases[i].Name = fmt.Sprintf("%s.", testCases[i].Name)
			}
			msg := new(dns.Msg)
			msg.SetQuestion(testCases[i].Name, testCases[i].Request)
			msg = msg.SetQuestion(testCases[i].Name, testCases[i].Request)
			var in *dns.Msg
			in, testCase.RTT, err = ResolveQUIC(msg, testCases[i].Server)
			assert.Nil(t, err)
			assert.NotNil(t, in)
		}

	}

}
