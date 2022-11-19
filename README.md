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

### Package Managers

- AUR: [awl-dns-git](https://aur.archlinux.org/packages/awl-dns-git)
- Debian/Ubuntu:

  ```sh
  # Add PGP key
  curl -sL https://packages.freecumextremist.finance/key.asc | sudo tee /usr/share/keyrings/awl.asc >/dev/null
  # Add repo
  echo "deb [signed-by=/usr/share/keyrings/awl.asc] https://packages.freecumextremist.finance/awl/apt awl main" | sudo tee /etc/apt/sources.list.d/awl.list >/dev/null
  sudo apt update
  sudo apt install awl
  ```

- Homebrew:

  ```sh
  brew install SamTherapy/tap/awl
  ```

- Scoop:

  ```pwsh
  scoop bucket add froth https://git.froth.zone/sam/scoop.git
  scoop install awl
  ```

### From source

Dependencies:

- [Go](https://go.dev/) >= 1.18
- GNU/BSD make or Plan 9 mk (if using the makefile/mkfile)
- [scdoc](https://git.sr.ht/~sircmpwn/scdoc) (optional, for man page)

Using `go install`:

```sh
go install git.froth.zone/sam/awl@latest
```

Using the makefile:

```sh
make && sudo make install
```

## Contributing

Send a [pull request](https://git.froth.zone/sam/awl/pulls) our way. Prefer
emails? Send a patch to the
[mailing list](https://lists.sr.ht/~sammefishe/awl-dev).

Found a bug or want a new feature? Create an issue
[here](https://git.froth.zone/sam/awl/issues).

### Licence

Revised BSD, See [LICENCE](./LICENCE)
