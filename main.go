// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	cli "git.froth.zone/sam/awl/cmd"
	"git.froth.zone/sam/awl/pkg/query"
	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
)

var version = "DEV"

func main() {
	if opts, code, err := run(os.Args); err != nil {
		// TODO: Make not ew
		if errors.Is(err, util.ErrNotError) || strings.Contains(err.Error(), "help requested") {
			os.Exit(0)
		} else {
			opts.Logger.Error(err)
			os.Exit(code)
		}
	}
}

func run(args []string) (opts *util.Options, code int, err error) {
	//nolint:gosec //Secure source not needed
	r := rand.New(rand.NewSource(time.Now().Unix()))

	opts, err = cli.ParseCLI(args, version)
	if err != nil {
		return opts, 1, fmt.Errorf("parse: %w", err)
	}

	var (
		resp          util.Response
		keepTracing   bool
		tempDomain    string
		tempQueryType uint16
	)

	for ok := true; ok; ok = keepTracing {
		if keepTracing {
			opts.Request.Name = tempDomain
			opts.Request.Type = tempQueryType
		} else {
			tempDomain = opts.Request.Name
			tempQueryType = opts.Request.Type

			// Override the query because it needs to be done
			opts.Request.Name = "."
			opts.Request.Type = dns.TypeNS
		}
		// Retry queries if a query fails
		for i := 0; i <= opts.Request.Retries; i++ {
			resp, err = query.CreateQuery(opts)
			if err == nil {
				keepTracing = opts.Trace && (!resp.DNS.Authoritative || (opts.Request.Name == "." && tempDomain != "."))

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

		if keepTracing {
			var records []dns.RR

			if opts.Request.Name == "." {
				records = resp.DNS.Answer
			} else {
				records = resp.DNS.Ns
			}

			want := func(rr dns.RR) bool {
				temp := strings.Split(rr.String(), "\t")

				return temp[len(temp)-2] == "NS"
			}

			i := 0

			for _, x := range records {
				if want(x) {
					records[i] = x
					i++
				}
			}

			records = records[:i]
			randomRR := records[r.Intn(len(records))]

			v := strings.Split(randomRR.String(), "\t")
			opts.Request.Server = strings.TrimSuffix(v[len(v)-1], ".")

			opts.TLS = false
			opts.HTTPS = false
			opts.QUIC = false

			opts.RD = false
			opts.Request.Port = 53
		}
	}

	return opts, 0, nil
}
