// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
)

func doQuery(c *cli.Context) error {
	var (
		err error
	)
	// load cli flags into options struct
	Options := query.Options{
		Logger:   util.InitLogger(c.Bool("debug")),
		Port:     c.Int("port"),
		IPv4:     c.Bool("4"),
		IPv6:     c.Bool("6"),
		DNSSEC:   c.Bool("dnssec"),
		Short:    c.Bool("short"),
		TCP:      c.Bool("tcp"),
		TLS:      c.Bool("tls"),
		HTTPS:    c.Bool("https"),
		QUIC:     c.Bool("quic"),
		Truncate: c.Bool("no-truncate"),
		AA:       c.Bool("aa"),
		TC:       c.Bool("tc"),
		Z:        c.Bool("z"),
		CD:       c.Bool("cd"),
		NoRD:     c.Bool("no-rd"),
		NoRA:     c.Bool("no-ra"),
		Reverse:  c.Bool("reverse"),
		Debug:    c.Bool("debug"),
	}
	Options.Answers, err = parseArgs(c.Args().Slice(), Options)
	if err != nil {
		Options.Logger.Error("Unable to parse args")
		return err
	}
	msg := new(dns.Msg)

	if Options.Reverse {
		if dns.TypeToString[Options.Answers.Request] == "A" {
			Options.Answers.Request = dns.StringToType["PTR"]
		}
		Options.Answers.Name, err = util.ReverseDNS(Options.Answers.Name, Options.Answers.Request)
		if err != nil {
			return err
		}
	}

	// if the domain is not canonical, make it canonical
	if !strings.HasSuffix(Options.Answers.Name, ".") {
		Options.Answers.Name = fmt.Sprintf("%s.", Options.Answers.Name)
	}
	msg.SetQuestion(Options.Answers.Name, Options.Answers.Request)
	// If port is not set, set it
	if Options.Port == 0 {
		if Options.TLS || Options.QUIC {
			Options.Port = 853
		} else {
			Options.Port = 53
		}
	}
	Options.Logger.Debug("setting any message flags")
	// Make this authoritative (does this do anything?)
	if Options.AA {
		Options.Logger.Debug("making message authorative")
		msg.Authoritative = true
	}
	// Set truncated flag (why)
	if Options.TC {
		msg.Truncated = true
	}
	// Set the zero flag if requested (does nothing)
	if Options.Z {
		Options.Logger.Debug("setting to zero")
		msg.Zero = true
	}
	// Disable DNSSEC validation
	if Options.CD {
		Options.Logger.Debug("disabling DNSSEC validation")
		msg.CheckingDisabled = true
	}
	// Disable wanting recursion
	if Options.NoRD {
		Options.Logger.Debug("disabling recursion")
		msg.RecursionDesired = false
	}
	// Disable recursion being available (I don't think this does anything)
	if Options.NoRA {
		msg.RecursionAvailable = false
	}
	// Set DNSSEC if requested
	if Options.DNSSEC {
		Options.Logger.Debug("using DNSSEC")
		msg.SetEdns0(1232, true)
	}

	resolver, err := query.LoadResolver(Options.Answers.Server, Options)
	if err != nil {
		return err
	}

	if Options.Debug {
		Options.Logger.SetLevel(3)
	}

	Options.Logger.Debug("Starting awl")

	var in = Options.Answers.DNS

	// Make the DNS request
	if Options.HTTPS {
		in, Options.Answers.RTT, err = resolver.LookUp(msg)
	} else if Options.QUIC {
		in, Options.Answers.RTT, err = resolver.LookUp(msg)
	} else {
		Options.Answers.Server = net.JoinHostPort(Options.Answers.Server, strconv.Itoa(Options.Port))
		d := new(dns.Client)

		// Set TCP/UDP, depending on flags
		if Options.TCP || Options.TLS {
			d.Net = "tcp"
		} else {
			Options.Logger.Debug("using udp")
			d.Net = "udp"
		}

		// Set IPv4 or IPv6, depending on flags
		switch {
		case Options.IPv4:
			d.Net += "4"
		case Options.IPv6:
			d.Net += "6"
		}

		// Add TLS, if requested
		if Options.TLS {
			d.Net += "-tls"
		}

		in, Options.Answers.RTT, err = d.Exchange(msg, Options.Answers.Server)
		if err != nil {
			return err
		}
		// If UDP truncates, use TCP instead (unless truncation is to be ignored)
		if in.MsgHdr.Truncated && !Options.Truncate {
			fmt.Printf(";; Truncated, retrying with TCP\n\n")
			d.Net = "tcp"
			switch {
			case Options.IPv4:
				d.Net += "4"
			case Options.IPv4:
				d.Net += "6"
			}
			in, Options.Answers.RTT, err = d.Exchange(msg, Options.Answers.Server)
		}
	}

	if err != nil {
		return err
	}

	if c.Bool("json") {
		json, err := json.MarshalIndent(in, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(json))
	} else {
		if !c.Bool("short") {
			// Print everything
			fmt.Println(in)
			fmt.Println(";; Query time:", Options.Answers.RTT)
			fmt.Println(";; SERVER:", Options.Answers.Server)
			fmt.Println(";; WHEN:", time.Now().Format(time.RFC1123Z))
			fmt.Println(";; MSG SIZE  rcvd:", in.Len())
		} else {
			// Print just the responses, nothing else
			for _, res := range in.Answer {
				temp := strings.Split(res.String(), "\t")
				fmt.Println(temp[len(temp)-1])
			}
		}
	}

	return nil
}
