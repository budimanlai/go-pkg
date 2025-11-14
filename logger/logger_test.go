package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// ============================================================================
// Helper Functions
// ============================================================================

// captureOutput captures stdout during function execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// ============================================================================
// Vardump Tests
// ============================================================================

func TestVardump(t *testing.T) {
	t.Run("simple_struct", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		output := captureOutput(func() {
			Vardump(Person{Name: "John", Age: 30})
		})

		if !strings.Contains(output, `"Name": "John"`) {
			t.Errorf("Expected output to contain Name field, got: %s", output)
		}
		if !strings.Contains(output, `"Age": 30`) {
			t.Errorf("Expected output to contain Age field, got: %s", output)
		}
	})

	t.Run("map", func(t *testing.T) {
		data := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		}

		output := captureOutput(func() {
			Vardump(data)
		})

		if !strings.Contains(output, "key1") {
			t.Error("Expected output to contain key1")
		}
		if !strings.Contains(output, "value1") {
			t.Error("Expected output to contain value1")
		}
	})

	t.Run("slice", func(t *testing.T) {
		data := []string{"apple", "banana", "cherry"}

		output := captureOutput(func() {
			Vardump(data)
		})

		if !strings.Contains(output, "apple") {
			t.Error("Expected output to contain apple")
		}
		if !strings.Contains(output, "banana") {
			t.Error("Expected output to contain banana")
		}
	})

	t.Run("nested_structure", func(t *testing.T) {
		type Address struct {
			City    string
			Country string
		}
		type User struct {
			Name    string
			Address Address
		}

		output := captureOutput(func() {
			Vardump(User{
				Name: "Alice",
				Address: Address{
					City:    "Jakarta",
					Country: "Indonesia",
				},
			})
		})

		if !strings.Contains(output, "Alice") {
			t.Error("Expected output to contain Alice")
		}
		if !strings.Contains(output, "Jakarta") {
			t.Error("Expected output to contain Jakarta")
		}
	})

	t.Run("empty_struct", func(t *testing.T) {
		type Empty struct{}

		output := captureOutput(func() {
			Vardump(Empty{})
		})

		if !strings.Contains(output, "{}") {
			t.Error("Expected output to contain empty object")
		}
	})
}

// ============================================================================
// Printf Tests
// ============================================================================

func TestPrintf(t *testing.T) {
	// Save original state
	originalShowOutput := ShowOutput
	defer func() { ShowOutput = originalShowOutput }()

	t.Run("basic_message", func(t *testing.T) {
		ShowOutput = true
		output := captureOutput(func() {
			Printf("Test message")
		})

		if !strings.Contains(output, "Test message") {
			t.Errorf("Expected 'Test message', got: %s", output)
		}
		// Should contain timestamp in format [YYYY-MM-DD HH:MM:SS]
		if !strings.Contains(output, "[20") {
			t.Error("Expected output to contain timestamp")
		}
	})

	t.Run("formatted_message", func(t *testing.T) {
		ShowOutput = true
		output := captureOutput(func() {
			Printf("User %s has %d points", "Alice", 100)
		})

		if !strings.Contains(output, "User Alice has 100 points") {
			t.Errorf("Expected formatted message, got: %s", output)
		}
	})

	t.Run("multiple_arguments", func(t *testing.T) {
		ShowOutput = true
		output := captureOutput(func() {
			Printf("Values: %d, %s, %v", 42, "test", true)
		})

		if !strings.Contains(output, "Values: 42, test, true") {
			t.Errorf("Expected formatted values, got: %s", output)
		}
	})

	t.Run("disabled_output", func(t *testing.T) {
		ShowOutput = false
		output := captureOutput(func() {
			Printf("This should not appear")
		})

		if output != "" {
			t.Errorf("Expected no output when ShowOutput is false, got: %s", output)
		}
	})

	t.Run("empty_message", func(t *testing.T) {
		ShowOutput = true
		output := captureOutput(func() {
			Printf("")
		})

		// Should still have timestamp
		if !strings.Contains(output, "[20") {
			t.Error("Expected timestamp even with empty message")
		}
	})
}

// ============================================================================
// PrintHex Tests
// ============================================================================

func TestPrintHex(t *testing.T) {
	// Save original state
	originalShowOutput := ShowOutput
	defer func() { ShowOutput = originalShowOutput }()

	t.Run("basic_hex_output", func(t *testing.T) {
		ShowOutput = true
		data := []byte{0x48, 0x65, 0x6c, 0x6c, 0x6f} // "Hello" in hex

		output := captureOutput(func() {
			PrintHex(data)
		})

		if !strings.Contains(output, "48656c6c6f") {
			t.Errorf("Expected hex output '48656c6c6f', got: %s", output)
		}
	})

	t.Run("empty_byte_slice", func(t *testing.T) {
		ShowOutput = true
		data := []byte{}

		output := captureOutput(func() {
			PrintHex(data)
		})

		// Should still show timestamp
		if !strings.Contains(output, "[20") {
			t.Error("Expected timestamp in output")
		}
	})

	t.Run("single_byte", func(t *testing.T) {
		ShowOutput = true
		data := []byte{0xFF}

		output := captureOutput(func() {
			PrintHex(data)
		})

		if !strings.Contains(output, "ff") {
			t.Errorf("Expected 'ff', got: %s", output)
		}
	})

	t.Run("disabled_output", func(t *testing.T) {
		ShowOutput = false
		data := []byte{0x01, 0x02, 0x03}

		output := captureOutput(func() {
			PrintHex(data)
		})

		if output != "" {
			t.Errorf("Expected no output when ShowOutput is false, got: %s", output)
		}
	})

	t.Run("multiple_bytes", func(t *testing.T) {
		ShowOutput = true
		data := []byte{0xDE, 0xAD, 0xBE, 0xEF}

		output := captureOutput(func() {
			PrintHex(data)
		})

		if !strings.Contains(output, "deadbeef") {
			t.Errorf("Expected 'deadbeef', got: %s", output)
		}
	})
}

// ============================================================================
// Debugf Tests
// ============================================================================

func TestDebugf(t *testing.T) {
	// Save original state
	originalShowDebug := ShowDebug
	defer func() { ShowDebug = originalShowDebug }()

	t.Run("basic_debug_message", func(t *testing.T) {
		ShowDebug = true
		output := captureOutput(func() {
			Debugf("Debug message")
		})

		if !strings.Contains(output, "DEBUG: Debug message") {
			t.Errorf("Expected DEBUG prefix, got: %s", output)
		}
		if !strings.Contains(output, "[20") {
			t.Error("Expected timestamp in debug output")
		}
	})

	t.Run("formatted_debug", func(t *testing.T) {
		ShowDebug = true
		output := captureOutput(func() {
			Debugf("Processing item %d of %d", 5, 10)
		})

		if !strings.Contains(output, "Processing item 5 of 10") {
			t.Errorf("Expected formatted message, got: %s", output)
		}
		if !strings.Contains(output, "DEBUG:") {
			t.Error("Expected DEBUG prefix")
		}
	})

	t.Run("disabled_debug", func(t *testing.T) {
		ShowDebug = false
		output := captureOutput(func() {
			Debugf("This debug should not appear")
		})

		if output != "" {
			t.Errorf("Expected no output when ShowDebug is false, got: %s", output)
		}
	})

	t.Run("debug_with_complex_format", func(t *testing.T) {
		ShowDebug = true
		output := captureOutput(func() {
			Debugf("Value: %v, Type: %T, Hex: %x", 42, 42, 42)
		})

		if !strings.Contains(output, "Value: 42") {
			t.Error("Expected value in output")
		}
		if !strings.Contains(output, "Type: int") {
			t.Error("Expected type in output")
		}
	})
}

// ============================================================================
// Errorf Tests
// ============================================================================

func TestErrorf(t *testing.T) {
	t.Run("basic_error_message", func(t *testing.T) {
		output := captureOutput(func() {
			Errorf("Something went wrong")
		})

		if !strings.Contains(output, "ERROR: Something went wrong") {
			t.Errorf("Expected ERROR prefix, got: %s", output)
		}
		if !strings.Contains(output, "[20") {
			t.Error("Expected timestamp in error output")
		}
	})

	t.Run("formatted_error", func(t *testing.T) {
		output := captureOutput(func() {
			Errorf("Failed to connect to %s:%d", "localhost", 3306)
		})

		if !strings.Contains(output, "Failed to connect to localhost:3306") {
			t.Errorf("Expected formatted error message, got: %s", output)
		}
		if !strings.Contains(output, "ERROR:") {
			t.Error("Expected ERROR prefix")
		}
	})

	t.Run("error_with_error_type", func(t *testing.T) {
		output := captureOutput(func() {
			err := io.EOF
			Errorf("Read failed: %v", err)
		})

		if !strings.Contains(output, "Read failed: EOF") {
			t.Errorf("Expected error message with EOF, got: %s", output)
		}
	})

	t.Run("multiple_placeholders", func(t *testing.T) {
		output := captureOutput(func() {
			Errorf("User %s failed login attempt %d", "john", 3)
		})

		if !strings.Contains(output, "User john failed login attempt 3") {
			t.Errorf("Expected formatted message, got: %s", output)
		}
	})
}

// ============================================================================
// Infof Tests
// ============================================================================

func TestInfof(t *testing.T) {
	t.Run("basic_info_message", func(t *testing.T) {
		output := captureOutput(func() {
			Infof("Application started")
		})

		if !strings.Contains(output, "INFO: Application started") {
			t.Errorf("Expected INFO prefix, got: %s", output)
		}
		if !strings.Contains(output, "[20") {
			t.Error("Expected timestamp in info output")
		}
	})

	t.Run("formatted_info", func(t *testing.T) {
		output := captureOutput(func() {
			Infof("Server listening on port %d", 8080)
		})

		if !strings.Contains(output, "Server listening on port 8080") {
			t.Errorf("Expected formatted message, got: %s", output)
		}
		if !strings.Contains(output, "INFO:") {
			t.Error("Expected INFO prefix")
		}
	})

	t.Run("info_with_multiple_args", func(t *testing.T) {
		output := captureOutput(func() {
			Infof("User %s logged in from %s at %s", "alice", "192.168.1.1", "10:30")
		})

		if !strings.Contains(output, "User alice logged in from 192.168.1.1 at 10:30") {
			t.Errorf("Expected formatted message, got: %s", output)
		}
	})

	t.Run("empty_info_message", func(t *testing.T) {
		output := captureOutput(func() {
			Infof("")
		})

		if !strings.Contains(output, "INFO:") {
			t.Error("Expected INFO prefix even with empty message")
		}
	})
}

// ============================================================================
// Global Flags Tests
// ============================================================================

func TestGlobalFlags(t *testing.T) {
	t.Run("showoutput_flag_controls_printf", func(t *testing.T) {
		originalShowOutput := ShowOutput
		defer func() { ShowOutput = originalShowOutput }()

		// Test enabled
		ShowOutput = true
		output1 := captureOutput(func() {
			Printf("Message 1")
		})
		if output1 == "" {
			t.Error("Expected output when ShowOutput is true")
		}

		// Test disabled
		ShowOutput = false
		output2 := captureOutput(func() {
			Printf("Message 2")
		})
		if output2 != "" {
			t.Error("Expected no output when ShowOutput is false")
		}
	})

	t.Run("showoutput_flag_controls_printhex", func(t *testing.T) {
		originalShowOutput := ShowOutput
		defer func() { ShowOutput = originalShowOutput }()

		data := []byte{0x01, 0x02}

		// Test enabled
		ShowOutput = true
		output1 := captureOutput(func() {
			PrintHex(data)
		})
		if output1 == "" {
			t.Error("Expected output when ShowOutput is true")
		}

		// Test disabled
		ShowOutput = false
		output2 := captureOutput(func() {
			PrintHex(data)
		})
		if output2 != "" {
			t.Error("Expected no output when ShowOutput is false")
		}
	})

	t.Run("showdebug_flag_controls_debugf", func(t *testing.T) {
		originalShowDebug := ShowDebug
		defer func() { ShowDebug = originalShowDebug }()

		// Test enabled
		ShowDebug = true
		output1 := captureOutput(func() {
			Debugf("Debug 1")
		})
		if output1 == "" {
			t.Error("Expected output when ShowDebug is true")
		}

		// Test disabled
		ShowDebug = false
		output2 := captureOutput(func() {
			Debugf("Debug 2")
		})
		if output2 != "" {
			t.Error("Expected no output when ShowDebug is false")
		}
	})

	t.Run("errorf_always_outputs", func(t *testing.T) {
		originalShowOutput := ShowOutput
		defer func() { ShowOutput = originalShowOutput }()

		// Error should always output regardless of ShowOutput
		ShowOutput = false
		output := captureOutput(func() {
			Errorf("Error message")
		})
		if output == "" {
			t.Error("Expected Errorf to always output")
		}
	})

	t.Run("infof_always_outputs", func(t *testing.T) {
		originalShowOutput := ShowOutput
		defer func() { ShowOutput = originalShowOutput }()

		// Info should always output regardless of ShowOutput
		ShowOutput = false
		output := captureOutput(func() {
			Infof("Info message")
		})
		if output == "" {
			t.Error("Expected Infof to always output")
		}
	})
}

// ============================================================================
// Timestamp Format Tests
// ============================================================================

func TestTimestampFormat(t *testing.T) {
	t.Run("printf_timestamp_format", func(t *testing.T) {
		ShowOutput = true
		output := captureOutput(func() {
			Printf("test")
		})

		// Check for timestamp pattern [YYYY-MM-DD HH:MM:SS]
		if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
			t.Error("Expected timestamp with brackets")
		}
		// Should contain colon for time
		if strings.Count(output[:25], ":") < 2 {
			t.Error("Expected timestamp with time colons")
		}
	})

	t.Run("debugf_timestamp_format", func(t *testing.T) {
		ShowDebug = true
		output := captureOutput(func() {
			Debugf("test")
		})

		if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
			t.Error("Expected timestamp with brackets in debug output")
		}
	})

	t.Run("errorf_timestamp_format", func(t *testing.T) {
		output := captureOutput(func() {
			Errorf("test")
		})

		if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
			t.Error("Expected timestamp with brackets in error output")
		}
	})

	t.Run("infof_timestamp_format", func(t *testing.T) {
		output := captureOutput(func() {
			Infof("test")
		})

		if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
			t.Error("Expected timestamp with brackets in info output")
		}
	})
}
