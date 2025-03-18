//go:build !windows

package entry

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

var uidCache = make(map[string]string)
var gidCache = make(map[string]string)

// userGroup retrieves the file owner and group names for the Entry.
// Falls back to UID/GID if names cannot be resolved.
func userGroup(e Entry) (string, string) {
	stat, ok := e.info.Sys().(*syscall.Stat_t)
	if !ok {
		fmt.Fprintf(os.Stderr, "warning: cannot get syscall.Stat_t for %s\n", e.info.Name())
		return "unknown", "unknown"
	}
	uid := fmt.Sprint(stat.Uid)
	usr, ok := uidCache[uid]
	if !ok {
		if u, err := user.LookupId(uid); err == nil {
			usr = u.Username
			uidCache[uid] = usr
		} else {
			fmt.Fprintf(os.Stderr, "warning: cannot resolve UID %s: %v\n", uid, err)
			usr = uid
		}
	}
	gid := fmt.Sprint(stat.Gid)
	group, ok := gidCache[gid]
	if !ok {
		if g, err := user.LookupGroupId(gid); err == nil {
			group = g.Name
			gidCache[gid] = group
		} else {
			fmt.Fprintf(os.Stderr, "warning: cannot resolve GID %s: %v\n", gid, err)
			group = gid
		}
	}
	return usr, group
}
