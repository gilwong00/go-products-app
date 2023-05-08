currency:
	protoc --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import protos/currency.proto

gen:
	protoc --proto_path=proto proto/*.proto  --go_out=. --go-grpc_out=.

server:
	go run main.go

.PHONY: currency server

