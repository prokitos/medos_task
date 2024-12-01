FROM golang:alpine AS builder

WORKDIR /build
ADD go.mod .
COPY . .
COPY /config /tempConfig
RUN go build -o . cmd/main.go

FROM alpine
WORKDIR /build
COPY --from=builder /build/main /build/main
COPY --from=builder /tempConfig /build/config
CMD ["./main"]