// SPDX-License-Identifier: BSD-3-Clause

package util_test

import (
	"testing"

	"dns.froth.zone/awl/pkg/util"
	"gotest.tools/v3/assert"
)

func TestSubnet(t *testing.T) {
	t.Parallel()

	subnet := []string{
		"0.0.0.0/0",
		"::0/0",
		"0",
		"127.0.0.1/32",
		"Invalid",
	}

	for _, test := range subnet {
		test := test

		t.Run(test, func(t *testing.T) {
			t.Parallel()
			err := util.ParseSubnet(test, new(util.Options))
			if err != nil {
				assert.ErrorContains(t, err, "invalid CIDR address")
			} else {
				assert.NilError(t, err)
			}
		})
	}
}
