package config

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	TomatickMementoDuration time.Duration
	ShortBreakDuration      time.Duration
	LongBreakDuration       time.Duration
	CyclesBeforeLongBreak   int
	MEMAIAPIToken           string
	ContextDir              string
	PerplexityAPIToken      string
	UserName                string
	WorkApps                []string
	Features                Features
}

type Features struct {
	BreakMonitoring bool
}

func LoadConfig() (*Config, error) {
	if err := validateEnvVars(); err != nil {
		return nil, err
	}

	pomoDuration, err := parseDurationEnv("POMODORO_DURATION", "25m")
	if err != nil {
		return nil, fmt.Errorf("invalid POMODORO_DURATION: %w", err)
	}

	shortBreak, err := parseDurationEnv("SHORT_BREAK_DURATION", "5m")
	if err != nil {
		return nil, fmt.Errorf("invalid SHORT_BREAK_DURATION: %w", err)
	}

	longBreak, err := parseDurationEnv("LONG_BREAK_DURATION", "15m")
	if err != nil {
		return nil, fmt.Errorf("invalid LONG_BREAK_DURATION: %w", err)
	}

	cycles, err := parseIntEnv("CYCLES_BEFORE_LONGBREAK", 4)
	if err != nil {
		return nil, fmt.Errorf("invalid CYCLES_BEFORE_LONGBREAK: %w", err)
	}

	contextDir := getEnvVar("TOMATICK_CONTEXT_DIR")
	if contextDir == "" {
		contextDir = getDefaultContextDir()
	}

	if err := ensureDirectoryExists(contextDir); err != nil {
		return nil, fmt.Errorf("failed to create context directory: %w", err)
	}

	// Get work apps from environment
	workApps := getWorkApps()

	// Determine available features based on OS
	features := Features{
		BreakMonitoring: runtime.GOOS == "darwin", // Only enable on macOS
	}

	return &Config{
		TomatickMementoDuration: pomoDuration,
		ShortBreakDuration:      shortBreak,
		LongBreakDuration:       longBreak,
		CyclesBeforeLongBreak:   cycles,
		MEMAIAPIToken:           getEnvVar("MEM_AI_API_TOKEN"),
		ContextDir:              contextDir,
		PerplexityAPIToken:      getEnvVar("PERPLEXITY_API_TOKEN"),
		UserName:                getEnvVar("USER_NAME"),
		WorkApps:                workApps,
		Features:                features,
	}, nil
}

// getWorkApps gets the list of work apps from environment variable
func getWorkApps() []string {
	defaultApps := []string{
		"Code",          // VS Code
		"Cursor",        // Cursor Editor
		"iTerm2",        // Terminal
		"Insomnia",      // API Testing
		"pgAdmin 4",     // PostgreSQL Admin
		"pgAdmin",       // PostgreSQL Admin
		"Chrome",        // Web Browser
		"Google Chrome", // Web Browser
		"Terminal",      // Built-in Terminal
	}

	workAppsEnv := getEnvVar("WORK_APPS")
	if workAppsEnv == "" {
		return defaultApps
	}

	// Split by comma and trim spaces
	customApps := strings.Split(workAppsEnv, ",")
	for i, app := range customApps {
		customApps[i] = strings.TrimSpace(app)
	}

	return customApps
}

func parseDurationEnv(key, defaultValue string) (time.Duration, error) {
	value := getEnvVar(key)
	if value == "" {
		value = defaultValue
	}
	return time.ParseDuration(value)
}

func parseIntEnv(key string, defaultValue int) (int, error) {
	value := getEnvVar(key)
	if value == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(value)
}
