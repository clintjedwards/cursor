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

// ExecuteJob is the implementation for how the plugins should respond to the client
func (m *GRPCServer) ExecuteJob(ctx context.Context, req *proto.Empty) (*proto.TestResponse, error) {
	message, err := m.Impl.ExecuteJob()
	return &proto.TestResponse{Message: message}, err
}
