package ux

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Spinner provides an animated loading indicator
type Spinner struct {
	frames  []string
	message string
	done    chan struct{}
	wg      sync.WaitGroup
	active  bool
	mu      sync.Mutex
}

// Spinner frames
var (
	SpinnerDots    = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	SpinnerLine    = []string{"-", "\\", "|", "/"}
	SpinnerBounce  = []string{"⠁", "⠂", "⠄", "⠂"}
	SpinnerGrow    = []string{"▁", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃"}
	SpinnerCircle  = []string{"◐", "◓", "◑", "◒"}
	SpinnerArrow   = []string{"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"}
	SpinnerDefault = SpinnerDots
)

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	frames := SpinnerDefault
	// Fall back to ASCII if not a TTY or colors disabled
	if !IsTTY() || NoColor {
		frames = SpinnerLine
	}
	return &Spinner{
		frames:  frames,
		message: message,
		done:    make(chan struct{}),
	}
}

// SetFrames sets custom spinner frames
func (s *Spinner) SetFrames(frames []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.frames = frames
}

// SetMessage updates the spinner message
func (s *Spinner) SetMessage(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.message = msg
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	if !IsTTY() || Quiet {
		// Non-interactive mode: just print message once
		if !Quiet {
			fmt.Fprintf(os.Stderr, "%s...\n", s.message)
		}
		return
	}

	s.mu.Lock()
	if s.active {
		s.mu.Unlock()
		return
	}
	s.active = true
	s.mu.Unlock()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		frame := 0
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-s.done:
				// Clear the spinner line
				fmt.Fprintf(os.Stderr, "\r\033[K")
				return
			case <-ticker.C:
				s.mu.Lock()
				f := s.frames[frame%len(s.frames)]
				msg := s.message
				s.mu.Unlock()

				if ColorEnabled() {
					fmt.Fprintf(os.Stderr, "\r%s%s%s %s", Cyan, f, Reset, msg)
				} else {
					fmt.Fprintf(os.Stderr, "\r%s %s", f, msg)
				}
				frame++
			}
		}
	}()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.mu.Lock()
	if !s.active {
		s.mu.Unlock()
		return
	}
	s.active = false
	s.mu.Unlock()

	close(s.done)
	s.wg.Wait()
}

// StopWithMessage stops and prints a final message
func (s *Spinner) StopWithMessage(msg string) {
	s.Stop()
	if !Quiet && IsTTY() {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}
}

// StopSuccess stops with a success message
func (s *Spinner) StopSuccess(msg string) {
	s.Stop()
	if !Quiet {
		if ColorEnabled() {
			fmt.Fprintf(os.Stderr, "%s✓%s %s\n", BoldGreen, Reset, msg)
		} else {
			fmt.Fprintf(os.Stderr, "[OK] %s\n", msg)
		}
	}
}

// StopFail stops with an error message
func (s *Spinner) StopFail(msg string) {
	s.Stop()
	if !Quiet {
		if ColorEnabled() {
			fmt.Fprintf(os.Stderr, "%s✗%s %s\n", BoldRed, Reset, msg)
		} else {
			fmt.Fprintf(os.Stderr, "[FAIL] %s\n", msg)
		}
	}
}

// ProgressBar displays a progress bar
type ProgressBar struct {
	total   int
	current int
	width   int
	message string
	mu      sync.Mutex
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, message string) *ProgressBar {
	return &ProgressBar{
		total:   total,
		width:   40,
		message: message,
	}
}

// SetWidth sets the progress bar width
func (p *ProgressBar) SetWidth(width int) {
	p.width = width
}

// Update updates the progress
func (p *ProgressBar) Update(current int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = current
	p.render()
}

// Increment increments progress by 1
func (p *ProgressBar) Increment() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current++
	p.render()
}

func (p *ProgressBar) render() {
	if !IsTTY() || Quiet {
		return
	}

	pct := float64(p.current) / float64(p.total)
	if pct > 1 {
		pct = 1
	}

	filled := int(pct * float64(p.width))
	empty := p.width - filled

	var bar string
	if ColorEnabled() {
		bar = fmt.Sprintf("%s%s%s%s",
			Green, strings.Repeat("█", filled),
			Dim+strings.Repeat("░", empty)+Reset,
			Reset)
	} else {
		bar = fmt.Sprintf("[%s%s]",
			strings.Repeat("=", filled),
			strings.Repeat(" ", empty))
	}

	fmt.Fprintf(os.Stderr, "\r%s %s %3.0f%% (%d/%d)",
		p.message, bar, pct*100, p.current, p.total)
}

// Done completes the progress bar
func (p *ProgressBar) Done() {
	if IsTTY() && !Quiet {
		fmt.Fprintf(os.Stderr, "\n")
	}
}
