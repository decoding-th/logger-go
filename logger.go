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

type logRequest struct {
	Request     interface{} `json:"request"`
	Data        interface{} `json:"data"`
	Message     string      `json:"message"`
	Level       string      `json:"level"`
	AppName     string      `json:"appName"`
	Version     string      `json:"version"`
	ServiceName string      `json:"serviceName"`
}

func New(token string, meta map[string]interface{}) *Logger {
	return &Logger{
		URL:       "https://gateway.decodx.co/send-log-message",
		AuthToken: token,
		Meta:      meta,
		Client:    &http.Client{Timeout: 5 * time.Second},
	}
}

func (l *Logger) log(level LogLevel, msg string, fields map[string]interface{}) {
	payload := logRequest{
		Message:     msg,
		Level:       string(level),
		AppName:     fmt.Sprintf("%v", l.Meta["appName"]),
		ServiceName: fmt.Sprintf("%v", l.Meta["serviceName"]),
		Version:     fmt.Sprintf("%v", l.Meta["version"]),
		Request:     fields["request"],
		Data:        fields["data"],
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Failed to marshal payload:", err)
		return
	}

	req, err := http.NewRequest("POST", l.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	if l.AuthToken == "" {
		fmt.Println("No auth token provided")
		return
	}
	req.Header.Set("Authorization", l.AuthToken)

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
