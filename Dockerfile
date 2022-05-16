# Multistaging building

ARG GO_VERSION=1.16.6

FROM golang:${GO_VERSION}-alpine as builder
RUN go env -w GOPROXY=direct
RUN apk add --no-cache git
RUN apk --no-cache add ca-certificates && update-ca-certificates

WORKDIR /src

COPY ./ ./
RUN go mod download

RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /app main.go

FROM scratch AS runner
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY .env ./
COPY --from=builder /app /app

EXPOSE 5050
ENTRYPOINT ["/app"]