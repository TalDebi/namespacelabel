package internal

import (
	"strings"
)

// IsManagementLabel check if label is a management label
func IsManagementLabel(label string) bool {
	return strings.HasPrefix(label, ManagementLabelPrefix)
}
