package logger

import (
	"fmt"
	"os"
	"strings"
)

const (
	ansiReset  = "\x1b[0m"
	ansiDim    = "\x1b[2m"
	ansiRed    = "\x1b[31m"
	ansiGreen  = "\x1b[32m"
	ansiYellow = "\x1b[33m"
	ansiCyan   = "\x1b[36m"
)

var debugEnabled bool

func Init() {
	debugEnabled = os.Getenv("MINIPAAS_LOGGER_LEVEL") == "DEBUG"
}

func Debug(message string) {
	if debugEnabled {
		write("•", ansiDim, message, ansiDim)
	}
}

func Info(message string) {
	write("»", ansiCyan, message, "")
}

func Success(message string) {
	write("✓", ansiGreen, message, "")
}

func Warn(message string) {
	write("!", ansiYellow, message, "")
}

func Error(message string, err error) {
	write("✗", ansiRed, withError(message, err), ansiRed)
}

func Command(args []string) {
	write("$", ansiDim, strings.Join(args, " "), ansiDim)
}

func Panic(message string) {
	write("✗", ansiRed, message, ansiRed)
	os.Exit(1)
}

func PanicErr(message string, err error) {
	write("✗", ansiRed, withError(message, err), ansiRed)
	os.Exit(1)
}

func write(symbol, symbolColor, message, messageColor string) {
	if isTerminal(os.Stderr) {
		fmt.Fprintf(os.Stderr, "%s%s%s %s%s%s\n", symbolColor, symbol, ansiReset, messageColor, message, ansiReset)
		return
	}
	fmt.Fprintf(os.Stderr, "%s %s\n", symbol, message)
}

func withError(message string, err error) string {
	if err == nil {
		return message
	}
	return message + ": " + err.Error()
}

func isTerminal(f *os.File) bool {
	info, err := f.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}
