# awl

`awl` is a command-line DNS client, much like
[`drill`](https://github.com/NLnetLabs/ldns),
[`dig`](https://bind9.readthedocs.io/en/v9_18_3/manpages.html#dig-dns-lookup-utility),
[`dog`](https://github.com/ogham/dog),
[`doggo`](https://github.com/mr-karan/doggo),
or [`q`](https://github.com/natesales/q)

This was made as my first major experiment with Go, so there are probably things that can be improved

The excellent [dns](https://github.com/miekg/dns) library for Go does most of the heavy
lifting.

## What works

- UDP
- TCP
- TLS
- HTTPS (maybe)
- QUIC (extreme maybe)

## What doesn't

- Your sanity after reading my awful code
- A motivation for making this after finding q and doggo

## What should change

- Make the CLI less abysmal (migrate to [cobra](https://github.com/spf13/cobra)?
  or just use stdlib's flags)
- Optimize everything
- Make the code less spaghetti
- Documentation, documentation, documentation