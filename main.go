package main

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	KEY_EVENT                = 1
	MOUSE_EVENT              = 2
	WINDOW_BUFFER_SIZE_EVENT = 4
	MENU_EVENT               = 8
	FOCUS_EVENT              = 16
)

var kernel32 = windows.NewLazyDLL("kernel32")
var readConsoleInput = kernel32.NewProc("ReadConsoleInputW")

type InputRecord struct {
	EventType uint16
	_         uint16
	Info      [8]uint16
}

type KeyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	UnicodeChar     uint16
	ControlKeyState uint32
}

type MouseEventRecord struct {
	X          int16
	Y          int16
	Button     uint32
	ControlKey uint32
	Event      uint32
}

type WindowBufferSizeRecord struct {
	X int16
	Y int16
}

func read(events []InputRecord) uint32 {
	var n uint32
	readConsoleInput.Call(
		uintptr(windows.Stdin),
		uintptr(unsafe.Pointer(&events[0])),
		uintptr(len(events)),
		uintptr(unsafe.Pointer(&n)))
	return n
}

func main() {
	fmt.Println("Hit ESCAPE key to stop.")
	for {
		var events [10]InputRecord

		n := read(events[:])
		for i := uint32(0); i < n; i++ {
			e := events[i]

			switch e.EventType {
			case KEY_EVENT:
				ee := (*KeyEventRecord)(unsafe.Pointer(&e.Info[0]))
				fmt.Printf("KeyEventRecord: %+v\n", ee)
				if ee.UnicodeChar == 27 {
					return
				}
			}
		}
	}
}
