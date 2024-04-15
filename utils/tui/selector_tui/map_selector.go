package selector_tui

import (
	"fmt"
	"gitee_cli/utils"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

type MapSelector struct {
	OptionsMap map[string]int
	Options    []string
	Promote    string
	Cursor     int
}

func NewMapSelector(optionsMap map[string]int, options []string, promote string, altScreen bool) *tea.Program {
	mapSelector := MapSelector{
		OptionsMap: optionsMap,
		Options:    options,
		Promote:    promote,
	}
	if altScreen {
		return tea.NewProgram(mapSelector, tea.WithAltScreen())
	} else {
		return tea.NewProgram(mapSelector)
	}
}

func (m MapSelector) Init() tea.Cmd {
	return nil
}

func (m MapSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			if m.Cursor < len(m.Options)-1 {
				m.Cursor++
			}
		case "ctrl+c", "q":
			m.Cursor = -1
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m MapSelector) SelectedValue() (int, error) {
	if m.Cursor == -1 {
		os.Exit(0)
	}
	option := m.Options[m.Cursor]

	if value, ok := m.OptionsMap[option]; ok {
		return value, nil
	}
	return 0, fmt.Errorf("无效的选项")
}

func (m MapSelector) View() string {
	promote := m.Promote + "\n"
	for i, option := range m.Options {
		cursor := " "
		if i == m.Cursor {
			cursor = utils.Green(">")
			option = utils.Yellow(option)
		}
		promote += fmt.Sprintf("%s %s\n", cursor, option)
	}
	return promote
}

func chunkOptions(options map[string]int) []string {
	_options := make([]string, 0)
	for key, _ := range options {
		_options = append(_options, key)
		if len(_options) == 16 {
			break
		}
	}
	return _options
}
