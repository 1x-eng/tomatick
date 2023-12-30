package markdown

import (
	"fmt"
	"strings"
	"time"
)

func FormatTasks(tasks []string) string {
	var md strings.Builder
	for _, task := range tasks {
		md.WriteString(fmt.Sprintf("- [ ] %s\n", task))
	}
	return md.String()
}

func FormatCycleSummary(tasks, reflections string) string {
	return fmt.Sprintf(
		"## Tomatick Cycle: %s\n\n### Tasks\n%s\n\n### Reflections\n%s\n***\n",
		time.Now().Format("02-01-2006 15:04"), tasks, reflections)
}
