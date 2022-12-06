// SPDX-License-Identifier: BSD-3-Clause

package cli

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"git.froth.zone/sam/awl/pkg/util"
	"github.com/miekg/dns"
	flag "github.com/stefansundin/go-zflag"
)

// ParseCLI parses arguments given from the CLI and passes them into an `Options`
// struct.
func ParseCLI(args []string, version string) (*util.Options, error) {
	flagSet := flag.NewFlagSet(args[0], flag.ContinueOnError)

	flagSet.Usage = func() {
		fmt.Println(`awl - drill, writ small

	Usage: awl name [@server] [record]
	 <name> domain, IP address, phone number
	 <record> defaults to A

	 Arguments may be in any order, including flags.
	 Dig-like +[no]commands are also supported, see dig(1) or dig -h

	Options:`)
		flagSet.PrintDefaults()
	}

	// CLI flags
	//
	// Remember, when adding a flag edit the manpage and the completions :)
	var (
		port  = flagSet.Int("port", 0, "`port` to make DNS query (default: 53 for UDP/TCP, 853 for TLS/QUIC)", flag.OptShorthand('p'), flag.OptDisablePrintDefault(true))
		query = flagSet.String("query", "", "domain name to `query` (default: .)", flag.OptShorthand('q'))
		class = flagSet.String("class", "IN", "DNS `class` to query", flag.OptShorthand('c'))
		qType = flagSet.String("qType", "", "`type` to query (default: A)", flag.OptShorthand('t'))

		ipv4    = flagSet.Bool("4", false, "force IPv4", flag.OptShorthand('4'))
		ipv6    = flagSet.Bool("6", false, "force IPv6", flag.OptShorthand('6'))
		reverse = flagSet.Bool("reverse", false, "do a reverse lookup", flag.OptShorthand('x'))

		timeout = flagSet.Float32("timeout", 5, "Timeout, in `seconds`")
		retry   = flagSet.Int("retries", 2, "number of `times` to retry")

		edns         = flagSet.Bool("no-edns", false, "disable EDNS entirely")
		ednsVer      = flagSet.Uint8("edns-ver", 0, "set EDNS version")
		dnssec       = flagSet.Bool("dnssec", false, "enable DNSSEC", flag.OptShorthand('D'))
		expire       = flagSet.Bool("expire", false, "set EDNS expire")
		nsid         = flagSet.Bool("nsid", false, "set EDNS NSID", flag.OptShorthand('n'))
		cookie       = flagSet.Bool("no-cookie", false, "disable sending EDNS cookie (default: cookie sent)")
		tcpKeepAlive = flagSet.Bool("keep-alive", false, "send EDNS TCP keep-alive")
		udpBufSize   = flagSet.Uint16("buffer-size", 1232, "set EDNS UDP buffer size", flag.OptShorthand('b'))
		mbzflag      = flagSet.String("zflag", "0", "set EDNS z-flag `value`")
		subnet       = flagSet.String("subnet", "", "set EDNS client subnet")
		padding      = flagSet.Bool("pad", false, "set EDNS padding")

		badCookie = flagSet.Bool("no-bad-cookie", false, "ignore BADCOOKIE EDNS responses (default: retry with correct cookie")
		truncate  = flagSet.Bool("no-truncate", false, "ignore truncation if a UDP request truncates (default: retry with TCP)")

		tcp      = flagSet.Bool("tcp", false, "use TCP")
		dnscrypt = flagSet.Bool("dnscrypt", false, "use DNSCrypt")
		tls      = flagSet.Bool("tls", false, "use DNS-over-TLS", flag.OptShorthand('T'))
		https    = flagSet.Bool("https", false, "use DNS-over-HTTPS", flag.OptShorthand('H'))
		quic     = flagSet.Bool("quic", false, "use DNS-over-QUIC", flag.OptShorthand('Q'))

		tlsHost  = flagSet.String("tls-host", "", "Server name to use for TLS verification")
		noVerify = flagSet.Bool("tls-no-verify", false, "Disable TLS cert verification")

		aaflag = flagSet.Bool("aa", false, "set/unset AA (Authoratative Answer) flag (default: not set)")
		adflag = flagSet.Bool("ad", false, "set/unset AD (Authenticated Data) flag (default: not set)")
		cdflag = flagSet.Bool("cd", false, "set/unset CD (Checking Disabled) flag (default: not set)")
		qrflag = flagSet.Bool("qr", false, "set/unset QR (QueRy) flag (default: not set)")
		rdflag = flagSet.Bool("rd", true, "set/unset RD (Recursion Desired) flag (default: set)", flag.OptDisablePrintDefault(true))
		raflag = flagSet.Bool("ra", false, "set/unset RA (Recursion Available) flag (default: not set)")
		tcflag = flagSet.Bool("tc", false, "set/unset TC (TrunCated) flag (default: not set)")
		zflag  = flagSet.Bool("z", false, "set/unset Z (Zero) flag (default: not set)", flag.OptShorthand('z'))

		short = flagSet.Bool("short", false, "print just the results", flag.OptShorthand('s'))
		json  = flagSet.Bool("json", false, "print the result(s) as JSON", flag.OptShorthand('j'))
		xml   = flagSet.Bool("xml", false, "print the result(s) as XML", flag.OptShorthand('X'))
		yaml  = flagSet.Bool("yaml", false, "print the result(s) as yaml", flag.OptShorthand('y'))

		noC     = flagSet.Bool("no-comments", false, "disable printing the comments")
		noQ     = flagSet.Bool("no-question", false, "disable printing the question section")
		noOpt   = flagSet.Bool("no-opt", false, "disable printing the OPT pseudosection")
		noAns   = flagSet.Bool("no-answer", false, "disable printing the answer section")
		noAuth  = flagSet.Bool("no-authority", false, "disable printing the authority section")
		noAdd   = flagSet.Bool("no-additional", false, "disable printing the additional section")
		noStats = flagSet.Bool("no-statistics", false, "disable printing the statistics section")

		verbosity   = flagSet.Int("verbosity", 1, "sets verbosity `level`", flag.OptShorthand('v'), flag.OptNoOptDefVal("2"))
		versionFlag = flagSet.Bool("version", false, "print version information", flag.OptShorthand('V'))
	)

	// Don't sort the flags when -h is given
	flagSet.SortFlags = false

	// Parse the flags
	if err := flagSet.Parse(args[1:]); err != nil {
		return &util.Options{Logger: util.InitLogger(*verbosity)}, fmt.Errorf("flag: %w", err)
	}

	// TODO: DRY, dumb dumb.
	mbz, err := strconv.ParseInt(*mbzflag, 0, 16)
	if err != nil {
		return &util.Options{Logger: util.InitLogger(*verbosity)}, fmt.Errorf("EDNS MBZ: %w", err)
	}

	opts := util.Options{
		Logger:      util.InitLogger(*verbosity),
		IPv4:        *ipv4,
		IPv6:        *ipv6,
		Short:       *short,
		TCP:         *tcp,
		DNSCrypt:    *dnscrypt,
		TLS:         *tls,
		TLSHost:     *tlsHost,
		TLSNoVerify: *noVerify,
		HTTPS:       *https,
		QUIC:        *quic,
		Truncate:    *truncate,
		BadCookie:   *badCookie,
		Reverse:     *reverse,
		JSON:        *json,
		XML:         *xml,
		YAML:        *yaml,
		HeaderFlags: util.HeaderFlags{
			AA: *aaflag,
			AD: *adflag,
			TC: *tcflag,
			Z:  *zflag,
			CD: *cdflag,
			QR: *qrflag,
			RD: *rdflag,
			RA: *raflag,
		},
		Request: util.Request{
			Type:    dns.StringToType[strings.ToUpper(*qType)],
			Class:   dns.StringToClass[strings.ToUpper(*class)],
			Name:    *query,
			Timeout: time.Duration(*timeout * float32(time.Second)),
			Retries: *retry,
			Port:    *port,
		},
		Display: util.Display{
			Comments:       !*noC,
			Question:       !*noQ,
			Opt:            !*noOpt,
			Answer:         !*noAns,
			Authority:      !*noAuth,
			Additional:     !*noAdd,
			Statistics:     !*noStats,
			TTL:            true,
			ShowClass:      true,
			ShowQuery:      false,
			HumanTTL:       false,
			UcodeTranslate: true,
		},
		EDNS: util.EDNS{
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
		HTTPSOptions: util.HTTPSOptions{
			Endpoint: "/dns-query",
			Get:      false,
		},
	}

	// TODO: DRY
	if *subnet != "" {
		if err = util.ParseSubnet(*subnet, &opts); err != nil {
			return &opts, fmt.Errorf("%w", err)
		}
	}

	opts.Logger.Info("POSIX flags parsed")
	opts.Logger.Debug(fmt.Sprintf("%+v", opts))

	if *versionFlag {
		fmt.Printf("awl version %s, built with %s\n", version, runtime.Version())

		return &opts, util.ErrNotError
	}

	// Parse all the arguments that don't start with - or --
	// This includes the dig-style (+) options
	err = ParseMiscArgs(flagSet.Args(), &opts)
	if err != nil {
		return &opts, err
	}

	opts.Logger.Info("Dig/Drill flags parsed")
	opts.Logger.Debug(fmt.Sprintf("%+v", opts))

	if opts.Request.Port == 0 {
		if opts.TLS || opts.QUIC {
			opts.Request.Port = 853
		} else {
			opts.Request.Port = 53
		}
	}

	opts.Logger.Info("Port set to", opts.Request.Port)

	// Set timeout to 0.5 seconds if set below 0.5
	if opts.Request.Timeout < (time.Second / 2) {
		opts.Request.Timeout = (time.Second / 2)
	}

	if opts.Request.Retries < 0 {
		opts.Request.Retries = 0
	}

	opts.Logger.Info("Options fully populated")
	opts.Logger.Debug(fmt.Sprintf("%+v", opts))

	return &opts, nil
}

var errNoArg = errors.New("no argument given")

type errInvalidArg struct {
	arg string
}

func (e *errInvalidArg) Error() string {
	return fmt.Sprintf("digflags: invalid argument %s", e.arg)
}
