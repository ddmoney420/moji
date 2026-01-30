// Package animate provides text animation effects including spinners, typewriter, scrolling, and fade effects.
//
// It offers a variety of preset animation sequences (spinners, dots, bounces) that can be played
// with text, as well as specialized effects like typewriter text, scrolling text, fade in, blinking,
// and matrix rain. Each animation is customizable with timing and display options.
//
// Example usage:
//
//	animate.Play(animate.SpinnerDots, "Loading...")
//	animate.Typewriter("Hello", 50*time.Millisecond)
//	animate.ScrollText("Message", 10)
//	animate.MatrixRain(10, 5*time.Second)
package animate
