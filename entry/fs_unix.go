//go:build !windows

package entry

import (
	"fmt"
	"log"
	"os/user"
	"syscall"
)

func FileOwner(e Entry) string {
	stat := e.info.Sys().(*syscall.Stat_t)
	u, err := user.LookupId(fmt.Sprint(stat.Uid))
	if err != nil {
		log.Fatal(err)
	}
	return u.Username
}

func FileGroup(e Entry) string {
	stat := e.info.Sys().(*syscall.Stat_t)
	g, err := user.LookupGroupId(fmt.Sprint(stat.Gid))
	if err != nil {
		log.Fatal(err)
	}
	return g.Name
}
