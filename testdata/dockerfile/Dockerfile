ARG BUILD_PLATFORM="linux/amd64"
ARG BUILDER_IMAGE="golang:1.20.6-alpine3.18"

FROM --platform=$BUILD_PLATFORM $BUILDER_IMAGE as builder

WORKDIR /
COPY . ./

# Get Signer plugin binary
ARG SIGNER_BINARY_LINK="https://d2hvyiie56hcat.cloudfront.net/linux/amd64/plugin/latest/notation-aws-signer-plugin.zip"
ARG SIGNER_BINARY_FILE="notation-aws-signer-plugin.zip"
RUN wget -O ${SIGNER_BINARY_FILE} ${SIGNER_BINARY_LINK} 
RUN apk update && \
    apk add unzip && \
    unzip -o ${SIGNER_BINARY_FILE}

# Build Go binary
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o kyverno-notation-aws .

FROM gcr.io/distroless/static:nonroot
WORKDIR /

# Notation home
ENV PLUGINS_DIR=/plugins

COPY --from=builder notation-com.amazonaws.signer.notation.plugin plugins/com.amazonaws.signer.notation.plugin/notation-com.amazonaws.signer.notation.plugin

COPY --from=builder kyverno-notation-aws kyverno-notation-aws
ENTRYPOINT ["/kyverno-notation-aws"]
