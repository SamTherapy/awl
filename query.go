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
	var err error
	var resp util.Response
	var isHTTPS bool
	resp.Answers, err = parseArgs(c.Args().Slice())
	if err != nil {
		return nil
	}
	port := c.Int("port")

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

	// Set the zero flag if requested (does nothing)
	if c.Bool("z") {
		msg.Zero = true
	}
	// Disable DNSSEC validation if enabled
	if c.Bool("cd") {
		msg.CheckingDisabled = true
	}

	if c.Bool("no-rd") {
		msg.RecursionDesired = false
	}

	if c.Bool("no-ra") {
		msg.RecursionAvailable = false
	}

	// Set DNSSEC if requested
	if c.Bool("dnssec") {
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
		if c.Bool("tcp") {
			d.Net = "tcp"
			if c.Bool("4") {
				d.Net = "tcp4"
			}
			if c.Bool("6") {
				d.Net = "tcp6"
			}
		} else {
			d.Net = "udp"
			if c.Bool("4") {
				d.Net = "udp4"
			}
			if c.Bool("6") {
				d.Net = "udp6"
			}
		}

		// This is apparently all it takes to enable DoT
		// TODO: Is it really?
		if c.Bool("tls") {
			d.Net = "tcp-tls"
		}

		in, resp.Answers.RTT, err = d.Exchange(msg, resp.Answers.Server)
		if err != nil {
			return err
		}
		// If UDP truncates, use TCP instead
		if !c.Bool("no-truncate") {
			if in.MsgHdr.Truncated {
				fmt.Printf(";; Truncated, retrying with TCP\n\n")
				d.Net = "tcp"
				if c.Bool("4") {
					d.Net = "tcp4"
				}
				if c.Bool("6") {
					d.Net = "tcp6"
				}
				in, resp.Answers.RTT, err = d.Exchange(msg, resp.Answers.Server)
			}
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
			fmt.Println(";; WHEN:", time.Now().Format(time.RFC1123))
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
