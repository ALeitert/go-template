## Build binary
FROM golang:1.24 as builder
WORKDIR /app
COPY ./ ./

RUN CGO_ENABLED=0 go build -o ./go-template -ldflags "-extldflags '-static'" cmd/main.go
RUN chmod 777 ./go-template

## Output Runner
FROM alpine
WORKDIR /app
COPY --from=builder /app/go-template /app/go-template
ENTRYPOINT /app/go-template
