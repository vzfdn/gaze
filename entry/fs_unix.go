//go:build !windows

package entry

import (
	"fmt"
	"os/user"
	"syscall"
)

// FileUserGroup retrieves the file owner and group names for the given Entry.
func FileUserGroup(e Entry) string {
	stat, ok := e.info.Sys().(*syscall.Stat_t)
	if !ok {
		return "0  0" 
	}
	
	uidStr := fmt.Sprint(stat.Uid)
	gidStr := fmt.Sprint(stat.Gid)

	owner := uidStr
	if u, err := user.LookupId(uidStr); err == nil {
		owner = u.Username
	}

	group := gidStr
	if g, err := user.LookupGroupId(gidStr); err == nil {
		group = g.Name
	}

	return fmt.Sprintf("%s  %s", owner, group)
}
