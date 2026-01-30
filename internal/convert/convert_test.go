package convert

import (
	"image"
	"image/color"
	"strings"
	"testing"
)

// createTestImage creates a simple test image
func createTestImage(w, h int, c color.Color) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}

// createGradientImage creates a horizontal gradient test image
func createGradientImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			gray := uint8(float64(x) / float64(w) * 255)
			img.Set(x, y, color.RGBA{gray, gray, gray, 255})
		}
	}
	return img
}

func TestDefaultOptions(t *testing.T) {
	opts := DefaultOptions()
	if opts.Width != 80 {
		t.Errorf("Default width = %d, want 80", opts.Width)
	}
	if opts.Charset == "" {
		t.Error("Default charset should not be empty")
	}
}

func TestGetCharset(t *testing.T) {
	tests := []struct {
		name   string
		expect string
	}{
		{"standard", " .:-=+*#%@"},
		{"simple", " .*#"},
		{"binary", " \u2588"},
		{"unknown", " .:-=+*#%@"}, // should default to standard
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCharset(tt.name)
			if result != tt.expect {
				t.Errorf("GetCharset(%q) = %q, want %q", tt.name, result, tt.expect)
			}
		})
	}
}

func TestListCharsets(t *testing.T) {
	charsets := ListCharsets()
	if len(charsets) < 20 {
		t.Errorf("ListCharsets should return at least 20 charsets, got %d", len(charsets))
	}

	// All listed charsets should exist in CharSets map
	for _, name := range charsets {
		if _, ok := CharSets[name]; !ok {
			t.Errorf("Listed charset %q not found in CharSets map", name)
		}
	}
}

func TestFromImageBlack(t *testing.T) {
	img := createTestImage(100, 50, color.Black)
	opts := Options{
		Width:   20,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	// Black image should use mostly dark characters (space)
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("Result should have lines")
	}

	// Count spaces - black image should be mostly spaces
	spaceCount := strings.Count(result, " ")
	totalChars := 0
	for _, line := range lines {
		totalChars += len(line)
	}
	if totalChars > 0 && float64(spaceCount)/float64(totalChars) < 0.8 {
		t.Error("Black image should produce mostly space characters")
	}
}

func TestFromImageWhite(t *testing.T) {
	img := createTestImage(100, 50, color.White)
	opts := Options{
		Width:   20,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	// White image should use dense characters (last few in charset)
	hasDenseChars := strings.ContainsAny(result, "#%@")
	if !hasDenseChars {
		t.Errorf("White image should contain dense characters (#, %%, @), got: %q", result)
	}
}

func TestFromImageGradient(t *testing.T) {
	img := createGradientImage(200, 100)
	opts := Options{
		Width:   40,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("Result should have lines")
	}

	// First column should be lighter, last column should be darker
	firstChar := rune(lines[0][0])
	lastLine := lines[0]
	lastChar := rune(lastLine[len(lastLine)-1])

	// Space is lightest, @ is densest
	charIndex := func(c rune) int {
		for i, ch := range " .:-=+*#%@" {
			if ch == c {
				return i
			}
		}
		return -1
	}

	if charIndex(firstChar) > charIndex(lastChar) {
		t.Error("Left side of gradient should be lighter than right side")
	}
}

func TestFromImageInvert(t *testing.T) {
	img := createTestImage(100, 50, color.Black)
	opts := Options{
		Width:   20,
		Charset: " .:-=+*#%@",
		Invert:  true,
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	// Inverted black should produce dense characters
	if !strings.Contains(result, "@") {
		t.Error("Inverted black image should contain @ characters")
	}
}

func TestFromImageColor(t *testing.T) {
	img := createTestImage(100, 50, color.RGBA{255, 0, 0, 255})
	opts := Options{
		Width:   20,
		Charset: " .:-=+*#%@",
		Color:   true,
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	// Should contain ANSI color codes
	if !strings.Contains(result, "\x1b[38;2;") {
		t.Error("Color mode should produce ANSI color codes")
	}
}

func TestFromImageEdgeDetection(t *testing.T) {
	// Create an image with a clear edge (left half black, right half white)
	img := image.NewRGBA(image.Rect(0, 0, 100, 50))
	for y := 0; y < 50; y++ {
		for x := 0; x < 100; x++ {
			if x < 50 {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	opts := Options{
		Width:      20,
		Charset:    " .:-=+*#%@",
		EdgeDetect: true,
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	// Edge detection should produce output
	if len(result) == 0 {
		t.Error("Edge detection should produce output")
	}
}

func TestFromImageDefaultCharset(t *testing.T) {
	img := createTestImage(50, 25, color.Gray{128})
	opts := Options{
		Width: 10,
		// No charset specified
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}
	if len(result) == 0 {
		t.Error("Should produce output with default charset")
	}
}

func TestFromImageAutoHeight(t *testing.T) {
	img := createTestImage(200, 100, color.Gray{128})
	opts := Options{
		Width:   40,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	// Height should be auto-calculated based on aspect ratio
	// 200x100 image at width 40 should be roughly 40/(200/100)/2 = 10 lines
	if len(lines) < 5 || len(lines) > 30 {
		t.Errorf("Auto height produced %d lines, expected 5-30 for a 2:1 aspect image at width 40", len(lines))
	}
}

func TestFromImageExplicitHeight(t *testing.T) {
	img := createTestImage(100, 100, color.Gray{128})
	opts := Options{
		Width:   20,
		Height:  10,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) != 10 {
		t.Errorf("Explicit height=10 produced %d lines", len(lines))
	}
}

func TestFromImageWidth(t *testing.T) {
	img := createTestImage(100, 50, color.Gray{128})
	opts := Options{
		Width:   30,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) > 0 {
		// Each line should be exactly width characters
		for i, line := range lines {
			if len(line) != 30 {
				t.Errorf("Line %d width = %d, want 30", i, len(line))
				break
			}
		}
	}
}

func TestFromImageZeroWidth(t *testing.T) {
	img := createTestImage(100, 50, color.Gray{128})
	opts := Options{
		Width:   0, // should default to 80
		Charset: " .:-=+*#%@",
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) > 0 && len(lines[0]) != 80 {
		t.Errorf("Zero width should default to 80, got %d", len(lines[0]))
	}
}

func TestSampleRegion(t *testing.T) {
	img := createTestImage(10, 10, color.RGBA{255, 128, 64, 255})

	r, g, b, brightness := sampleRegion(img, 0, 0, 5, 5)
	if r != 255 || g != 128 || b != 64 {
		t.Errorf("sampleRegion color = (%d,%d,%d), want (255,128,64)", r, g, b)
	}
	if brightness <= 0 || brightness >= 1 {
		t.Errorf("sampleRegion brightness = %f, want between 0 and 1", brightness)
	}

	// Zero-size region
	_, _, _, brightness = sampleRegion(img, 0, 0, 0, 0)
	// Should handle gracefully (defaults to 1x1)
	if brightness < 0 {
		t.Error("Zero region should not return negative brightness")
	}
}

func TestCharSetsNotEmpty(t *testing.T) {
	for name, cs := range CharSets {
		if len(cs) < 2 {
			t.Errorf("Charset %q should have at least 2 characters, got %d", name, len(cs))
		}
	}
}

// Test parallel processing
func TestFromImageParallel(t *testing.T) {
	img := createTestImage(600, 600, color.Gray{128})
	opts := Options{
		Width:   40,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImageParallel(img, opts)
	if err != nil {
		t.Fatalf("FromImageParallel error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("Result should have lines")
	}

	// Check width matches
	for i, line := range lines {
		if len(line) != 40 {
			t.Errorf("Line %d width = %d, want 40", i, len(line))
			break
		}
	}
}

// TestFromImageParallelBlack verifies parallel processing produces correct output for solid colors
func TestFromImageParallelBlack(t *testing.T) {
	img := createTestImage(1000, 1000, color.Black)
	opts := Options{
		Width:   50,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImageParallel(img, opts)
	if err != nil {
		t.Fatalf("FromImageParallel error: %v", err)
	}

	// Black image should use mostly dark characters (space)
	spaceCount := strings.Count(result, " ")
	totalChars := 0
	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	for _, line := range lines {
		totalChars += len(line)
	}
	if totalChars > 0 && float64(spaceCount)/float64(totalChars) < 0.8 {
		t.Error("Parallel: Black image should produce mostly space characters")
	}
}

// TestFromImageParallelWhite verifies parallel processing handles bright regions
func TestFromImageParallelWhite(t *testing.T) {
	img := createTestImage(1000, 1000, color.White)
	opts := Options{
		Width:   50,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImageParallel(img, opts)
	if err != nil {
		t.Fatalf("FromImageParallel error: %v", err)
	}

	// White image should use dense characters
	hasDenseChars := strings.ContainsAny(result, "#%@")
	if !hasDenseChars {
		t.Errorf("Parallel: White image should contain dense characters (#, %%, @)")
	}
}

// TestFromImageParallelGradient verifies parallel processing handles gradient correctly
func TestFromImageParallelGradient(t *testing.T) {
	img := createGradientImage(800, 800)
	opts := Options{
		Width:   40,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImageParallel(img, opts)
	if err != nil {
		t.Fatalf("FromImageParallel error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("Parallel gradient: Result should have lines")
	}

	// Verify output matches sequential version
	optionsSeq := Options{
		Width:   40,
		Charset: " .:-=+*#%@",
	}
	// Force sequential by using small threshold
	oldThreshold := parallelConfig.threshold
	parallelConfig.threshold = 999999999 // Set impossibly high threshold
	sequential, errSeq := FromImage(img, optionsSeq)
	parallelConfig.threshold = oldThreshold

	if errSeq != nil {
		t.Fatalf("Sequential FromImage error: %v", errSeq)
	}

	if result != sequential {
		t.Error("Parallel result differs from sequential result")
	}
}

// TestFromImageParallelMatchesSequential verifies parallel and sequential produce identical output
func TestFromImageParallelMatchesSequential(t *testing.T) {
	// Create a large test image to ensure parallelization is used
	img := createGradientImage(1000, 1000)
	opts := Options{
		Width:   80,
		Height:  40,
		Charset: " .:-=+*#%@",
	}

	// Get parallel result
	parallelResult, err := FromImageParallel(img, opts)
	if err != nil {
		t.Fatalf("FromImageParallel error: %v", err)
	}

	// Get sequential result by forcing low threshold
	oldThreshold := parallelConfig.threshold
	parallelConfig.threshold = 999999999
	sequentialResult, err := FromImage(img, opts)
	parallelConfig.threshold = oldThreshold

	if err != nil {
		t.Fatalf("FromImage error: %v", err)
	}

	if parallelResult != sequentialResult {
		t.Error("Parallel result does not match sequential result")
		// Show first difference for debugging
		for i := 0; i < len(parallelResult) && i < len(sequentialResult); i++ {
			if parallelResult[i] != sequentialResult[i] {
				t.Logf("First difference at byte %d: parallel=%q sequential=%q", i, parallelResult[i], sequentialResult[i])
				break
			}
		}
	}
}

// TestFromImageParallelColor verifies parallel processing with color enabled
func TestFromImageParallelColor(t *testing.T) {
	img := createTestImage(600, 600, color.RGBA{255, 0, 0, 255})
	opts := Options{
		Width:   40,
		Charset: " .:-=+*#%@",
		Color:   true,
	}

	result, err := FromImageParallel(img, opts)
	if err != nil {
		t.Fatalf("FromImageParallel with color error: %v", err)
	}

	// Should contain ANSI color codes
	if !strings.Contains(result, "\x1b[38;2;") {
		t.Error("Parallel color mode should produce ANSI color codes")
	}
}

// TestFromImageParallelEdgeDetection verifies parallel processing with edge detection
func TestFromImageParallelEdgeDetection(t *testing.T) {
	// Create high contrast image
	img := image.NewRGBA(image.Rect(0, 0, 600, 600))
	for y := 0; y < 600; y++ {
		for x := 0; x < 600; x++ {
			if x < 300 {
				img.Set(x, y, color.Black)
			} else {
				img.Set(x, y, color.White)
			}
		}
	}

	opts := Options{
		Width:      40,
		Charset:    " .:-=+*#%@",
		EdgeDetect: true,
	}

	result, err := FromImageParallel(img, opts)
	if err != nil {
		t.Fatalf("FromImageParallel with edge detection error: %v", err)
	}

	if len(result) == 0 {
		t.Error("Parallel edge detection should produce output")
	}
}

// TestSetWorkerCount verifies worker count configuration
func TestSetWorkerCount(t *testing.T) {
	SetWorkerCount(4)
	if parallelConfig.workerCount != 4 {
		t.Errorf("SetWorkerCount(4) failed, got %d", parallelConfig.workerCount)
	}

	SetWorkerCount(1)
	if parallelConfig.workerCount != 1 {
		t.Errorf("SetWorkerCount(1) failed, got %d", parallelConfig.workerCount)
	}

	// Test invalid values
	SetWorkerCount(0)
	if parallelConfig.workerCount < 1 {
		t.Errorf("SetWorkerCount(0) should clamp to 1, got %d", parallelConfig.workerCount)
	}

	SetWorkerCount(-5)
	if parallelConfig.workerCount < 1 {
		t.Errorf("SetWorkerCount(-5) should clamp to 1, got %d", parallelConfig.workerCount)
	}
}

// TestSetParallelThreshold verifies threshold configuration
func TestSetParallelThreshold(t *testing.T) {
	oldThreshold := parallelConfig.threshold
	defer func() { parallelConfig.threshold = oldThreshold }()

	SetParallelThreshold(1000000)
	if parallelConfig.threshold != 1000000 {
		t.Errorf("SetParallelThreshold(1000000) failed, got %d", parallelConfig.threshold)
	}

	SetParallelThreshold(1)
	if parallelConfig.threshold != 1 {
		t.Errorf("SetParallelThreshold(1) failed, got %d", parallelConfig.threshold)
	}

	// Test invalid values
	SetParallelThreshold(0)
	if parallelConfig.threshold < 1 {
		t.Errorf("SetParallelThreshold(0) should clamp to 1, got %d", parallelConfig.threshold)
	}
}

// TestFromImageAutoParallel verifies that FromImage automatically uses parallel for large images
func TestFromImageAutoParallel(t *testing.T) {
	// Set a low threshold to trigger parallelization
	oldThreshold := parallelConfig.threshold
	defer func() { parallelConfig.threshold = oldThreshold }()

	SetParallelThreshold(100 * 100) // 10000 pixels

	// Create an image larger than threshold
	img := createTestImage(600, 600, color.Gray{128})
	opts := Options{
		Width:   40,
		Charset: " .:-=+*#%@",
	}

	result, err := FromImage(img, opts)
	if err != nil {
		t.Fatalf("FromImage with auto-parallel error: %v", err)
	}

	lines := strings.Split(strings.TrimRight(result, "\n"), "\n")
	if len(lines) == 0 {
		t.Fatal("Result should have lines")
	}
}

// BenchmarkFromImageSequential benchmarks sequential processing
func BenchmarkFromImageSequential(b *testing.B) {
	// Force sequential processing by setting very high threshold
	oldThreshold := parallelConfig.threshold
	parallelConfig.threshold = 999999999
	defer func() { parallelConfig.threshold = oldThreshold }()

	img := createTestImage(1000, 1000, color.Gray{128})
	opts := Options{
		Width:   80,
		Charset: " .:-=+*#%@",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromImage(img, opts)
	}
}

// BenchmarkFromImageParallel benchmarks parallel processing
func BenchmarkFromImageParallel(b *testing.B) {
	img := createTestImage(1000, 1000, color.Gray{128})
	opts := Options{
		Width:   80,
		Charset: " .:-=+*#%@",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromImageParallel(img, opts)
	}
}

// BenchmarkFromImageLarge benchmarks on larger image (to show parallelization benefit)
func BenchmarkFromImageLarge(b *testing.B) {
	// Force parallel processing
	oldThreshold := parallelConfig.threshold
	SetParallelThreshold(100 * 100)
	defer func() { parallelConfig.threshold = oldThreshold }()

	img := createTestImage(2000, 2000, color.Gray{128})
	opts := Options{
		Width:   100,
		Charset: " .:-=+*#%@",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = FromImage(img, opts)
	}
}
