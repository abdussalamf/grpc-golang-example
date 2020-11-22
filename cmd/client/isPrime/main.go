package main

import (
	"context"
	"log"
	"simple-grpc/pkg/proto/math"

	"google.golang.org/grpc"
)

func main() {
	log.Println("Client 1 running ...")

	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	response, err := math.NewMathClient(conn).IsPrime(context.Background(), &math.Request{Num: int32(1)})

	log.Println(response)
	log.Println(err)
}
