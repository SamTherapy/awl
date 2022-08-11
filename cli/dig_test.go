// SPDX-License-Identifier: BSD-3-Clause

package cli_test

import (
	"testing"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/util"
	"gotest.tools/v3/assert"
)

func FuzzDig(f *testing.F) {
	f.Log("ParseDig Fuzzing")

	seeds := []string{
		"aaflag", "aaonly", "noaaflag", "noaaonly",
		"adflag", "noadflag",
		"cdflag", "nocdflag",
		"qrflag", "noqrflag",
		"raflag", "noraflag",
		"rdflag", "recurse", "nordflag", "norecurse",
		"tcflag", "notcflag",
		"zflag", "nozflag",
		"qr", "noqr",
		"ttlunits", "nottlunits",
		"ttlid", "nottlid",
		"dnssec", "nodnssec",
		"edns", "edns=a", "edns=0", "noedns",
		"expire", "noexpire",
		"ednsflags", "ednsflags=\"", "ednsflags=1", "noednsflags",
		"subnet=0.0.0.0/0", "subnet=::0/0", "subnet=b", "subnet=0", "subnet",
		"cookie", "nocookeie",
		"keepopen", "keepalive", "nokeepopen", "nokeepalive",
		"nsid", "nonsid",
		"padding", "nopadding",
		"bufsize=512", "bufsize=a", "bufsize",
		"time=5", "timeout=a", "timeout",
		"retry=a", "retry=3", "retry",
		"tries=2", "tries=b", "tries",
		"tcp", "vc", "notcp", "novc",
		"ignore", "noignore",
		"tls", "notls",
		"dnscrypt", "nodnscrypt",
		"https", "nohttps",
		"quic", "noquic",
		"short", "noshort",
		"identify", "noidentify",
		"json", "nojson",
		"xml", "noxml",
		"yaml", "noyaml",
		"comments", "nocomments",
		"question", "noquestion",
		"opt", "noopt",
		"answer", "noanswer",
		"authority", "noauthority",
		"additional", "noadditional",
		"stats", "nostats",
		"all", "noall",
		"idnout", "noidnout",
		"invalid",
	}

	for _, tc := range seeds {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		opts := new(util.Options)
		opts.Logger = util.InitLogger(0)
		if err := cli.ParseDig(orig, opts); err != nil {
			assert.ErrorContains(t, err, "digflags:")
		}
	})
}
