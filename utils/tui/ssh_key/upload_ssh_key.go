package tui

import (
	"fmt"
	"gitee_cli/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func InitialUploadSSHKeyTui(fileList []string) *tea.Program {
	return tea.NewProgram(UploadSSHKeyTui{
		FileList: fileList,
	})
}

type UploadSSHKeyTui struct {
	FileList []string
	Cursor   int
}

func (s UploadSSHKeyTui) Init() tea.Cmd {
	return nil
}

func (s UploadSSHKeyTui) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.Cursor > 0 {
				s.Cursor--
			}
		case "down", "j":
			if s.Cursor < len(s.FileList)-1 {
				s.Cursor++
			}
		case "ctrl+c", "q":
			s.Cursor = -1
			return s, tea.Quit
		case "enter":
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s UploadSSHKeyTui) View() string {
	promote := "请选择要上传的 SSH 公钥\n"
	for i, file := range s.FileList {
		cursor := " "
		if i == s.Cursor {
			cursor = utils.Green(">")
			file = utils.Yellow(file)
		}
		promote += fmt.Sprintf("%s %s\n", cursor, file)
	}
	return promote
}
