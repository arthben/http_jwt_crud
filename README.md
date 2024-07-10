# http_jwt_crud

## This repo demonstrate
- Go HTTP ( with middleware )
- Structure Log ( slog package )
- JWT Authentication ( Bearer )
- Signature to validate request
- Simple CRUD with SQLX package

## Table structure
```sql
CREATE TABLE todos (
	id varchar(64),
	title varchar(256),
	detail text,
	created_date timestamp,
	updated_date timestamp,
	st_completed char(1),
	completed_date timestamp,
	primary key(id)
) engine=Innodb;
```

## Swagger ( API Documentation )
Install swagger with Golang
```console
go get -u github.com/swaggo/swag
```

To generate docs, run this command
```console
$ swag init -g api/handlers/handlers.go
```