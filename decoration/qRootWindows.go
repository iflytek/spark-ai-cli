package decoration

import (
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

type windows string

var (
	qInteractionWindows windows = "qInteractionWindows"
	qCommandRunWindows  windows = "qCommandRunWindows"
)

// 主模型，包含子模型
type qRootWindows struct {
	timer            timer.Model
	interactionModel *qInteractionModel
	qCommandRunModel qCommandRunModel
	wShow            windows
	shellScripts     string
}

// change window cmd
func changeCommandWindowsCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return ChangeScreen(qCommandRunWindows, message)
	}
}

func ChangeScreen(w windows, message string) tea.Msg {
	return changeScreenMsg{windowsName: w, message: message}
}

// changeScreenMsg is message that signals to change the screen.
type changeScreenMsg struct {
	windowsName windows
	message     string
}

func NewQRootModel(question string) qRootWindows {

	m := qRootWindows{}

	m.interactionModel = NewQInteractionModel(question)
	m.timer = timer.NewWithInterval(time.Second*60, time.Millisecond) //TODO config
	m.wShow = qInteractionWindows                                     //default

	return m
}

func (m qRootWindows) Init() tea.Cmd {
	// 返回初始化命令
	return tea.Batch(
		m.interactionModel.Init(),
	)
}

func (m qRootWindows) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd1 tea.Cmd

	switch m.wShow {
	case qInteractionWindows:
		var m1 tea.Model
		m1, cmd1 = m.interactionModel.Update(msg)
		m.interactionModel = m1.(*qInteractionModel)

	case qCommandRunWindows:
		var m1 tea.Model
		m1, cmd1 = m.qCommandRunModel.Update(msg)
		m.qCommandRunModel = m1.(qCommandRunModel)
	}

	switch msg.(type) {
	case changeScreenMsg:
		tmpMsg := msg.(changeScreenMsg)
		m.wShow = tmpMsg.windowsName
		m.qCommandRunModel = NewCommandRunModel(tmpMsg.message)
		return m, nil
	}

	return m, tea.Batch(cmd1) // 使用tea.Batch来合并命令
}

func (m qRootWindows) View() string {
	// 组合子模型的视图
	switch m.wShow {
	case qCommandRunWindows:
		return m.qCommandRunModel.View()
	case qInteractionWindows:
		return m.interactionModel.View()
	}
	return ""
}
