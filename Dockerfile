FROM golang:1.22.2-alpine3.18 AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go get gorm.io/driver/postgres gorm.io/gorm

RUN go build -o /app/meeting-api ./cmd/api

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /app/meeting-api /app/meeting-api 
EXPOSE 8080

CMD ["/app/meeting-api"]

