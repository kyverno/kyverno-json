# Install

You can install the pre-compiled binary (in several ways), or compile from source.
We also provide a [GitHub action](#github-action) to easily install Kyverno-JSON in your workflows.

## Install the pre-compiled binary

### Homebrew tap

**add tap:**

```bash
brew tap kyverno/kyverno-json https://github.com/kyverno/kyverno-json
```

**install kyverno-json:**

```bash
brew install kyverno/kyverno-json/kyverno-json
```

### Manually

Download the pre-compiled binaries for your system from the [releases page](https://github.com/kyverno/kyverno-json/releases) and copy them to the desired location.

## Using `go install`

You can install with `go install` with:

```bash
go install github.com/kyverno/kyverno-json@latest
```
## Running with Docker 

Kyverno-JSON is also available as a Docker image which you can pull and run:

```bash
docker pull ghcr.io/kyverno/kyverno-json:<version>
```

!!! info

    Since kyverno-JSON relies on files for its operation (like ValidatingPolicy definitions), you will need to bind mount the necessary directories when running it via Docker.


```bash
$ docker run --rm                       \
    -v /path/on/host:/path/in/container \
    ghcr.io/kyverno/kyverno-json:<version>  \
    <kyverno-json-command>
```


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