package image

import (
	"os"
	"unsafe"

	"golang.org/x/sys/windows"
)

var getCurrentConsoleFontEx = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetCurrentConsoleFontEx")
var getConsoleWindow = windows.NewLazySystemDLL("kernel32.dll").NewProc("GetConsoleWindow")
var isWindowVisible = windows.NewLazySystemDLL("user32.dll").NewProc("IsWindowVisible")

type consoleFontInfo struct {
	Size   uint32
	Index  uint32
	Width  int16
	Height int16
	Family uint32
	Weight uint32
	Face   [32]uint16
}

// ConPTY font metrics describe the hidden console, not the host terminal.
// Only use this API for a visible local console window.
func systemCellPixels() (int, int) {
	if os.Getenv("WT_SESSION") != "" || os.Getenv("TERM_PROGRAM") != "" || os.Getenv("SSH_CONNECTION") != "" {
		return 0, 0
	}
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return 0, 0
	}
	visible, _, _ := isWindowVisible.Call(hwnd)
	if visible == 0 {
		return 0, 0
	}
	info := consoleFontInfo{}
	info.Size = uint32(unsafe.Sizeof(info))
	ok, _, _ := getCurrentConsoleFontEx.Call(os.Stdout.Fd(), 0, uintptr(unsafe.Pointer(&info)))
	if ok == 0 {
		return 0, 0
	}
	return int(info.Width), int(info.Height)
}
