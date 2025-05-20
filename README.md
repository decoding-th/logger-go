# logger-go

A minimal Graylog logger for Go (Golang), sending structured logs over GELF HTTP.  
Supports metadata, log levels, and Basic Auth headers for secure logging.

---

## ✨ Features

-   ✅ GELF 1.1 format
-   ✅ Info, Warn, Error, Debug log levels
-   ✅ Custom metadata support (`service`, `appName`, `version`, etc.)
-   ✅ Sends logs over HTTPS with optional Basic Auth (base64)
-   ✅ Lightweight and production-ready

---

## 📦 Installation

```bash
go get github.com/decoding-th/logger-go
```

---

## 🚀 Usage

```go
import (
  "encoding/base64"
  "github.com/decoding-th/logger-go/logger"
)

func main() {
  token := base64.StdEncoding.EncodeToString([]byte("a:b"))

  log := logger.New(token, map[string]interface{}{
    "service":  "api",
    "appName":  "pairing-content",
    "version":  "1.0.0",
  })

  log.Info("User logged in", map[string]interface{}{
    "userId": "abc123",
  })
}
```

---

## 📤 Output Format (GELF)

```json
{
    "version": "1.1",
    "host": "pairing-content",
    "short_message": "User logged in",
    "level": 6,
    "_service": "api",
    "_appName": "pairing-content",
    "_version": "1.0.0",
    "_userId": "abc123"
}
```

---

## 🔒 Authorization

This library supports **Basic Auth** via a base64 token (`Authorization: Basic ...`).

Use this to generate token:

```go
base64.StdEncoding.EncodeToString([]byte("username:password"))
```

---

## 📘 Log Levels

| Method        | Graylog Level |
| ------------- | ------------- |
| `log.Debug()` | 7             |
| `log.Info()`  | 6             |
| `log.Warn()`  | 4             |
| `log.Error()` | 3             |

---

## 🛡 License

MIT — by the decoding.co.th team
