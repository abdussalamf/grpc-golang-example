package math

import (
	"context"
	"io"
	"log"
	"math/big"
	"simple-grpc/pkg/proto/math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct{}

type Server interface {
	MultiplyByTen(srv math.Math_MultiplyByTenServer) error
	IsPrime(ctx context.Context, req *math.Request) (*math.IsPrimeResponse, error)
}

func NewServer() Server {
	return &server{}
}

func (*server) MultiplyByTen(srv math.Math_MultiplyByTenServer) error {
	var multi int32

	ctx := srv.Context()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
		}

		req, err := srv.Recv()
		if err == io.EOF {
			log.Println("Exit")
			return nil
		}

		if err != nil {
			log.Printf("Receiving error: %v", err)
			continue
		}

		multi = req.Num * 10

		resp := math.Response{
			Result: multi,
		}
		if err := srv.Send(&resp); err != nil {
			log.Printf("Sending error: %v", err)
		}
		log.Printf("Multiplied result of %v: %v", req.Num, multi)
	}
}

func (*server) IsPrime(ctx context.Context, req *math.Request) (*math.IsPrimeResponse, error) {

	if req.Num == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "Number %v is not in prime numbers", req.Num)
	}

	boolReturn := big.NewInt(int64(req.Num)).ProbablyPrime(0)

	if boolReturn {
		return &math.IsPrimeResponse{
			IsPrime: true,
		}, nil
	}

	return &math.IsPrimeResponse{
		IsPrime: false,
	}, nil
}
