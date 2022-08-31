// SPDX-License-Identifier: BSD-3-Clause
//go:build plan9

package conf_test

import (
	"testing"

	"git.froth.zone/sam/awl/conf"
	"gotest.tools/v3/assert"
)

func TestGetPlan9Config(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "plan9" {
		t.Skip("Not running Plan 9, skipping")
	}

	ndbs := []struct {
		in   string
		want string
	}{
		{`ip=192.168.122.45 ipmask=255.255.255.0 ipgw=192.168.122.1
	sys=chog9
	dns=192.168.122.1`, "192.168.122.1"},
		{`ipnet=murray-hill ip=135.104.0.0 ipmask=255.255.0.0
	dns=135.104.10.1
	ntp=ntp.cs.bell-labs.com
	ipnet=plan9 ip=135.104.9.0 ipmask=255.255.255.0
	ntp=oncore.cs.bell-labs.com
	smtp=smtp1.cs.bell-labs.com
	ip=135.104.9.6 sys=anna dom=anna.cs.bell-labs.com
	smtp=smtp2.cs.bell-labs.com`, "135.104.10.1"},
	}

	for _, ndb := range ndbs {
		// Go is a little quirky
		ndb := ndb
		t.Run(ndb.want, func(t *testing.T) {
			t.Parallel()
			act, err := conf.GetPlan9Config(ndb.in)
			assert.NilError(t, err)
			assert.Equal(t, ndb.want, act.Servers[0])
		})
	}

	invalid := `sys = spindle
	dom=spindle.research.bell-labs.com
	bootf=/mips/9powerboot
	ip=135.104.117.32 ether=080069020677
	proto=il`

	act, err := conf.GetPlan9Config(invalid)
	assert.ErrorContains(t, err, "no DNS servers found")
	assert.Assert(t, act == nil)
}
