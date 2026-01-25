package batch

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ddmoney420/moji/internal/convert"
)

// Result holds the result of converting a single image
type Result struct {
	Path   string
	Output string
	Error  error
}

// ConvertImages converts multiple images concurrently
func ConvertImages(patterns []string, opts convert.Options, workers int) []Result {
	// Collect all files matching patterns
	var files []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		files = append(files, matches...)
	}

	if len(files) == 0 {
		return nil
	}

	if workers <= 0 {
		workers = 4
	}

	// Create channels
	jobs := make(chan string, len(files))
	results := make(chan Result, len(files))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range jobs {
				art, err := convert.FromFile(path, opts)
				results <- Result{
					Path:   path,
					Output: art,
					Error:  err,
				}
			}
		}()
	}

	// Send jobs
	for _, f := range files {
		jobs <- f
	}
	close(jobs)

	// Wait and collect
	go func() {
		wg.Wait()
		close(results)
	}()

	var out []Result
	for r := range results {
		out = append(out, r)
	}

	return out
}

// SaveResults saves batch results to files
func SaveResults(results []Result, outDir string, format string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	for _, r := range results {
		if r.Error != nil {
			fmt.Fprintf(os.Stderr, "Skipping %s: %v\n", r.Path, r.Error)
			continue
		}

		base := filepath.Base(r.Path)
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		outPath := filepath.Join(outDir, name+"."+format)

		if err := os.WriteFile(outPath, []byte(r.Output), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write %s: %v\n", outPath, err)
			continue
		}
	}

	return nil
}

// PrintResults prints results in a gallery format
func PrintResults(results []Result, showPath bool) {
	for i, r := range results {
		if r.Error != nil {
			fmt.Fprintf(os.Stderr, "Error converting %s: %v\n", r.Path, r.Error)
			continue
		}

		if showPath {
			fmt.Printf("=== %s ===\n", r.Path)
		}
		fmt.Println(r.Output)
		if i < len(results)-1 {
			fmt.Println()
		}
	}
}
