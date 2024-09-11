package main

import (
	"fmt"
	"os/exec"
	"time"
	"user32/win"
)

func main() {
	// 打开记事本应用程序
	cmd := exec.Command("notepad.exe", "example.txt")
	err := cmd.Start()
	if err != nil {
		fmt.Println("无法打开记事本:", err)
		return
	}

	// 等待记事本打开
	time.Sleep(2 * time.Second)

	// 获取记事本窗口句柄
	hwnd := win.FindWindow("", "example.txt - Notepad")
	if hwnd == 0 {
		fmt.Println("查找记事本窗口失败")
		return
	}

	fmt.Println("记事本窗口句柄:", hwnd)

	// 寻找记事本 编辑条目
	hwnd1 := win.FindWindowEx(hwnd, 0, "NotepadTextBox", "")
	if hwnd1 == 0 {
		fmt.Println("hwnd1失败")
		return
	}

	hwndEdit := win.FindWindowEx(hwnd1, 0, "RichEditD2DPT", "")

	fmt.Println("记事本编辑条目句柄:", hwndEdit)

	win.InputText(hwndEdit, "Hello, World!")

	// text := win.GetWindowText(hwndEdit)
	// fmt.Println("记事本编辑条目文本:", text)

	// // 点击记事本编辑条目
	// win.ClickMouse(hwndEdit)
}
