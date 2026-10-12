package util

import (
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	certificateUser32 = syscall.NewLazyDLL("user32.dll")
	promptFocusMutex  sync.Mutex
	promptWindow      uintptr
	// Allocate the Windows callback once rather than on each focus attempt.
	promptWindowCallback = syscall.NewCallback(func(hwnd, _ uintptr) uintptr {
		var pid uint32
		certificateUser32.NewProc("GetWindowThreadProcessId").Call(hwnd, uintptr(unsafe.Pointer(&pid)))
		if pid != uint32(os.Getpid()) {
			return 1
		}
		visible, _, _ := certificateUser32.NewProc("IsWindowVisible").Call(hwnd)
		if visible == 0 {
			return 1
		}
		var class [64]uint16
		certificateUser32.NewProc("GetClassNameW").Call(hwnd, uintptr(unsafe.Pointer(&class[0])), uintptr(len(class)))
		if syscall.UTF16ToString(class[:]) != "#32770" {
			return 1
		}
		promptWindow = hwnd
		return 0
	})
)

// Does not approve, dismiss, or modify the native security warning.
func FocusCertificatePrompt() bool {
	promptFocusMutex.Lock()
	defer promptFocusMutex.Unlock()
	if !certificatePromptActive.Load() {
		return false
	}
	promptWindow = 0
	certificateUser32.NewProc("EnumWindows").Call(promptWindowCallback, 0)
	if promptWindow == 0 {
		return false
	}
	certificateUser32.NewProc("ShowWindow").Call(promptWindow, 9) // SW_RESTORE
	foreground, _, _ := certificateUser32.NewProc("SetForegroundWindow").Call(promptWindow)
	return foreground != 0
}

// Focus the native confirmation without answering it. Used for install and removal.
func beginCertificatePrompt() func() {
	certificatePromptActive.Store(true)
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(250 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				if FocusCertificatePrompt() {
					return
				}
			}
		}
	}()
	return func() { close(done); certificatePromptActive.Store(false) }
}
