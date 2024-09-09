package chat

import (
	"fmt"
)

// State 接口定义了所有状态必须实现的方法
type State interface {
	HandleInput(input string) State
	Display() string
}

// StartState 是初始状态
type StartState struct{}

func (s *StartState) HandleInput(input string) State {
	if input == "1" {
		return &HelpState{}
	} else if input == "2" {
		return &EchoState{}
	}
	return s
}

func (s *StartState) Display() string {
	return "欢迎使用聊天机器人！\n1. 帮助\n2. 回声\n请输入选项："
}

// HelpState 是帮助状态
type HelpState struct{}

func (h *HelpState) HandleInput(input string) State {
	return &StartState{}
}

func (h *HelpState) Display() string {
	return "这是帮助信息。返回主菜单请输入任意键。"
}

// EchoState 是回声状态
type EchoState struct{}

func (e *EchoState) HandleInput(input string) State {
	fmt.Println("回声:", input)
	return &StartState{}
}

func (e *EchoState) Display() string {
	return "请输入要回显的文本："
}
