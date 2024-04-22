FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./lawang-project ./main.go
 
 
FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/lawang-project .
EXPOSE 8080
ENTRYPOINT ["./lawang-project"]