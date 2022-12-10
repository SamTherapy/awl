# SPDX-License-Identifier: BSD-3-Clause
function __fish_complete_awl
    set -l token (commandline -ct)
    switch $token
        case '+tries=*' '+retry=*' '+time=*' '+bufsize=*' '+edns=*'
            printf '%s\n' $token(seq 0 255)
        case '-v=*'
            printf '%s\n' $token(seq -1 3)
    end
end

complete -c awl -x -a "(__fish_print_hostnames) A AAAA AFSDB APL CAA CDNSKEY CDS CERT CNAME DHCID DLV DNAME DNSKEY DS HIP IPSECKEY KEY KX LOC MX NAPTR NS NSEC NSEC3 NSEC3PARAM PTR RRSIG RP SIG SOA SRV SSHFP TA TKEY TLSA TSIG TXT URI"
complete -c awl -x -a "@(__fish_print_hostnames)"

complete -f -c awl -s 4 -d 'Use IPv4 query transport only'
complete -f -c awl -s 6 -d 'Use IPv6 query transport only'
complete -c awl -s c -l class -x -a 'IN CH HS QCLASS' -d 'Specify query class'
complete -c awl -s p -l port  -x -d 'Specify port number'
complete -c awl -s q -l query -x -a "(__fish_print_hostnames)" -d 'Query domain'
complete -c awl -s t -l qType -x -a 'A AAAA AFSDB APL CAA CDNSKEY CDS CERT CNAME DHCID DLV DNAME DNSKEY DS HIP IPSECKEY KEY KX LOC MX NAPTR NS NSEC NSEC3 NSEC3PARAM PTR RRSIG RP SIG SOA SRV SSHFP TA TKEY TLSA TSIG TXT URI' -d 'Specify query type'
complete -c awl -l timeout -x -d 'Set timeout'
complete -c awl -l retries -x -d 'Set number of query retries'
complete -c awl -l no-edns -x -d 'Disable EDNS'
complete -f -c awl -l tcp -a '+vc +novc +tcp +notcp' -d 'TCP mode'
complete -f -c awl -l dnscrypt -a '+dnscrypt +nodnscrypt' -d 'Use DNSCrypt'
complete -c awl -s T -l tls -a '+tls +notls' -d 'Use DNS-over-TLS'
complete -c awl -s H -l https -a '+https +nohttps' -d 'Use DNS-over-HTTPS'
complete -c awl -s Q -l quic -a '+quic +noquic'  -d 'Use DNS-over-QUIC'

complete -c awl -s j -l json -a '+json +nojson' -d 'Print as JSON'
complete -c awl -s j -l xml -a '+xml +noxml' -d 'Print as XML'
complete -c awl -s j -l yaml -a '+yaml +noyaml' -d 'Print as YAML'

complete -c awl -s x -l reverse -x -d 'Reverse lookup'
complete -f -c awl -s h -l help -d 'Print help and exit'
complete -f -c awl -s V -l version -d 'Print version and exit'


# complete -f -c awl -a '+search +nosearch' -d 'Set whether to use searchlist'
# complete -f -c awl -a '+showsearch +noshowsearch' -d 'Search with intermediate results'
complete -f -c awl -a '+recurse +norecurse' -d 'Recursive mode'
complete -f -c awl -l no-truncate -a '+ignore +noignore' -d 'Dont revert to TCP for TC responses.'
# complete -f -c awl -a '+fail +nofail' -d 'Dont try next server on SERVFAIL'
# complete -f -c awl -a '+besteffort +nobesteffort' -d 'Try to parse even illegal messages'
complete -f -c awl -a '+aaonly +noaaonly' -d 'Set AA flag in query (+[no]aaflag)'
complete -f -c awl -a '+adflag +noadflag' -d 'Set AD flag in query'
complete -f -c awl -a '+cdflag +nocdflag' -d 'Set CD flag in query'
complete -f -c awl -a '+cl +nocl' -d 'Control display of class in records'
# complete -f -c awl -a '+cmd +nocmd' -d 'Control display of command line'
complete -f -c awl -a '+comments +nocomments' -d 'Control display of comment lines'
complete -f -c awl -a '+question +noquestion' -d 'Control display of question'
complete -f -c awl -a '+answer +noanswer' -d 'Control display of answer'
complete -f -c awl -a '+authority +noauthority' -d 'Control display of authority'
complete -f -c awl -a '+additional +noadditional' -d 'Control display of additional'
complete -f -c awl -a '+stats +nostats' -d 'Control display of statistics'
complete -f -c awl -s s -l short -a '+short +noshort' -d 'Disable everything except short form of answer'
complete -f -c awl -a '+ttlid +nottlid' -d 'Control display of ttls in records'
complete -f -c awl -a '+all +noall' -d 'Set or clear all display flags'
complete -f -c awl -a '+qr +noqr' -d 'Print question before sending'
# complete -f -c awl -a '+nssearch +nonssearch' -d 'Search all authoritative nameservers'
complete -f -c awl -a '+identify +noidentify' -d 'ID responders in short answers'
complete -f -c awl -a '+trace +notrace' -d 'Trace delegation down from root'
complete -f -c awl -l dnssec -a '+dnssec +nodnssec +do +nodo' -d 'Request DNSSEC records'
complete -f -c awl -a '+nsid +nonsid' -d 'Request Name Server ID'
# complete -f -c awl -a '+multiline +nomultiline' -d 'Print records in an expanded format'
# complete -f -c awl -a '+onesoa +noonesoa' -d 'AXFR prints only one soa record'

complete -f -c awl -a '+tries=' -d 'Set number of UDP attempts'
complete -f -c awl -a '+retry=' -d 'Set number of UDP retries'
complete -f -c awl -a '+time=' -d 'Set query timeout'
complete -f -c awl -a '+bufsize=' -d 'Set EDNS0 Max UDP packet size'
complete -f -c awl -a '+ndots=' -d 'Set NDOTS value'
complete -f -c awl -a '+edns=' -d 'Set EDNS version'

complete -c awl -a '(__fish_complete_awl)'
