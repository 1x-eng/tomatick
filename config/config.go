package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	TomatickMementoDuration time.Duration
	ShortBreakDuration      time.Duration
	LongBreakDuration       time.Duration
	MEMAIAPIToken           string
}

func LoadConfig() (*Config, error) {
	pomoDuration, _ := time.ParseDuration("25m")
	shortBreak, _ := time.ParseDuration("5m")
	longBreak, _ := time.ParseDuration("15m")

	if val, exists := os.LookupEnv("POMODORO_DURATION"); exists {
		if d, err := time.ParseDuration(val); err == nil {
			fmt.Println("POMODORO_DURATION", d)
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

	return &Config{
		TomatickMementoDuration: pomoDuration,
		ShortBreakDuration:      shortBreak,
		LongBreakDuration:       longBreak,
		MEMAIAPIToken:           os.Getenv("MEM_AI_API_TOKEN"),
	}, nil
}
