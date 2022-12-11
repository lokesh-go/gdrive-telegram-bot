package scopes

import (
	"google.golang.org/api/drive/v3"
)

// Get ...
func Get() (scope string) {
	// Drive scope
	driveScope := drive.DriveScope

	// Returns
	return driveScope
}
