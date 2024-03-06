package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	gotReposSuccessMsg []repo
	gotReposErrMsg     error
)

// 用于反序列化当前的仓库名
type repo struct {
	Name string `json:"name"`
}

const reposURL = "https://api.github.com/orgs/charmbracelet/repos"

// 这里非常简单，就是获取charmbracelet仓库底下的所有项目
func getRepos() tea.Msg {
	req, err := http.NewRequest(http.MethodGet, reposURL, nil)
	if err != nil {
		return gotReposErrMsg(err)
	}

	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return gotReposErrMsg(err)
	}
	defer resp.Body.Close() // nolint: errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return gotReposErrMsg(err)
	}

	var repos []repo

	err = json.Unmarshal(data, &repos)
	if err != nil {
		return gotReposErrMsg(err)
	}

	return gotReposSuccessMsg(repos)
}

type model struct {
	// 文本框输入模型
	textInput textinput.Model
}

func initialModel() model {
	ti := textinput.New()
	ti.Prompt = "charmbracelet/"
	ti.Placeholder = "repo..."
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	ti.ShowSuggestions = true
	return model{textInput: ti}
}

// 用于初始化数据，有些数据是通过动态的方式获取的
func (m model) Init() tea.Cmd {
	// textinput.Blink用于执行光标闪烁
	return tea.Batch(getRepos, textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case gotReposSuccessMsg:
		var suggestions []string
		for _, r := range msg {
			suggestions = append(suggestions, r.Name)
		}
		m.textInput.SetSuggestions(suggestions)
	}

	// 之所以在这里没有看到处理ctrl+n, ctrl+p，是因为textinput.Model已经自动处理了

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return fmt.Sprintf(
		"What’s your favorite Charm repository?\n\n%s\n\n%s\n",
		m.textInput.View(),
		"(tab to complete, ctrl+n/ctrl+p to cycle through suggestions, esc to quit)",
	)
}
