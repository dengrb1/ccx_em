//go:build windows

package appdirs

import (
	"os"
	"path/filepath"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var kernel32 = windows.NewLazySystemDLL("kernel32.dll")

var procGetCurrentPackageFamilyName = kernel32.NewProc("GetCurrentPackageFamilyName")

func windowsPackagedRoamingDir() string {
	familyName, ok := currentPackageFamilyName()
	if !ok {
		return ""
	}
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return ""
	}
	return filepath.Join(localAppData, "Packages", familyName, "LocalCache", "Roaming")
}

func currentPackageFamilyName() (string, bool) {
	if err := procGetCurrentPackageFamilyName.Find(); err != nil {
		return "", false
	}

	var length uint32
	err := getCurrentPackageFamilyName(&length, nil)
	if err == windows.APPMODEL_ERROR_NO_PACKAGE || length == 0 {
		return "", false
	}
	if err != windows.ERROR_INSUFFICIENT_BUFFER {
		return "", false
	}

	buffer := make([]uint16, length)
	err = getCurrentPackageFamilyName(&length, &buffer[0])
	if err != windows.ERROR_SUCCESS {
		return "", false
	}
	familyName := strings.TrimRight(windows.UTF16ToString(buffer), "\x00")
	return familyName, familyName != ""
}

func getCurrentPackageFamilyName(length *uint32, familyName *uint16) windows.Errno {
	ret, _, _ := procGetCurrentPackageFamilyName.Call(
		uintptr(unsafe.Pointer(length)),
		uintptr(unsafe.Pointer(familyName)),
	)
	return windows.Errno(ret)
}
