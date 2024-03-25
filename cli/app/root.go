package app

import (
	"SparkAICLI/i18n"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"os"
	"strings"
)

var sparkCmd = &cobra.Command{
	Use:   "spark",
	Short: "spark ai cli is a next gen shell with llm",
	Long:  `spark ai cli is a next gen shell with llm`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ASK Question:", args)
		question := strings.Join(args, " ")
		questionCommand(question)
	},
}

func Execute() {
	//i18n init
	//TODO config
	i18n.Init(language.English)

	if err := sparkCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
