FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod ./
COPY . .

RUN go build -o /out/backend ./cmd

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /out/backend /app/backend
EXPOSE 8080
ENV PORT=8080
CMD ["/app/backend"]
