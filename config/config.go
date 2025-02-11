package config

import (
	"fmt"
	"strconv"
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

	return &Config{
		TomatickMementoDuration: pomoDuration,
		ShortBreakDuration:      shortBreak,
		LongBreakDuration:       longBreak,
		CyclesBeforeLongBreak:   cycles,
		MEMAIAPIToken:           getEnvVar("MEM_AI_API_TOKEN"),
		ContextDir:              contextDir,
		PerplexityAPIToken:      getEnvVar("PERPLEXITY_API_TOKEN"),
		UserName:                getEnvVar("USER_NAME"),
	}, nil
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
