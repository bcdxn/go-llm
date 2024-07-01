package shared

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

type LLM interface {
	GetPluginProtocolVersion() string
	GetModels() []string
}

type LLMRPC struct {
	client *rpc.Client
}

func (l *LLMRPC) GetPluginProtocolVersion() string {
	var resp string
	err := l.client.Call("Plugin.GetPluginProtocolVersion", new(interface{}), &resp)
	if err != nil {
		// TODO: return err
		panic(err)
	}

	return resp
}

func (l *LLMRPC) GetModels() string {
	var resp string
	err := l.client.Call("Plugin.GetModels", new(interface{}), &resp)
	if err != nil {
		// TODO: return err
		panic(err)
	}

	return resp
}

type LLMRPCServer struct {
	Impl LLM
}

type LLMPlugin struct {
	Impl LLM
}

func (p *LLMPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &LLMRPCServer{Impl: p.Impl}, nil
}

func (LLMPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &LLMRPC{client: c}, nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var DefaultHandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "llm",
}
