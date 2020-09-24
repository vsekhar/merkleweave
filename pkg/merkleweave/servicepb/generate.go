package servicepb

//go:generate protoc -I ../../../proto --go-grpc_out=../../.. --go-grpc_opt=module=github.com/vsekhar/merkleweave --go_out=../../.. --go_opt=module=github.com/vsekhar/merkleweave ../../../proto/service.proto
