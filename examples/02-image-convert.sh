#!/usr/bin/env bash

##############################################################################
# Example 2: Image to ASCII Conversion
#
# This example demonstrates how to convert images into ASCII art using
# different character sets, rendering modes, and color options.
#
# Key Topics:
#   - PNG/JPEG to ASCII conversion
#   - Character set selection (blocks, shade, braille, boxes)
#   - Width/height customization
#   - Edge detection for line art
#   - Terminal graphics protocols (Sixel, Kitty, iTerm2)
#   - Watch mode for real-time updates
#   - Color preservation
#
# Note: This example uses test images. You can replace with your own.
#
##############################################################################

set -e

# Color codes for section headers
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 2: Image to ASCII Conversion${RESET}\n"

# ============================================================================
# Setup: Create a simple test image if available
# ============================================================================
echo -e "${HEADER}Setup${RESET}"
echo "$DIVIDER"

# Check if moji can work with images
if moji --help | grep -q "convert"; then
    echo "✓ moji convert command is available"
else
    echo "✗ moji convert not available. This example requires image support."
    echo "  Make sure moji is fully built with image support enabled."
    exit 1
fi

# Try to find a test image in common locations
TEST_IMAGE=""
for img in /tmp/test.png ~/Pictures/test.png ./test.png; do
    if [ -f "$img" ]; then
        TEST_IMAGE="$img"
        break
    fi
done

if [ -z "$TEST_IMAGE" ]; then
    echo "Note: No test image found. This example shows commands that would work"
    echo "      with actual image files. Replace 'photo.jpg' with your image path."
    echo ""
    echo "Download test images or provide your own:"
    echo "  - PNG format: image.png"
    echo "  - JPEG format: photo.jpg"
    echo "  - GIF format: animation.gif"
    echo ""
fi

# ============================================================================
# 1. Basic Image Conversion
# ============================================================================
echo -e "${HEADER}1. Basic Image Conversion${RESET}"
echo "Converting images with default settings:"
echo "$DIVIDER"
cat << 'EOF'
# Convert to ASCII with standard characters
moji convert photo.jpg

# Specify output width (default 80)
moji convert photo.jpg --width 120

# Specify both width and height
moji convert photo.jpg --width 100 --height 50
EOF
echo ""

# ============================================================================
# 2. Character Set Selection
# ============================================================================
echo -e "${HEADER}2. Character Set Selection${RESET}"
echo "Different character sets provide different visual styles:"
echo "$DIVIDER"
cat << 'EOF'
# Standard ASCII characters (numbers and symbols)
moji convert photo.jpg --charset standard

# Block characters (better gradient representation)
moji convert photo.jpg --charset blocks

# Shade characters (smooth gradients)
moji convert photo.jpg --charset shade

# Braille characters (high detail)
moji convert photo.jpg --charset braille

# Box drawing characters (geometric)
moji convert photo.jpg --charset boxes

# Decorative characters (artistic)
moji convert photo.jpg --charset decorative

# View all available character sets
moji list-charsets
EOF
echo ""

# ============================================================================
# 3. Color and Inversion Options
# ============================================================================
echo -e "${HEADER}3. Color and Inversion${RESET}"
echo "Preserve colors or invert brightness:"
echo "$DIVIDER"
cat << 'EOF'
# Convert with color preservation (ANSI colors in terminal)
moji convert photo.jpg --color

# Invert brightness (useful for light backgrounds)
moji convert photo.jpg --invert

# Combine color and inversion
moji convert photo.jpg --color --invert

# Combine with specific character set
moji convert photo.jpg --charset blocks --color
EOF
echo ""

# ============================================================================
# 4. Edge Detection
# ============================================================================
echo -e "${HEADER}4. Edge Detection${RESET}"
echo "Detect and highlight edges for line-art style conversion:"
echo "$DIVIDER"
cat << 'EOF'
# Convert with edge detection
moji convert photo.jpg --edge

# Combine with color for edge lines
moji convert photo.jpg --edge --color

# Use braille for detailed edges
moji convert photo.jpg --edge --charset braille

# Useful for converting diagrams, sketches, or creating outlines
EOF
echo ""

# ============================================================================
# 5. Terminal Graphics Protocols
# ============================================================================
echo -e "${HEADER}5. Terminal Graphics Protocols${RESET}"
echo "Render images in supported terminal emulators:"
echo "$DIVIDER"
cat << 'EOF'
# Auto-detect and use best protocol for your terminal
moji convert photo.jpg --protocol auto

# Force specific protocol
moji convert photo.jpg --protocol sixel     # xterm, foot, mlterm, mintty
moji convert photo.jpg --protocol kitty     # Kitty terminal
moji convert photo.jpg --protocol iterm2    # iTerm2 (macOS)

# These protocols show actual images instead of ASCII
# Useful for preview or high-quality terminal display

# Check terminal capabilities
moji term
EOF
echo ""

# ============================================================================
# 6. Batch Processing
# ============================================================================
echo -e "${HEADER}6. Batch Processing${RESET}"
echo "Convert multiple images at once:"
echo "$DIVIDER"
cat << 'EOF'
# Batch convert all JPEG files in a directory
moji batch "*.jpg" --width 80

# Batch convert with specific output directory
moji batch "photos/*.jpg" --width 100

# Process with pattern and save to different format
moji batch "images/**/*.png" --width 60

# The output is saved with .ascii extension by default
EOF
echo ""

# ============================================================================
# 7. Watch Mode
# ============================================================================
echo -e "${HEADER}7. Watch Mode - Live Updates${RESET}"
echo "Auto-update conversion when file changes:"
echo "$DIVIDER"
cat << 'EOF'
# Watch image file and auto-re-render on changes
moji convert photo.jpg --watch

# Watch with specific width
moji convert photo.jpg --width 100 --watch

# Watch with character set selection
moji convert photo.jpg --charset blocks --watch

# Useful for:
#   - Real-time preview while editing image
#   - Rapid iteration on conversion settings
#   - Slideshow effect with file updates
EOF
echo ""

# ============================================================================
# 8. Saving Output
# ============================================================================
echo -e "${HEADER}8. Saving Output${RESET}"
echo "Save conversions to files:"
echo "$DIVIDER"
cat << 'EOF'
# Save to text file
moji convert photo.jpg --output result.txt

# Save with width specification
moji convert photo.jpg --width 100 --output large.txt

# Piping to other commands
moji convert photo.jpg | less

# Pipe to file
moji convert photo.jpg > output.txt

# Combine with filters
moji convert photo.jpg --color | moji filter neon - > styled.txt
EOF
echo ""

# ============================================================================
# 9. Advanced Combinations
# ============================================================================
echo -e "${HEADER}9. Advanced Combinations${RESET}"
echo "Complex workflows combining multiple features:"
echo "$DIVIDER"
cat << 'EOF'
# High-detail conversion with edge detection and color
moji convert photo.jpg --charset braille --edge --color --width 120

# Create ASCII version from URL
moji convert photo.jpg --url "https://example.com/image.jpg" --width 80

# Watch and colorize output
moji convert sketch.jpg --edge --watch | moji filter fire -

# Batch process with specific settings
for img in *.jpg; do
    moji convert "$img" --width 100 --charset blocks --output "${img%.jpg}.ascii"
done
EOF
echo ""

# ============================================================================
# 10. Performance Tips
# ============================================================================
echo -e "${HEADER}10. Performance Tips${RESET}"
echo "$DIVIDER"
cat << 'EOF'
# Larger widths = slower processing
# --width 200 is slower than --width 80

# Character sets impact quality and speed:
#   - standard: fastest
#   - blocks: medium
#   - braille: slower but more detail

# Edge detection adds processing time
# Color preservation adds processing time

# For large batches, consider:
moji batch "*.jpg" --max-workers 4 --width 100

# Profile with time command
time moji convert large.jpg --width 100
EOF
echo ""

# ============================================================================
# 11. Troubleshooting
# ============================================================================
echo -e "${HEADER}11. Troubleshooting${RESET}"
echo "$DIVIDER"
cat << 'EOF'
# Check supported image formats
moji convert --help

# Supported formats:
#   - PNG (.png)
#   - JPEG (.jpg, .jpeg)
#   - GIF (.gif)
#   - BMP (.bmp)
#   - WebP (.webp)

# Check terminal color support
moji doctor

# Terminal protocol issues?
moji term

# Debug conversion
moji convert photo.jpg --verbose --width 50
EOF
echo ""

# ============================================================================
# Command Reference
# ============================================================================
echo -e "${HEADER}12. Quick Command Reference${RESET}"
echo "$DIVIDER"
cat << 'EOF'
# Basic conversion
moji convert photo.jpg

# With width
moji convert photo.jpg --width 120

# With character set
moji convert photo.jpg --charset blocks

# With color
moji convert photo.jpg --color

# With edge detection
moji convert photo.jpg --edge

# With protocol
moji convert photo.jpg --protocol kitty

# Watch mode
moji convert photo.jpg --watch

# Save to file
moji convert photo.jpg --output result.txt

# Batch processing
moji batch "*.jpg" --width 100

# List all character sets
moji list-charsets

# Show help
moji convert --help
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Try these commands with your own images to see the results!"
