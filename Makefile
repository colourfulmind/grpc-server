.PHONY: build
build:
	go build -o team00 ./cmd/team00/main.go

.PHONY: run
run: build
	./team00 --config=./config/local.yaml

.PHONY: clean
clean:
	rm -rf team00

.PHONY: generate_files
generate_files:
	protoc -I protos/proto protos/proto/transmitter/transmitter.proto --go_out=./protos/gen/go --go_opt=paths=source_relative --go-grpc_out=./protos/gen/go --go-grpc_opt=paths=source_relative


.DEFAULT_GOAL := build