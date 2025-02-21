//go:build !windows

package entry

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

// fileUserGroup retrieves the file owner and group names for the given Entry.
// Falls back to UID/GID if names cannot be resolved.
func fileUserGroup(e Entry) (string, string) {
	stat, ok := e.info.Sys().(*syscall.Stat_t)
	if !ok {
		fmt.Fprintf(os.Stderr, "warning: cannot get syscall.Stat_t for %s\n", e.info.Name())
		return "unknown", "unknown"
	}

	uid := fmt.Sprint(stat.Uid)
	usr := uid
	if u, err := user.LookupId(uid); err == nil {
		usr = u.Username
	}

	gid := fmt.Sprint(stat.Gid)
	group := gid
	if g, err := user.LookupGroupId(gid); err == nil {
		group = g.Name
	}

	return usr, group
}
