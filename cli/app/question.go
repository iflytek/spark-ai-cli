package app

import (
	"SparkAICLI/decoration"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func questionCommand(question string) {

	//TODO config

	//workflows
	m := decoration.NewQRootModel(question)
	p := tea.NewProgram(m)

	//m := decoration.NewCommandRunModel("test")
	if _, err := p.Run(); err != nil {
		fmt.Println("Couldn't start program:", err)
		os.Exit(1)
	}

}

func init() {
	sparkCmd.AddCommand(questionCmd)
}

var questionCmd = &cobra.Command{
	Use:   "q",
	Short: "ask spark ai cli question.",
	Long:  `ask question of spark`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ASK Question:", args)
		question := strings.Join(args, " ")
		questionCommand(question)
	},
}
