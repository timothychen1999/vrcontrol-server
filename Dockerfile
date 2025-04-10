FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /app/
#Copy static files
COPY --from=builder /app/static /app/static

EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["./main"]