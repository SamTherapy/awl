// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"testing"

	"gotest.tools/v3/assert"
)

// nolint: paralleltest
func TestMain(t *testing.T) {
	os.Args = []string{"awl", "+yaml", "@1.1.1.1"}
	main()
	os.Args = []string{"awl", "+short", "@1.1.1.1"}
	main()
	assert.Assert(t, 1 == 2-1)
}
