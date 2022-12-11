# Build the manager binary
FROM golang:1.19 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY config/ config/
COPY kube/ kube/
COPY score/ score/
COPY weather/ weather/

RUN CGO_ENABLED=0 go build -a -o ecos main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM ubuntu:jammy
WORKDIR /
COPY --from=builder /workspace/ecos .

# USER 65532:65532
RUN apt update && apt install -y curl

ENTRYPOINT ["/ecos"]
