# vw

`vw` is an opinionated command-line interface for Bitwarden and Vaultwarden.

It wraps the official Bitwarden CLI (`bw`) and provides a friendlier workflow for common terminal usage:

```sh
vw setup --server https://vaultwarden.example.com
vw login
vw unlock
vw sync
vw copy github
vw code github
vw new github --user me@example.com --url https://github.com
```

`vw` does not reimplement the Bitwarden/Vaultwarden protocol or cryptography. It delegates that work to `bw`.

## Why

The official `bw` CLI is powerful, but many day-to-day actions require JSON, item IDs, templates, `bw encode`, and extra shell scripting.

`vw` provides a smaller, more convenient interface for common workflows:

```sh
vw copy github
vw get github
vw user github
vw url github
vw code github
vw new github --user me@example.com --url https://github.com
```

## Relationship to `bw`

`vw` uses the official Bitwarden CLI internally.

This means `vw` works with:

- Bitwarden cloud
- Self-hosted Bitwarden
- Vaultwarden

For Vaultwarden, configure the server URL during setup:

```sh
vw setup --server https://vaultwarden.example.com
```

Under the hood, `vw` calls commands such as:

```sh
bw config server ...
bw login
bw unlock
bw sync
bw list items
bw get item
bw get password
bw get totp
bw create item
```

## Installing from a bundled release

A bundled release contains both:

```text
vw
bw
```

You do not need to install `bw` separately when using a bundled release.

Extract the archive:

```sh
tar -xzf vw-*.tar.gz
cd vw-*
```

Install to the default prefix:

```sh
scripts/install.sh
```

By default this installs:

```text
~/.local/bin/vw
~/.local/bin/bw
```

Make sure this directory is in your `PATH`:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

To install somewhere else:

```sh
scripts/install.sh --prefix /usr/local
```

To install only `vw` and skip the bundled `bw`:

```sh
scripts/install.sh --no-bw
```

For non-interactive installation:

```sh
scripts/install.sh --yes
```

## Installing from source

Clone the repository:

```sh
git clone https://github.com/comfucios/vw.git
cd vw
```

Build:

```sh
make build
```

Install:

```sh
mkdir -p ~/.local/bin
cp bin/vw ~/.local/bin/vw
```

Then either install `bw` yourself or let `vw` download a managed copy:

```sh
vw bootstrap-bw
```

## Managed Bitwarden CLI

`vw` can use a managed copy of `bw`.

To download it:

```sh
vw bootstrap-bw
```

Managed `bw` is stored under:

```text
~/.local/share/vw/bin/bw
```

`vw` resolves the `bw` binary in this order:

```text
1. configured bw path
2. managed/bundled bw
3. bw found in PATH
```

You can check what `vw` sees with:

```sh
vw doctor
```

## Initial setup

For Vaultwarden:

```sh
vw setup --server https://vaultwarden.example.com
```

For Bitwarden cloud:

```sh
vw setup
```

Then log in:

```sh
vw login
```

Unlock your vault:

```sh
vw unlock
```

Sync:

```sh
vw sync
```

## Common commands

Copy a password:

```sh
vw copy github
```

Copy a username:

```sh
vw copy github --field user
```

Copy a URL:

```sh
vw copy github --field url
```

Copy a TOTP code:

```sh
vw copy github --field code
```

Print a password:

```sh
vw get github
```

Print a username:

```sh
vw user github
```

Print a URL:

```sh
vw url github
```

Print a TOTP code:

```sh
vw code github
```

List matching items:

```sh
vw list github
```

Create a new login item:

```sh
vw new github --user me@example.com --url https://github.com
```

Create a new login item interactively:

```sh
vw new
```

Lock the vault/session:

```sh
vw lock
```

## Configuration

Show the config path:

```sh
vw config path
```

Show config values:

```sh
vw config get
```

Set a config value:

```sh
vw config set server_url https://vaultwarden.example.com
```

Supported environment variables use the `VW_` prefix.

Examples:

```sh
VW_SERVER_URL=https://vaultwarden.example.com
VW_BW_PATH=/custom/path/to/bw
```

## Uninstalling

From an extracted release bundle:

```sh
scripts/uninstall.sh
```

This removes:

```text
~/.local/bin/vw
~/.local/bin/bw
~/.local/share/vw/bin/bw
```

User config and cache data are kept by default.

To remove config, cache, state, and managed data too:

```sh
scripts/uninstall.sh --purge
```

For non-interactive uninstall:

```sh
scripts/uninstall.sh --purge --yes
```

If you use the OS keyring session backend, run this before uninstalling:

```sh
vw lock
```

Otherwise, remove the `vw` keyring entry manually from your system keychain/keyring.

## Development

Build:

```sh
make build
```

Test:

```sh
make test
```

Run:

```sh
go run ./cmd/vw --help
```

Format:

```sh
make fmt
```

Vet:

```sh
make vet
```

## Creating a bundled package

Build `vw` first:

```sh
make build
```

Place a compatible `bw` binary at:

```text
bin/bw
```

or pass it explicitly:

```sh
BW_PATH=/path/to/bw VERSION=0.1.0 scripts/package-with-bw.sh
```

The bundled archive will be created under:

```text
dist/
```

The bundled archive layout is:

```text
vw-<version>-<os>-<arch>/
  bin/
    vw
    bw
  scripts/
    install.sh
    uninstall.sh
  README.md
  LICENSE
  NOTICE
  VERSION
  BUNDLE
```

## CI

The repository uses GitHub Actions for CI.

On pull requests and pushes to `main`, CI runs:

```sh
go mod download
gofmt check
go vet ./...
go test -race ./...
go build ./cmd/vw
```

You can run the same checks locally with:

```sh
make ci
```

## Releases

Releases are automated with Release Please and GoReleaser.

The release flow is:

```text
1. Merge Conventional Commit messages into main
2. Release Please opens or updates a release PR
3. Merge the release PR
4. Release Please creates a GitHub release and tag
5. GoReleaser builds binaries and uploads release artifacts
```

Use Conventional Commits:

```text
feat: add interactive item editor
fix: handle missing bw session
docs: update install instructions
chore: update dependencies
```

`feat:` creates a minor release.

`fix:` creates a patch release.

Breaking changes create a major release when marked with `!` or a `BREAKING CHANGE:` footer.

Plain vw releases are built by GoReleaser:

```sh
goreleaser release --snapshot --clean
```

or

```sh
make snapshot
```

Bundled releases that include `bw` are produced separately with:

```sh
VERSION=0.1.0 BW_PATH=/path/to/bw scripts/package-with-bw.sh
```

This distinction is intentional:

```sh
GoReleaser -> builds vw release artifacts
package-with-bw.sh -> creates optional bundles containing vw + bw
```

## Security notes

`vw` is a convenience wrapper around `bw`.

It does not implement Bitwarden cryptography itself.

Session handling should be treated carefully. Use:

```sh
vw lock
```

when you want to clear the active session.

Do not commit secrets, config files containing credentials, or generated release bundles to the repository.

## License

Apache License 2.0. See [LICENSE](LICENSE) and [NOTICE](NOTICE).
