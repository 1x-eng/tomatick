package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type LLMConfig struct {
	Provider   string // e.g., "openai", "anthropic"
	Model      string // e.g., "gpt-4", "claude-3"
	APIKey     string
	Parameters map[string]any // Provider-specific parameters
}

type CopilotConfig struct {
	Enabled     bool
	DefaultMode string // "lead" or "follow"
	LLM         LLMConfig
}

type Config struct {
	TomatickMementoDuration time.Duration
	ShortBreakDuration      time.Duration
	LongBreakDuration       time.Duration
	CyclesBeforeLongBreak   int
	MEMAIAPIToken           string
	ContextDir              string
	Copilot                 CopilotConfig
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

	// Load Copilot Configuration
	copilotConfig := CopilotConfig{
		Enabled:     getEnvBool("COPILOT_ENABLED", true),
		DefaultMode: getEnvString("COPILOT_MODE", "follow"),
		LLM: LLMConfig{
			Provider: getEnvString("LLM_PROVIDER", "openai"),
			Model:    getEnvString("LLM_MODEL", "gpt-4"),
			APIKey:   os.Getenv("OPENAI_API_KEY"),
			Parameters: map[string]any{
				"temperature":       getEnvFloat("LLM_TEMPERATURE", 0.7),
				"max_tokens":        getEnvInt("LLM_MAX_TOKENS", 150),
				"top_p":             getEnvFloat("LLM_TOP_P", 1.0),
				"presence_penalty":  getEnvFloat("LLM_PRESENCE_PENALTY", 0.0),
				"frequency_penalty": getEnvFloat("LLM_FREQUENCY_PENALTY", 0.0),
			},
		},
	}

	return &Config{
		TomatickMementoDuration: pomoDuration,
		ShortBreakDuration:      shortBreak,
		LongBreakDuration:       longBreak,
		CyclesBeforeLongBreak:   cyclesBeforeLongBreak,
		MEMAIAPIToken:           os.Getenv("MEM_AI_API_TOKEN"),
		ContextDir:              contextDir,
		Copilot:                 copilotConfig,
	}, nil
}

// Helper functions for environment variables
func getEnvBool(key string, defaultVal bool) bool {
	if val, exists := os.LookupEnv(key); exists {
		b, err := strconv.ParseBool(val)
		if err == nil {
			return b
		}
	}
	return defaultVal
}

func getEnvString(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}

func getEnvFloat(key string, defaultVal float64) float64 {
	if val, exists := os.LookupEnv(key); exists {
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f
		}
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
