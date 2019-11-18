package plugin

import (
	"context"

	"github.com/clintjedwards/cursor/plugin/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// GRPCClient represents the implementation for a client that can talk to plugins
type GRPCClient struct{ client proto.CursorPluginClient }

// GRPCClient is the client implementation that allows our host to send RPCs to plugins
func (p *CursorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewCursorPluginClient(c)}, nil
}

// Below are wrappers for how plugins should respond to the RPC in question
// They are all pretty simple since the general flow is to just call the implementation
// of the rpc method for that specific plugin and return the result

// ExecuteTask calls ExecuteTask on the plugin through the GRPC client
func (m *GRPCClient) ExecuteTask(request *proto.ExecuteTaskRequest) (*proto.ExecuteTaskResponse, error) {
	response, err := m.client.ExecuteTask(context.Background(), request)
	if err != nil {
		return &proto.ExecuteTaskResponse{}, err
	}
	return response, nil
}

// GetPipelineInfo calls GetPipelineInfo on the plugin through the GRPC client
func (m *GRPCClient) GetPipelineInfo(request *proto.GetPipelineInfoRequest) (*proto.GetPipelineInfoResponse, error) {
	response, err := m.client.GetPipelineInfo(context.Background(), request)
	if err != nil {
		return &proto.GetPipelineInfoResponse{}, err
	}
	return response, nil
}
