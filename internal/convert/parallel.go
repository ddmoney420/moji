package convert

import (
	"fmt"
	"image"
	"math"
	"runtime"
	"strings"
	"sync"
)

// ParallelConfig holds configuration for parallel processing
type ParallelConfig struct {
	workerCount int
	threshold   int // Minimum image size (width*height) to enable parallelization
}

// Global parallel config
var parallelConfig = ParallelConfig{
	workerCount: runtime.NumCPU(),
	threshold:   500 * 500, // Default threshold: 500x500 = 250,000 pixels
}

// SetWorkerCount sets the number of workers for parallel processing
func SetWorkerCount(n int) {
	if n < 1 {
		n = 1
	}
	parallelConfig.workerCount = n
}

// SetParallelThreshold sets the minimum pixel count to enable parallelization
func SetParallelThreshold(pixels int) {
	if pixels < 1 {
		pixels = 1
	}
	parallelConfig.threshold = pixels
}

// rowResult holds a processed row of ASCII art
type rowResult struct {
	rowIndex int
	content  string
	err      error
}

// FromImageParallel converts an image to ASCII art using parallel processing
// It splits the image into horizontal bands and processes each band in a separate worker
func FromImageParallel(img image.Image, opts Options) (string, error) {
	if opts.Charset == "" {
		opts.Charset = CharSets["standard"]
	}

	bounds := img.Bounds()
	imgWidth := bounds.Max.X - bounds.Min.X
	imgHeight := bounds.Max.Y - bounds.Min.Y

	// Calculate target dimensions
	targetWidth := opts.Width
	if targetWidth <= 0 {
		targetWidth = 80
	}

	aspectRatio := float64(imgWidth) / float64(imgHeight)
	targetHeight := opts.Height
	if targetHeight <= 0 {
		targetHeight = int(float64(targetWidth) / aspectRatio / 2.0)
	}

	// Determine optimal worker count
	workers := parallelConfig.workerCount
	if workers > targetHeight {
		workers = targetHeight
	}
	if workers < 1 {
		workers = 1
	}

	// For very small heights, don't parallelize
	minRowsPerWorker := 10
	if targetHeight < workers*minRowsPerWorker {
		return fromImageSequential(img, opts)
	}

	// Calculate sampling steps
	stepX := float64(imgWidth) / float64(targetWidth)
	stepY := float64(imgHeight) / float64(targetHeight)

	// Precompute edges if needed
	var edges [][]float64
	if opts.EdgeDetect {
		edges = detectEdges(img, targetWidth, targetHeight, stepX, stepY)
	}

	// Create result channel with buffer to avoid blocking workers
	resultsChan := make(chan rowResult, targetHeight)
	var wg sync.WaitGroup

	// Distribute rows to workers
	rowQueue := make(chan int, workers)

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processRowWorker(img, opts, bounds, targetWidth, targetHeight, stepX, stepY, edges, rowQueue, resultsChan)
		}()
	}

	// Queue all rows
	go func() {
		for y := 0; y < targetHeight; y++ {
			rowQueue <- y
		}
		close(rowQueue)
	}()

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results and assemble in order
	results := make(map[int]string)
	var firstErr error

	for result := range resultsChan {
		if result.err != nil && firstErr == nil {
			firstErr = result.err
		}
		results[result.rowIndex] = result.content
	}

	if firstErr != nil {
		return "", firstErr
	}

	// Assemble rows in correct order
	var output strings.Builder
	for y := 0; y < targetHeight; y++ {
		if content, ok := results[y]; ok {
			output.WriteString(content)
			output.WriteString("\n")
		}
	}

	return output.String(), nil
}

// processRowWorker processes rows from the queue
func processRowWorker(img image.Image, opts Options, bounds image.Rectangle, targetWidth, targetHeight int, stepX, stepY float64, edges [][]float64, rowQueue chan int, results chan rowResult) {
	chars := []rune(opts.Charset)
	numChars := len(chars)

	for y := range rowQueue {
		var rowContent strings.Builder

		for x := 0; x < targetWidth; x++ {
			// Sample the image region
			sampleX := int(float64(x)*stepX) + bounds.Min.X
			sampleY := int(float64(y)*stepY) + bounds.Min.Y

			// Ensure within bounds
			if sampleX >= bounds.Max.X {
				sampleX = bounds.Max.X - 1
			}
			if sampleY >= bounds.Max.Y {
				sampleY = bounds.Max.Y - 1
			}

			var brightness float64
			var r8, g8, b8 uint8

			if opts.EdgeDetect && edges != nil {
				brightness = edges[y][x]
			} else {
				r8, g8, b8, brightness = sampleRegion(img, sampleX, sampleY, int(stepX), int(stepY))
			}

			if opts.Invert {
				brightness = 1.0 - brightness
			}

			// Map brightness to character
			charIdx := int(brightness * float64(numChars-1))
			if charIdx >= numChars {
				charIdx = numChars - 1
			}
			if charIdx < 0 {
				charIdx = 0
			}

			char := chars[charIdx]

			if opts.Color && !opts.EdgeDetect {
				// ANSI 24-bit color
				rowContent.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm%c\x1b[0m", r8, g8, b8, char))
			} else {
				rowContent.WriteRune(char)
			}
		}

		results <- rowResult{
			rowIndex: y,
			content:  rowContent.String(),
			err:      nil,
		}
	}
}

// detectEdgesParallel applies Sobel edge detection with parallel processing
func detectEdgesParallel(img image.Image, width, height int, stepX, stepY float64) [][]float64 {
	bounds := img.Bounds()

	// First, create grayscale version at target resolution (sequential)
	gray := make([][]float64, height)
	for y := 0; y < height; y++ {
		gray[y] = make([]float64, width)
		for x := 0; x < width; x++ {
			sampleX := int(float64(x)*stepX) + bounds.Min.X
			sampleY := int(float64(y)*stepY) + bounds.Min.Y

			if sampleX >= bounds.Max.X {
				sampleX = bounds.Max.X - 1
			}
			if sampleY >= bounds.Max.Y {
				sampleY = bounds.Max.Y - 1
			}

			_, _, _, brightness := sampleRegion(img, sampleX, sampleY, int(stepX), int(stepY))
			gray[y][x] = brightness
		}
	}

	// Apply Sobel operator in parallel
	edges := make([][]float64, height)
	for y := 0; y < height; y++ {
		edges[y] = make([]float64, width)
	}

	// Determine optimal worker count for edge detection
	workers := parallelConfig.workerCount
	if workers > height {
		workers = height
	}

	rowQueue := make(chan int, workers)
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for y := range rowQueue {
				for x := 0; x < width; x++ {
					if x == 0 || x == width-1 || y == 0 || y == height-1 {
						edges[y][x] = 0
						continue
					}

					gx := -gray[y-1][x-1] + gray[y-1][x+1] +
						-2*gray[y][x-1] + 2*gray[y][x+1] +
						-gray[y+1][x-1] + gray[y+1][x+1]

					gy := -gray[y-1][x-1] - 2*gray[y-1][x] - gray[y-1][x+1] +
						gray[y+1][x-1] + 2*gray[y+1][x] + gray[y+1][x+1]

					magnitude := math.Sqrt(gx*gx + gy*gy)
					edges[y][x] = 1.0 - math.Min(magnitude*2, 1.0)
				}
			}
		}()
	}

	go func() {
		for y := 0; y < height; y++ {
			rowQueue <- y
		}
		close(rowQueue)
	}()

	wg.Wait()

	return edges
}
