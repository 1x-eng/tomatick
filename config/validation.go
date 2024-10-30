package config

import (
	"fmt"
	"os"
	"strings"
)

type EnvVar struct {
	Name        string
	Description string
	Required    bool
}

var requiredEnvVars = []EnvVar{
	{
		Name:        "MEM_AI_API_TOKEN",
		Description: "API token for Mem.ai integration",
		Required:    true,
	},
	{
		Name:        "PERPLEXITY_API_TOKEN",
		Description: "API token for Perplexity AI integration",
		Required:    true,
	},
	{
		Name:        "TOMATICK_CONTEXT_DIR",
		Description: "Directory for storing context files",
		Required:    false, // We have a default value
	},
	{
		Name:        "POMODORO_DURATION",
		Description: "Duration of each Pomodoro cycle (e.g., 25m)",
		Required:    false, // We have a default value
	},
	{
		Name:        "SHORT_BREAK_DURATION",
		Description: "Duration of short breaks (e.g., 5m)",
		Required:    false, // We have a default value
	},
	{
		Name:        "LONG_BREAK_DURATION",
		Description: "Duration of long breaks (e.g., 15m)",
		Required:    false, // We have a default value
	},
	{
		Name:        "CYCLES_BEFORE_LONGBREAK",
		Description: "Number of cycles before a long break",
		Required:    false, // We have a default value
	},
}

func validateEnvVars() error {
	var missingVars []string

	for _, env := range requiredEnvVars {
		if env.Required && strings.TrimSpace(getEnvVar(env.Name)) == "" {
			missingVars = append(missingVars, fmt.Sprintf("- %s: %s", env.Name, env.Description))
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required environment variables:\n%s\n\nPlease set these environment variables and try again", strings.Join(missingVars, "\n"))
	}

	return nil
}

func getEnvVar(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}
