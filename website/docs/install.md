# Install

You can install the pre-compiled binary (in several ways), or compile from source.

## Install using `go install`

You can install with `go install` with:

```bash
go install github.com/kyverno/kyverno-json@latest
```

## Manually

Download the pre-compiled binaries from the [releases page](https://github.com/kyverno/kyverno-json/releases) and copy them to the desired location.

## Build from the source code

**clone the repository:**

```bash
git clone github.com/kyverno/kyverno-json
```

**build the binaries:**

```bash
cd kyverno-json
```

```bash
make build
```

**verify the build:**

```bash
./kyverno-json version
```