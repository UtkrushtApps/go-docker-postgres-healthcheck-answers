FROM golang:1.21 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/api /app
EXPOSE 8080
# Pass environment variables at container runtime
ENTRYPOINT ["/app/api"]