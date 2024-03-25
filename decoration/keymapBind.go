package decoration

import (
	"github.com/charmbracelet/bubbles/key"
)

type keymap struct {
	start key.Binding
	stop  key.Binding
	reset key.Binding
	quit  key.Binding
	left  key.Binding
	right key.Binding
	up    key.Binding
	down  key.Binding
	enter key.Binding
	table key.Binding
}

var km keymap

func init() {
	km = keymap{
		start: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "start"),
		),
		stop: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "stop"),
		),
		reset: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reset"),
		),
		quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl-c", "quit"),
		),
		left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("left | h", "left"),
		),
		right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("right | l", "right"),
		),
		down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down | j", "down"),
		),
		up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up | k", "up"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "enter"),
		),
		table: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "change"),
		),
	}
}
