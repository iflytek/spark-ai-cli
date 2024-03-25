package decoration

import (
	"SparkAICLI/i18n"
	"SparkAICLI/llm"
	"SparkAICLI/prompt"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	SparkRequest windowsStatus = iota
	ReviseTextInput
	Menu
	Revise
)

type qInteractionModel struct {
	shellResultChan        chan string
	shellExplanationResult chan string
	question               string
	currentStatus          windowsStatus

	resultFinishChan       chan bool
	shellExplanationFinish chan bool

	reviseResultChan   chan string
	reviseResultFinish chan bool

	llmFinish           bool
	explanationFinish   bool
	isRevise            bool
	reviseRequire       string
	reviseRequireResult string
	reviseRequireFinish bool

	isFirstResult        bool
	shellSuggestData     string
	shellExplanationData string
	choicesMenu          []string // items on the to-do list
	choice               string
	cursor               int              // which to-do list item our cursor is pointing at
	selected             map[int]struct{} // which to-do items are selected

	spinner   spinner.Model
	timer     timer.Model
	textInput textinput.Model
	keymap    keymap
	help      help.Model

	width    int
	height   int
	quitting bool
}

func NewQInteractionModel(question string) *qInteractionModel {

	m := &qInteractionModel{
		shellResultChan:        make(chan string, 10),
		resultFinishChan:       make(chan bool, 10),
		shellExplanationResult: make(chan string, 10),
		shellExplanationFinish: make(chan bool, 10),
		reviseResultChan:       make(chan string, 10),
		reviseResultFinish:     make(chan bool, 10),
		question:               question,
		isFirstResult:          true,
		choicesMenu:            []string{i18n.Copy, i18n.RunI18N, i18n.ReviseI18N, i18n.CancelI18N}, // Menu list
		selected:               make(map[int]struct{}),
		keymap:                 km,
		spinner:                spinner.New(),
		help:                   help.New(),
		timer:                  timer.NewWithInterval(time.Second*60, time.Millisecond),
		currentStatus:          SparkRequest,
	}
	m.textInput = textinput.New()
	m.textInput.Placeholder = "Describe a shell command or ask question"
	m.textInput.Focus()

	var purple = lipgloss.NewStyle().Foreground(lipgloss.Color("#800080"))
	m.spinner.Style = purple

	return m
}

func (m *qInteractionModel) Init() tea.Cmd {
	//return m.timer.Init()
	workflows(m.question, m.shellResultChan, m.shellExplanationResult, m.resultFinishChan, m.shellExplanationFinish)

	return tea.Batch(
		m.spinner.Tick,
		textinput.Blink,
		m.timer.Init(),
	)
}

func handleLlmResult(m *qInteractionModel) {
	select {
	case data, ok := <-m.shellResultChan:
		if ok {
			if m.isFirstResult {
				data = shellPostProcess(data)
				m.isFirstResult = false
			}
			m.shellSuggestData += data
		}
	case finish, ok := <-m.resultFinishChan:
		if ok {
			m.llmFinish = finish
		}
	case finish, ok := <-m.shellExplanationFinish:
		if ok {
			m.explanationFinish = finish
			m.currentStatus = Menu
		}
	case data, ok := <-m.shellExplanationResult:
		if ok {
			m.shellExplanationData += data
		}
	case finish, ok := <-m.reviseResultFinish:
		if ok {
			m.reviseRequireFinish = finish
			m.currentStatus = Menu
			//update suggest shell content
			m.shellSuggestData = m.reviseRequireResult
		}
	case data, ok := <-m.reviseResultChan:
		if ok {
			m.reviseRequireResult += data
		}
	default:

	}

}

func (m *qInteractionModel) handleTextInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.enter):
			// Send the choice on the channel and exit.
			m.reviseRequire = m.textInput.Value()
			ReviseWorkflows(m.reviseRequire, m.shellSuggestData, m.reviseResultChan, m.reviseResultFinish)
			m.isRevise = true

			m.currentStatus = Revise
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m *qInteractionModel) handleMenuInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.enter):
			// Send the choice on the channel and exit.
			if m.llmFinish {
				return m.handleEnter()
			}
		case key.Matches(msg, m.keymap.down):
			m.cursor++
			if m.cursor >= len(m.choicesMenu) {
				m.cursor = 0
			}
		case key.Matches(msg, m.keymap.up):
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choicesMenu) - 1
			}
		}
	}
	return m, nil
}

func (m *qInteractionModel) handleSparkRequest(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil

}

func (m *qInteractionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	handleLlmResult(m)

	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		m.keymap.stop.SetEnabled(m.timer.Running())
		m.keymap.start.SetEnabled(!m.timer.Running())
		return m, cmd

	case timer.TimeoutMsg:
		m.quitting = true
		return m, tea.Quit
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
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
		}
	}

	switch m.currentStatus {
	case SparkRequest:
		return m.handleSparkRequest(msg)
	case ReviseTextInput:
		return m.handleTextInput(msg)
	case Menu:
		return m.handleMenuInput(msg)
	case Revise:

	}

	return m, cmd
}

func (m *qInteractionModel) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.up,
		m.keymap.down,
		m.keymap.enter,
		m.keymap.quit,
	})
}

func (m *qInteractionModel) View() string {
	s := strings.Builder{}
	if !m.explanationFinish {
		s.WriteString(m.spinner.View() + gold.Render(" Ask Spark ...") + "\n")
	}
	s.WriteString(i18n.YourScripts + "\n")
	if m.timer.Timedout() {
		s.WriteString("All done!")
	}
	s.WriteString(m.shellSuggestData)

	// shell suggestion finish,show explanation
	if m.llmFinish {
		s.WriteString("\n" + i18n.Explanation + "\n")
		s.WriteString(m.shellExplanationData)
	}

	//Revise Input show
	if m.currentStatus == ReviseTextInput {
		s.WriteString("\n" + i18n.ReviseMenu + "\n")
		s.WriteString(m.textInput.View())
		//TODO show menu
	}

	//Revise input finish content show
	if m.isRevise {
		s.WriteString("\n" + i18n.ReviseMenu + "\n")
		s.WriteString(m.reviseRequire + "\n")
		s.WriteString(i18n.ReviseResult + "\n")
		s.WriteString(m.reviseRequireResult + "\n")
	}

	//TODO add emoji
	if m.currentStatus == Menu {
		s.WriteString("\n" + i18n.MenuI18n + "\n")
		for i := 0; i < len(m.choicesMenu); i++ {
			if m.cursor == i {
				s.WriteString("(•) ")
			} else {
				s.WriteString("( ) ")
			}
			s.WriteString(m.choicesMenu[i])
			s.WriteString("\n")
		}
	}

	if !m.quitting {
		//s = "Exiting in " + s
		s.WriteString(m.helpView())
	}

	return lipgloss.NewStyle().Width(m.width).Render(s.String())
}

func workflows(question string, shellResultChan, shellExplanationResult chan string, resultFinishChan, shellExplanationFinish chan bool) {
	shellSuggestionPrompt := prompt.GetCommandPrompt(question)
	//llm request
	go func() {
		//get shell command suggestion
		tmpResult := llm.LLMRequest(shellSuggestionPrompt, shellResultChan, resultFinishChan)

		explanationPrompt := prompt.GetCommandExplanationPrompt(tmpResult)
		llm.LLMRequest(explanationPrompt, shellExplanationResult, shellExplanationFinish)

	}()

}

func shellPostProcess(cmd string) string {
	cmd = strings.TrimPrefix(cmd, "$")
	cmd = strings.TrimPrefix(cmd, " ")
	return cmd
}

func ReviseWorkflows(requirement, shell string, resultChan chan string, resultFinish chan bool) {
	shellSuggestionPrompt := prompt.GetReviseCommandPrompt(requirement, shell)
	//llm request
	go func() {
		//get shell command suggestion
		llm.LLMRequest(shellSuggestionPrompt, resultChan, resultFinish)

	}()

}

func (m *qInteractionModel) handleEnter() (tea.Model, tea.Cmd) {
	m.choice = m.choicesMenu[m.cursor]
	switch m.choice {
	case i18n.RunI18N:
		sendCmd2Shell(m.shellSuggestData)
		return m, tea.Quit
		//return m, changeCommandWindowsCmd(m.shellSuggestData)
	case i18n.Copy:
		m.runCopyCommand()
		return m, tea.Quit
	case i18n.ReviseI18N:
		m.currentStatus = ReviseTextInput
		//set default value for multiple uses
		m.reviseRequireResult = ""
		m.textInput.SetValue("")
		m.reviseRequireFinish = false
		m.explanationFinish = false
		m.isRevise = false
	case i18n.CancelI18N:
		return m, tea.Quit
	}
	return m, nil
}

func (m *qInteractionModel) runCopyCommand() {
	if clipboard.Unsupported {
		panic(`Linux, Unix (requires 'xclip' or 'xsel' command to be installed)`)
	}
	clipboard.WriteAll(m.shellSuggestData)
	fmt.Println("Copy Finish,Enjoy")
}

func sendCmd2Shell(cmd string) {
	//postprocess
	cmd = strings.Replace(cmd, "```", "", -1)
	cmd = strings.Replace(cmd, "```bash", "", -1)
	if strings.HasPrefix(cmd, "$") {
		cmd = strings.TrimPrefix(cmd, "$")
	}

	fmt.Println(cmd)
	cmdObj := exec.Command("bash", "-c", cmd)
	output, err := cmdObj.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(i18n.RunResult)
	fmt.Println(string(output))
}
