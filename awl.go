package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/c-robinson/iplib"
	"github.com/miekg/dns"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/idna"
)

// The basic structure of a DNS request
type request struct {
	server  string // The server to make the DNS request from
	request uint16 // The type of request
	name    string // The domain name to make a DNS request for
}

func main() {
	// Custom version string
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("%s version %s, built with %s\n", c.App.Name, c.App.Version, runtime.Version())
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:  "v",
		Usage: "show version",
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:  "h",
		Usage: "show this help",
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
				Usage:       "port to make DNS query",
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
				Usage:   "use DNS-over-HTTPS (NOT FULLY COMPLETE)",
			},
			&cli.BoolFlag{
				Name:    "quic",
				Aliases: []string{"Q"},
				Usage:   "use DNS-over-QUIC (NOT YET IMPLEMENTED)",
			},
			&cli.BoolFlag{
				Name:  "no-truncate",
				Usage: "Ignore truncation if a UDP request truncates (default: retry with TCP)",
			},
			&cli.BoolFlag{
				Name:    "reverse",
				Aliases: []string{"x"},
				Usage:   "do a reverse lookup",
			},
		},
		Action: func(c *cli.Context) error {
			var err error

			req := parseArgs(c.Args().Slice())
			// Set DNS-over-TLS, if enabled
			port := c.Int("port")

			// If port is not set, set it
			if port == 0 {
				if c.Bool("tls") || c.Bool("quic") {
					port = 853
				} else {
					port = 53
				}
			}

			if !c.Bool("https") {
				req.server = net.JoinHostPort(req.server, strconv.Itoa(port))
			} else {
				req.server = "https://" + req.server
			}

			// Process the IP/Phone number so a PTR/NAPTR can be done
			if c.Bool("reverse") {
				if dns.TypeToString[req.request] == "A" {
					req.request = dns.StringToType["PTR"]
				}
				req.name, err = reverseDNS(req.name, dns.TypeToString[req.request])
				if err != nil {
					return err
				}
			}

			// if the domain is not canonical, make it canonical
			if !strings.HasSuffix(req.name, ".") {
				req.name = fmt.Sprintf("%s.", req.name)
			}

			msg := new(dns.Msg)
			msg.SetQuestion(req.name, req.request)

			// Set DNSSEC if requested
			if c.Bool("dnssec") {
				msg.SetEdns0(1232, true)
			}

			var (
				in  *dns.Msg
				rtt time.Duration
			)

			// Make the DNS request
			if c.Bool("https") {
				in, err = resolveHTTPS(msg, req.server)
			} else if c.Bool("quic") {
				return fmt.Errorf("quic: not yet implemented")
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
				in, rtt, err = d.Exchange(msg, req.server)

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
						in, rtt, err = d.Exchange(msg, req.server)
					}
				}
			}

			if err != nil {
				return err
			}

			if c.Bool("json") {
				json, _ := json.Marshal(in)
				fmt.Println(string(json))
			} else {
				if !c.Bool("short") {
					// Print everything
					fmt.Println(in)
					fmt.Println(";; Query time:", rtt)
					fmt.Println(";; SERVER:", req.server)
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseArgs(args []string) request {
	var (
		server string
		req    uint16
		name   string
	)
	for _, arg := range args {
		// If it starts with @, it's a DNS server
		if strings.HasPrefix(arg, "@") {
			server = strings.Split(arg, "@")[1]
			continue
		}
		// If there's a dot, it's a name
		if strings.Contains(arg, ".") {
			name, _ = idna.ToUnicode(arg)
			continue
		}
		// If it's a request, it's a request (duh)
		if r, ok := dns.StringToType[strings.ToUpper(arg)]; ok {
			req = r
			continue
		}

		//else, assume it's a name
		name, _ = idna.ToUnicode(arg)
	}

	// If nothing was set, set a default
	if name == "" {
		name = "."
		if req == 0 {
			req = dns.StringToType["NS"]
		}
	} else {
		if req == 0 {
			req = dns.StringToType["A"]
		}
	}
	if server == "" {
		resolv, err := dns.ClientConfigFromFile("/etc/resolv.conf")
		if err != nil { // Query Google by default, needed for Windows since the DNS library doesn't support Windows
			// TODO: Actually find where windows stuffs its dns resolvers
			server = "8.8.4.4"
		} else {
			server = resolv.Servers[0]
		}
	}

	return request{server: server, request: req, name: name}
}

func reverseDNS(dom string, q string) (string, error) {
	if q == "PTR" {
		if strings.Contains(dom, ".") {
			// It's an IPv4 address
			ip := net.ParseIP(dom)
			if ip != nil {
				return iplib.IP4ToARPA(ip), nil
			} else {
				return "", errors.New("error: Could not parse IPv4 address")
			}

		} else if strings.Contains(dom, ":") {
			// It's an IPv6 address
			ip := net.ParseIP(dom)
			if ip != nil {
				return iplib.IP6ToARPA(ip), nil
			} else {
				return "", errors.New("error: Could not parse IPv6 address")
			}
		}
	}
	return "", errors.New("error: -x flag given but no IP found")
}

func resolveHTTPS(msg *dns.Msg, server string) (*dns.Msg, error) {
	httpR := &http.Client{}
	buf, err := msg.Pack()
	if err != nil {
		return nil, err
	}
	query := server + "?dns=" + base64.RawURLEncoding.EncodeToString(buf)
	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/dns-message")

	res, err := httpR.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad HTTP Request: %d", res.StatusCode)
	}

	fullRes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	response := dns.Msg{}
	err = response.Unpack(fullRes)
	if err != nil {
		return nil, err
	}

	return &response, nil

}
