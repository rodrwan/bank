clean-proto cp:
	@rm -rf ./pkg/pb/*

proto p: clean-proto
	@echo "Generate proto"
	@protoc -I ./proto --go_out=${GOPATH}/src  --go-grpc_out=${GOPATH}/src  ./proto/*.proto

.PHONY: proto clean-proto