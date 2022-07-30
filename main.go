// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"
	"git.froth.zone/sam/awl/query"
)

var version = "DEV"

func main() {
	opts, err := cli.ParseCLI(version)
	if err != nil {
		// TODO: Make not ew
		if errors.Is(err, cli.ErrNotError) || strings.Contains(err.Error(), "help requested") {
			os.Exit(0)
		}
		opts.Logger.Error(err)
		os.Exit(1)
	}

	var resp helpers.Response

	// Retry queries if a query fails
	for i := 0; i < opts.Request.Retries; i++ {
		resp, err = query.CreateQuery(opts)
		if err == nil {
			break
		} else {
			opts.Logger.Warn("Retrying request, error:", err)
		}
	}

	// Query failed, make it fail
	if err != nil {
		opts.Logger.Error(err)
		os.Exit(9)
	}

	var str string
	if opts.JSON || opts.XML || opts.YAML {
		str, err = query.PrintSpecial(resp.DNS, opts)
		if err != nil {
			opts.Logger.Error("Special print:", err)
			os.Exit(10)
		}
	} else {
		str = query.ToString(resp, opts)
	}

	fmt.Println(str)
}
