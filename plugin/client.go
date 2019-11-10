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

// ExecuteJob calls ExecuteJob on the plugin through the GRPC client
func (m *GRPCClient) ExecuteJob() (string, error) {
	ret, err := m.client.ExecuteJob(context.Background(), &proto.Empty{})
	if err != nil {
		return "", err
	}
	return ret.Message, nil
}
