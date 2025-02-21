//go:build windows

package entry

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

var sidCache = make(map[string]string)

// fileUserGroup retrieves the file's owner and group names for the given Entry.
// Falls back to SID strings if names cannot be resolved.
func fileUserGroup(e Entry) (string, string) {
	path := filepath.Join(e.path, e.info.Name())
	securityFlags := windows.OWNER_SECURITY_INFORMATION | windows.GROUP_SECURITY_INFORMATION

	sd, err := windows.GetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.SECURITY_INFORMATION(securityFlags),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: cannot get security info for %s: %v\n", path, err)
		return "unknown", "unknown"
	}

	owner := "unknown"
	if ownerSid, _, err := sd.Owner(); err == nil && ownerSid != nil {
		owner = sidToName(ownerSid)
	}

	group := "unknown"
	if groupSid, _, err := sd.Group(); err == nil && groupSid != nil {
		group = sidToName(groupSid)
	}

	return owner, group
}

// sidToName converts a Windows SID into a human-readable account name.
func sidToName(sid *windows.SID) string {
	if sid == nil {
		return "unknown"
	}

	// Check cache first
	sidStr := sid.String()
	if name, ok := sidCache[sidStr]; ok {
		return name
	}

	name, _, _, err := sid.LookupAccount("")
	if err != nil {
		return sidStr
	}

	sidCache[sidStr] = name
	return name
}
