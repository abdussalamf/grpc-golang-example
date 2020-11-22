package math

import (
	"context"
	"log"
	"net"
	"simple-grpc/pkg/proto/math"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	math.RegisterMathServer(server, NewServer())

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
