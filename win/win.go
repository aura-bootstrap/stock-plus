package win

import (
	"syscall"
	"unicode/utf16"
	"user32/user32"
)

func String(s string) *uint16 {
	if s == "" {
		return nil
	}
	wordString, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return wordString
}

func FindWindow(className, windowName string) syscall.Handle {
	return user32.FindWindowW(String(className), String(windowName))
}

func FindWindowEx(hwndParent, hwndChildAfter syscall.Handle, className, windowName string) syscall.Handle {
	return user32.FindWindowExW(hwndParent, hwndChildAfter, String(className), String(windowName))
}

func GetWindowText(hwnd syscall.Handle) string {
	const nMaxCount = 4096
	buf := make([]uint16, nMaxCount)
	user32.GetWindowTextW(hwnd, &buf[0], nMaxCount)
	s := syscall.UTF16ToString(buf)
	return s
}

func ClickMouse(hwnd syscall.Handle) {
	user32.PostMessageA(hwnd, user32.WM_LBUTTONDOWN, user32.MK_LBUTTON, 0)
	user32.PostMessageA(hwnd, user32.WM_LBUTTONUP, user32.MK_LBUTTON, 0)
}

// https://www.eolink.com/news/post/1063.html

func InputText(hwnd syscall.Handle, text string) {
	for _, r := range text {
		utf16Chars := utf16.Encode([]rune{r})
		for _, char := range utf16Chars {
			// PostMessage(hwnd, WM_LBUTTONDOWN, 0,mX + mY * 65536)。
			user32.PostMessageA(hwnd, user32.WM_KEYDOWN, uintptr(char), 0)
			user32.PostMessageA(hwnd, user32.WM_KEYUP, uintptr(char), 0)
		}
	}
}

func GetDlgItem(hwndParent syscall.Handle, controlID int) syscall.Handle {
	return user32.GetDlgItem(hwndParent, controlID)
}
