package output

import (
	"fmt"
	"strings"
)

const DefaultFormat = "yaml"

func NormalizeFormat(format string) string {
	if strings.TrimSpace(format) == "" {
		return DefaultFormat
	}

	return format
}

func ValidateFormat(format string) error {
	format = NormalizeFormat(format)

	switch format {
	case "yaml":
		return nil
	default:
		return fmt.Errorf("invalid output format %q: supported format is yaml", format)
	}
}
