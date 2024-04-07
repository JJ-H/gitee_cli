package utils

import (
	"fmt"
	"gitee_cli/config"
	"github.com/AlecAivazis/survey/v2"
	"os"
)

func InitialEditor(message string, defaultContent string) *survey.Editor {
	return &survey.Editor{
		Editor:        config.Conf.DefaultEditor,
		Default:       fmt.Sprint(defaultContent),
		AppendDefault: true,
		Message:       message,
		FileName:      "*.md",
	}
}

func ReadFromEditor(editor *survey.Editor, content string) string {
	survey.EditorQuestionTemplate = `
{{- if .ShowHelp }}{{- color .Config.Icons.Help.Format }}{{ .Config.Icons.Help.Text }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color .Config.Icons.Question.Format }}{{ .Config.Icons.Question.Text }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ .Config.HelpInput }} for help]{{color "reset"}} {{end}}
  {{- color "cyan"}}[Enter to launch editor] {{color "reset"}}
{{- end}}`
	if err := survey.AskOne(editor, &content); err != nil {
		os.Exit(0)
	}
	return content
}

func ReadFromInput(message, content string) string {
	prompt := &survey.Input{
		Message: message,
	}
	if err := survey.AskOne(prompt, &content); err != nil {
		os.Exit(0)
	}
	return content
}
