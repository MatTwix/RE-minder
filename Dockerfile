FROM golang:1.23-alpine AS server-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /main

FROM alpine:3.19
WORKDIR /app
COPY --from=server-builder /main .
COPY .env* ./
EXPOSE ${PORT:-3000}
CMD [ "./main" ]