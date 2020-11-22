package math

import (
	"context"
	"log"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"simple-grpc/pkg/proto/math"
)

func TestIsPrime_Server(t *testing.T) {
	tests := []struct {
		name    string
		value   int32
		res     *math.IsPrimeResponse
		errCode codes.Code
		errMsg  string
	}{
		{
			"Valid request with true return",
			3,
			&math.IsPrimeResponse{IsPrime: true},
			codes.OK,
			"",
		},
		{
			"Valid request with false return",
			4,
			&math.IsPrimeResponse{IsPrime: false},
			codes.OK,
			"",
		},
		{
			"Invalid request with 1 value",
			1,
			nil,
			codes.InvalidArgument,
			"Number 1 is not in prime numbers",
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := math.NewMathClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := &math.Request{Num: tt.value}

			resp, err := client.IsPrime(ctx, request)

			if resp != nil {
				if resp.IsPrime != tt.res.IsPrime {
					t.Error("response: expected", tt.res.IsPrime, "received", resp.IsPrime)
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Code() != tt.errCode {
						t.Error("error code: expected", codes.InvalidArgument, "received", er.Code())
					}
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}
