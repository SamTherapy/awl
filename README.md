<!-- markdownlint-disable MD033 -->
# <img src="./docs/img/awl-text.png" width="50%" title="awl logo" alt="awl">

> awl *(noun)*: A pointed tool for making small holes in wood or leather

A command-line DNS lookup tool that supports DNS queries over UDP, TCP, TLS, HTTPS, DNSCrypt, and QUIC.

[![Gitea Release](https://img.shields.io/gitea/v/release/sam/awl?gitea_url=https%3A%2F%2Fgit.froth.zone&display_name=release&style=for-the-badge)](https://git.froth.zone/sam/awl)
[![Last Commit](https://img.shields.io/gitea/last-commit/sam/awl?gitea_url=https%3A%2F%2Fgit.froth.zone&style=for-the-badge)](https://git.froth.zone/sam/awl/commits/branch/master)
[![License](https://img.shields.io/github/license/samtherapy/awl?style=for-the-badge)](https://spdx.org/licenses/BSD-3-Clause.html)
[![Go Report](https://goreportcard.com/badge/dns.froth.zone/awl?style=for-the-badge)](https://goreportcard.com/report/dns.froth.zone/awl)

Awl is designed to be a drop-in replacement for [dig](https://bind9.readthedocs.io/en/v9_18_3/manpages.html#dig-dns-lookup-utility).

## Examples

```shell
# Query a domain over UDP
awl example.com

# Query a domain over HTTPS, print only the results
awl example.com +https --short

# Query a domain over TLS, print as JSON
awl example.com +tls +json
```

For more and the usage, see the [manpage](https://git.froth.zone/sam/awl/wiki/awl.1).

## Installing

On any platform, with [Go](https://go.dev) installed, run the following command to install:

```shell
go install dns.froth.zone/awl@latest
```

### Packaging

Alternatively, many package managers are supported:

<details>
<summary>Linux</summary>

#### Distro-specific

<details>
<summary>Alpine Linux</summary>

Provided by [Gitea packages](https://git.froth.zone/sam/-/packages/alpine/awl-dns) \
***Any distro that uses apk should also work***

```shell
# Add the repository
echo "https://git.froth.zone/api/packages/sam/alpine/edge/main" | tee -a /etc/apk/repositories
# Get the signing key
curl -JO https://git.froth.zone/api/packages/sam/alpine/key --output-dir /etc/apk/keys
# Install
apk add awl-dns
```

</details>

<details>
<summary>Arch</summary>

AUR package available as [awl-dns-git](https://aur.archlinux.org/packages/awl-dns-git/)

```shell
yay -S awl-dns-git ||
paru -S awl-dns-git
```

</details>

<details>
<summary>Debian / Ubuntu</summary>

Provided by [Gitea packages](https://git.froth.zone/sam/-/packages/debian/awl-dns/) \
***Any distro that uses deb/dpkg should also work***

```shell
# Install the repository and GPG keys
curl -JO https://git.froth.zone/packaging/-/packages/debian/git-froth-zone-debian/1-0/files/5937
sudo dpkg -i git-froth-zone-debian_1-0_all.deb
rm git-froth-zone-debian_1-0_all.deb
# Update and install
sudo apt update
sudo apt install awl-dns
```

</details>

<details>
<summary>Fedora / RHEL / SUSE</summary>

Provided by [Gitea packages](https://git.froth.zone/sam/-/packages/rpm/awl-dns/) \
***Any distro that uses rpm/dnf might also work, I've never tried it***

```shell
# Add the repository
dnf config-manager --add-repo https://git.froth.zone/api/packages/sam/rpm.repo ||
zypper addrepo https://git.froth.zone/api/packages/sam/rpm.repo
# Install
dnf install awl-dns ||
zypper install awl-dns
```

</details>

<details>
<summary>Gentoo</summary>

```shell
# Add the ebuild repository
eselect repository add froth-zone git https://git.froth.zone/packaging/portage.git
emaint sync -r froth-zone
# Install
emerge -av net-dns/awl
```

</details>

#### Distro-agnostic


<details>
<summary><a href="https://brew.sh" nofollow>Homebrew</a></summary>

```shell
brew install SamTherapy/tap/awl
```

</details>
<details>
<summary>Snap</summary>

Snap package available as [awl-dns](https://snapcraft.io/awl-dns)

```shell
snap install awl-dns ||
sudo snap install awl-dns
```

</details>
</details>
<hr />
<details>
<summary>macOS</summary>

<details open>
<summary><a href="https://brew.sh" nofollow>Homebrew</a></summary>

```shell
brew install SamTherapy/tap/awl
```

</details>
</details>
<hr />
<details>
<summary>Windows</summary>

<details open>
<summary><a href="https://scoop.sh" nofollow>Scoop</a></summary>

```pwsh
scoop bucket add froth https://git.froth.zone/packages/scoop.git
scoop install awl
```

</details>
</details>

## Contributing

Please see the [CONTRIBUTING.md](./docs/CONTRIBUTING.md) file for more information.

TL;DR: If you like the project, spread the word! If you want to contribute, [use the issue tracker](https://git.froth.zone/sam/awl/issues) or [open a pull request](https://git.froth.zone/sam/awl/pulls).
Want to use email instead? Use our [mailing list](https://lists.sr.ht/~sammefishe/awl-devel)!

### Mirrors

The canonical repository is located on [my personal Forgejo instance](https://git.froth.zone/sam/awl). \
Official mirrors are located on [GitHub](https://github.com/SamTherapy/awl) and [SourceHut](https://git.sr.ht/~sammefishe/awl/).
Contributions are accepted on all mirrors, but the Forgejo instance is preferred.

## License

[BSD-3-Clause](https://spdx.org/licenses/BSD-3-Clause.html)

### Credits

- Awl image taken from [Wikimedia Commons](https://commons.wikimedia.org/wiki/File:Awl.tif), imaged is licensed CC0.
