package math

import (
	"context"
	"io"
	"log"
	"testing"

	"google.golang.org/grpc"

	"simple-grpc/pkg/proto/math"
)

func TestMultiplyTen_Server(t *testing.T) {

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := math.NewMathClient(conn)

	stream, err := client.MultiplyByTen(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	ctx := stream.Context()
	done := make(chan bool)

	// first goroutine sends random increasing numbers to stream
	// and closes int after 10 iterations
	go func() {
		req := math.Request{Num: 5}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("can not send %v", err)
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
				// close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			val := resp.Result
			if val != 50 {
				t.Error("Result value unexpected", val)
			}
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
}
