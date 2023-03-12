package helpers

import "strings"

func IsValidString(value string, err error) bool {
	return err == nil && len(strings.TrimSpace(value)) > 0
}
