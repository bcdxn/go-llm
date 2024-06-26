package chat

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
)

// The application state for a Bubble Tea component
type model struct {
	viewport       viewport.Model
	textarea       textarea.Model
	conversationId string
	messages       []string
	model          string
	width          int // window width
	height         int // window height
}
