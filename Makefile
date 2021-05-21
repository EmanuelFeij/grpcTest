proto:

	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    protos/company/company.proto


	

install:
	go get github.com/lib/pq
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
	go get google.golang.org/grpc
	go get google.golang.org/grpc/codes
	go get google.golang.org/grpc/status