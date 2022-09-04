#compdef awl
# SPDX-License-Identifier: BSD-3-Clause

local curcontext="$curcontext" state line expl
local -a alts args
[[ -prefix + ]] && args=(
  '*+'{no,}'tcp[use TCP instead of UDP for queries]'
  '*+'{no,}'ignore[ignore truncation in UDP responses]'
  '*+'{no,}'tls[use DNS-over-TLS for queries]'
  '*+'{no,}'dnscrypt[use DNSCrypt for queries]'
  '*+'{no,}'https[use DNS-over-HTTPS for queries]'
  '*+'{no,}'quic[use DNS-over-QUIC for queries]'
  '*+'{no,}'aaonly[set aa flag in the query]'
  '*+'{no,}'additional[print additional section of a reply]'
  '*+'{no,}'adflag[set the AD (authentic data) bit in the query]'
  '*+'{no,}'cdflag[set the CD (checking disabled) bit in the query]'
  '*+'{no,}'cookie[add a COOKIE option to the request]'
  '*+edns=[specify EDNS version for query]:version (0-255)'
  '*+noedns[clear EDNS version to be sent]'
  '*+ednsflags=[set EDNS flags bits]:flags'
  # '*+ednsopt=[specify EDNS option]:code point'
  '*+noedns[clear EDNS options to be sent]'
  '*+'{no,}'expire[send an EDNS Expire option]'
  # '*+'{no,}'idnin[set processing of IDN domain names on input]'
  '*+'{no,}'idnout[set conversion of IDN puny code on output]'
  '*+'{no,}'keepalive[request EDNS TCP keepalive]'
  '*+'{no,}'keepopen[keep TCP socket open between queries]'
  '*+'{no,}'recurse[set the RD (recursion desired) bit in the query]'
  # '*+'{no,}'nssearch[search all authoritative nameservers]'
  # '*+'{no,}'trace[trace delegation down from root]'
  # '*+'{no,}'cmd[print initial comment in output]'
  '*+'{no,}'short[print terse output]'
  '*+'{no,}'identify[print IP and port of responder]'
  '*+'{no,}'comments[print comment lines in output]'
  '*+'{no,}'stats[print statistics]'
  '*+padding[set padding block size]:size [0]'
  '*+'{no,}'qr[print query as it was sent]'
  '*+'{no,}'question[print question section of a query]'
  '*+'{no,}'raflag[set RA flag in the query]'
  '*+'{no,}'answer[print answer section of a reply]'
  '*+'{no,}'authority[print authority section of a reply]'
  '*+'{no,}'all[set all print/display flags]'
  '*+'{no,}'subnet=[send EDNS client subnet option]:addr/prefix-length'
  '*+'{no,}'tcflag[set TC flag in the query]'
  '*+time=[set query timeout]:timeout (seconds) [1]'
  '*+timeout=[set query timeout]:timeout (seconds) [1]'
  '*+tries=[specify number of UDP query attempts]:tries [3]'
  '*+retry=[specify number of UDP query retries]:retries [2]'
  # '*+'{no,}'rrcomments[set display of per-record comments]'
  # '*+ndots=[specify number of dots to be considered absolute]:dots'
  '*+bufsize=[specify UDP buffer size]:size (bytes)'
  '*+'{no,}'dnssec[enable DNSSEC]'
  '*+'{no,}'nsid[include EDNS name server ID request in query]'
  '*+'{no,}'class[display the class whening printing the answer]'
  '*+'{no,}'ttlid[display the TTL whening printing the record]'
  '*+'{no,}'ttlunits[display the TTL in human-readable units]'
  # '*+'{no,}'unknownformat[print RDATA in RFC 3597 "unknown" format]'
  '*+'{no,}'json[present the results as JSON]'
  '*+'{no,}'xml[present the results as XML]'
  '*+'{no,}'yaml[present the results as YAML]'
  '*+'{no,}'zflag[set Z flag in query]'
)
# TODO: Add the regular (POSIX/GNU) flags
_arguments -s -C $args \
  '(- *)-'{h,-help}'[display help information]' \
  '(- *)-'{V,-version}'[display version information]' \
  '-'{v,-verbosity}'=+[set verbosity to custom level]:verbosity:compadd -M "m\:{\-1-3}={\-1-3}" - \-1 0 1 2 3' \
  '-'{v,-verbosity}'+[set verbosity to info]' \
  '*-'{p,-port}'+[specify port number]:port:_ports' \
  '*-'{q,-query}'+[specify host name to query]:host:_hosts' \
  '*-'{c,-class}'+[specify class]:class:compadd -M "m\:{a-z}={A-Z}" - IN CS CH HS' \
  '*-'{t,-qType}'+[specify type]:type:_dns_types' \
  '*-4+[force IPv4 only]' \
  '*-6+[force IPv6 only]' \
  '*-'{x,-reverse}'+[reverse lookup]' \
  '*--timeout+[timeout in seconds]:number [1]' \
  '*--retry+[specify number of UDP query retries]:number [2]' \
  '*--no-edns+[disable EDNS]' \
  '*--edns-ver+[specify EDNS version for query]:version (0-255) [0]' \
  '*-'{D,-dnssec}'+[enable DNSSEC]' \
  '*--expire+[send EDNS expire]' \
  '*-'{n,-nsid}'+[include EDNS name server ID request in query]' \
  '*--no-cookie+[disable sending EDNS cookie]' \
  '*--keep-alive+[request EDNS TCP keepalive]' \
  '*-'{b,-buffer-size}'+[specify UDP buffer size]:size (bytes) [1232]' \
  '*--zflag+[set EDNS z-flag]:decimal, hex or octal [0]' \
  '*--subnet+[set EDNS client subnet]:addr/prefix-length' \
  '*--no-truncate+[ignore truncation in UDP responses]' \
  '*--tcp+[use TCP instead of UDP for queries]' \
  '*--dnscrypt+[use DNSCrypt for queries]' \
  '*-'{T,-tls}'+[use DNS-over-TLS for queries]' \
  '*-'{H,-https}'+[use DNS-over-HTTPS for queries]' \
  '*-'{Q,-quic}'+[use DNS-over-QUIC for queries]' \
  '*--tls-no-verify+[disable TLS verification]' \
  '*--tls-host+[set TLS lookup hostname]:host:_hosts' \
  '*-'{s,-short}'+[print terse output]' \
  '*-'{j,-json}'+[present the results as JSON]' \
  '*-'{x,-xml}'+[present the results as XML]' \
  '*-'{y,-yaml}'+[present the results as YAML]' \
  '*: :->args' && ret=0

if [[ -n $state ]]; then
  if compset -P @; then
    _wanted hosts expl 'DNS server' _hosts && ret=0;
  else
    _alternative 'hosts:host:_hosts' 'types:query type:_dns_types' && ret=0
  fi
fi

return ret