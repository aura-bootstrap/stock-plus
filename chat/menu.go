package chat

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/bootstrap-library/stock-plus/function"
	"github.com/bootstrap-library/stock-plus/telegram"
)

type MenuState struct{}

func (s *MenuState) String() string {
	return "MenuState"
}

func (s *MenuState) EnterState() string {
	var b bytes.Buffer
	b.WriteString("请选择一个选项：\n")
	for i, item := range menu.items {
		b.WriteString(fmt.Sprintf("%d. %s\n", i+1, item.Name))
	}

	return b.String()
}

func (s *MenuState) HandleInput(input string, sender telegram.MessageSender) State {
	i, _ := strconv.Atoi(input)
	item := menu.GetItem(i - 1)
	if item == nil {
		sender("无效的选项")
		return s
	}

	sender(fmt.Sprintf("正在运行：%s", item.Name))
	return NewFunctionState(item.Function)
}

func (s *MenuState) LeaveState() string {
	return ""
}

type Menu struct {
	items []*MenuItem
}

type MenuItem struct {
	Name     string
	Function function.Function
}

var menu = &Menu{}

func (m *Menu) Register(item *MenuItem) {
	m.items = append(m.items, item)
}

func (m *Menu) GetItem(id int) *MenuItem {
	if id >= 0 && id < len(m.items) {
		return m.items[id]
	}
	return nil
}

func init() {
	menu.Register(&MenuItem{
		Name:     "检查Ping延迟",
		Function: &function.CheckPingFunction{},
	})
}
