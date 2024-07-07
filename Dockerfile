FROM golang:1.22.5-alpine3.19 AS builder

WORKDIR /build
COPY . .
RUN go mod download
RUN go build -o ./http_jwt_crud ./cmd/http_jwt_crud/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/http_jwt_crud ./http_jwt_crud
COPY --from=builder /build/configs ./configs

CMD ["/http_jwt_crud"]

