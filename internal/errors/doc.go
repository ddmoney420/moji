// Package errors provides standardized error handling with error codes and suggestions.
//
// It defines a MojiError type that includes error codes, user-friendly messages, and contextual
// suggestions for fixing common issues. The package includes preset constructors for common
// error scenarios like missing files or invalid input.
//
// Example usage:
//
//	err := errors.FileNotFound("config.yaml")
//	err := errors.InvalidInput("font name", "shadow")
//	err := errors.PermissionDenied("/etc/config")
//	if mojiErr, ok := err.(*errors.MojiError); ok {
//		fmt.Println(mojiErr.Suggest())
//	}
package errors
