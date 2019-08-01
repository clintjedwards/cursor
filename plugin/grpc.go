package plugin

import (
	"context"

	"github.com/clintjedwards/cursor/plugin/proto"
	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// GRPCClient is an implementation of KV that talks over RPC.
type GRPCClient struct{ client proto.CursorPluginClient }

func (m *GRPCClient) ExecuteJob() error {
	_, err := m.client.ExecuteJob(context.Background(), &proto.Empty{})
	return err
}

type GRPCServer struct {
	Impl PipelinePlugin
}

func (m *GRPCServer) ExecuteJob(
	ctx context.Context,
	req *proto.Empty) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.ExecuteJob()
}

func (p *CursorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterCursorPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *CursorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewCursorPluginClient(c)}, nil
}
