FROM golang:1.20.4-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /app/authenticationExec ./cmd/api

# Final stage
FROM alpine:latest

RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/authenticationExec .

CMD ["./authenticationExec"]