package chat

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/bootstrap-library/stock-plus/function"
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

func (s *MenuState) HandleInput(input string) (State, string) {
	i, _ := strconv.Atoi(input)
	item := menu.GetItem(i - 1)
	if item == nil {
		return s, "无效的选项"
	}

	return NewFunctionState(item.Function), item.Name
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
