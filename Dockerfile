FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN go build -ldflags="-s -w" -o /app/ui-platform-backend-service ./cmd/app/main.go


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Greenwich /usr/share/zoneinfo/Greenwich
ENV TZ Greenwich

WORKDIR /app
COPY --from=builder /app/ui-platform-backend-service /app/ui-platform-backend-service

EXPOSE 8080

CMD ["./ui-platform-backend-service"]