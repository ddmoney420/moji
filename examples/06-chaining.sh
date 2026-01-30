#!/usr/bin/env bash

##############################################################################
# Example 6: Pipeline Chaining
#
# This example demonstrates how to combine multiple moji features in
# powerful workflows using shell pipes and composition patterns.
#
# Key Topics:
#   - Piping output between moji commands
#   - Combining banners with effects and filters
#   - Using stdin/stdout effectively
#   - Complex multi-step transformations
#   - Shell redirection tricks
#
##############################################################################

set -e

# Color codes for section headers
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 6: Pipeline Chaining${RESET}\n"

# ============================================================================
# 1. Basic Piping - Banner to Filter
# ============================================================================
echo -e "${HEADER}1. Basic Piping: Banner → Filter${RESET}"
echo "Apply filters to banner output:"
echo "$DIVIDER"

echo "Banner piped through Rainbow filter:"
moji banner "PIPED" --font slant | moji filter rainbow -
echo ""

echo "Banner piped through Fire filter:"
moji banner "FLAMES" --font shadow | moji filter fire -
echo ""

echo "Banner piped through Neon filter:"
moji banner "GLOW" --font block | moji filter neon -
echo ""

# ============================================================================
# 2. Effect Piping - Text to Filter
# ============================================================================
echo -e "${HEADER}2. Effect Piping: Effect → Filter${RESET}"
echo "Apply filters to effect outputs:"
echo "$DIVIDER"

echo "Bubble effect with Rainbow filter:"
moji effect bubble "BUBBLE" | moji filter rainbow -
echo ""

echo "Bold text with Neon filter:"
moji effect bold "BOLD" | moji filter neon -
echo ""

echo "Script with Fire filter:"
moji effect script "SCRIPT" | moji filter fire -
echo ""

# ============================================================================
# 3. Triple Pipelines - Three Transformations
# ============================================================================
echo -e "${HEADER}3. Triple Pipelines: Effect → Filter → Piped Output${RESET}"
echo "Chain three or more transformations:"
echo "$DIVIDER"

echo "Banner → Effect → Filter:"
moji banner "TRIPLE" --font slant | moji effect bold - | moji filter rainbow -
echo ""

echo "Effect → Filter → Echo:"
echo "CHAINED" | moji effect bubble - | moji filter neon -
echo ""

echo "Banner → Custom → Filter:"
moji banner "STYLED" --font shadow --style fire | moji filter border -
echo ""

# ============================================================================
# 4. Using stdin with Pipes
# ============================================================================
echo -e "${HEADER}4. Effective stdin Usage${RESET}"
echo "Work with piped input effectively:"
echo "$DIVIDER"

cat << 'EOF'
# Direct echo piping
echo "Text" | moji effect bold - | moji filter rainbow -

# Multi-line piping
echo -e "Line 1\nLine 2\nLine 3" | moji filter neon -

# From variable
text="My Text"
echo "$text" | moji effect bubble - | moji filter fire -

# From file
cat myfile.txt | moji filter matrix -

# Command substitution
echo "$(whoami)" | moji effect smallcaps - | moji filter neon -

# Pipe from another command
date | moji filter rainbow -

# Complex: Output file through transformation
ls -la | moji filter matrix - | less
EOF
echo ""

# ============================================================================
# 5. Combining Multiple Files
# ============================================================================
echo -e "${HEADER}5. Multi-File Processing${RESET}"
echo "Process multiple inputs in sequence:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Process multiple text files with effects

files=("file1.txt" "file2.txt" "file3.txt")

for file in "${files[@]}"; do
    echo "=== $file ==="
    cat "$file" | moji filter rainbow - | moji effect bold -
    echo ""
done

# Or with different effects per file
echo "File 1:"
cat file1.txt | moji filter fire - | moji effect bold -

echo "File 2:"
cat file2.txt | moji filter ice - | moji effect italic -

echo "File 3:"
cat file3.txt | moji filter neon - | moji effect bubble -
EOF
echo ""

# ============================================================================
# 6. Creating Formatted Documentation
# ============================================================================
echo -e "${HEADER}6. Formatted Documentation Pipelines${RESET}"
echo "Create beautiful documentation with pipelines:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Create a formatted README section

echo "Creating formatted output:"
echo ""

# Title
echo "Title:" && moji banner "MY PROJECT" --font shadow | moji filter rainbow -
echo ""

# Subtitle with effects
echo "Subtitle:"
echo "A Magical Tool" | moji effect script - | moji filter rainbow -
echo ""

# Section headers
echo "Features:" && echo "=======..." | moji filter neon -
echo ""

# Bullet points with effects
echo "  ✨ Feature 1" | moji filter rainbow -
echo "  ✨ Feature 2" | moji filter rainbow -
echo "  ✨ Feature 3" | moji filter rainbow -
echo ""

# Code block (monospaced)
echo "Example:"
echo "$ command" | moji effect monospace -
echo ""

# Footer
echo "Footer:" && moji banner "END" --font slant | moji filter rainbow -
EOF
echo ""

# ============================================================================
# 7. Conditional Piping
# ============================================================================
echo -e "${HEADER}7. Conditional Piping${RESET}"
echo "Pipe different transformations based on conditions:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Conditional effect application

text="DYNAMIC"
effect_type="rainbow"  # or "fire", "neon", "ice"

case "$effect_type" in
    rainbow)
        echo "$text" | moji filter rainbow -
        ;;
    fire)
        echo "$text" | moji filter fire -
        ;;
    neon)
        echo "$text" | moji filter neon -
        ;;
    *)
        echo "$text" | moji filter matrix -
        ;;
esac

# Or with effects
text2="STYLED"
if [ "$1" = "fancy" ]; then
    echo "$text2" | moji effect bubble - | moji filter rainbow -
else
    echo "$text2" | moji effect bold - | moji filter neon -
fi
EOF
echo ""

# ============================================================================
# 8. Saving Piped Output
# ============================================================================
echo -e "${HEADER}8. Saving Pipeline Results${RESET}"
echo "Capture and save piped transformations:"
echo "$DIVIDER"

cat << 'EOF'
# Save to file
moji banner "HEADER" --font shadow | moji filter neon - > styled_header.txt

# Save with redirection
echo "Content" | moji effect bubble - | moji filter rainbow - > output.txt

# Save multiple sections
{
    moji banner "TITLE" --font shadow | moji filter rainbow -
    echo ""
    echo "Content here" | moji filter neon -
    echo ""
    moji banner "END" --font shadow | moji filter rainbow -
} > complete_file.txt

# Append to file
echo "New line" | moji filter fire - >> existing.txt

# Copy final result to clipboard
moji banner "COPY" --font slant | moji filter neon - | moji effect bold - --copy
EOF
echo ""

# ============================================================================
# 9. Piping to Other Commands
# ============================================================================
echo -e "${HEADER}9. Piping to Other Unix Tools${RESET}"
echo "Feed moji output into other commands:"
echo "$DIVIDER"

cat << 'EOF'
# Pipe to less (pagination)
moji banner "LONG OUTPUT" --font shadow | moji filter rainbow - | less

# Pipe to wc (count characters)
moji banner "COUNT" --font slant | wc -c

# Pipe to grep (filter lines)
moji banner "SEARCH" | moji filter rainbow - | grep "pattern"

# Pipe to sort
echo -e "zebra\napple\nbanana" | moji effect bold - | sort

# Pipe to head/tail
moji banner "MULTILINE" --font shadow | head -5

# Pipe to tee (display and save)
moji effect bubble "SAVE" | tee output.txt

# Combine with other tools
ls | moji filter matrix - | head -10
EOF
echo ""

# ============================================================================
# 10. Nested Pipelines (Functions)
# ============================================================================
echo -e "${HEADER}10. Reusable Pipeline Functions${RESET}"
echo "Create helper functions for common pipelines:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Define reusable pipeline functions

# Fancy header function
fancy_header() {
    local text="$1"
    moji banner "$text" --font shadow | moji filter rainbow -
}

# Styled output function
styled_output() {
    local text="$1"
    local style="${2:-rainbow}"
    echo "$text" | moji filter "$style" - | moji effect bold -
}

# Warning banner function
warning_banner() {
    local text="$1"
    echo "$text" | moji effect bubble - | moji filter fire -
}

# Usage:
fancy_header "MY PROJECT"
styled_output "Important Notice" "neon"
warning_banner "ALERT"

# Or with variables
title="Dynamic Title"
fancy_header "$title"
EOF
echo ""

# ============================================================================
# 11. Performance Optimization
# ============================================================================
echo -e "${HEADER}11. Pipeline Performance Tips${RESET}"
echo "$DIVIDER"

cat << 'EOF'
1. Pipeline Order:
   - Put fastest operations first
   - Banner generation is usually fast
   - Filters are fast
   - File I/O can be slower

2. Avoid Redundant Operations:
   # Bad: Multiple passes
   data | moji filter rainbow - > temp.txt
   cat temp.txt | moji effect bold - > final.txt

   # Good: Single pipeline
   data | moji filter rainbow - | moji effect bold - > final.txt

3. Use Direct Piping:
   # Faster: Direct pipe
   cat input.txt | moji filter neon -

   # Slower: Intermediate files
   cat input.txt > temp.txt
   moji filter neon < temp.txt

4. Batch Operations:
   # Instead of looping with echo
   { echo "line1"; echo "line2"; } | moji filter rainbow -

5. Avoid Unnecessary Filters:
   # Only chain filters you need
   # Each filter adds slight overhead
EOF
echo ""

# ============================================================================
# 12. Real-World Examples
# ============================================================================
echo -e "${HEADER}12. Real-World Pipeline Examples${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# 1. System Info Display
{
    moji banner "SYSTEM INFO" --font shadow | moji filter matrix -
    echo ""
    moji sysinfo | moji filter cyan -
} > sysinfo.txt

# 2. Project README
{
    moji banner "PROJECT" --font shadow | moji filter rainbow -
    echo ""
    echo "Description" | moji effect bold -
    echo ""
    moji banner "INSTALLATION" --font slant | moji filter neon -
    cat install.txt
    echo ""
    moji banner "USAGE" --font slant | moji filter neon -
    cat usage.txt
} > README.md

# 3. Log File Styling
tail -f app.log | moji filter matrix - | moji effect monospace -

# 4. Directory Tree with Style
moji tree . --depth 2 | moji filter rainbow - | less

# 5. Process Listing
ps aux | moji filter matrix - | head -10

# 6. Git Status Display
git status | moji effect bold - | moji filter neon -

# 7. Error Message
{
    moji banner "ERROR" --font block | moji filter fire -
    echo "Description of error" | moji effect bold -
} >&2

# 8. Progress Indicator
for i in {1..5}; do
    echo "Step $i..." | moji filter rainbow -
    sleep 1
done

# 9. Interactive Menu
{
    moji banner "MENU" --font shadow | moji filter rainbow -
    echo "1. Option One" | moji filter neon -
    echo "2. Option Two" | moji filter neon -
    echo "3. Option Three" | moji filter neon -
} | less
EOF
echo ""

# ============================================================================
# 13. Quick Reference
# ============================================================================
echo -e "${HEADER}13. Quick Pipeline Patterns${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# Basic patterns
command | moji filter FILTER -
command | moji effect EFFECT -
command | moji effect EFFECT - | moji filter FILTER -

# Banner patterns
moji banner TEXT | moji filter FILTER -
moji banner TEXT | moji effect EFFECT -
moji banner TEXT --option VALUE | moji filter FILTER -

# Multiple steps
command | moji effect EFFECT1 - | moji effect EFFECT2 - | moji filter FILTER -

# Saving
command | moji filter FILTER - > output.txt

# Appending
command | moji effect EFFECT - >> output.txt

# Display and save
command | moji filter FILTER - | tee output.txt

# Piping to other tools
command | moji filter FILTER - | sort | uniq
command | moji effect EFFECT - | grep "pattern"
command | moji filter FILTER - | less
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Combine moji commands in creative ways to build powerful workflows!"
