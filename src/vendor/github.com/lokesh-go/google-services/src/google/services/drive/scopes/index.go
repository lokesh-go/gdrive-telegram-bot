package scopes

import (
	"google.golang.org/api/drive/v3"
)

// Get ...
func Get() (scope string) {
	// Drive read only scope
	readOnly := drive.DriveReadonlyScope

	// Returns
	return readOnly
}
