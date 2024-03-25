package decoration

import (
	"SparkAICLI/i18n"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 100

	columnWidth = 30
)

type windowsStatus int

const (
	ShellEdit windowsStatus = iota
	Confirm
)

var (
	gold   = lipgloss.NewStyle().Foreground(lipgloss.Color("#B8860B"))
	purple = lipgloss.NewStyle().Foreground(lipgloss.Color("#800080"))
	subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
)

var (
	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("57")).
			Foreground(lipgloss.Color("230"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("238"))

	endOfBufferStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("235"))

	focusedPlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("99"))

	focusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("238"))

	blurredBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.HiddenBorder())

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFF7DB")).
			Background(lipgloss.Color("#888B7E")).
			Padding(0, 3).
			MarginTop(1)

	activeButtonStyle = buttonStyle.Copy().
				Foreground(lipgloss.Color("#FFF7DB")).
				Background(lipgloss.Color("#F25D94")).
				MarginRight(2).
				Underline(true)

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#874BFD")).
			Padding(1, 0).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)
)

type qCommandRunModel struct {
	textArea      textarea.Model
	currentStatus windowsStatus

	commandData string
	choicesMenu []string // items on the to-do list
	cursor      int      // which to-do list item our cursor is pointing at

	spinner spinner.Model
	timer   timer.Model
	keymap  keymap
	help    help.Model

	width    int
	height   int
	quitting bool
}

func newTextarea() textarea.Model {
	t := textarea.New()
	t.SetHeight(2)
	t.SetWidth(width)
	t.Placeholder = "Your scripts"
	t.ShowLineNumbers = true
	t.Cursor.Style = cursorStyle
	t.FocusedStyle.Placeholder = focusedPlaceholderStyle
	t.BlurredStyle.Placeholder = placeholderStyle
	t.FocusedStyle.CursorLine = cursorLineStyle
	t.FocusedStyle.Base = focusedBorderStyle
	t.BlurredStyle.Base = blurredBorderStyle
	t.FocusedStyle.EndOfBuffer = endOfBufferStyle
	t.BlurredStyle.EndOfBuffer = endOfBufferStyle
	t.KeyMap.DeleteWordBackward.SetEnabled(false)
	t.KeyMap.LineNext = key.NewBinding(key.WithKeys("down"))
	t.KeyMap.LinePrevious = key.NewBinding(key.WithKeys("up"))
	t.Blur()
	return t
}

func NewCommandRunModel(command string) qCommandRunModel {

	m := qCommandRunModel{
		commandData: command,
		// Our to-do list is a grocery list
		choicesMenu: []string{i18n.RunI18N, i18n.CancelI18N},
		keymap:      km,
		spinner:     spinner.New(),
		help:        help.New(),
		timer:       timer.NewWithInterval(time.Second*10, time.Millisecond),
		textArea:    newTextarea(),
	}
	m.spinner.Style = purple
	m.textArea.Focus()
	m.textArea.SetValue(command)
	m.currentStatus = ShellEdit

	return m
}

func (m qCommandRunModel) Init() tea.Cmd {
	//return m.timer.Init()
	return tea.Batch(
		m.spinner.Tick,
		textarea.Blink,
		m.timer.Init(),
	)
}

func (m qCommandRunModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.currentStatus {
	case ShellEdit:
		m.textArea, cmd = m.textArea.Update(msg)
		m.keymap.right.SetEnabled(false)
		m.keymap.left.SetEnabled(false)
	case Confirm:
		m.keymap.right.SetEnabled(true)
		m.keymap.left.SetEnabled(true)
	}

	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.width = msg.Width
		m.height = msg.Height
		return m, cmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.enter):
			// Send the choice on the channel and exit.
			choice := m.choicesMenu[m.cursor]
			switch choice {
			case i18n.RunI18N:

			case i18n.CancelI18N:
				return m, tea.Quit
			}
		case key.Matches(msg, m.keymap.right):
			m.cursor++
			if m.cursor >= len(m.choicesMenu) {
				m.cursor = 0
			}
		case key.Matches(msg, m.keymap.left):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choicesMenu) - 1
			}
		case key.Matches(msg, m.keymap.table):
			m.currentStatus = toggleStatus(m.currentStatus)
		}
	}

	return m, cmd
}

func (m qCommandRunModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.left,
		m.keymap.right,
		m.keymap.enter,
		m.keymap.table,
		m.keymap.quit,
	})
}

func (m qCommandRunModel) View() string {
	s := strings.Builder{}

	var mode string
	switch m.currentStatus {
	case ShellEdit:
		mode = lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("编辑模式")
	case Confirm:
		mode = lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render("执行模式")
	}
	var okButton, cancelButton string

	okButton = buttonStyle.Render(i18n.RunI18N)
	cancelButton = buttonStyle.Render(i18n.CancelI18N)

	switch m.choicesMenu[m.cursor] {
	case i18n.RunI18N:
		okButton = activeButtonStyle.Render("(•)" + i18n.RunI18N)
	case i18n.CancelI18N:
		cancelButton = activeButtonStyle.Render("(•)" + i18n.CancelI18N)
	}
	text := lipgloss.NewStyle().Width(50).Render(m.textArea.View())
	question := lipgloss.NewStyle().Width(50).Align(lipgloss.Center).Render(i18n.CommandRunMenu)
	buttons := lipgloss.JoinHorizontal(lipgloss.Top, okButton, cancelButton)
	ui := lipgloss.JoinVertical(lipgloss.Center, mode, text, question, buttons)
	dialog := lipgloss.Place(width, 20,
		lipgloss.Center, lipgloss.Center,
		dialogBoxStyle.Render(ui),
		lipgloss.WithWhitespaceChars("Spark "),
		lipgloss.WithWhitespaceForeground(subtle),
	)

	s.WriteString(dialog + "\n\n")

	s.WriteString(m.helpView())

	return s.String()
}

func toggleStatus(currentStatus windowsStatus) windowsStatus {
	if currentStatus == ShellEdit {
		return Confirm
	}
	return ShellEdit
}
