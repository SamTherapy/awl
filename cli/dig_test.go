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
		"dnssec", "nodnssec",
		"tcp", "vc", "notcp", "novc",
		"ignore", "noignore",
		"tls", "notls",
		"dnscrypt", "nodnscrypt",
		"https", "nohttps",
		"quic", "noquic",
		"short", "noshort",
		"json", "nojson",
		"xml", "noxml",
		"yaml", "noyaml",
		"question", "noquestion",
		"answer", "noanswer",
		"authority", "noauthority",
		"additional", "noadditional",
		"stats", "nostats",
		"all", "noall",
		"invalid",
	}
	for _, tc := range seeds {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		opts := new(cli.Options)
		opts.Logger = util.InitLogger(0)
		err := cli.ParseDig(orig, opts)
		if err != nil {
			assert.ErrorContains(t, err, "unknown flag")
		}
	})
}
