package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-plugin/examples/basic/shared"
)

type PluginListItem struct {
	Name string
	Path string
}

func Find() ([]PluginListItem, error) {
	var plugins []PluginListItem
	home, err := os.UserHomeDir()
	if err != nil {
		return plugins, err
	}

	gobin := filepath.Join(home, "go", "bin")
	bins, err := os.ReadDir(gobin)
	if err != nil {
		return plugins, err
	}

	r := regexp.MustCompile("^go-llm-plugin-(.*)$")

	for _, f := range bins {
		if r.MatchString(f.Name()) {
			plugins = append(plugins, PluginListItem{
				Name: r.FindStringSubmatch(f.Name())[1],
				Path: filepath.Join(gobin, f.Name()),
			})
		}
	}

	return plugins, nil
}

func Load(pluginPath string) error {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(pluginPath),
		Logger:          logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("greeter")
	if err != nil {
		return err
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	greeter := raw.(shared.Greeter)
	fmt.Println(greeter.Greet())

	return nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"greeter": &shared.GreeterPlugin{},
}

type LlmModel struct {
	Id      string
	Aliases []string
	Name    string
}

type LlmPlugin interface {
	GetModels() []LlmModel
}
