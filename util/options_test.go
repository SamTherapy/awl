// SPDX-License-Identifier: BSD-3-Clause

package util_test

import (
	"testing"

	"git.froth.zone/sam/awl/util"
	"gotest.tools/v3/assert"
)

func TestSubnet(t *testing.T) {
	t.Parallel()

	subnet := []string{
		"0.0.0.0/0",
		"::0/0",
		"0",
		"127.0.0.1/32",
	}

	for _, test := range subnet {
		test := test
		t.Run(test, func(t *testing.T) {
			t.Parallel()
			err := util.ParseSubnet(test, new(util.Options))
			assert.NilError(t, err)
		})
	}
}

func TestInvalidSub(t *testing.T) {
	t.Parallel()

	err := util.ParseSubnet("1", new(util.Options))
	assert.ErrorContains(t, err, "invalid CIDR address")
}
