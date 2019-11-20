package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/clintjedwards/cursor/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// initClientConn creates a client connection for cmd utilities to make requests
// make sure to close upon using
func initClientConn() *grpc.ClientConn {
	config, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to read configuration: %v", err)
	}

	creds, err := credentials.NewClientTLSFromFile(config.TLSCertPath, "")
	if err != nil {
		log.Fatalf("failed to get certificate: %v", err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(creds))

	hostPortTuple := strings.Split(config.Master.GRPCURL, ":")

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", hostPortTuple[0], hostPortTuple[1]), opts...)
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}

	return conn
}

func generateClientContext() context.Context {
	config, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to read configuration: %v", err)
	}
	md := metadata.Pairs("Authorization", "Bearer "+config.CommandLine.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	return ctx
}
