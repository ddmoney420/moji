// Package progress provides progress bars and spinners for terminal output.
//
// It implements 8+ bar styles and 11 spinner types for showing progress and activity.
// The package includes utilities for ETA calculation and rate measurement.
//
// Example usage:
//
//	bar := progress.Bar(current, total)
//	spinner := progress.Spinner(progress.SpinnerDots)
//	eta := progress.ETA(elapsed, processed, total)
package progress
