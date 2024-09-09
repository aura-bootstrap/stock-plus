package chat

import "github.com/bootstrap-library/stock-plus/function"

type FunctionState struct {
	function.Function
	input  chan string
	output chan string
	sender func(string)
}

func NewFunctionState(f function.Function) *FunctionState {
	return &FunctionState{
		Function: f,
		input:    make(chan string),
		output:   make(chan string),
	}
}

func (s *FunctionState) String() string {
	return s.Function.String()
}

func (s *FunctionState) EnterState() string {
	go func() {
		for text := range s.output {
			bot.sender(text)
		}
	}()
	go func() {
		defer close(s.output)
		s.Function.Main(s.input, s.output)
	}()

	return ""
}

func (s *FunctionState) HandleInput(input string) (State, string) {
	s.input <- input
	return s, ""
}

func (s *FunctionState) LeaveState() string {
	return ""
}
