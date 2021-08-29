package main

import (
	"syscall"
)

type (
	HANDLE uintptr
	WORD   uint16
	DWORD  uint32
)

const (
	STD_OUTPUT_HANDLE    = 0xFFFFFFF5
	FOREGROUND_BLUE      = 0x01
	FOREGROUND_GREEN     = 0x02
	FOREGROUND_RED       = 0x04
	FOREGROUND_INTENSITY = 0x08
	BACKGROUND_BLUE      = 0x10
	BACKGROUND_GREEN     = 0x20
	BACKGROUND_RED       = 0x40
	BACKGROUND_INTENSITY = 0x80
)

var (
	modkernel32                 = syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle            = modkernel32.NewProc("GetStdHandle")
	procSetConsoleTextAttribute = modkernel32.NewProc("SetConsoleTextAttribute")
)

func init() {
	textbackground(0x07)
}

func GetStdHandle(nStdHandle DWORD) HANDLE {
	ret, _, _ := procGetStdHandle.Call(uintptr(nStdHandle))
	return HANDLE(ret)
}

func SetConsoleTextAttribute(hConsoleOutput HANDLE, wAttributes WORD) bool {
	ret, _, _ := procSetConsoleTextAttribute.Call(
		uintptr(hConsoleOutput),
		uintptr(wAttributes),
	)
	return ret != 0
}

func textbackground(color int) {
	hOut := GetStdHandle(STD_OUTPUT_HANDLE)
	SetConsoleTextAttribute(hOut, WORD(color))
}

func resettextbackground() {
	textbackground(0x07)
}
