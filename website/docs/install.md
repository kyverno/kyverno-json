# Install

You can install the pre-compiled binary (in several ways), or compile from source.

## Install using `go install`

You can install with `go install` with:

```bash
go install github.com/kyverno/kyverno-json@latest
```

## Manually

```
Download the pre-compiled binaries from the [releases page](https://github.com/kyverno/kyverno-json/releases) and copy them to the desired location.
```

## Compile from sources

**clone:**

```bash
git clone github.com/kyverno/kyverno-json
```

**get the dependencies:**

```bash
go mod tidy
```

**build:**

```bash
make build
```

**verify it works:**

```bash
./kyverno-json version
```