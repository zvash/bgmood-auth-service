dburl="postgresql://root:123@127.0.0.1:5432/bgmood_auth?sslmode=disable"

postgres:
	docker run --name postgres15 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=123 -p 5432:5432 -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root bgmood_auth

mu:
	migrate -path internal/db/migration -database $(dburl) -verbose up

mu1:
	migrate -path internal/db/migration -database $(dburl) -verbose up 1

md:
	migrate -path internal/db/migration -database $(dburl) -verbose down

md1:
	migrate -path internal/db/migration -database $(dburl) -verbose down 1

mr:
	make md && make mu

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

proto:
	rm -rf internal/authpb/* 2>/dev/null
	protoc --go_out=internal/authpb --go_opt=paths=source_relative --go-grpc_out=internal/authpb --go-grpc_opt=paths=source_relative --grpc-gateway_out=internal/authpb --grpc-gateway_opt=paths=source_relative --proto_path=../bgmood-protos/auth-service ../bgmood-protos/auth-service/*.proto

.PHONY: postgres createdb mu mu1 md md1 mr sqlc test server proto