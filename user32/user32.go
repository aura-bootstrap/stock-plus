package user32

import (
	"syscall"
	"unsafe"
)

const (
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	MK_LBUTTON     = 0x0001
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
)

var (
	user32             = syscall.NewLazyDLL("user32.dll")
	procFindWindowW    = user32.NewProc("FindWindowW")
	procFindWindowExW  = user32.NewProc("FindWindowExW")
	procGetWindowTextW = user32.NewProc("GetWindowTextW")
	procPostMessageA   = user32.NewProc("PostMessageA")
	procGetDlgItem     = user32.NewProc("GetDlgItem")
)

func FindWindowW(className, windowName *uint16) syscall.Handle {
	ret, _, _ := procFindWindowW.Call(uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowName)))
	return syscall.Handle(ret)
}

func FindWindowExW(hwndParent, hwndChildAfter syscall.Handle, className, windowName *uint16) syscall.Handle {
	ret, _, _ := procFindWindowExW.Call(uintptr(hwndParent), uintptr(hwndChildAfter), uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(windowName)))
	return syscall.Handle(ret)
}

func GetWindowTextW(hwnd syscall.Handle, lpString *uint16, nMaxCount int) int {
	ret, _, _ := procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(lpString)), uintptr(nMaxCount))
	return int(ret)
}

func PostMessageA(hwnd syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procPostMessageA.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return ret
}

func GetDlgItem(hwndParent syscall.Handle, controlID int) syscall.Handle {
	ret, _, _ := procGetDlgItem.Call(uintptr(hwndParent), uintptr(controlID))
	return syscall.Handle(ret)
}
