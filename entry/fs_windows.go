//go:build windows

package entry

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// FileUserGroup retrieves the file's owner and group names for the given Entry.
func FileUserGroup(e Entry) string {
	securityFlags := windows.OWNER_SECURITY_INFORMATION | windows.GROUP_SECURITY_INFORMATION
	sd, err := windows.GetNamedSecurityInfo(
		filepath.Clean(e.info.Name()),
		windows.SE_FILE_OBJECT,
		windows.SECURITY_INFORMATION(securityFlags),
	)
	if err != nil {
		return "unknown  unknown"  
	}

	owner := "unknown"
	if ownerSid, _, err := sd.Owner(); err == nil && ownerSid != nil {
		owner = sidToName(ownerSid)
	} else if ownerSid != nil {
		owner = ownerSid.String()  
	}

	group := "unknown"
	if groupSid, _, err := sd.Group(); err == nil && groupSid != nil {
		group = sidToName(groupSid)
	} else if groupSid != nil {
		group = groupSid.String()  
	}
	
	return fmt.Sprintf("%s  %s", owner, group)
}

// sidToName converts a Windows SID into a human-readable account name ("DOMAIN\User").
func sidToName(sid *windows.SID) string {
	if sid == nil {
		return "unknown"  
	}

	name, domain, _, err := sid.LookupAccount("")
	if err != nil {
		return sid.String() 
	}
	return fmt.Sprintf("%s\\%s", domain, name)
}
