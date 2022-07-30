// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	main()
	os.Args = []string{"awl", "+yaml", "@dns.froth.zone"}
	main()
}
