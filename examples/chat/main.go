package main

// A simple program demonstrating the text area component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
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
	errMsg error
)

type model struct {
	viewport    viewport.Model // 用于展示消息
	messages    []string       // 保存当前输入的消息
	textarea    textarea.Model // 用于输入文本
	senderStyle lipgloss.Style // 用于设置文本的样式
	err         error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 80 // 限制字符输入数量

	ta.SetWidth(30) // 设置文本框的输入宽度为30
	ta.SetHeight(3) // 设置文本框的高度为3

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false // 设置不展示行号

	vp := viewport.New(30, 5)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:    ta, // 用于文本框的展示
		messages:    []string{},
		viewport:    vp, // 用于展示数据
		senderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		err:         nil,
	}
}

func (m model) Init() tea.Cmd {
	// 设置文本框的光标闪烁
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	// 需要更新文本框输入组件  TODO 其实内部组件并不需要关心消息，只是有些内部组件需要关心一些快捷键做出响应
	m.textarea, tiCmd = m.textarea.Update(msg)
	// 需要更新展示区域
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			// 这里退出的时候完全可以进行手动渲染，其实就是把数据打印到控制台
			fmt.Println(m.textarea.Value())
			return m, tea.Quit
		case tea.KeyEnter:
			// 保存消息
			m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			// 设置展示文本框需要展示的消息
			m.viewport.SetContent(strings.Join(m.messages, "\n"))
			// 重置输入文本框
			m.textarea.Reset()
			// 如果展示文本框中需要展示的消息超过了最大的展示行数，那么直接显示最后面的几条消息
			m.viewport.GotoBottom()
		}

	// We handle errors just like any other message
	// TODO 为啥收到错误消息的时候不需要进行处理？
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(), // 调用展示组件进行渲染
		m.textarea.View(), // 调用文本框输入组件进行渲染。
	) + "\n\n"
}
