package manager

import (
	"encoding/gob"
	"log"
	"os"
	"os/exec"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mitchellh/go-homedir"
	"github.com/stangirard/yatas/plugins/commons"
)

func RunPlugin(pluginInput commons.Plugin, c *commons.Config) []commons.Tests {
	// Create an hclog.Logger
	gob.Register(map[string]interface{}{})
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	// We're a host! Start by launching the plugin process.
	homeDir, _ := homedir.Expand("~/.yatas.d/plugins/")
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Cmd:             exec.Command(homeDir + "/" + pluginInput.Source + "/" + pluginInput.Version + "/yatas-" + pluginInput.Name),
		Logger:          logger,
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(pluginInput.Name)
	if err != nil {
		log.Fatal(err)
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	yatasPlugin := raw.(commons.Yatas)

	return yatasPlugin.Run(c)
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  2,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.

var pluginMap = map[string]plugin.Plugin{
	"aws": &commons.YatasPlugin{},
}
