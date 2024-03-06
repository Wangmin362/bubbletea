package main

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/*
This example assumes an existing understanding of commands and messages. If you
haven't already read our tutorials on the basics of Bubble Tea and working
with commands, we recommend reading those first.

Find them at:
https://github.com/charmbracelet/bubbletea/tree/master/tutorials/commands
https://github.com/charmbracelet/bubbletea/tree/master/tutorials/basics
*/

// sessionState is used to track which model is focused
type sessionState uint

const (
	defaultTime              = time.Minute
	timerView   sessionState = iota // 0
	spinnerView                     // 1
)

var (
	// Available spinners
	spinners = []spinner.Spinner{ // 旋转器不同的样式
		spinner.Line,
		spinner.Dot,
		spinner.MiniDot,
		spinner.Jump,
		spinner.Pulse,
		spinner.Points,
		spinner.Globe,
		spinner.Moon,
		spinner.Monkey,
	}
	// 非聚焦的样式
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
		// 聚焦样式
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type mainModel struct {
	state   sessionState  // 记录当前光标指向的位置
	timer   timer.Model   // 内嵌了时间组件
	spinner spinner.Model // 内嵌了旋转器组件
	index   int           // 记录旋转器当前需要使用的样式索引
}

func newModel(timeout time.Duration) mainModel {
	m := mainModel{state: timerView}
	m.timer = timer.New(timeout)
	m.spinner = spinner.New()
	return m
}

func (m mainModel) Init() tea.Cmd {
	// start the timer and spinner on program start
	return tea.Batch(m.timer.Init(), m.spinner.Tick)
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			// 按住tab键，控制当前聚焦在哪个文本框里面
			if m.state == timerView {
				m.state = spinnerView
			} else {
				m.state = timerView
			}
		case "n":
			// 如果是计时器，那么重置计时器
			if m.state == timerView {
				m.timer = timer.New(defaultTime)    // 重置定时器
				cmds = append(cmds, m.timer.Init()) // 添加命令
			} else { // 否则，肯定就是旋转器了
				m.Next()                            // 选择下一个旋转器的样式
				m.resetSpinner()                    // 重置旋转器的样式
				cmds = append(cmds, m.spinner.Tick) // 添加命令
			}
		}

		// 接下来需要渲染当前选中的组件，切换之后，手动更新一次组件
		switch m.state {
		// update whichever model is focused
		case spinnerView:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		default:
			m.timer, cmd = m.timer.Update(msg)
			cmds = append(cmds, cmd)
		}
	case spinner.TickMsg:
		// 更新旋转器，实际上旋转器之所以看起来在动，就是通过这里看更新的
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	case timer.TickMsg:
		// 更新定时器，这里看到计时器在动也是因为这里一直在更新
		m.timer, cmd = m.timer.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	var s string
	model := m.currentFocusedModel()
	if m.state == timerView { // 计时器
		s += lipgloss.JoinHorizontal(
			lipgloss.Top,
			focusedModelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), // 注意这里聚焦窗口的使用，把计时器包裹在里面
			modelStyle.Render(m.spinner.View()))                          // 注意这里
	} else { // 旋转器
		s += lipgloss.JoinHorizontal(lipgloss.Top,
			modelStyle.Render(fmt.Sprintf("%4s", m.timer.View())), // 非聚焦
			focusedModelStyle.Render(m.spinner.View()))            // 聚焦窗口把旋转器包裹了起来
	}
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • n: new %s • q: exit\n", model))
	return s
}

func (m mainModel) currentFocusedModel() string {
	if m.state == timerView {
		return "timer"
	}
	return "spinner"
}

func (m *mainModel) Next() {
	if m.index == len(spinners)-1 { // 从头开始展示
		m.index = 0
	} else {
		m.index++
	}
}

func (m *mainModel) resetSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[m.index]
}

func main() {
	p := tea.NewProgram(newModel(defaultTime))

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
