package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

var (
	// 用于设置原色
	color = termenv.EnvColorProfile().Color
	// 关键字颜色
	keyword = termenv.Style{}.Foreground(color("204")).Background(color("235")).Styled
	// 帮助颜色
	help = termenv.Style{}.Foreground(color("241")).Styled
)

type model struct {
	altscreen bool
	quitting  bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			m.quitting = true
			return m, tea.Quit
		case " ": // 通过空格切换是进入全屏还是退出全屏
			var cmd tea.Cmd
			if m.altscreen {
				cmd = tea.ExitAltScreen // 退出全屏
			} else {
				cmd = tea.EnterAltScreen // 进入全屏
			}
			m.altscreen = !m.altscreen
			return m, cmd
		}
	}
	return m, nil
}

func (m model) View() string {
	// 如果退出了，打印下面这段话
	if m.quitting {
		return "Bye!\n"
	}

	const (
		altscreenMode = " altscreen mode " // 全屏模式
		inlineMode    = " inline mode "    // 命令行模式
	)

	var mode string
	if m.altscreen {
		mode = altscreenMode
	} else {
		mode = inlineMode
	}

	return fmt.Sprintf("\n\n  You're in %s\n\n\n", keyword(mode)) +
		help("  space: switch modes • q: exit\n")
}

func main() {
	if _, err := tea.NewProgram(model{}).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
