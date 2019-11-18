package plugin

import (
	"context"

	"github.com/clintjedwards/cursor/plugin/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// GRPCServer is the implementation that allows the plugin to respond to requests from the host
type GRPCServer struct {
	Impl PipelineDefinition
}

// GRPCServer is the server implementation that allows our plugins to recieve RPCs
func (p *CursorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCursorPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

// Below are wrappers for how plugins should respond to the RPC in question
// They are all pretty simple since the general flow is to just call the implementation
// of the rpc method for that specific plugin and return the result

// ExecuteTask executes a single task on a plugin
func (m *GRPCServer) ExecuteTask(ctx context.Context, request *proto.ExecuteTaskRequest) (*proto.ExecuteTaskResponse, error) {
	response, err := m.Impl.ExecuteTask(request)
	return response, err
}

// GetPipelineInfo executes a single task on a plugin
func (m *GRPCServer) GetPipelineInfo(ctx context.Context, request *proto.GetPipelineInfoRequest) (*proto.GetPipelineInfoResponse, error) {
	response, err := m.Impl.GetPipelineInfo(request)
	return response, err
}
