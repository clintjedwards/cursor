package client

import (
	"context"
	"fmt"

	"github.com/clintjedwards/cursor/api"
	"google.golang.org/grpc"
)

type CursorClient struct {
	Conn *grpc.ClientConn
}

func (client *CursorClient) Connect(hostname, port string) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure()) //TODO: Get rid of this

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", hostname, port), opts...)
	if err != nil {
		return err
	}

	client.Conn = conn

	return nil
}

func (client *CursorClient) Close() error {
	return client.Conn.Close()
}

func (client *CursorClient) CreatePipeline(request *api.CreatePipelineRequest) (*api.CreatePipelineResponse, error) {
	grpcClient := api.NewCursorMasterClient(client.Conn)

	pipeline, err := grpcClient.CreatePipeline(context.Background(), request)
	if err != nil {
		return &api.CreatePipelineResponse{}, err
	}

	return pipeline, nil
}
