package errors

import (
	"errors"
	"fmt"

	"github.com/ddmoney420/moji/internal/ux"
)

// Error code constants for programmatic error handling
const (
	CodeUnknown       = 0
	CodeFileNotFound  = 1
	CodeInvalidInput  = 2
	CodeInvalidConfig = 3
	CodePermission    = 4
	CodeIO            = 5
	CodeNotSupported  = 6
	CodeAlreadyExists = 7
)

// MojiError is the standardized error type for moji
type MojiError struct {
	Code       int
	Message    string
	Suggestion string
	Err        error // Wrapped error for error chain
}

// Error implements the error interface
func (e *MojiError) Error() string {
	if e.Err != nil && e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

// Unwrap returns the wrapped error for use with errors.Is and errors.As
func (e *MojiError) Unwrap() error {
	return e.Err
}

// WithSuggestion adds a suggestion to the error for chaining
func (e *MojiError) WithSuggestion(s string) *MojiError {
	e.Suggestion = s
	return e
}

// WithCode sets the error code for chaining
func (e *MojiError) WithCode(c int) *MojiError {
	e.Code = c
	return e
}

// Format returns a formatted error message suitable for CLI output
// Includes ANSI colors for better visibility
func (e *MojiError) Format() string {
	var output string

	// Main error message in bold red
	if ux.NoColor {
		output = fmt.Sprintf("Error: %s\n", e.Error())
	} else {
		output = fmt.Sprintf("%sError: %s%s\n", ux.BoldRed, e.Error(), ux.Reset)
	}

	// Add suggestion in cyan if present
	if e.Suggestion != "" {
		if ux.NoColor {
			output += fmt.Sprintf("Suggestion: %s\n", e.Suggestion)
		} else {
			output += fmt.Sprintf("%sSuggestion: %s%s\n", ux.Cyan, e.Suggestion, ux.Reset)
		}
	}

	return output
}

// New creates a new MojiError with a message
func New(msg string) *MojiError {
	return &MojiError{
		Code:    CodeUnknown,
		Message: msg,
	}
}

// Wrap wraps an existing error with additional context
func Wrap(err error, msg string) *MojiError {
	if err == nil {
		return New(msg)
	}
	return &MojiError{
		Code:    CodeUnknown,
		Message: msg,
		Err:     err,
	}
}

// Is provides error matching functionality compatible with errors.Is()
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// Predefined error constructors for common cases

// FileNotFound creates an error for missing files
func FileNotFound(filename string) *MojiError {
	return &MojiError{
		Code:       CodeFileNotFound,
		Message:    fmt.Sprintf("file not found: %s", filename),
		Suggestion: fmt.Sprintf("Check that the path '%s' is correct and the file exists", filename),
	}
}

// InvalidInput creates an error for invalid user input
func InvalidInput(details string) *MojiError {
	return &MojiError{
		Code:       CodeInvalidInput,
		Message:    fmt.Sprintf("invalid input: %s", details),
		Suggestion: "Check your input and try again with valid parameters",
	}
}

// InvalidConfig creates an error for configuration issues
func InvalidConfig(details string) *MojiError {
	return &MojiError{
		Code:       CodeInvalidConfig,
		Message:    fmt.Sprintf("invalid configuration: %s", details),
		Suggestion: "Check your configuration file and ensure all required fields are present",
	}
}

// PermissionDenied creates an error for permission issues
func PermissionDenied(resource string) *MojiError {
	return &MojiError{
		Code:       CodePermission,
		Message:    fmt.Sprintf("permission denied: %s", resource),
		Suggestion: fmt.Sprintf("You may need elevated privileges to access '%s'", resource),
	}
}

// IOError creates an error for I/O operations
func IOError(operation string, err error) *MojiError {
	return &MojiError{
		Code:       CodeIO,
		Message:    fmt.Sprintf("I/O error during %s", operation),
		Suggestion: "Check your disk space and file permissions",
		Err:        err,
	}
}

// NotSupported creates an error for unsupported operations
func NotSupported(feature string) *MojiError {
	return &MojiError{
		Code:       CodeNotSupported,
		Message:    fmt.Sprintf("%s is not supported on this system", feature),
		Suggestion: "Check the documentation for supported features and alternatives",
	}
}

// AlreadyExists creates an error for existing resources
func AlreadyExists(resource string) *MojiError {
	return &MojiError{
		Code:       CodeAlreadyExists,
		Message:    fmt.Sprintf("%s already exists", resource),
		Suggestion: fmt.Sprintf("Remove or rename the existing '%s' and try again", resource),
	}
}
