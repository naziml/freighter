all:
	protoc --go_out=./freighter --go_opt=paths=source_relative --go-grpc_out=./freighter --go-grpc_opt=paths=source_relative proto/freighter.proto