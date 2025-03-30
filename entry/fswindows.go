//go:build windows

package entry

import (
	"path/filepath"

	"golang.org/x/sys/windows"
)

// userGroup retrieves the file owner and group names for the Entry.
// Falls back to SID strings if names cannot be resolved.
func userGroup(e Entry) (string, string) {
	sidCache := make(map[string]string) // Local cache; SIDs rarely repeat in a directory
	path := filepath.Join(e.path, e.info.Name())
	securityFlags := windows.OWNER_SECURITY_INFORMATION | windows.GROUP_SECURITY_INFORMATION
	sd, err := windows.GetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.SECURITY_INFORMATION(securityFlags),
	)
	if err != nil {
		return "unknown", "unknown"
	}
	owner := "unknown"
	if ownerSid, _, err := sd.Owner(); err == nil && ownerSid != nil {
		owner = sidToName(ownerSid, sidCache)
	}
	group := "unknown"
	if groupSid, _, err := sd.Group(); err == nil && groupSid != nil {
		group = sidToName(groupSid, sidCache)
	}
	return owner, group
}

// sidToName converts a Windows SID to a human-readable account name.
// Uses the provided cache to store resolved names.
func sidToName(sid *windows.SID, cache map[string]string) string {
	if sid == nil {
		return "unknown"
	}
	sidStr := sid.String()
	if name, ok := cache[sidStr]; ok {
		return name
	}
	name, _, _, err := sid.LookupAccount("")
	if err != nil {
		return sidStr
	}
	cache[sidStr] = name
	return name
}
