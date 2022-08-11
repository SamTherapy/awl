awl(1)

# NAME

awl - DNS lookup tool

# SYNOPSIS

_awl_ [ _OPTIONS_ ] _name_ [ _@server_ ] [ _type_ ], where

_name_ is the query to make (_example: froth.zone_)++
_@server_ is the server to query (_example: dns.froth.zone_)++
_type_ is the DNS resource type (_example: AAAA_)

# DESCRIPTION

_awl_ (*a*wls *w*ant *l*icorice) is a simple tool designed to make DNS queries,
much like the venerable _dig_(1). An awl is a tool used to make small holes,
typically used in leatherworking.

_awl_ is designed to be a more "modern" version of _drill_(1) by including
some more recent RFCs and output options. _awl_ is still heavily
Work-In-Progress so some features may get added or removed.

# OPTIONS

Dig-like +[no]flags are supported, see dig(1)

_-D_, _--dnssec_, _+dnssec_
	Enable DNSSEC. This needs to be manually enabled.

_-v_ _value_
	Set verbosity (currently WIP)

_-V_
	Print the version and exit.

_-h_
	Show a "short" help message.

## Query Options

_-4_
	Only make query over IPv4

_-6_
	Only make query over IPv6

_-p_, _--port_ _port_
	Sets the port to query. Default ports listed below.
	
	- _53_ for _UDP_ and _TCP_
	- _853_ for _TLS_ and _QUIC_
	- _443_ for _HTTPS_

_-q_, _--query_ _domain_
	Domain to query (eg. example.com)

_-c_, _--class_ _class_
	DNS class to query (eg. IN, CH)

_-t_, _--qType_ _type_
	DNS type to query (eg. A, NS)

_--no-truncate_, _+ignore_
	Ignore UDP truncation (by default, awl _retries with TCP_)

_--tcp_, _+tcp_, _+vc_
	Use TCP for the query (see _RFC 7766_)

_--dnscrypt_, _+dnscrypt_
	Use DNSCrypt

_-T_, _--tls_, _+tls_
	Use DNS-over-TLS, implies _--tcp_ (see _RFC 7858_)

_--tls-host_ _string_
	Set hostname to use for TLS certificate validation.
	Default is the name of the domain when querying over TLS, and empty for IPs.

_--tls-no-verify_
	Ignore TLS validation when performing a DNS query.

_-H_. _--https_, _+https_
	Use DNS-over-HTTPS (see _RFC 8484_)

_-Q_. _--quic_, _+quic_
	Use DNS-over-QUIC (see _RFC 9250_)

_-x_, _--reverse_
	Do a reverse lookup. Sets default _type_ to PTR.++
_awl_ automatically makes an IP or phone number canonical.

_--timeout_ _seconds_, _+timeout=__seconds_
	Set the timeout period. Floating point numbers are accepted.++
0.5 seconds is the minimum.

_--retries_ _int_, _+tries_=_int_, _+ retry_=_int_
	Set the number of retries.
	Retry is one more than tries, dig style

## DNS Flags

_--aa=[false]_, _+[no]aaflag_
	(Set, Unset) AA (Authoritative Answer) flag

_--ad=[false]_, _+[no]adflag_
	(Set, Unset) AD (Authenticated Data) flag

_--tc=[false]_, _+[no]tcflag_
	(Set, Unset) TC (TrunCated) flag

_-z=[false]_, _+[no]zflag_
	(Set, Unset) Z (Zero) flag

_--cd=[false]_, _+[no]cdflag_
	(Set, Unset) CD (Checking Disabled) flag

_--qr=[false]_, _+[no]qrflag_
	(Set, Unset) QR (QueRy) flag

_--rd=[true]_, _+[no]rdflag_
	(Set, Unset) RD (Recursion Desired) flag

_--ra=[false]_, _+[no]raflag_
	(Set, Unset) RA (Recursion Available) flag

## Output Display

_--no-question_, _+noquestion_
	Do not display the Question section

_--no-answer_, _+noanswer_
	Do not display the Answer section

_--no-answer_, _+noanswer_
	Do not display the Answer section

_--no-authority_, _+noauthority_
	Do not display the Authority section

_--no-additional_, _+noadditional_
	Do not display the Additional section

_--no-statistics_, _+nostats_
	Do not display the Statistics (additional comments) section

## Output Formats

_-j_, _--json_, _+json_
	Print the query results as JSON.

_-X_, _--xml_, _+xml_
	Print the query results as XML.

_-y_, _--yaml_, _+yaml_
	Print the query results as YAML.

_-s_, _--short_, _+short_
	Print just the address of the answer.

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

_drill_(1), _dig_(1), the many DNS RFCs
