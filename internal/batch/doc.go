// Package batch provides concurrent batch processing for image conversions.
//
// It implements efficient parallel image conversion using worker pools to handle multiple files
// simultaneously. The package manages file I/O, progress reporting, and error handling for
// large-scale image processing tasks.
//
// Example usage:
//
//	results := batch.ConvertImages(files, opts)
//	batch.SaveResults(results, outputDir)
//	batch.PrintResults(results)
package batch
