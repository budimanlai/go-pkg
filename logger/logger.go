package logger

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

var (
	ShowOutput = true
	ShowDebug  = true
)

// Vardump prints a formatted JSON representation of the given value to standard output.
// It uses json.MarshalIndent with 2-space indentation for readable output.
// Any marshaling errors are silently ignored.
//
// Parameters:
//   - v: any value to be printed as formatted JSON
//
// Example:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//	user := User{Name: "John", Age: 30}
//	Vardump(user)
//	// Output:
//	// {
//	//   "Name": "John",
//	//   "Age": 30
//	// }
func Vardump(v any) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

// Printf formats and prints a log message with a timestamp prefix.
// The message is only printed if ShowOutput is true.
// The format string and arguments follow the same conventions as fmt.Sprintf.
// Each log entry is prefixed with the current timestamp in "2006-01-02 15:04:05" format.
//
// Parameters:
//   - format: A format string following fmt.Sprintf conventions
//   - args: Variadic arguments to be formatted according to the format string
//
// Example:
//
//	Printf("User %s logged in at %d", username, loginTime)
func Printf(format string, args ...interface{}) {
	if ShowOutput {
		text := fmt.Sprintf(format, args...)
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] %s\n", now, text)
	}
}

// PrintHex prints the hexadecimal representation of the provided byte slice to stdout.
// The output includes a timestamp in the format "2006-01-02 15:04:05" followed by the
// hex-encoded data. The function only produces output if the ShowOutput flag is set to true.
func PrintHex(data []byte) {
	if ShowOutput {
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] %s\n", now, hex.EncodeToString(data))
	}
}

// Debugf formats and prints a debug message with timestamp if ShowDebug is enabled.
// The message is formatted according to the format specifier and arguments provided.
// Output format: [YYYY-MM-DD HH:MM:SS] DEBUG: <formatted message>
//
// Parameters:
//   - format: A format string following fmt.Sprintf conventions
//   - args: Variable number of arguments to be formatted according to the format string
//
// The function will only produce output when the global ShowDebug flag is set to true.
func Debugf(format string, args ...interface{}) {
	if ShowDebug {
		text := fmt.Sprintf(format, args...)
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] DEBUG: %s\n", now, text)
	}
}

// Fatalf logs a formatted fatal message with timestamp and terminates the program with exit code 1.
// The function formats the message according to the format specifier and arguments provided,
// prepends it with the current timestamp in "2006-01-02 15:04:05" format and "FATAL" level,
// then calls log.Fatalf which prints the message and exits the program.
//
// Parameters:
//   - format: A format string following fmt.Printf conventions
//   - args: Variable number of arguments to be formatted according to the format string
//
// Note: This function does not return as it terminates the program execution.
func Fatalf(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	now := time.Now().Format("2006-01-02 15:04:05")
	log.Fatalf("[%s] FATAL: %s\n", now, text)
}

// Errorf logs an error message with formatted arguments.
// It formats the message using fmt.Sprintf with the provided format string and arguments,
// then prints it to standard output with an ERROR prefix and current timestamp.
// The timestamp format is "2006-01-02 15:04:05".
//
// Parameters:
//   - format: A format string following fmt.Sprintf conventions
//   - args: Variadic arguments to be formatted according to the format string
//
// Example:
//
//	logger.Errorf("failed to connect to database: %v", err)
func Errorf(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] ERROR: %s\n", now, text)
}

// Infof logs an informational message with formatted output.
// It formats the message according to the format specifier and arguments,
// prefixes it with a timestamp in "2006-01-02 15:04:05" format and "INFO" level,
// then outputs to standard output.
//
// Parameters:
//   - format: A format string following fmt.Sprintf conventions
//   - args: Variable arguments to be formatted according to the format string
//
// Example:
//
//	Infof("User %s logged in at %d", username, loginTime)
func Infof(format string, args ...interface{}) {
	text := fmt.Sprintf(format, args...)
	now := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] INFO: %s\n", now, text)
}
