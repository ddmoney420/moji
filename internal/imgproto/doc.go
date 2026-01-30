// Package imgproto provides terminal image protocol support.
//
// It handles multiple terminal image display protocols including Sixel, Kitty, iTerm2, WezTerm,
// and Terminology. The package auto-detects terminal capabilities and uses the appropriate
// protocol for rendering images in the terminal.
//
// Example usage:
//
//	caps := imgproto.Detect()
//	if caps.SixelSupported {
//		imgproto.RenderSixel(img)
//	}
//	imgproto.Render(img) // Auto-detect and render
package imgproto
