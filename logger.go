package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Logger struct {
	URL       string
	AuthToken string
	Meta      map[string]interface{}
	Client    *http.Client
}

type LogLevel string

const (
	LevelDebug LogLevel = "debug"
	LevelInfo  LogLevel = "info"
	LevelWarn  LogLevel = "warn"
	LevelError LogLevel = "error"
)

func New(token string, meta map[string]interface{}) *Logger {
	return &Logger{
		URL:       "https://logs.decodx.co/gelf",
		AuthToken: token,
		Meta:      meta,
		Client:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (l *Logger) log(level LogLevel, msg string, fields map[string]interface{}) {
	payload := map[string]interface{}{
		"version":       "1.1",
		"host":          l.Meta["appName"],
		"short_message": msg,
		"level":         l.mapLevel(level),
		"timestamp":     float64(time.Now().UnixNano()) / 1e9,
	}

	for k, v := range l.Meta {
		payload["_"+k] = v
	}
	for k, v := range fields {
		payload["_"+k] = v
	}

	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", l.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if l.AuthToken != "" {
		req.Header.Set("Authorization", "Basic "+l.AuthToken)
	}

	resp, err := l.Client.Do(req)
	if err != nil {
		fmt.Println("Failed to send log:", err)
		return
	}
	defer resp.Body.Close()
}

func (l *Logger) Info(msg string, fields map[string]interface{})  { l.log(LevelInfo, msg, fields) }
func (l *Logger) Warn(msg string, fields map[string]interface{})  { l.log(LevelWarn, msg, fields) }
func (l *Logger) Error(msg string, fields map[string]interface{}) { l.log(LevelError, msg, fields) }
func (l *Logger) Debug(msg string, fields map[string]interface{}) { l.log(LevelDebug, msg, fields) }

func (l *Logger) mapLevel(level LogLevel) int {
	switch level {
	case LevelDebug:
		return 7
	case LevelInfo:
		return 6
	case LevelWarn:
		return 4
	case LevelError:
		return 3
	default:
		return 6
	}
}
