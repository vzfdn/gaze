package entry

import (
	"fmt"
	"log"
	"os/user"
	"syscall"
	"unsafe"
)

func UnixOwner(e Entry) string {
	stat := e.info.Sys().(*syscall.Stat_t)
	u, err := user.LookupId(fmt.Sprint(stat.Uid))
	if err != nil {
		log.Fatal(err)
	}
	return u.Name
}

func UnixGroup(e Entry) string {
	stat := e.info.Sys().(*syscall.Stat_t)
	g, err := user.LookupGroupId(fmt.Sprint(stat.Gid))
	if err != nil {
		log.Fatal(err)
	}
	return g.Name
}

func GetTerminalWidth() int {
	ts := &struct {
		Row uint16
		Col uint16
		X   uint16
		Y   uint16
	}{}
	retCode, _, _ := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ts)),
	)
	if int(retCode) == -1 {
		return 0
	}
	return int(ts.Col)
}
