// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/internal/helpers"
	"git.froth.zone/sam/awl/util"
	"github.com/miekg/dns"
	flag "github.com/stefansundin/go-zflag"
)

// ParseCLI parses arguments given from the CLI and passes them into an `Options`
// struct.
func ParseCLI(version string) (Options, error) {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	flag.Usage = func() {
		fmt.Println(`awl - drill, writ small

	Usage: awl name [@server] [record]
	 <name> domain, IP address, phone number
	 <record> defaults to A
	  
	 Arguments may be in any order, including flags.
	 Dig-like +[no]commands are also supported, see dig(1) or dig -h
		  
	Options:`)
		flag.PrintDefaults()
	}

	// CLI flag
	var (
		port  = flag.Int("port", 0, "`port` to make DNS query (default: 53 for UDP/TCP, 853 for TLS/QUIC)", flag.OptShorthand('p'), flag.OptDisablePrintDefault(true))
		query = flag.String("query", "", "domain name to `query` (default: .)", flag.OptShorthand('q'))
		class = flag.String("class", "IN", "DNS `class` to query", flag.OptShorthand('c'))
		qType = flag.String("qType", "", "`type` to query (default: A)", flag.OptShorthand('t'))

		ipv4    = flag.Bool("4", false, "force IPv4", flag.OptShorthandStr("4"))
		ipv6    = flag.Bool("6", false, "force IPv6", flag.OptShorthand('6'))
		reverse = flag.Bool("reverse", false, "do a reverse lookup", flag.OptShorthand('x'))

		timeout = flag.Float32("timeout", 1, "Timeout, in `seconds`")
		retry   = flag.Int("retries", 2, "number of `times` to retry")

		edns         = flag.Bool("no-edns", false, "disable EDNS entirely")
		ednsVer      = flag.Uint8("edns-ver", 0, "set EDNS version")
		expire       = flag.Bool("expire", false, "set EDNS expire")
		dnssec       = flag.Bool("dnssec", false, "enable DNSSEC", flag.OptShorthand('D'))
		nsid         = flag.Bool("nsid", false, "set EDNS NSID", flag.OptShorthand('n'))
		cookie       = flag.Bool("no-cookie", false, "disable sending EDNS cookie (default: cookie sent)")
		tcpKeepAlive = flag.Bool("keep-alive", false, "send EDNS TCP keep-alive")
		udpBufSize   = flag.Uint16("buffer-size", 1232, "set EDNS UDP buffer size", flag.OptShorthand('b'))
		zflag        = flag.String("zflag", "0", "set EDNS z-flag")
		subnet       = flag.String("subnet", "", "set EDNS subnet")
		padding      = flag.Bool("pad", false, "set EDNS padding")

		truncate = flag.Bool("no-truncate", false, "ignore truncation if a UDP request truncates (default: retry with TCP)")

		tcp      = flag.Bool("tcp", false, "use TCP")
		dnscrypt = flag.Bool("dnscrypt", false, "use DNSCrypt")
		tls      = flag.Bool("tls", false, "use DNS-over-TLS", flag.OptShorthand('T'))
		https    = flag.Bool("https", false, "use DNS-over-HTTPS", flag.OptShorthand('H'))
		quic     = flag.Bool("quic", false, "use DNS-over-QUIC", flag.OptShorthand('Q'))

		aa = flag.Bool("aa", false, "set/unset AA (Authoratative Answer) flag (default: not set)")
		ad = flag.Bool("ad", false, "set/unset AD (Authenticated Data) flag (default: not set)")
		cd = flag.Bool("cd", false, "set/unset CD (Checking Disabled) flag (default: not set)")
		qr = flag.Bool("qr", false, "set/unset QR (QueRy) flag (default: not set)")
		rd = flag.Bool("rd", true, "set/unset RD (Recursion Desired) flag (default: set)", flag.OptDisablePrintDefault(true))
		ra = flag.Bool("ra", false, "set/unset RA (Recursion Available) flag (default: not set)")
		tc = flag.Bool("tc", false, "set/unset TC (TrunCated) flag (default: not set)")
		z  = flag.Bool("z", false, "set/unset Z (Zero) flag (default: not set)", flag.OptShorthand('z'))

		short = flag.Bool("short", false, "print just the results", flag.OptShorthand('s'))
		json  = flag.Bool("json", false, "print the result(s) as JSON", flag.OptShorthand('j'))
		xml   = flag.Bool("xml", false, "print the result(s) as XML", flag.OptShorthand('X'))
		yaml  = flag.Bool("yaml", false, "print the result(s) as yaml", flag.OptShorthand('y'))

		noC     = flag.Bool("no-comments", false, "disable printing the comments")
		noQ     = flag.Bool("no-question", false, "disable printing the question section")
		noOpt   = flag.Bool("no-opt", false, "disable printing the OPT pseudosection")
		noAns   = flag.Bool("no-answer", false, "disable printing the answer section")
		noAuth  = flag.Bool("no-authority", false, "disable printing the authority section")
		noAdd   = flag.Bool("no-additional", false, "disable printing the additional section")
		noStats = flag.Bool("no-statistics", false, "disable printing the statistics section")

		verbosity   = flag.Int("verbosity", 1, "sets verbosity `level`", flag.OptShorthand('v'), flag.OptNoOptDefVal("2"))
		versionFlag = flag.Bool("version", false, "print version information", flag.OptShorthand('V'))
	)

	// Don't sort the flags when -h is given
	flag.CommandLine.SortFlags = false

	// Parse the flags
	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		return Options{Logger: util.InitLogger(*verbosity)}, fmt.Errorf("flag parse error: %w", err)
	}

	// TODO: DRY, dumb dumb.
	mbz, err := strconv.ParseInt(*zflag, 0, 16)
	if err != nil {
		return Options{Logger: util.InitLogger(*verbosity)}, fmt.Errorf("EDNS MBZ error: %w", err)
	}

	opts := Options{
		Logger:    util.InitLogger(*verbosity),
		Port:      *port,
		IPv4:      *ipv4,
		IPv6:      *ipv6,
		Short:     *short,
		TCP:       *tcp,
		DNSCrypt:  *dnscrypt,
		TLS:       *tls,
		HTTPS:     *https,
		QUIC:      *quic,
		Truncate:  *truncate,
		ShowQuery: false,
		AA:        *aa,
		AD:        *ad,
		TC:        *tc,
		Z:         *z,
		CD:        *cd,
		QR:        *qr,
		RD:        *rd,
		RA:        *ra,
		Reverse:   *reverse,
		HumanTTL:  false,
		ShowTTL:   true,
		JSON:      *json,
		XML:       *xml,
		YAML:      *yaml,
		Request: helpers.Request{
			Type:    dns.StringToType[strings.ToUpper(*qType)],
			Class:   dns.StringToClass[strings.ToUpper(*class)],
			Name:    *query,
			Timeout: time.Duration(*timeout * float32(time.Second)),
			Retries: *retry,
		},
		Display: Displays{
			Comments:   !*noC,
			Question:   !*noQ,
			Opt:        !*noOpt,
			Answer:     !*noAns,
			Authority:  !*noAuth,
			Additional: !*noAdd,
			Statistics: !*noStats,
		},
		EDNS: EDNS{
			EnableEDNS: !*edns,
			Cookie:     !*cookie,
			DNSSEC:     *dnssec,
			BufSize:    *udpBufSize,
			Version:    *ednsVer,
			Expire:     *expire,
			KeepOpen:   *tcpKeepAlive,
			Nsid:       *nsid,
			ZFlag:      uint16(mbz & 0x7FFF),
			Padding:    *padding,
		},
	}

	// TODO: DRY
	if *subnet != "" {
		err := parseSubnet(*subnet, &opts)
		if err != nil {
			return opts, err
		}
	}

	opts.Logger.Info("POSIX flags parsed")
	opts.Logger.Debug(fmt.Sprintf("%+v", opts))

	if *versionFlag {
		fmt.Printf("awl version %s, built with %s\n", version, runtime.Version())
		return opts, ErrNotError
	}

	// Parse all the arguments that don't start with - or --
	// This includes the dig-style (+) options
	err = ParseMiscArgs(flag.Args(), &opts)
	if err != nil {
		return opts, err
	}
	opts.Logger.Info("Dig/Drill flags parsed")
	opts.Logger.Debug(fmt.Sprintf("%+v", opts))

	if opts.Port == 0 {
		if opts.TLS || opts.QUIC {
			opts.Port = 853
		} else {
			opts.Port = 53
		}
	}
	opts.Logger.Info("Port set to", opts.Port)

	// Set timeout to 0.5 seconds if set below 0.5
	if opts.Request.Timeout < (time.Second / 2) {
		opts.Request.Timeout = (time.Second / 2)
	}

	if opts.Request.Retries < 0 {
		opts.Request.Retries = 0
	}

	return opts, nil
}
