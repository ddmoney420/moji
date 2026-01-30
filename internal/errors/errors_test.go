package errors

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/ddmoney420/moji/internal/ux"
)

func TestMojiErrorError(t *testing.T) {
	tests := []struct {
		name     string
		err      *MojiError
		expected string
	}{
		{
			name:     "message only",
			err:      &MojiError{Message: "test error"},
			expected: "test error",
		},
		{
			name:     "with wrapped error",
			err:      &MojiError{Message: "context", Err: errors.New("root cause")},
			expected: "context: root cause",
		},
		{
			name:     "wrapped error only",
			err:      &MojiError{Err: errors.New("root cause")},
			expected: "root cause",
		},
		{
			name:     "empty error",
			err:      &MojiError{},
			expected: "unknown error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Error() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestMojiErrorUnwrap(t *testing.T) {
	rootErr := errors.New("root error")
	err := &MojiError{Message: "context", Err: rootErr}

	unwrapped := err.Unwrap()
	if unwrapped != rootErr {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, rootErr)
	}
}

func TestWithSuggestion(t *testing.T) {
	err := New("test").WithSuggestion("try this")

	if err.Suggestion != "try this" {
		t.Errorf("WithSuggestion() set suggestion to %q, want %q", err.Suggestion, "try this")
	}
	// Test chaining returns same error
	if err.Message != "test" {
		t.Errorf("WithSuggestion() changed message to %q", err.Message)
	}
}

func TestWithCode(t *testing.T) {
	err := New("test").WithCode(42)

	if err.Code != 42 {
		t.Errorf("WithCode() set code to %d, want 42", err.Code)
	}
	// Test chaining returns same error
	if err.Message != "test" {
		t.Errorf("WithCode() changed message to %q", err.Message)
	}
}

func TestFormat(t *testing.T) {
	tests := []struct {
		name             string
		err              *MojiError
		shouldContain    []string
		shouldNotContain []string
	}{
		{
			name:          "with suggestion",
			err:           New("test error").WithSuggestion("try this"),
			shouldContain: []string{"Error: test error", "Suggestion: try this"},
		},
		{
			name:             "without suggestion",
			err:              New("test error"),
			shouldContain:    []string{"Error: test error"},
			shouldNotContain: []string{"Suggestion:"},
		},
		{
			name:          "with wrapped error",
			err:           Wrap(errors.New("root"), "context"),
			shouldContain: []string{"Error: context: root"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Format()

			for _, substr := range tt.shouldContain {
				if !strings.Contains(result, substr) {
					t.Errorf("Format() missing %q\nGot: %s", substr, result)
				}
			}

			for _, substr := range tt.shouldNotContain {
				if strings.Contains(result, substr) {
					t.Errorf("Format() should not contain %q\nGot: %s", substr, result)
				}
			}
		})
	}
}

func TestNew(t *testing.T) {
	err := New("test message")

	if err.Message != "test message" {
		t.Errorf("New() message = %q, want %q", err.Message, "test message")
	}
	if err.Code != CodeUnknown {
		t.Errorf("New() code = %d, want %d", err.Code, CodeUnknown)
	}
	if err.Err != nil {
		t.Errorf("New() wrapped error = %v, want nil", err.Err)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		msg         string
		expectedMsg string
	}{
		{
			name:        "wrap error",
			err:         errors.New("root"),
			msg:         "context",
			expectedMsg: "context: root",
		},
		{
			name:        "wrap nil error",
			err:         nil,
			msg:         "test",
			expectedMsg: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Wrap(tt.err, tt.msg)
			if result.Error() != tt.expectedMsg {
				t.Errorf("Wrap() = %q, want %q", result.Error(), tt.expectedMsg)
			}
		})
	}
}

func TestIs(t *testing.T) {
	targetErr := errors.New("target")
	err := Wrap(targetErr, "context")

	if !Is(err, targetErr) {
		t.Error("Is() should find wrapped error")
	}
}

func TestFileNotFound(t *testing.T) {
	err := FileNotFound("test.txt")

	if err.Code != CodeFileNotFound {
		t.Errorf("FileNotFound() code = %d, want %d", err.Code, CodeFileNotFound)
	}
	if !strings.Contains(err.Message, "test.txt") {
		t.Errorf("FileNotFound() message should contain filename: %q", err.Message)
	}
	if !strings.Contains(err.Suggestion, "test.txt") {
		t.Errorf("FileNotFound() suggestion should contain filename: %q", err.Suggestion)
	}
}

func TestInvalidInput(t *testing.T) {
	err := InvalidInput("expected number")

	if err.Code != CodeInvalidInput {
		t.Errorf("InvalidInput() code = %d, want %d", err.Code, CodeInvalidInput)
	}
	if !strings.Contains(err.Message, "expected number") {
		t.Errorf("InvalidInput() message should contain details")
	}
	if err.Suggestion == "" {
		t.Error("InvalidInput() should have a suggestion")
	}
}

func TestInvalidConfig(t *testing.T) {
	err := InvalidConfig("missing field")

	if err.Code != CodeInvalidConfig {
		t.Errorf("InvalidConfig() code = %d, want %d", err.Code, CodeInvalidConfig)
	}
	if !strings.Contains(err.Message, "missing field") {
		t.Errorf("InvalidConfig() message should contain details")
	}
	if err.Suggestion == "" {
		t.Error("InvalidConfig() should have a suggestion")
	}
}

func TestPermissionDenied(t *testing.T) {
	err := PermissionDenied("/etc/passwd")

	if err.Code != CodePermission {
		t.Errorf("PermissionDenied() code = %d, want %d", err.Code, CodePermission)
	}
	if !strings.Contains(err.Message, "/etc/passwd") {
		t.Errorf("PermissionDenied() message should contain resource")
	}
	if !strings.Contains(err.Suggestion, "/etc/passwd") {
		t.Errorf("PermissionDenied() suggestion should contain resource")
	}
}

func TestIOError(t *testing.T) {
	rootErr := errors.New("disk full")
	err := IOError("write", rootErr)

	if err.Code != CodeIO {
		t.Errorf("IOError() code = %d, want %d", err.Code, CodeIO)
	}
	if !strings.Contains(err.Message, "write") {
		t.Errorf("IOError() message should contain operation")
	}
	if err.Err != rootErr {
		t.Errorf("IOError() should wrap the error")
	}
	if err.Suggestion == "" {
		t.Error("IOError() should have a suggestion")
	}
}

func TestNotSupported(t *testing.T) {
	err := NotSupported("WASM rendering")

	if err.Code != CodeNotSupported {
		t.Errorf("NotSupported() code = %d, want %d", err.Code, CodeNotSupported)
	}
	if !strings.Contains(err.Message, "WASM rendering") {
		t.Errorf("NotSupported() message should contain feature")
	}
	if err.Suggestion == "" {
		t.Error("NotSupported() should have a suggestion")
	}
}

func TestAlreadyExists(t *testing.T) {
	err := AlreadyExists("config.yaml")

	if err.Code != CodeAlreadyExists {
		t.Errorf("AlreadyExists() code = %d, want %d", err.Code, CodeAlreadyExists)
	}
	if !strings.Contains(err.Message, "config.yaml") {
		t.Errorf("AlreadyExists() message should contain resource")
	}
	if !strings.Contains(err.Suggestion, "config.yaml") {
		t.Errorf("AlreadyExists() suggestion should contain resource")
	}
}

func TestMojiErrorChaining(t *testing.T) {
	// Test that methods can be chained
	err := New("base").
		WithCode(CodeFileNotFound).
		WithSuggestion("check the path")

	if err.Code != CodeFileNotFound {
		t.Errorf("chaining WithCode failed")
	}
	if err.Suggestion != "check the path" {
		t.Errorf("chaining WithSuggestion failed")
	}
	if err.Message != "base" {
		t.Errorf("chaining changed base message")
	}
}

func TestErrorInterface(t *testing.T) {
	// Ensure MojiError implements the error interface
	var _ error = (*MojiError)(nil)

	// Test that it works with error handling
	err := New("test")
	if err.Error() == "" {
		t.Error("Error() should not return empty string")
	}
}

func TestFormatWithoutColor(t *testing.T) {
	// Save original NoColor setting
	originalNoColor := ux.NoColor
	defer func() { ux.NoColor = originalNoColor }()

	ux.NoColor = true
	err := New("test error").WithSuggestion("try this")
	result := err.Format()

	// Should not contain ANSI codes when NoColor is true
	if strings.Contains(result, "\033[") {
		t.Error("Format() should not contain ANSI codes when NoColor is true")
	}

	if !strings.Contains(result, "Error: test error") {
		t.Error("Format() should contain error message")
	}
	if !strings.Contains(result, "Suggestion: try this") {
		t.Error("Format() should contain suggestion")
	}
}

func TestFormatWithColor(t *testing.T) {
	// Save original NoColor setting
	originalNoColor := ux.NoColor
	defer func() { ux.NoColor = originalNoColor }()

	ux.NoColor = false
	err := New("test error").WithSuggestion("try this")
	result := err.Format()

	// Should contain ANSI codes when NoColor is false
	if !strings.Contains(result, "\033[") {
		t.Error("Format() should contain ANSI codes when NoColor is false")
	}
}

func TestCodeConstants(t *testing.T) {
	// Verify all code constants are defined and unique
	codes := map[int]string{
		CodeUnknown:       "CodeUnknown",
		CodeFileNotFound:  "CodeFileNotFound",
		CodeInvalidInput:  "CodeInvalidInput",
		CodeInvalidConfig: "CodeInvalidConfig",
		CodePermission:    "CodePermission",
		CodeIO:            "CodeIO",
		CodeNotSupported:  "CodeNotSupported",
		CodeAlreadyExists: "CodeAlreadyExists",
	}

	if len(codes) != 8 {
		t.Errorf("expected 8 unique codes, got %d", len(codes))
	}
}

func TestMultipleWraps(t *testing.T) {
	// Test wrapping multiple times
	rootErr := errors.New("root")
	wrapped1 := Wrap(rootErr, "level1")
	wrapped2 := Wrap(wrapped1, "level2")

	if !strings.Contains(wrapped2.Error(), "level2") {
		t.Error("outermost message should be present")
	}
	// Check that we can still find the root error
	if !Is(wrapped2, rootErr) {
		t.Error("should be able to find root error in chain")
	}
}

func TestErrorMessages(t *testing.T) {
	tests := []struct {
		name    string
		creator func() *MojiError
		check   func(*MojiError) bool
	}{
		{
			name:    "FileNotFound has code",
			creator: func() *MojiError { return FileNotFound("test.txt") },
			check:   func(e *MojiError) bool { return e.Code == CodeFileNotFound },
		},
		{
			name:    "InvalidInput has code",
			creator: func() *MojiError { return InvalidInput("details") },
			check:   func(e *MojiError) bool { return e.Code == CodeInvalidInput },
		},
		{
			name:    "InvalidConfig has code",
			creator: func() *MojiError { return InvalidConfig("details") },
			check:   func(e *MojiError) bool { return e.Code == CodeInvalidConfig },
		},
		{
			name:    "PermissionDenied has code",
			creator: func() *MojiError { return PermissionDenied("resource") },
			check:   func(e *MojiError) bool { return e.Code == CodePermission },
		},
		{
			name:    "IOError has code",
			creator: func() *MojiError { return IOError("read", errors.New("test")) },
			check:   func(e *MojiError) bool { return e.Code == CodeIO },
		},
		{
			name:    "NotSupported has code",
			creator: func() *MojiError { return NotSupported("feature") },
			check:   func(e *MojiError) bool { return e.Code == CodeNotSupported },
		},
		{
			name:    "AlreadyExists has code",
			creator: func() *MojiError { return AlreadyExists("resource") },
			check:   func(e *MojiError) bool { return e.Code == CodeAlreadyExists },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.creator()
			if !tt.check(err) {
				t.Errorf("%s failed", tt.name)
			}
		})
	}
}

func TestErrorFormatting(t *testing.T) {
	// Test that Format() produces valid output
	tests := []struct {
		name string
		err  *MojiError
	}{
		{"simple", New("error")},
		{"with suggestion", New("error").WithSuggestion("suggestion")},
		{"with code", New("error").WithCode(42)},
		{"wrapped", Wrap(errors.New("root"), "context")},
		{"wrapped with suggestion", Wrap(errors.New("root"), "context").WithSuggestion("suggestion")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Format()
			if result == "" {
				t.Error("Format() returned empty string")
			}
			if !strings.Contains(result, "Error:") {
				t.Error("Format() should contain 'Error:'")
			}
		})
	}
}

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(fmt.Sprintf("error %d", i))
	}
}

func BenchmarkWrap(b *testing.B) {
	baseErr := errors.New("root")
	for i := 0; i < b.N; i++ {
		Wrap(baseErr, fmt.Sprintf("context %d", i))
	}
}

func BenchmarkFileNotFound(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FileNotFound(fmt.Sprintf("file%d.txt", i))
	}
}

func BenchmarkFormat(b *testing.B) {
	err := New("test error").WithSuggestion("try this")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err.Format()
	}
}
