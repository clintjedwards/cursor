package client

import (
	"context"
	"fmt"

	"github.com/clintjedwards/cursor/api"
	"google.golang.org/grpc"
)

// CursorClient represents a cursor master client object
type CursorClient struct {
	Conn *grpc.ClientConn
}

// Connect opens a grpc connection to given host:port
// make sure to close connection once finished with Close() function
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

// Close disconnects an open grpc connection
func (client *CursorClient) Close() error {
	return client.Conn.Close()
}

// CreatePipeline creates a new cursor pipeline
func (client *CursorClient) CreatePipeline(request *api.CreatePipelineRequest) (*api.CreatePipelineResponse, error) {
	grpcClient := api.NewCursorMasterClient(client.Conn)

	pipeline, err := grpcClient.CreatePipeline(context.Background(), request)
	if err != nil {
		return &api.CreatePipelineResponse{}, err
	}

	return pipeline, nil
}

// ListPipelines returns a list of all pipelines that exist for a cursor master
func (client *CursorClient) ListPipelines(request *api.ListPipelinesRequest) (*api.ListPipelinesResponse, error) {
	grpcClient := api.NewCursorMasterClient(client.Conn)

	pipelines, err := grpcClient.ListPipelines(context.Background(), request)
	if err != nil {
		return &api.ListPipelinesResponse{}, err
	}

	return pipelines, nil
}
