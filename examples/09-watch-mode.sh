#!/usr/bin/env bash

##############################################################################
# Example 9: Watch Mode
#
# This example demonstrates live development with real-time updates
# as files change. Watch mode enables rapid iteration and testing.
#
# Key Topics:
#   - Banner generation with watch mode
#   - Image conversion with file watching
#   - Auto-refresh for rapid iteration
#   - Development workflows
#   - Real-time testing and debugging
#
##############################################################################

set -e

# Color codes for section headers
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 9: Watch Mode${RESET}\n"

# ============================================================================
# 1. Understanding Watch Mode
# ============================================================================
echo -e "${HEADER}1. What is Watch Mode?${RESET}"
echo "Watch mode monitors files for changes and auto-updates output:"
echo "$DIVIDER"

cat << 'EOF'
# Watch mode features:
# - Monitors input file for changes
# - Automatically re-renders output when file changes
# - Useful for live preview and testing
# - Speeds up development workflow
# - Works with most moji commands

# Supported in:
# - moji banner (with text input file)
# - moji convert (for images)
# - moji gradient (with text input)
# - Many other commands
EOF
echo ""

# ============================================================================
# 2. Basic Watch Mode Usage
# ============================================================================
echo -e "${HEADER}2. Basic Watch Mode for Image Conversion${RESET}"
echo "Real-time ASCII conversion as image changes:"
echo "$DIVIDER"

cat << 'EOF'
# Watch a single image file
moji convert photo.jpg --watch

# Watch with specific width
moji convert photo.jpg --width 100 --watch

# Watch with character set
moji convert photo.jpg --charset blocks --watch

# Watch with color preservation
moji convert photo.jpg --color --watch

# Watch with edge detection
moji convert photo.jpg --edge --watch

# How to use:
# 1. Edit the image in an image editor
# 2. Save the changes
# 3. moji automatically detects the change
# 4. ASCII output updates in real-time
# 5. Press Ctrl+C to stop watching
EOF
echo ""

# ============================================================================
# 3. Watch Mode with Text Files
# ============================================================================
echo -e "${HEADER}3. Watch Mode for Text Processing${RESET}"
echo "Monitor text files and apply effects in real-time:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Real-time text processing with watch mode

# Create a text file to monitor
echo "My text here" > watch_me.txt

# Watch and apply effect (manual polling approach)
while true; do
    clear
    echo "=== Live Text Processing ==="
    echo "File: watch_me.txt"
    echo ""
    echo "Bold version:"
    cat watch_me.txt | moji effect bold -
    echo ""
    echo "Rainbow version:"
    cat watch_me.txt | moji filter rainbow -
    echo ""
    echo "Neon version:"
    cat watch_me.txt | moji filter neon -
    echo ""
    echo "Watching for changes... (Ctrl+C to stop)"
    sleep 2
done
EOF
echo ""

# ============================================================================
# 4. Watch Mode in Development
# ============================================================================
echo -e "${HEADER}4. Development Workflow with Watch Mode${RESET}"
echo "Use watch mode during development for rapid iteration:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Development workflow example

# Create a project structure
mkdir -p project/{images,output}

# In one terminal: Watch and convert images
cd project
moji convert images/screenshot.png --width 100 --watch

# In another terminal: Edit your image
# - Make changes in image editor
# - Save file
# - Watch moji window automatically updates!

# Use case: Screenshot documentation
# 1. Take screenshot
# 2. Watch mode displays ASCII version immediately
# 3. Refine screenshot
# 4. Updated ASCII shows changes in real-time
# 5. No need to re-run commands
EOF
echo ""

# ============================================================================
# 5. Watch Mode with Multiple Files
# ============================================================================
echo -e "${HEADER}5. Watching Multiple Files${RESET}"
echo "Monitor several files simultaneously:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Watch multiple files in parallel

# Create test images
files=("image1.jpg" "image2.jpg" "image3.png")

# Start watch sessions in parallel
for file in "${files[@]}"; do
    (
        echo "Watching: $file"
        moji convert "$file" --width 80 --watch
    ) &
done

# All watch windows run independently
# Press Ctrl+C in each to stop

# In practice, use separate terminal tabs/windows:
# Tab 1: moji convert image1.jpg --watch
# Tab 2: moji convert image2.jpg --watch
# Tab 3: moji convert image3.jpg --watch
EOF
echo ""

# ============================================================================
# 6. Watch Mode with Logging
# ============================================================================
echo -e "${HEADER}6. Watch Mode with Logging${RESET}"
echo "Keep logs of watch mode activity:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Watch mode with logging

LOG_FILE="watch_$(date +%Y%m%d_%H%M%S).log"

{
    echo "Watch Mode Session"
    echo "Started: $(date)"
    echo "File: photo.jpg"
    echo "Width: 100"
    echo "==============================="
    echo ""
} > "$LOG_FILE"

# Run watch mode, capture to log
{
    while true; do
        timestamp=$(date +"%Y-%m-%d %H:%M:%S")
        echo "[$timestamp] Checking for changes..."
        moji convert photo.jpg --width 100
        echo ""
        sleep 5  # Check every 5 seconds
    done
} >> "$LOG_FILE" 2>&1 &

WATCH_PID=$!

echo "Watch mode running (PID: $WATCH_PID)"
echo "Log: $LOG_FILE"
echo "Press enter to stop..."
read

kill $WATCH_PID
echo "Watch session ended at $(date)" >> "$LOG_FILE"
EOF
echo ""

# ============================================================================
# 7. Watch Mode with Version Control
# ============================================================================
echo -e "${HEADER}7. Watch Mode for Documentation Updates${RESET}"
echo "Auto-generate ASCII when docs are edited:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Auto-update ASCII art in documentation

# Setup: Create doc template
cat > doc_template.txt << 'DOC'
# My Project

## Screenshot

```
[ASCII art will be inserted here]
```

## Features
- Feature 1
- Feature 2
DOC

# Watch mode polling approach
watch_and_update() {
    local image="screenshot.png"
    local doc="README.md"

    while true; do
        # Check if image changed
        if [ "$image" -nt "${image}.timestamp" ]; then
            echo "Image changed, updating documentation..."

            # Generate ASCII
            ascii_output=$(moji convert "$image" --width 80)

            # Update documentation (in real project, parse and replace section)
            echo "$ascii_output" > temp_ascii.txt

            # Could use sed/awk to insert into README
            echo "Documentation updated"

            touch "${image}.timestamp"
        fi

        sleep 5
    done
}

# Usage
watch_and_update &
EOF
echo ""

# ============================================================================
# 8. Real-Time Preview Scripts
# ============================================================================
echo -e "${HEADER}8. Real-Time Preview Dashboard${RESET}"
echo "Create a live updating dashboard:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Live dashboard with multiple watch streams

monitor_dashboard() {
    while true; do
        clear
        echo "┌─────────────────────────────────────────┐"
        echo "│     Live Conversion Dashboard           │"
        echo "└─────────────────────────────────────────┘"
        echo ""

        echo "Image: photo.jpg (width: 80)"
        echo "────────────────────────────────────────"
        moji convert photo.jpg --width 80 | head -10
        echo "... (truncated)"
        echo ""

        echo "Image: screenshot.png (width: 100)"
        echo "────────────────────────────────────────"
        moji convert screenshot.png --width 100 | head -10
        echo "... (truncated)"
        echo ""

        echo "Last updated: $(date +%H:%M:%S)"
        echo "Refreshing every 5 seconds... (Ctrl+C to stop)"

        sleep 5
    done
}

monitor_dashboard
EOF
echo ""

# ============================================================================
# 9. Watch Mode for Testing
# ============================================================================
echo -e "${HEADER}9. Using Watch Mode for Testing${RESET}"
echo "Test moji features with live changes:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Testing workflow with watch mode

# Setup test environment
mkdir -p test_images
mkdir -p test_output

# Create test image
echo "Test image ready"

# Terminal 1: Watch conversion
echo "Starting watch mode for testing..."
moji convert test_images/test.png --width 100 --watch

# While that runs:
# - Edit image in graphics program
# - Save changes
# - Watch moji updates immediately
# - Verify output looks correct
# - No need to run command again

# Test different parameters:
# Stop watch (Ctrl+C)
moji convert test_images/test.png --width 120 --watch  # Different width
moji convert test_images/test.png --charset blocks --watch  # Different charset
moji convert test_images/test.png --edge --watch  # Edge detection

# This rapid iteration helps find best settings
EOF
echo ""

# ============================================================================
# 10. Watch Mode Performance
# ============================================================================
echo -e "${HEADER}10. Performance Considerations${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# Watch mode performance tips:

1. File size matters:
   - Small images (< 500KB): instant updates
   - Medium images (500KB - 2MB): quick updates (< 1 sec)
   - Large images (> 2MB): may take longer (1-5 sec)
   - Test with your typical file sizes

2. Processing settings:
   - Simple settings: very fast
   - Edge detection: adds processing time
   - Color preservation: slight overhead
   - Larger width: more processing

3. System resources:
   - Watch mode uses minimal memory
   - CPU usage spikes during conversion
   - Suitable for most modern systems
   - Don't run too many watch sessions simultaneously

4. Optimization:
   - Use reasonable width (80-120)
   - Choose simple character sets for faster feedback
   - Disable features you don't need
   - Test with your file types first

# Example: Optimized watch
moji convert photo.jpg --width 80 --watch  # Fast
moji convert photo.jpg --width 200 --edge --color --watch  # Slower
EOF
echo ""

# ============================================================================
# 11. Watch Mode Troubleshooting
# ============================================================================
echo -e "${HEADER}11. Troubleshooting Watch Mode${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# Problem: Watch mode not detecting changes
Solution:
  - Ensure you're saving the file (not just editing in memory)
  - Some editors buffer writes; try explicit save/flush
  - Check file permissions

# Problem: Watch mode shows old content
Solution:
  - Wait a moment for the file to fully save
  - Some systems have file sync delays
  - Try restarting watch mode

# Problem: Watch mode is too slow
Solution:
  - Reduce image width (use --width 80 instead of 200)
  - Use simpler character set (--charset standard)
  - Disable edge detection if not needed
  - Use faster system or close other apps

# Problem: Watch mode exits unexpectedly
Solution:
  - Check the error message
  - Verify file still exists
  - Check disk space
  - Restart watch mode

# Debugging watch mode:
moji convert photo.jpg --watch --verbose  # If supported
moji doctor  # Check system compatibility
EOF
echo ""

# ============================================================================
# 12. Use Cases and Workflows
# ============================================================================
echo -e "${HEADER}12. Real-World Use Cases${RESET}"
echo "$DIVIDER"

cat << 'EOF'
1. Screenshot to ASCII Documentation:
   - Take screenshot
   - moji convert screenshot.png --watch
   - Edit screenshot if needed
   - ASCII updates automatically
   - Copy and paste into docs

2. Image Processing Pipeline:
   - Start moji in one terminal
   - Run image editor in another
   - Make changes, save
   - See ASCII result immediately
   - Refine until satisfied

3. Thumbnail Preview:
   - Watch mode shows live thumbnails
   - Useful for batch image review
   - Watch each image as you process it

4. Tutorial/Demo Recording:
   - Record yourself using watch mode
   - Show file changing and ASCII updating
   - Great for demonstrations

5. Design Iteration:
   - Design in graphics software
   - Watch moji for immediate feedback
   - Iterate quickly
   - Determine best ASCII representation

6. Testing Conversions:
   - Test different settings with watch mode
   - No need to retype commands
   - Quick A/B comparison
   - Find optimal parameters

7. Continuous Preview:
   - Background process watching files
   - Dashboard showing multiple conversions
   - Good for monitoring batch processing
EOF
echo ""

# ============================================================================
# 13. Advanced Watch Patterns
# ============================================================================
echo -e "${HEADER}13. Advanced Watch Patterns${RESET}"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Advanced watch mode patterns

# Pattern 1: Watch multiple conversions
watch_multi() {
    local image="$1"

    echo "Standard:"
    moji convert "$image" --width 80 &
    pid1=$!

    echo "With Color:"
    moji convert "$image" --width 80 --color &
    pid2=$!

    echo "Edge Detect:"
    moji convert "$image" --width 80 --edge &
    pid3=$!

    wait $pid1 $pid2 $pid3
}

# Pattern 2: Conditional watch
watch_if_needed() {
    local file="$1"
    local last_mtime=0

    while true; do
        mtime=$(stat -f%m "$file" 2>/dev/null)

        if [ "$mtime" != "$last_mtime" ]; then
            clear
            echo "Updated: $(date)"
            moji convert "$file" --width 100
            last_mtime=$mtime
        fi

        sleep 1
    done
}

# Pattern 3: Watch with notification
watch_with_notify() {
    local file="$1"

    while true; do
        echo "Watching: $file"
        moji convert "$file" --width 100

        # Notify when done (macOS)
        osascript -e 'display notification "ASCII conversion complete"'

        sleep 5
    done
}

# Usage:
# watch_multi photo.jpg
# watch_if_needed image.png &
# watch_with_notify screenshot.jpg
EOF
echo ""

# ============================================================================
# 14. Quick Command Reference
# ============================================================================
echo -e "${HEADER}14. Quick Command Reference${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# Basic watch
moji convert photo.jpg --watch

# Watch with width
moji convert photo.jpg --width 100 --watch

# Watch with charset
moji convert photo.jpg --charset blocks --watch

# Watch with color
moji convert photo.jpg --color --watch

# Watch with edge
moji convert photo.jpg --edge --watch

# Watch with multiple options
moji convert photo.jpg --width 80 --color --edge --watch

# To stop watch mode
# Press: Ctrl+C

# Manual monitoring (polling)
while true; do
    clear
    moji convert photo.jpg --width 100
    sleep 5
done

# Help
moji convert --help
# Look for --watch option
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Use watch mode for rapid iteration and live previewing of conversions!"
