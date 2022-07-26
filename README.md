# awl

[![Build Status](https://ci.git.froth.zone/api/badges/sam/awl/status.svg)](https://ci.git.froth.zone/sam/awl)

`awl` is a command-line DNS client, much like
[`drill`](https://github.com/NLnetLabs/ldns),
[`dig`](https://bind9.readthedocs.io/en/v9_18_3/manpages.html#dig-dns-lookup-utility),
[`dog`](https://github.com/ogham/dog),
[`doggo`](https://github.com/mr-karan/doggo), or
[`q`](https://github.com/natesales/q).

`awl` is designed to be a drop-in replacement for the venerable dig, but support
newer RFC query types, such as DNS-over-HTTPS and DNS-over-QUIC.

## Usage

- [Feature wiki](https://git.froth.zone/sam/awl/wiki/Supported)
- [Manpage](https://git.froth.zone/sam/awl/wiki/awl.1)


## Building and installing

### From releases

Grab a prebuilt binary from the
[release](https://git.froth.zone/sam/awl/releases) section.

### From source

Dependencies:

- Go >= 1.18
- GNU/BSD make or Plan 9 mk
- [scdoc](https://git.sr.ht/~sircmpwn/scdoc) (optional, for manpage)

Make sure to recursively clone the repo:

```sh
git clone --recursive https://git.froth.zone/sam/awl
```

Using the makefile:

```sh
make
sudo make install
```

Alternatively, using `go install`:

```sh
go install git.froth.zone/sam/awl@latest
```

## Contributing

Send a [pull request](https://git.froth.zone/sam/awl/pulls) our way. Prefer
emails? Send a patch to the
[mailing list](https://lists.sr.ht/~sammefishe/awl-dev).

Found a bug or want a new feature? Create an issue
[here](https://git.froth.zone/sam/awl/issues).

### License

See [LICENSE](./LICENSE)
