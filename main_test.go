// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"
	"testing"
)

func TestMain(t *testing.T) {
	os.Args = []string{"awl", "+yaml", "@1.1.1.1"}
	main()
}
