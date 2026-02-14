//go:build windows

package actions

import (
	"fmt"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

// SHFILEOPSTRUCTW is the Go representation of the Windows SHFILEOPSTRUCT.
type shFileOpStruct struct {
	hwnd                  uintptr
	wFunc                 uint32
	pFrom                 *uint16
	pTo                   *uint16
	fFlags                uint16
	fAnyOperationsAborted int32
	hNameMappings         uintptr
	lpszProgressTitle     *uint16
}

const (
	foDelete          = 0x0003
	fofAllowUndo      = 0x0040
	fofNoConfirmation = 0x0010
	fofNoErrorUI      = 0x0400
	fofSilent         = 0x0004
)

var (
	shell32              = windows.NewLazySystemDLL("shell32.dll")
	procSHFileOperationW = shell32.NewProc("SHFileOperationW")
)

// trash moves a file to the Windows Recycle Bin using SHFileOperationW.
func (e *Executor) trash(src string) error {
	absPath, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// SHFileOperationW requires double-null-terminated string
	from, err := windows.UTF16FromString(absPath)
	if err != nil {
		return fmt.Errorf("failed to encode path: %w", err)
	}
	// Append extra null terminator
	from = append(from, 0)

	op := shFileOpStruct{
		wFunc:  foDelete,
		pFrom:  &from[0],
		fFlags: fofAllowUndo | fofNoConfirmation | fofNoErrorUI | fofSilent,
	}

	ret, _, _ := procSHFileOperationW.Call(uintptr(unsafe.Pointer(&op)))
	if ret != 0 {
		return fmt.Errorf("SHFileOperationW failed with code %d", ret)
	}

	return nil
}
