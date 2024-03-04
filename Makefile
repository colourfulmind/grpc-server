.PHONY: run-server
run-server:
	go build -o server ./cmd/server/main.go
	./server --config=./config/config.yaml

.PHONY: run-client
run-client:
	go build -o client ./cmd/client/main.go
	./client -k 1.2 --config=./config/config.yaml

.PHONY: clean
clean:
	rm -rf server client

.PHONY: generate_files
generate_files:
	protoc -I protos/proto protos/proto/transmitter/transmitter.proto --go_out=./protos/gen/go --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go --go-grpc_opt=paths=source_relative


.DEFAULT_GOAL := run-server