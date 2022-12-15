awl(1)
; SPDX-License-Identifier: BSD-3-Clause

# NAME

awl - DNS lookup tool

# SYNOPSIS

*awl* [ _OPTIONS_ ] _name_ [ _@server_ ] [ _type_ ], where

_name_ is the query to make (example: froth.zone)++
_@server_ is the server to query (example: dns.froth.zone)++
_type_ is the DNS resource type (example: AAAA)

# DESCRIPTION

*awl* (*a*wls *w*ant *l*icorice) is a simple tool designed to make DNS queries,
much like the venerable *dig*(1). An awl is a tool used to make small holes,
typically used in leatherworking.

*awl* is designed to be a more "modern" version of *drill*(1) by including
some more recent RFCs and output options.

When no arguments are given, *awl* will perform an _NS_ query on the root ('_._').

When a nameserver is not given, *awl* will query a random system nameserver.
If one cannot be found, *awl* will query the localhost.

# OPTIONS

*-4*
	Force only IPv4

*-6*
	Force only IPv6

*-c*, *--class* _class_
	DNS class to query (eg. IN, CH)
	The default is IN.

*-h*
	Show a "short" help message.

*-p*, *--port* *port*
	Sets the port to query. Default ports listed below.
	- _53_ for *UDP* and *TCP*
	- _853_ for *TLS* and *QUIC*
	- _443_ for *HTTPS*

*-q*, *--query* _domain_
	Explicitly set a domain to query (eg. example.com)

*-t*, *--qType* _type_
	Explicitly set a DNS type to query (eg. A, AAAA, NS)
	The default is A.

*-v*[=_int_]
	Set verbosity of output
	Accepted values are as follows:
	- _0_: Only log errors.
	- _1_: Log warnings. *This is the default.*
	- _2_: Log information *Default when specifying just* _-v_.
	- _3_: Log information useful for debugging.

	Setting a value lower than 0 disables logging entirely.

	By default, specifying just *-v* sets the verbosity to 2 (info).

*-x*, *--reverse*
	Do a reverse lookup. Sets default *type* to PTR.
	*awl* automatically makes an IP or phone number canonical.

*-V*
	Print the version and exit.

# QUERY OPTIONS

Anything in [brackets] is optional.
Many options are inherited from *dig*(1).

*--aa*[=_bool_], *+*[no]*aaflag*, *+*[no]*aaonly*
	Sets the AA (Authoritative Answer) flag.

*--ad*[=_bool_], *+*[no]*adflag*
	Sets the AD (Authenticated Data) flag.

*--no-additional*, *+*[no]*additional*
	Toggle the display of the Additional section.

*--no-answer*, *+*[no]*answer*
	Toggle the display of the Answer section.

*--no-authority*, *+*[no]*authority*
	Toggle the display of the Authority section.

*--no-bad-cookie*, *+*[no]*badcookie*
	\[Do not\] ignore BADCOOKIE responses

*--buffer-size* _int_, *+bufize*=_int_
	Set the UDP message buffer size, using EDNS.
	Max is 65535, minimum is zero.
	The default value is 1232.

*--cd*[=_bool_], *+*[no]*cdflag*
	(Set, Unset) CD (Checking Disabled) flag.

*--no-cookie*, *+*[no]*cookie*[=_string_]
	Send an EDNS cookie.
	This is enabled by default with a random string.

*-D*, *--dnssec*, *+dnssec*, *+do*
	Request DNSSEC records as well.
	This sets the DNSSEC OK bit (DO)

*--dnscrypt*, *+*[no]*dnscrypt*
	Use DNSCrypt.

*--expire*. *+*[no]*expire*
	Send an EDNS Expire.


*--edns-ver*, *+edns*[=_int_]
	Enable EDNS and set EDNS version.
	The maximum value is 255, and the minimum (default) value is 0.

*--no-edns*, *+noedns*
	Disable EDNS.

*-H*, *--https*, *+*[no]*https*[=_endpoint_], *+*[no]*https-post*[=_endpoint_]
	Use DNS-over-HTTPS (see RFC 8484).
	The default endpoint is _/dns-query_

*+*[no]*https-get*[=_endpoint_]
	Use an HTTP GET instead of an HTTP POST when making a DNS-over-HTTPS query.

*+*[no]*idnout*
	Converts [or leaves] punycode on output.
	Input is automatically translated to punycode.

*--no-truncate*, *+ignore*
	Ignore UDP truncation (by default, awl *retries with TCP*).

*-j*, *--json*, *+*[no]*json*
	Print the query results as JSON.
	The result is *not* in compliance with RFC 8427.

*--keep-alive*, *+*[no]*keepalive*, *+*[no]*keepopen*
	Send an EDNS keep-alive.
	This does nothing unless using TCP.

*--nsid*, *+*[no]*nsid*
	Send an EDNS name server ID request.

*--qr*[=_bool_], *+*[no]*qrflag*
	Sets the QR (QueRy) flag.

*--no-question*, *+*[no]*question*
	Toggle the display of the Question section.

*-Q*. *--quic*, *+*[no]*quic*
	Use DNS-over-QUIC (see RFC 9250).

*-s*, *--short*, *+*[no]*short*
	Print just the address of the answer.

*--no-statistics*, *+*[no]*stats*
	Toggle the display of the Statistics (additional comments) section.

*--subnet* _ip_[_/prefix_], *+*[no]*subnet*[=_ip_[_/prefix_]]
	Send an EDNS Client Subnet option with the specified address.

	Like *dig*(1), setting the IP to _0.0.0.0/0_, _::/0_ or _0_ will signal the resolver to not use any client information when returning the query.

*--tc*[=_bool_], *+*[no]*tcflag*
	Sets the TC (TrunCated) flag

*--tcp*, *+*[no]*tcp*, *+*[no]*vc*
	Use TCP for the query (see RFC 7766).

*--timeout* _seconds_, *+timeout*=_seconds_
	Set the timeout period. Floating point numbers are accepted.
	0.5 seconds is the minimum.

*-T*, *--tls*, *+*[no]*tls*
	Use DNS-over-TLS, implies *--tcp* (see RFC 7858)

*--tls-host* _string_
	Set hostname to use for TLS certificate validation.
	Default is the name of the domain when querying over TLS, and empty for IPs.

*--tls-no-verify*
	Ignore TLS validation when performing a DNS query.

*--trace*, *+trace*
	Trace the path of the query from the root, acting like its own resolver.
	This option enables DNSSEC.
	When *@server* is specified, this will only affect the initial query.

*--retries* _int_, *+tries*=_int_, *+retry*=_int_
	Set the number of retries.
	Retry is one more than tries, dig style.

*-X*, *--xml*, *+*[no]*xml*
	Print the query results as XML.

*-y*, *--yaml*, *+*[no]*yaml*
	Print the query results as YAML.

*-z*[=_bool_], *+*[no]*zflag*
	Sets the Z (Zero) flag.

*--zflag* _int_, *+ednsflags*=_int_
	Set the must-be-zero EDNS flags.
	Decimal, hexadecimal and octal are supported.
	Trying to set DO will be ignored.

# EXIT STATUS

The exit code is 0 when a query is successfully made and received.
This includes SERVFAILs, NOTIMPL among others.

# EXAMPLES

```
awl grumbulon.xyz -j +cd
```

Run a query of your local resolver for the A records of grumbulon.xyz, print
them as JSON and disable DNSSEC verification.

```
awl +short example.com AAAA @1.1.1.1
```

Query 1.1.1.1 for the AAAA records of example.com, print just the answers

```
awl -xT PTR 8.8.4.4 @dns.google
```

Query dns.google over TLS for the PTR record to the IP address 8.8.4.4

# SEE ALSO

*drill*(1), *dig*(1)

# STANDARDS

RFC 1034,1035 (UDP), 7766 (TCP), 7858 (TLS), 8484 (HTTPS), 9230 (QUIC)

Probably more, _https://www.statdns.com/rfc_

# BUGS

Full parity with *dig*(1) is not complete.

This man page is probably not complete.

Likely numerous more, report them either to the tracker
_https://git.froth.zone/sam/awl/issues_ or via email
_~sammefishe/awl-dev@lists.sr.ht_
