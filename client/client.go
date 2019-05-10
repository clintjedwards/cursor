package client

import "google.golang.org/grpc"

type cursorClient struct {
	Conn *grpc.ClientConn
}
