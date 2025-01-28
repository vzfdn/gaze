//go:build windows

package entry

import (
	"syscall"
	"unsafe"
)

var (
	modadvapi32         = syscall.NewLazyDLL("advapi32.dll")
	procGetSecurityInfo = modadvapi32.NewProc("GetSecurityInfo")
)

// Get file owner using Windows SID
func FileOwner(e Entry) string {
	path, err := syscall.UTF16PtrFromString(e.info.Name())
	if err != nil {
		return ""
	}

	var sid *syscall.SID
	var secDesc uintptr

	err = syscall.GetNamedSecurityInfo(
		path,
		syscall.SE_FILE_OBJECT,
		syscall.OWNER_SECURITY_INFORMATION,
		&sid,
		nil,
		nil,
		nil,
		&secDesc,
	)

	if err != nil {
		return ""
	}

	sidStr, _ := sid.String()
	return sidStr
}

// Get file group using Windows SID
func FileGroup(e Entry) string {
	path, err := syscall.UTF16PtrFromString(e.info.Name())
	if err != nil {
		return ""
	}

	var sid *syscall.SID
	var secDesc uintptr

	err = syscall.GetNamedSecurityInfo(
		path,
		syscall.SE_FILE_OBJECT,
		syscall.GROUP_SECURITY_INFORMATION,
		nil,
		&sid,
		nil,
		nil,
		&secDesc,
	)

	if err != nil {
		return ""
	}

	sidStr, _ := sid.String()
	return sidStr
}
