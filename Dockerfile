FROM node:22-alpine AS web-builder

WORKDIR /src/web
COPY web/package*.json ./
RUN npm ci
COPY web ./
RUN npm run build

FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
COPY --from=web-builder /src/internal/server/web/dist ./internal/server/web/dist
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/schemaforge ./cmd/schemaforge

FROM alpine:3.20

RUN adduser -D -H -u 10001 schemaforge
COPY --from=builder /out/schemaforge /usr/local/bin/schemaforge
USER schemaforge
EXPOSE 8989

ENTRYPOINT ["schemaforge"]
CMD ["ui", "--addr", "0.0.0.0:8989"]
