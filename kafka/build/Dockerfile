FROM golang:1.16 as builder

ENV GO111MODULE="on"
ENV CGO_ENABLED="0"
ENV GOOS="linux"

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN go build -o ./app ./cmd/app

FROM gcr.io/distroless/static
COPY --from=builder /app/app /usr/bin/app
ENTRYPOINT ["/usr/bin/app"]
