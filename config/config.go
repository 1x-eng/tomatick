package config

import (
	"fmt"
	"os"
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
}

func LoadConfig() (*Config, error) {
	pomoDuration, _ := time.ParseDuration("25m")
	shortBreak, _ := time.ParseDuration("5m")
	longBreak, _ := time.ParseDuration("15m")
	cyclesBeforeLongBreak := 4

	if val, exists := os.LookupEnv("POMODORO_DURATION"); exists {
		if d, err := time.ParseDuration(val); err == nil {
			pomoDuration = d
		}
	}
	if val, exists := os.LookupEnv("SHORT_BREAK_DURATION"); exists {
		if d, err := time.ParseDuration(val); err == nil {
			shortBreak = d
		}
	}
	if val, exists := os.LookupEnv("LONG_BREAK_DURATION"); exists {
		if d, err := time.ParseDuration(val); err == nil {
			longBreak = d
		}
	}
	if val, exists := os.LookupEnv("CYCLES_BEFORE_LONGBREAK"); exists {
		if n, err := strconv.Atoi(val); err == nil {
			cyclesBeforeLongBreak = n
		}
	}

	contextDir := os.Getenv("TOMATICK_CONTEXT_DIR")
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
		CyclesBeforeLongBreak:   cyclesBeforeLongBreak,
		MEMAIAPIToken:           os.Getenv("MEM_AI_API_TOKEN"),
		ContextDir:              contextDir,
	}, nil
}
