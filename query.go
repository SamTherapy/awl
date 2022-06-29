package main

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/logawl"
	"git.froth.zone/sam/awl/query"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
)

func doQuery(c *cli.Context) error {
	var (
		err     error
		resp    util.Response
		isHTTPS bool
		Logger  = logawl.New() //init logger
	)

	resp.Answers, err = parseArgs(c.Args().Slice())
	if err != nil {
		Logger.Error("Unable to parse args")
		return err
	}
	port := c.Int("port")
	if c.Bool("debug") {
		Logger.SetLevel(3)
	}

	Logger.Debug("Starting awl")
	// If port is not set, set it
	if port == 0 {
		if c.Bool("tls") || c.Bool("quic") {
			port = 853
		} else {
			port = 53
		}
	}

	if c.Bool("https") || strings.HasPrefix(resp.Answers.Server, "https://") {
		// add https:// if it doesn't already exist
		if !strings.HasPrefix(resp.Answers.Server, "https://") {
			resp.Answers.Server = "https://" + resp.Answers.Server
		}
		isHTTPS = true
	} else {
		resp.Answers.Server = net.JoinHostPort(resp.Answers.Server, strconv.Itoa(port))
	}

	// Process the IP/Phone number so a PTR/NAPTR can be done
	if c.Bool("reverse") {
		if dns.TypeToString[resp.Answers.Request] == "A" {
			resp.Answers.Request = dns.StringToType["PTR"]
		}
		resp.Answers.Name, err = util.ReverseDNS(resp.Answers.Name, dns.TypeToString[resp.Answers.Request])
		if err != nil {
			return err
		}
	}

	// if the domain is not canonical, make it canonical
	if !strings.HasSuffix(resp.Answers.Name, ".") {
		resp.Answers.Name = fmt.Sprintf("%s.", resp.Answers.Name)
	}

	msg := new(dns.Msg)

	msg.SetQuestion(resp.Answers.Name, resp.Answers.Request)

	// Make this authoritative (does this do anything?)
	if c.Bool("aa") {
		msg.Authoritative = true
	}
	// Set truncated flag (why)
	if c.Bool("tc") {
		msg.Truncated = true
	}
	// Set the zero flag if requested (does nothing)
	if c.Bool("z") {
		Logger.Debug("Setting message to zero")
		msg.Zero = true
	}
	// Disable DNSSEC validation
	if c.Bool("cd") {
		msg.CheckingDisabled = true
	}
	// Disable wanting recursion
	if c.Bool("no-rd") {
		msg.RecursionDesired = false
	}
	// Disable recursion being available (I don't think this does anything)
	if c.Bool("no-ra") {
		msg.RecursionAvailable = false
	}
	// Set DNSSEC if requested
	if c.Bool("dnssec") {
		Logger.Debug("Using DNSSEC")
		msg.SetEdns0(1232, true)
	}

	var in *dns.Msg

	// Make the DNS request
	if isHTTPS {
		in, resp.Answers.RTT, err = query.ResolveHTTPS(msg, resp.Answers.Server)
	} else if c.Bool("quic") {
		in, resp.Answers.RTT, err = query.ResolveQUIC(msg, resp.Answers.Server)
	} else {

		d := new(dns.Client)

		// Set TCP/UDP, depending on flags
		if c.Bool("tcp") || c.Bool("tls") {
			d.Net = "tcp"
		} else {
			d.Net = "udp"
		}

		// Set IPv4 or IPv6, depending on flags
		switch {
		case c.Bool("4"):
			d.Net += "4"
		case c.Bool("6"):
			d.Net += "6"
		}

		// Add TLS, if requested
		if c.Bool("tls") {
			d.Net += "-tls"
		}

		in, resp.Answers.RTT, err = d.Exchange(msg, resp.Answers.Server)
		if err != nil {
			return err
		}
		// If UDP truncates, use TCP instead (unless truncation is to be ignored)
		if in.MsgHdr.Truncated && !c.Bool("no-truncate") {
			fmt.Printf(";; Truncated, retrying with TCP\n\n")
			d.Net = "tcp"
			switch {
			case c.Bool("4"):
				d.Net += "4"
			case c.Bool("6"):
				d.Net += "6"
			}
			in, resp.Answers.RTT, err = d.Exchange(msg, resp.Answers.Server)
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
			fmt.Println(";; Query time:", resp.Answers.RTT)
			fmt.Println(";; SERVER:", resp.Answers.Server)
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
