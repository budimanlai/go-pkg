package logger

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

var (
	ShowOutput = true
	ShowDebug  = true
)

func Vardump(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func Printf(format string, args ...interface{}) {
	if ShowOutput {
		text := fmt.Sprintf(format, args...)
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] %s\n", now, text)
	}
}

func PrintHex(data []byte) {
	if ShowOutput {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] %s\n", now, hex.EncodeToString(data))
	}
}

func Debugf(format string, args ...interface{}) {
	if ShowDebug {
		text := fmt.Sprintf(format, args...)
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] %s\n", now, text)
	}
}
