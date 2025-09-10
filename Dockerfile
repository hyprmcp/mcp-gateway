FROM golang:1.25 AS builder
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY main.go main.go
COPY cmd/ cmd/
COPY config/ config/
COPY htmlresponse/ htmlresponse/
COPY jsonrpc/ jsonrpc/
COPY log/ log/
COPY oauth/ oauth/
COPY proxy/ proxy/
COPY webhook/ webhook/
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
  go build -a -o mcp-gateway -ldflags="-s -w" .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=builder /workspace/mcp-gateway .
USER 65532:65532
ENTRYPOINT [ "/mcp-gateway" ]
