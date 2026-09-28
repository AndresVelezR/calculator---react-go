FROM golang:1.27.1-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /calculator ./cmd/server

FROM alpine:3.23

RUN addgroup -S calculator && adduser -S calculator -G calculator
COPY --from=build /calculator /usr/local/bin/calculator
USER calculator
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/calculator"]
