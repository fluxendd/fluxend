package pkg

import (
	"fmt"
)

func ConvertBytesToKiloBytes(bytes int) int {
	return bytes / 1024
}

func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func FormatPercentage(used, total uint64) string {
	if total == 0 {
		return "0%"
	}
	percentage := float64(used) / float64(total) * 100
	return fmt.Sprintf("%.1f%%", percentage)
}
