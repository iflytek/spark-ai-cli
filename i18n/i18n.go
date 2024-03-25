package i18n

import (
	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	MenuI18n       string
	RunI18N        string
	Copy           string
	ReviseI18N     string
	CancelI18N     string
	YourScripts    string
	Explanation    string
	CommandRunMenu string
	ReviseMenu     string
	ReviseResult   string
	RunResult      string
)

func Init(tag language.Tag) {
	bundle := i18n.NewBundle(tag)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.LoadMessageFile("active.zh.toml")

	loc := i18n.NewLocalizer(bundle, "en")

	RunI18N = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Run",
			Other: "Run",
		},
	})

	Copy = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Copy",
			Other: "Copy",
		},
	})

	Explanation = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Explanation",
			Other: "Explanation:",
		},
	})

	YourScripts = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "YourScripts",
			Other: "Your Scripts:",
		},
	})

	ReviseI18N = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Revise",
			Other: "Revise",
		},
	})

	CancelI18N = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Cancel",
			Other: "Cancel",
		},
	})

	MenuI18n = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "Menu",
			Other: "What would you do with this script?",
		},
	})

	CommandRunMenu = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "CommandRunMenu",
			Other: "Are you sure you want to run this scripts?",
		},
	})

	ReviseMenu = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ReviseMenu",
			Other: "What would you like spark to change this scripts?",
		},
	})

	ReviseResult = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ReviseResult",
			Other: "Your new scripts:",
		},
	})

	RunResult = loc.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "RunResult",
			Other: "Run Results:",
		},
	})

}
