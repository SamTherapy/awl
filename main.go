// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	cli "git.froth.zone/sam/awl/cmd"
	"git.froth.zone/sam/awl/pkg/query"
	"git.froth.zone/sam/awl/pkg/util"
)

var version = "DEV"

func main() {
	if opts, code, err := run(os.Args); err != nil {
		// TODO: Make not ew
		if errors.Is(err, cli.ErrNotError) || strings.Contains(err.Error(), "help requested") {
			os.Exit(0)
		} else {
			opts.Logger.Error(err)
			os.Exit(code)
		}
	}
}

func run(args []string) (opts *util.Options, code int, err error) {
	opts, err = cli.ParseCLI(args, version)
	if err != nil {
		return opts, 1, fmt.Errorf("parse: %w", err)
	}

	var resp util.Response

	// Retry queries if a query fails
	for i := 0; i <= opts.Request.Retries; i++ {
		resp, err = query.CreateQuery(opts)
		if err == nil {
			break
		} else if i != opts.Request.Retries {
			opts.Logger.Warn("Retrying request, error:", err)
		}
	}

	// Query failed, make it fail
	if err != nil {
		return opts, 9, fmt.Errorf("query: %w", err)
	}

	var str string
	if opts.JSON || opts.XML || opts.YAML {
		str, err = query.PrintSpecial(resp, opts)
		if err != nil {
			return opts, 10, fmt.Errorf("format print: %w", err)
		}
	} else {
		str, err = query.ToString(resp, opts)
		if err != nil {
			return opts, 15, fmt.Errorf("standard print: %w", err)
		}
	}

	fmt.Println(str)

	return opts, 0, nil
}
