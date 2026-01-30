#!/usr/bin/env bash

##############################################################################
# Example 1: Basic Banner Creation
#
# This example demonstrates the fundamentals of ASCII art banner generation
# with moji. Learn how to create banners using different fonts, colors,
# styles, and borders.
#
# Key Topics:
#   - Simple banner creation
#   - Font selection (slant, shadow, block, 3d, graffiti)
#   - Color styles (rainbow, fire, ice, neon, matrix, glitch)
#   - Border styles (single, double, round, bold, ascii)
#   - Text alignment
#
##############################################################################

set -e

# Color codes for section headers in output
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 1: Basic Banner Creation${RESET}\n"

# ============================================================================
# 1. Simple Banner - Default Style
# ============================================================================
echo -e "${HEADER}1. Simple Banner (Default)${RESET}"
echo "Command: moji banner 'HELLO'"
echo "$DIVIDER"
moji banner "HELLO"
echo ""

# ============================================================================
# 2. Different Fonts - Visual Variety
# ============================================================================
echo -e "${HEADER}2. Different Fonts${RESET}"
echo "Font styles showcase:"
echo "$DIVIDER"

echo "Slant Font:"
moji banner "SLANT" --font slant
echo ""

echo "Shadow Font:"
moji banner "SHADOW" --font shadow
echo ""

echo "Block Font:"
moji banner "BLOCK" --font block
echo ""

echo "3D Font:"
moji banner "3D" --font 3d
echo ""

echo "Graffiti Font:"
moji banner "ART" --font graffiti
echo ""

# ============================================================================
# 3. Color Styles - Add Visual Interest
# ============================================================================
echo -e "${HEADER}3. Color Styles${RESET}"
echo "Different color schemes:"
echo "$DIVIDER"

echo "Rainbow Style:"
moji banner "RAINBOW" --font shadow --style rainbow
echo ""

echo "Fire Style:"
moji banner "FIRE" --font shadow --style fire
echo ""

echo "Ice Style:"
moji banner "ICE" --font shadow --style ice
echo ""

echo "Neon Style:"
moji banner "NEON" --font shadow --style neon
echo ""

echo "Matrix Style:"
moji banner "MATRIX" --font shadow --style matrix
echo ""

# ============================================================================
# 4. Border Styles - Frame Your Banner
# ============================================================================
echo -e "${HEADER}4. Border Styles${RESET}"
echo "Different border options:"
echo "$DIVIDER"

echo "Bold Border:"
moji banner "BOLD" --font slant --border bold
echo ""

echo "Double Border:"
moji banner "DOUBLE" --font slant --border double
echo ""

echo "Round Border:"
moji banner "ROUND" --font slant --border round
echo ""

echo "Ascii Border:"
moji banner "ASCII" --font slant --border ascii
echo ""

echo "Stars Border:"
moji banner "STARS" --font slant --border stars
echo ""

# ============================================================================
# 5. Text Alignment - Position Control
# ============================================================================
echo -e "${HEADER}5. Text Alignment${RESET}"
echo "Alignment options:"
echo "$DIVIDER"

echo "Left Aligned:"
moji banner "LEFT" --font slant --border bold --align left
echo ""

echo "Center Aligned (default):"
moji banner "CENTER" --font slant --border bold --align center
echo ""

echo "Right Aligned:"
moji banner "RIGHT" --font slant --border bold --align right
echo ""

# ============================================================================
# 6. Combining Features - Multi-Attribute Banners
# ============================================================================
echo -e "${HEADER}6. Combining Features - Creating Professional Banners${RESET}"
echo "Using multiple attributes together:"
echo "$DIVIDER"

echo "Professional Header:"
moji banner "PROJECT ALPHA" --font shadow --style neon --border double --align center
echo ""

echo "Section Divider:"
moji banner "SECTION" --font slant --style fire --border bold --align center
echo ""

echo "Attention Grabber:"
moji banner "WARNING" --font 3d --style fire --border stars --align center
echo ""

# ============================================================================
# 7. Command Reference
# ============================================================================
echo -e "${HEADER}7. Quick Command Reference${RESET}"
echo "$DIVIDER"
cat << 'EOF'
# List all available fonts
moji list-fonts

# Basic usage
moji banner "TEXT"

# With specific font
moji banner "TEXT" --font shadow

# With color style
moji banner "TEXT" --style rainbow

# With border
moji banner "TEXT" --border double

# With alignment
moji banner "TEXT" --align left

# Combining options
moji banner "TEXT" --font shadow --style neon --border double --align center

# Save to file
moji banner "TEXT" > banner.txt

# Export as PNG (with colors)
moji banner "TEXT" --output banner.png

# Copy to clipboard
moji banner "TEXT" --copy

# Show help
moji banner --help
EOF
echo ""

# ============================================================================
# Tips
# ============================================================================
echo -e "${HEADER}Tips & Tricks${RESET}"
echo "$DIVIDER"
cat << 'EOF'
1. Use --font shadow for chunky, bold letters
2. Use --font slant for italicized, modern look
3. Combine fire or ice styles with bold borders for impact
4. Use --border round for friendly, approachable headers
5. Try different font+style combinations to find your style
6. Use --align to match surrounding text layout
7. Export as PNG for presentations and documentation
8. Test colors with --style before final version
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Try modifying the commands above to create your own banners."
