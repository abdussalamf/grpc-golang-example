package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"time"

	"simple-grpc/pkg/proto/math"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := math.NewMathClient(conn)
	stream, err := client.MultiplyByTen(context.Background())
	if err != nil {
		log.Fatalf("open stream error %v", err)
	}

	var val int32

	ctx := stream.Context()
	done := make(chan bool)

	// first goroutine sends random increasing numbers to stream
	// and closes int after 10 iterations
	go func() {
		for i := 1; i < 10; i++ {
			rnd := int32(rand.Intn(i))
			req := math.Request{Num: rnd}
			if err := stream.Send(&req); err != nil {
				log.Fatalf("can not send %v", err)
			}
			log.Printf("%d sent", req.Num)
			time.Sleep(time.Millisecond * 200)
		}

		if err := stream.CloseSend(); err != nil {
			log.Println(err)
		}
	}()

	// second goroutine receives data from stream
	// and saves result in max variable
	//
	// if stream is finished it closes done channel
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			val = resp.Result
			log.Printf("Multiplied result received: %d", val)
		}
	}()

	// third goroutine closes done channel
	// if context is done
	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err)
		}
		close(done)
	}()

	<-done
	log.Printf("finished with value=%d", val)
}
