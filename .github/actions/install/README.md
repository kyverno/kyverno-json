# kyverno-json-installer GitHub Action

This action enables you to install `kyverno-json`.

For a quick start guide on the usage of `kyverno-json`, please refer to https://kyverno.github.io/kyverno-json.

# Usage

This action currently supports GitHub-provided Linux, macOS and Windows runners (self-hosted runners may not work).

Add the following entry to your Github workflow YAML file:

```yaml
uses: kyverno/kyverno-json/.github/actions/install@main
with:
  release: 'v0.0.1' # optional
```

Example using a pinned version:

```yaml
jobs:
  example:
    runs-on: ubuntu-latest

    permissions: {}

    name: Install kyverno-json
    steps:
      - name: Install kyverno-json
        uses: kyverno/kyverno-json/.github/actions/install@main
        with:
          release: 'v0.0.1'
      - name: Check install
        run: kyverno-json version
```

Example using the default version:

```yaml
jobs:
  example:
    runs-on: ubuntu-latest

    permissions: {}

    name: Install kyverno-json
    steps:
      - name: Install kyverno-json
        uses: kyverno/kyverno-json/.github/actions/install@main
      - name: Check install
        run: kyverno-json version
```

Example using [cosign](https://github.com/sigstore/cosign) verification:

```yaml
jobs:
  example:
    runs-on: ubuntu-latest

    permissions: {}

    name: Install kyverno-json
    steps:
      - name: Install Cosign
        uses: sigstore/cosign-installer@v3.1.1
      - name: Install kyverno-json
        uses: kyverno/kyverno-json/.github/actions/install@main
        with:
          verify: true
      - name: Check install
        run: kyverno-json version
```

If you want to install `kyverno-json` from its main version by using `go install` under the hood, you can set `release` as `main`.
Once you did that, `kyverno-json` will be installed via `go install` which means that please ensure that go is installed.

Example of installing `kyverno-json` via `go install`:

```yaml
jobs:
  example:
    runs-on: ubuntu-latest

    permissions: {}

    name: Install kyverno-json via go install
    steps:
      - name: Install go
        uses: actions/setup-go@v4
        with:
          go-version: '1.20'
          check-latest: true
      - name: Install kyverno-json
        uses: kyverno/kyverno-json/.github/actions/install@main
        with:
          release: main
      - name: Check install
        run: kyverno-json version
```

### Optional Inputs

The following optional inputs:

| Input | Description |
| --- | --- |
| `release` | `kyverno-json` version to use instead of the default. |
| `install-dir` | directory to place the `kyverno-json` binary into instead of the default (`$HOME/.kyverno-json`). |
| `use-sudo` | set to `true` if `install-dir` location requires sudo privs. Defaults to false. |
| `verify` | set to `true` to enable [cosign](https://github.com/sigstore/cosign) verification of the downloaded archive. |

## Security

Should you discover any security issues, please refer to Kyverno's [security process](https://github.com/kyverno/kyverno/blob/main/SECURITY.md)