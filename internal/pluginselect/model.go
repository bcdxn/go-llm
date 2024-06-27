package pluginselect

import (
	"github.com/bcdxn/go-llm/internal/config"
	"github.com/bcdxn/go-llm/internal/plugins"
	"github.com/charmbracelet/bubbles/list"
)

// The application state for a Bubble Tea component
type model struct {
	plugins  []plugins.PluginListItem
	list     list.Model
	selected plugins.PluginListItem
	width    int // window width
	height   int // window height
	cfg      config.Config
}
