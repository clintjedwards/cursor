package plugin

import (
	"github.com/hashicorp/go-plugin"
)

// This file contains structures that both the plugin and the plugin host has to implement

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "CURSOR_PLUGIN",
	MagicCookieValue: "cKykOnGDBJ",
}

// PipelineDefinition is the interface in which both the plugin and the host has to implement
type PipelineDefinition interface {
	ExecuteJob() (string, error)
}

// CursorPlugin is just a wrapper so we implement the correct go-plugin interface
// it allows us to serve/consume the plugin
type CursorPlugin struct {
	plugin.Plugin
	Impl PipelineDefinition
}
