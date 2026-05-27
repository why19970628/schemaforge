FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/schemaforge ./cmd/schemaforge

FROM alpine:3.20

RUN adduser -D -H -u 10001 schemaforge
COPY --from=builder /out/schemaforge /usr/local/bin/schemaforge
USER schemaforge
EXPOSE 8989

ENTRYPOINT ["schemaforge"]
CMD ["ui", "--addr", "0.0.0.0:8989"]
