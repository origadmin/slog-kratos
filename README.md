# slog-kratos

[![Go Report Card](https://goreportcard.com/badge/github.com/origadmin/slog-kratos)](https://goreportcard.com/report/github.com/origadmin/slog-kratos) [![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/origadmin/slog-kratos) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A simple and efficient adapter to use Go's standard structured logging library (`slog`) as a logger within the [Kratos](https://github.com/go-kratos/kratos) framework.

This library implements the `log.Logger` interface from the Kratos project, allowing you to seamlessly integrate `slog` into your Kratos applications for powerful, structured, and level-based logging.

## ✨ Features

- **Seamless Integration**: Implements the `log.Logger` interface for direct use in Kratos.
- **Structured Logging**: Leverage the full power of `slog` for structured, key-value pair logging.
- **Level-based Logging**: Supports different log levels (Debug, Info, Warn, Error).
- **Customizable**: Easily configure the underlying `slog.Handler` (e.g., `slog.TextHandler`, `slog.JSONHandler`) to control log format and output.
- **High Performance**: Built on Go's standard library, ensuring minimal overhead.

## 🚀 Installation

To install `slog-kratos`, use `go get`:

```sh
go get github.com/origadmin/slog-kratos
```

## 💡 Usage

Integrating the logger into your Kratos application is straightforward. You can create a new `slog` logger and pass it to the Kratos application.

Here is a basic example:

```go
package main

import (
	"context"
	"log/slog"
	"os"

	kratoslog "github.com/go-kratos/kratos/v2/log"
	slogkratos "github.com/origadmin/slog-kratos"
)

func main() {
	// 1. Create a standard slog handler (e.g., JSONHandler)
	slogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Set your desired log level
	})

	// 2. Create the slog-kratos logger adapter
	logger := slogkratos.NewLogger(slogHandler)

	// 3. Use it as a Kratos logger
	kratosLogger := kratoslog.NewHelper(logger)

	// 4. Log messages with key-value pairs
	kratosLogger.WithContext(context.Background()).Info("message", "hello world", "user", "kratos")
	kratosLogger.WithContext(context.Background()).Warnw("key1", "value1", "key2", 123)
	kratosLogger.Error("This is an error message")

	// Example Output (JSON):
	// {"time":"2023-10-27T10:00:00.000Z","level":"INFO","msg":"hello world","user":"kratos"}
	// {"time":"2023-10-27T10:00:00.000Z","level":"WARN","msg":"","key1":"value1","key2":123}
	// {"time":"2023-10-27T10:00:00.000Z","level":"ERROR","msg":"This is an error message"}
}
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue.

Please see [CONTRIBUTING.md](.github/CONTRIBUTING.md) for guidelines.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
