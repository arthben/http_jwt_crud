OS      = linux
ARCH    = amd64
OUTPUT  = http_jwt_crud

run:
	swag init -g api/handlers/handlers.go --markdownFiles docs
	go run cmd/http_jwt_crud/main.go

build:
	swag init -g api/handlers/handlers.go --markdownFiles docs
	CC=${CC} GOOS=${OS} GOARCH=${ARCH} CGO_ENABLED=1 go build --ldflags=${LDFLAGS} -v -buildvcs=false -o ${OUTPUT} cmd/http_jwt_crud/main.go