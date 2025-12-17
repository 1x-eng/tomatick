package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// HTTPDispatcher implements the Dispatcher interface for HTTP webhooks
type HTTPDispatcher struct {
	urls       []string
	client     *http.Client
	workerPool chan struct{} // Semaphore to limit concurrent dispatches
	logger     *log.Logger
	logFile    *os.File
	wg         sync.WaitGroup
}

type dispatchLog struct {
	Timestamp string      `json:"timestamp"`
	Event     EventType   `json:"event"`
	Target    string      `json:"target"`
	Attempt   int         `json:"attempt"`
	Status    string      `json:"status"`
	Code      int         `json:"code,omitempty"`
	Error     string      `json:"error,omitempty"`
	Duration  string      `json:"duration"`
}

// NewHTTPDispatcher creates a new dispatcher with the given webhook URLs and log path
func NewHTTPDispatcher(urls []string, logDir string) *HTTPDispatcher {
	// Ensure log directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Warning: Failed to create webhook log directory: %v\n", err)
	}

	logPath := filepath.Join(logDir, "webhooks.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Warning: Failed to open webhook log file: %v\n", err)
	}

	return &HTTPDispatcher{
		urls: urls,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		workerPool: make(chan struct{}, 10),
		logger:     log.New(f, "", 0),
		logFile:    f,
	}
}

// Dispatch sends an event to all configured webhooks asynchronously
func (d *HTTPDispatcher) Dispatch(eventType EventType, data map[string]string) {
	if len(d.urls) == 0 {
		return
	}

	payload := EventPayload{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		d.logAttempt(eventType, "marshal_error", 0, "Failed", 0, err, 0)
		return
	}

	for _, url := range d.urls {
		d.wg.Add(1)
		go d.sendWithRetry(url, eventType, body)
	}
}

func (d *HTTPDispatcher) sendWithRetry(url string, eventType EventType, body []byte) {
	defer d.wg.Done()

	// Acquire worker slot
	d.workerPool <- struct{}{}
	defer func() { <-d.workerPool }()

	const maxRetries = 3
	baseDelay := 1 * time.Second

	var lastErr error
	var statusCode int

	for attempt := 1; attempt <= maxRetries+1; attempt++ {
		start := time.Now()
		statusCode, lastErr = d.sendOnce(url, body)
		duration := time.Since(start)

		if lastErr == nil && statusCode >= 200 && statusCode < 300 {
			d.logAttempt(eventType, url, attempt, "Success", statusCode, nil, duration)
			return
		}

		// Log failure for this attempt
		status := "Failed"
		if attempt <= maxRetries {
			status = "Retrying"
		}
		d.logAttempt(eventType, url, attempt, status, statusCode, lastErr, duration)

		if attempt <= maxRetries {
			// Exponential backoff: 1s, 2s, 4s
			sleepDuration := baseDelay * time.Duration(math.Pow(2, float64(attempt-1)))
			time.Sleep(sleepDuration)
		}
	}
}

// Wait blocks until all active webhooks have been sent or exhausted their retries
func (d *HTTPDispatcher) Wait() {
	d.wg.Wait()
}

func (d *HTTPDispatcher) sendOnce(url string, body []byte) (int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Tomatick-Webhook-Dispatcher/1.0")

	resp, err := d.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func (d *HTTPDispatcher) logAttempt(event EventType, target string, attempt int, status string, code int, err error, duration time.Duration) {
	if d.logger == nil {
		return
	}

	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	logEntry := dispatchLog{
		Timestamp: time.Now().Format(time.RFC3339),
		Event:     event,
		Target:    target,
		Attempt:   attempt,
		Status:    status,
		Code:      code,
		Error:     errStr,
		Duration:  duration.String(),
	}

	jsonBytes, _ := json.Marshal(logEntry)
	d.logger.Println(string(jsonBytes))
}

// Close gracefully closes the log file
func (d *HTTPDispatcher) Close() {
	if d.logFile != nil {
		d.logFile.Close()
	}
}
