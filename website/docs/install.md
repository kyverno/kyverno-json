# Install

You can install the pre-compiled binary (in several ways), or compile from source.

## Using `go install`

You can install with `go install` with:

```bash
go install github.com/kyverno/kyverno-json@latest
```

## Download binary

Download the pre-compiled binaries from the [releases page](https://github.com/kyverno/kyverno-json/releases) and copy them to the desired location.

## Build from the source code

**clone the repository:**

```bash
git clone https://github.com/kyverno/kyverno-json.git
```

**build the binaries:**

```bash
cd kyverno-json
go mod tidy
make build
```

**verify the build:**

```bash
./kyverno-json version
```