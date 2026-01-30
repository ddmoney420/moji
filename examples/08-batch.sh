#!/usr/bin/env bash

##############################################################################
# Example 8: Batch Processing
#
# This example demonstrates how to process multiple images or files in
# parallel using moji's batch processing capabilities and shell loops.
#
# Key Topics:
#   - Batch image conversion
#   - Parallel processing with pattern matching
#   - Output organization and file naming
#   - Loop-based batch processing
#   - Handling multiple file types
#
##############################################################################

set -e

# Color codes for section headers
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 8: Batch Processing${RESET}\n"

# ============================================================================
# 1. Basic Batch Image Conversion
# ============================================================================
echo -e "${HEADER}1. Basic Batch Image Conversion${RESET}"
echo "Convert multiple images at once:"
echo "$DIVIDER"

cat << 'EOF'
# Batch convert all JPEG files
moji batch "*.jpg"

# Batch convert with specific width
moji batch "*.jpg" --width 100

# Batch with specific charset
moji batch "*.png" --charset blocks

# Batch with color preservation
moji batch "*.jpg" --color

# Batch with output pattern
moji batch "images/*.jpg" --width 80

# Batch convert all image formats
moji batch "**/*.{jpg,png,gif}" --width 100
EOF
echo ""

# ============================================================================
# 2. Output Organization
# ============================================================================
echo -e "${HEADER}2. Organizing Batch Output${RESET}"
echo "Save batch results to specific directories:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Organize batch output into directories

# Create output directory
mkdir -p ascii_output

# Batch convert and save to output dir
for img in *.jpg; do
    output_file="ascii_output/${img%.jpg}.ascii"
    moji convert "$img" --width 100 --output "$output_file"
done

# Or with subdirectories
mkdir -p {small_width,large_width,edge_detect}

for img in *.jpg; do
    base="${img%.jpg}"
    moji convert "$img" --width 60 --output "small_width/$base.ascii"
    moji convert "$img" --width 120 --output "large_width/$base.ascii"
    moji convert "$img" --width 100 --edge --output "edge_detect/$base.ascii"
done

# Create a manifest
{
    echo "Batch Processing Results"
    echo "========================"
    echo "Generated: $(date)"
    echo ""
    ls -lh small_width/
    echo ""
    ls -lh large_width/
    echo ""
    ls -lh edge_detect/
} > results_manifest.txt
EOF
echo ""

# ============================================================================
# 3. Processing with Different Settings
# ============================================================================
echo -e "${HEADER}3. Batch with Multiple Processing Options${RESET}"
echo "Generate variants with different settings:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Generate multiple variants of each image

mkdir -p {standard,color,edge,blocks}

for img in *.jpg; do
    base="${img%.jpg}"
    echo "Processing: $img"

    # Standard ASCII
    moji convert "$img" --width 80 --output "standard/$base.ascii"

    # With color
    moji convert "$img" --width 80 --color --output "color/$base.ascii"

    # Edge detection
    moji convert "$img" --width 80 --edge --output "edge/$base.ascii"

    # Block characters
    moji convert "$img" --width 80 --charset blocks --output "blocks/$base.ascii"
done

echo "Processing complete!"
echo "Results saved in subdirectories"
EOF
echo ""

# ============================================================================
# 4. Batch Processing with Progress
# ============================================================================
echo -e "${HEADER}4. Batch Processing with Progress Indicator${RESET}"
echo "Track progress while processing multiple files:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Batch processing with progress

total=$(ls *.jpg 2>/dev/null | wc -l)
current=0

mkdir -p output

for img in *.jpg; do
    current=$((current + 1))
    percent=$((current * 100 / total))

    echo "[$percent%] Processing: $img ($current/$total)"
    moji convert "$img" --width 100 --output "output/${img%.jpg}.ascii"
done

echo ""
echo "✓ Complete! Processed $current files"
EOF
echo ""

# ============================================================================
# 5. Selective Batch Processing
# ============================================================================
echo -e "${HEADER}5. Selective Batch Processing${RESET}"
echo "Process only files matching criteria:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Process only images larger than specific size

min_size=100000  # 100KB minimum

mkdir -p large_images

for img in *.jpg; do
    size=$(stat -f%z "$img" 2>/dev/null || stat -c%s "$img")

    if [ $size -gt $min_size ]; then
        echo "Processing large image: $img"
        moji convert "$img" --width 120 --output "large_images/${img%.jpg}.ascii"
    else
        echo "Skipping small image: $img"
    fi
done

# Or process by modification date
mkdir -p today_images

for img in $(find . -maxdepth 1 -name "*.jpg" -mtime -1); do
    moji convert "$img" --width 100 --output "today_images/$(basename "$img" .jpg).ascii"
done
EOF
echo ""

# ============================================================================
# 6. Batch with Error Handling
# ============================================================================
echo -e "${HEADER}6. Robust Batch Processing${RESET}"
echo "Handle errors gracefully during batch operations:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Batch processing with error handling

mkdir -p output
mkdir -p failed

total=0
success=0
failed_count=0

for img in *.jpg; do
    total=$((total + 1))

    if moji convert "$img" --width 100 --output "output/${img%.jpg}.ascii" 2>/dev/null; then
        success=$((success + 1))
        echo "✓ $img"
    else
        failed_count=$((failed_count + 1))
        echo "✗ $img (failed)"
        cp "$img" "failed/$img"  # Save failed images
    fi
done

# Summary
echo ""
echo "================================"
echo "Batch Processing Summary"
echo "================================"
echo "Total:   $total"
echo "Success: $success"
echo "Failed:  $failed_count"
echo "Success Rate: $((success * 100 / total))%"

if [ $failed_count -gt 0 ]; then
    echo ""
    echo "Failed images copied to 'failed/' directory"
fi
EOF
echo ""

# ============================================================================
# 7. Batch Text Effects
# ============================================================================
echo -e "${HEADER}7. Batch Processing Text Files${RESET}"
echo "Apply effects to multiple text files:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Apply effects to multiple text files

mkdir -p {bold,neon,rainbow}

for txt in *.txt; do
    base="${txt%.txt}"

    # Apply bold effect
    cat "$txt" | moji effect bold - > "bold/$base.txt"

    # Apply neon filter
    cat "$txt" | moji filter neon - > "neon/$base.txt"

    # Apply rainbow filter
    cat "$txt" | moji filter rainbow - > "rainbow/$base.txt"
done

echo "Batch text processing complete!"
EOF
echo ""

# ============================================================================
# 8. Parallel Batch Processing
# ============================================================================
echo -e "${HEADER}8. Parallel Batch Processing${RESET}"
echo "Process multiple files concurrently:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Parallel batch processing (speeds up large batches)

mkdir -p output

# Method 1: Using xargs
ls *.jpg | xargs -P 4 -I {} moji convert {} --width 100 --output "output/{}.ascii"

# Method 2: Using GNU parallel (if available)
parallel moji convert {} --width 100 --output output/{}.ascii ::: *.jpg

# Method 3: Background jobs with wait
for img in *.jpg; do
    (
        echo "Processing: $img"
        moji convert "$img" --width 100 --output "output/${img%.jpg}.ascii"
    ) &

    # Limit concurrent jobs to 4
    if [ $(($(jobs -r | wc -l))) -ge 4 ]; then
        wait -n
    fi
done
wait

echo "Parallel processing complete!"

# Note: Parallel processing uses more system resources
# Good for many files, but not necessary for small batches
EOF
echo ""

# ============================================================================
# 9. Batch Processing with Logging
# ============================================================================
echo -e "${HEADER}9. Batch Processing with Detailed Logging${RESET}"
echo "Keep detailed logs of batch operations:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Batch processing with logging

LOG_FILE="batch_processing_$(date +%Y%m%d_%H%M%S).log"
mkdir -p output

{
    echo "Batch Processing Log"
    echo "Started: $(date)"
    echo "========================================"
    echo ""
} > "$LOG_FILE"

for img in *.jpg; do
    base="${img%.jpg}"
    output_file="output/$base.ascii"

    {
        echo "Processing: $img"
        echo "Output: $output_file"
        echo "Time: $(date)"
    } >> "$LOG_FILE"

    if moji convert "$img" --width 100 --output "$output_file" 2>&1 | tee -a "$LOG_FILE"; then
        echo "Status: Success" >> "$LOG_FILE"
    else
        echo "Status: Failed" >> "$LOG_FILE"
    fi

    echo "" >> "$LOG_FILE"
done

{
    echo "========================================"
    echo "Completed: $(date)"
} >> "$LOG_FILE"

echo "Log saved to: $LOG_FILE"
cat "$LOG_FILE"
EOF
echo ""

# ============================================================================
# 10. Batch Conversion with Verification
# ============================================================================
echo -e "${HEADER}10. Batch Processing with Verification${RESET}"
echo "Verify batch results quality:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Batch processing with verification

mkdir -p output
mkdir -p verify

for img in *.jpg; do
    base="${img%.jpg}"
    output_file="output/$base.ascii"

    echo "Processing: $img"
    moji convert "$img" --width 100 --output "$output_file"

    # Verify output was created and has content
    if [ -s "$output_file" ]; then
        echo "✓ Output verified: $(wc -l < "$output_file") lines"
        cp "$output_file" "verify/$base.ascii"
    else
        echo "✗ Output verification failed"
        rm -f "$output_file"
    fi
done

echo ""
echo "Verified files in 'verify/' directory"
EOF
echo ""

# ============================================================================
# 11. Batch with Custom Naming
# ============================================================================
echo -e "${HEADER}11. Custom Naming in Batch Processing${RESET}"
echo "Organize output with custom naming schemes:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Batch processing with custom naming

mkdir -p output

for img in *.jpg; do
    base="${img%.jpg}"

    # Option 1: Add timestamp
    output="output/${base}_$(date +%s).ascii"
    moji convert "$img" --width 100 --output "$output"

    # Option 2: Add suffix for processing type
    moji convert "$img" --width 80 --output "output/${base}_80wide.ascii"
    moji convert "$img" --width 120 --output "output/${base}_120wide.ascii"
    moji convert "$img" --width 80 --color --output "output/${base}_color.ascii"

    # Option 3: Categorize by size
    size=$(ls -lh "$img" | awk '{print $5}')
    output="output/${base}_${size}.ascii"
    moji convert "$img" --width 100 --output "$output"
done

echo "Batch complete with custom naming"
EOF
echo ""

# ============================================================================
# 12. Batch Processing Functions
# ============================================================================
echo -e "${HEADER}12. Reusable Batch Functions${RESET}"
echo "Create helper functions for batch operations:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Reusable batch processing functions

# Function: Batch convert images
batch_convert_images() {
    local pattern="$1"
    local width="${2:-80}"
    local output_dir="${3:-output}"

    mkdir -p "$output_dir"

    for img in $pattern; do
        [ -f "$img" ] || continue
        base="${img%.*}"
        moji convert "$img" --width "$width" --output "$output_dir/$base.ascii"
    done
}

# Function: Batch apply effects
batch_apply_effects() {
    local pattern="$1"
    local effect="$2"
    local output_dir="${3:-output}"

    mkdir -p "$output_dir"

    for file in $pattern; do
        [ -f "$file" ] || continue
        base="${file%.*}"
        cat "$file" | moji effect "$effect" - > "$output_dir/$base.txt"
    done
}

# Function: Batch apply filters
batch_apply_filters() {
    local pattern="$1"
    local filter="$2"
    local output_dir="${3:-output}"

    mkdir -p "$output_dir"

    for file in $pattern; do
        [ -f "$file" ] || continue
        base="${file%.*}"
        cat "$file" | moji filter "$filter" - > "$output_dir/$base.txt"
    done
}

# Usage:
# batch_convert_images "*.jpg" 100
# batch_apply_effects "*.txt" bold
# batch_apply_filters "*.txt" rainbow output_filtered
EOF
echo ""

# ============================================================================
# 13. Quick Command Reference
# ============================================================================
echo -e "${HEADER}13. Quick Command Reference${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# Basic batch conversion
moji batch "*.jpg"

# With specific width
moji batch "*.jpg" --width 100

# With character set
moji batch "*.png" --charset blocks

# With color
moji batch "*.jpg" --color

# With output pattern
moji batch "images/*.jpg" --width 80

# Loop-based batch (more control)
for img in *.jpg; do
    moji convert "$img" --width 100 --output "output/${img%.jpg}.ascii"
done

# With progress
for img in *.jpg; do
    echo "Processing: $img"
    moji convert "$img" --width 100 --output "output/${img%.jpg}.ascii"
done

# Parallel (if xargs available)
ls *.jpg | xargs -P 4 -I {} moji convert {} --width 100 --output "output/{}.ascii"

# Save log
moji batch "*.jpg" 2>&1 | tee batch.log

# Help
moji batch --help
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Use batch processing to handle large collections of images efficiently!"
