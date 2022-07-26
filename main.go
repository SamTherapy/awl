// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/cli"
	"git.froth.zone/sam/awl/internal/helpers"
	"git.froth.zone/sam/awl/query"

	"gopkg.in/yaml.v2"
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
			opts.Logger.Warn("Retrying request, error", err)
		}
	}

	// Query failed, make it fail
	if err != nil {
		opts.Logger.Error(err)
		os.Exit(9)
	}
	switch {
	case opts.JSON:
		opts.Logger.Info("Printing as JSON")
		json, err := json.MarshalIndent(resp.DNS, "", "  ")
		if err != nil {
			opts.Logger.Error(err)
			os.Exit(10)
		}
		fmt.Println(string(json))
	case opts.XML:
		opts.Logger.Info("Printing as XML")
		xml, err := xml.MarshalIndent(resp.DNS, "", "  ")
		if err != nil {
			opts.Logger.Error(err)
			os.Exit(10)
		}
		fmt.Println(string(xml))
	case opts.YAML:
		opts.Logger.Info("Printing as YAML")
		yaml, err := yaml.Marshal(resp.DNS)
		if err != nil {
			opts.Logger.Error(err)
			os.Exit(10)
		}
		fmt.Println(string(yaml))
	default:
		if !opts.Short {
			// Print everything

			if !opts.Display.Question {
				resp.DNS.Question = nil
				opts.Logger.Info("Disabled question display")
			}
			if !opts.Display.Answer {
				resp.DNS.Answer = nil
				opts.Logger.Info("Disabled answer display")
			}
			if !opts.Display.Authority {
				resp.DNS.Ns = nil
				opts.Logger.Info("Disabled authority display")
			}
			if !opts.Display.Additional {
				resp.DNS.Extra = nil
				opts.Logger.Info("Disabled additional display")
			}

			fmt.Println(resp.DNS)

			if opts.Display.Statistics {
				fmt.Println(";; Query time:", resp.RTT)

				// Add extra information to server string
				var extra string
				switch {
				case opts.TCP:
					extra = ":" + strconv.Itoa(opts.Port) + " (TCP)"
				case opts.TLS:
					extra = ":" + strconv.Itoa(opts.Port) + " (TLS)"
				case opts.HTTPS, opts.DNSCrypt:
					extra = ""
				case opts.QUIC:
					extra = ":" + strconv.Itoa(opts.Port) + " (QUIC)"
				default:
					extra = ":" + strconv.Itoa(opts.Port) + " (UDP)"
				}

				fmt.Println(";; SERVER:", opts.Request.Server+extra)
				fmt.Println(";; WHEN:", time.Now().Format(time.RFC1123Z))
				fmt.Println(";; MSG SIZE  rcvd:", resp.DNS.Len())
			}

		} else {
			// Print just the responses, nothing else
			for _, res := range resp.DNS.Answer {
				temp := strings.Split(res.String(), "\t")
				fmt.Println(temp[len(temp)-1])
			}
		}
	}
}
