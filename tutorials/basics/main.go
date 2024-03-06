package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	// 光标，需要记录用户当前的光标在哪里，不然我怎么知道当前在哪个位置；显然
	// 并非所有的场景都需要光标，需要根据渲染情况、交互情况来决定是否需要光标。
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func (m *model) String() string {
	var res []string
	for k := range m.selected {
		res = append(res, m.choices[k])
	}

	return strings.Join(res, ",")
}

func initialModel() model {
	// 初始化模型，需要根据需求决定当前应该渲染什么
	return model{
		cursor: 2, // 光标默认指向第一个
		// 初始化数据
		choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		// 初始化的时候默认没有任何一个选中
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// TODO 这玩意有啥用？
	return tea.SetWindowTitle("Grocery List")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// 用于判断是否是按键事件
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// TODO 这里再退出的时候能不能清除屏幕，打印用户选择的东西
			return m, tea.Quit
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(m.choices) - 1
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(m.choices) {
				m.cursor = 0
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				// 如果二次选中，就是取消选择
				delete(m.selected, m.cursor)
			} else {
				// 如果是第一次选中，那就是选择
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

// 用于渲染表单，用户需要根据自己设计的modle渲染表单；一旦model修改了，bubbletea就会执行View，重新渲染表单
func (m model) View() string {
	s := "What should we buy at the market?\n\n"

	// 根据modle
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i { // 说明光标在这个位置
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok { // 说明用户选中了这个选项
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n"

	return s
}

func main() {
	// 第一步：初始化一个模型，模型主要有两个作用：
	// 1、构造模型数据，让bubbletea把数据渲染在命令行上，让用户可以按照我们想要的方式输入数据
	// 2、接收用户的输入，即在bubbletea渲染之后，用户就需要根据自己的实际需求填写数据
	mm := initialModel()

	// 第二步：实例化程序
	p := tea.NewProgram(&mm)

	// 第三步：运行程序，实际上可以理解为渲染表单，然后死循环接收用户的输入；直到用户输入了合法的数据才退出
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	// 第四步：使用数据，实际上使用数据才是我们的根本目的，不然用户数据的数据就没有任何意义
	fmt.Printf("your choose: %s\n", mm.String())
}
