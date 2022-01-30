package grpc

import (
	"google.golang.org/grpc"
	"log"
)

func SetupClient() *grpc.ClientConn {
	conn, err := grpc.Dial("127.0.0.1:2021", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	return conn
}
