OS      = linux
ARCH    = amd64
OUTPUT  = http_jwt_crud

run:
	swag init -g api/handlers/handlers.go --markdownFiles docs
	go run cmd/http_jwt_crud/main.go

build:
	swag init -g api/handlers/handlers.go --markdownFiles docs
	GOOS=${OS} GOARCH=${ARCH} CGO_ENABLED=0 go build --ldflags=${LDFLAGS} -v -buildvcs=false -o ${OUTPUT} cmd/http_jwt_crud/main.go

test:
	cd tests; go test -v -count=1