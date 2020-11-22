package main

import (
	"log"
	"net"

	srv "simple-grpc/internal/math"
	"simple-grpc/pkg/proto/math"

	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	math.RegisterMathServer(s, srv.NewServer())

	log.Println("starting gRPC server...")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
