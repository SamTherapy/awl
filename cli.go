package main

import (
	"fmt"
	"runtime"
	"strings"

	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/idna"
)

// Do all the magic CLI crap
func prepareCLI() *cli.App {
	// Custom version string
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s, built with %s\n", c.App.Name, c.App.Version, runtime.Version())
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "v",
		Usage: "show version and exit",
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "h",
		Usage: "show this help and exit",
	}

	// Hack to get rid of the annoying default on the CLI
	oldFlagStringer := cli.FlagStringer
	cli.FlagStringer = func(f cli.Flag) string {
		return strings.TrimSuffix(oldFlagStringer(f), " (default: false)")
	}

	cli.AppHelpTemplate = `{{.Name}} - {{.Usage}}
	
		Usage: {{.HelpName}} name [@server] [type]
			<name>	can be a name or an IP address
			<type>	defaults to A
	
			arguments can be in any order
		{{if .VisibleFlags}}
		Options:
			{{range .VisibleFlags}}{{.}}
			{{end}}{{end}}`
	app := &cli.App{
		Name:    "awl",
		Usage:   "drill, writ small",
		Version: "v0.2.0",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Usage:       "`port` to make DNS query",
				DefaultText: "53 over plain TCP/UDP, 853 over TLS or QUIC, and 443 over HTTPS",
			},
			&cli.BoolFlag{
				Name:  "4",
				Usage: "force IPv4",
			},
			&cli.BoolFlag{
				Name:  "6",
				Usage: "force IPv6",
			},
			&cli.BoolFlag{
				Name:    "dnssec",
				Aliases: []string{"D"},
				Usage:   "enable DNSSEC",
			},
			&cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"j"},
				Usage:   "return the result(s) as JSON",
			},
			&cli.BoolFlag{
				Name:    "short",
				Aliases: []string{"s"},
				Usage:   "print just the results, equivalent to dig +short",
			},
			&cli.BoolFlag{
				Name:    "tcp",
				Aliases: []string{"t"},
				Usage:   "use TCP (default: use UDP)",
			},
			&cli.BoolFlag{
				Name:    "tls",
				Aliases: []string{"T"},
				Usage:   "use DNS-over-TLS",
			},
			&cli.BoolFlag{
				Name:    "https",
				Aliases: []string{"H"},
				Usage:   "use DNS-over-HTTPS",
			},
			&cli.BoolFlag{
				Name:    "quic",
				Aliases: []string{"Q"},
				Usage:   "use DNS-over-QUIC",
			},
			&cli.BoolFlag{
				Name:  "no-truncate",
				Usage: "ignore truncation if a UDP request truncates (default: retry with TCP)",
			},
			&cli.BoolFlag{
				Name:  "z",
				Usage: "set Z (zero) flag (default: not set) [THIS DOESN'T DO ANYTHING]",
			},
			&cli.BoolFlag{
				Name:  "cd",
				Usage: "set CD (CheckingDisabled) flag (default: not set)",
			},
			&cli.BoolFlag{
				Name:  "no-rd",
				Usage: "UNset RD (RecursionDesired) flag (default: set)",
			},
			&cli.BoolFlag{
				Name:  "no-ra",
				Usage: "UNset RA (RecursionAvailable) flag (default: set)",
			},
			&cli.BoolFlag{
				Name:    "reverse",
				Aliases: []string{"x"},
				Usage:   "do a reverse lookup",
			},
		},
		Action: doQuery,
	}
	return app
}

// Parse the wildcard arguments, drill style
func parseArgs(args []string) (util.Answers, error) {
	var (
		resp util.Response
	)
	for _, arg := range args {
		// If it starts with @, it's a DNS server
		if strings.HasPrefix(arg, "@") {
			resp.Answers.Server = strings.Split(arg, "@")[1]
			continue
		}
		// If there's a dot, it's a name
		if strings.Contains(arg, ".") {
			resp.Answers.Name, _ = idna.ToUnicode(arg)
			continue
		}
		// If it's a request, it's a request (duh)
		if r, ok := dns.StringToType[strings.ToUpper(arg)]; ok {
			resp.Answers.Request = r
			continue
		}

		//else, assume it's a name
		var err error
		resp.Answers.Name, err = idna.ToUnicode(arg)
		if err != nil {
			return util.Answers{}, err
		}
	}

	// If nothing was set, set a default
	if resp.Answers.Name == "" {
		resp.Answers.Name = "."
		if resp.Answers.Request == 0 {
			resp.Answers.Request = dns.StringToType["NS"]
		}
	} else {
		if resp.Answers.Request == 0 {
			resp.Answers.Request = dns.StringToType["A"]
		}
	}
	if resp.Answers.Server == "" {
		resolv, err := dns.ClientConfigFromFile("/etc/resolv.conf")
		if err != nil { // Query Google by default, needed for Windows since the DNS library doesn't support Windows
			// TODO: Actually find where windows stuffs its dns resolvers
			resp.Answers.Server = "8.8.4.4"
		} else {
			resp.Answers.Server = resolv.Servers[0]
		}
	}

	return util.Answers{Server: resp.Answers.Server, Request: resp.Answers.Request, Name: resp.Answers.Name}, nil
}
