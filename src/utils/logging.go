package utils

import (
	"log/slog"
)

const ERROR_HEADER = "Syntax Error"

func PrintParserErrors(logger *slog.Logger, errors []string) {
	for _, msg := range errors {
		logger.Error(msg, "type", "syntax")
	}
}
