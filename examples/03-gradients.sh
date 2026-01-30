#!/usr/bin/env bash

##############################################################################
# Example 3: Gradient Effects
#
# This example demonstrates how to apply beautiful color gradients to text
# and banners. Gradients add visual interest and professional styling.
#
# Key Topics:
#   - Built-in gradient themes
#   - Gradient modes (horizontal, vertical, diagonal)
#   - Per-line gradients
#   - Combining gradients with banners
#   - Custom color palettes
#
##############################################################################

set -e

# Color codes for section headers
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 3: Gradient Effects${RESET}\n"

# ============================================================================
# 1. Built-in Gradient Themes
# ============================================================================
echo -e "${HEADER}1. Built-in Gradient Themes${RESET}"
echo "Apply beautiful color gradients to text:"
echo "$DIVIDER"

echo "Rainbow Gradient:"
moji gradient "RAINBOW" --theme rainbow
echo ""

echo "Neon Gradient:"
moji gradient "NEON" --theme neon
echo ""

echo "Fire Gradient:"
moji gradient "FIRE" --theme fire
echo ""

echo "Ice Gradient:"
moji gradient "ICE" --theme ice
echo ""

echo "Matrix Gradient:"
moji gradient "MATRIX" --theme matrix
echo ""

echo "Sunset Gradient:"
moji gradient "SUNSET" --theme sunset
echo ""

echo "Ocean Gradient:"
moji gradient "OCEAN" --theme ocean
echo ""

echo "Dracula Gradient:"
moji gradient "DRACULA" --theme dracula
echo ""

echo "Vaporwave Gradient:"
moji gradient "VAPORWAVE" --theme vaporwave
echo ""

echo "C64 (Retro) Gradient:"
moji gradient "RETRO" --theme c64
echo ""

echo "Pastel Gradient:"
moji gradient "PASTEL" --theme pastel
echo ""

# ============================================================================
# 2. Gradient Modes
# ============================================================================
echo -e "${HEADER}2. Gradient Modes${RESET}"
echo "Different ways to apply gradients:"
echo "$DIVIDER"

echo "Horizontal Gradient (left to right):"
moji gradient "HORIZONTAL TEXT HERE" --theme rainbow --mode horizontal
echo ""

echo "Vertical Gradient (top to bottom per-line):"
moji gradient "VERTICAL
TEXT
HERE" --theme fire --mode vertical
echo ""

echo "Diagonal Gradient:"
moji gradient "DIAGONAL TEXT" --theme ice --mode diagonal
echo ""

# ============================================================================
# 3. Combining Gradients with Banners
# ============================================================================
echo -e "${HEADER}3. Combining Gradients with Banners${RESET}"
echo "Apply gradients to ASCII art banners:"
echo "$DIVIDER"

echo "Banner with Rainbow Gradient:"
moji banner "GRADIENT" --font shadow --gradient rainbow
echo ""

echo "Banner with Fire Gradient:"
moji banner "FIRE" --font shadow --gradient fire
echo ""

echo "Banner with Neon Gradient:"
moji banner "NEON" --font shadow --gradient neon
echo ""

echo "Banner with Ocean Gradient:"
moji banner "OCEAN" --font shadow --gradient ocean
echo ""

# ============================================================================
# 4. Gradient with Effects
# ============================================================================
echo -e "${HEADER}4. Gradients Combined with Border Effects${RESET}"
echo "Professional styled output:"
echo "$DIVIDER"

echo "Gradient Banner with Border:"
moji banner "STYLED" --font shadow --gradient sunset --border bold
echo ""

echo "Gradient with Multiple Fonts:"
moji banner "MODERN" --font slant --gradient vaporwave --border double
echo ""

echo "Gradient with Alignment:"
moji banner "CENTERED" --font block --gradient dracula --border round --align center
echo ""

# ============================================================================
# 5. Per-Line Gradient Effect
# ============================================================================
echo -e "${HEADER}5. Per-Line Gradient Effect${RESET}"
echo "Gradient applied per line (vertical mode):"
echo "$DIVIDER"

cat << 'EOF'
# Create multi-line text with per-line gradients
text="Line One
Line Two
Line Three"

# Apply gradient per line
echo "$text" | moji gradient - --theme rainbow --mode vertical
EOF
echo ""

# ============================================================================
# 6. Theme Customization
# ============================================================================
echo -e "${HEADER}6. Interactive Theme Selection${RESET}"
echo "Explore themes with interactive mode:"
echo "$DIVIDER"

cat << 'EOF'
# Interactive TUI mode shows all themes live
moji interactive

# Then:
# 1. Navigate to the "Gradients" tab
# 2. See all themes with live preview
# 3. Adjust gradient mode in real-time
# 4. Copy selected theme to clipboard

# List all available themes
moji list-themes
EOF
echo ""

# ============================================================================
# 7. Piping with Gradients
# ============================================================================
echo -e "${HEADER}7. Piping Text Through Gradients${RESET}"
echo "Apply gradients to any piped input:"
echo "$DIVIDER"

cat << 'EOF'
# Gradient any text input
echo "Hello, colorful world!" | moji gradient - --theme rainbow

# Multiple lines
echo -e "First line\nSecond line\nThird line" | moji gradient - --theme fire

# From file
cat input.txt | moji gradient - --theme ocean

# Pipe from other commands
whoami | moji gradient - --theme neon
date | moji gradient - --theme pastel
EOF
echo ""

# ============================================================================
# 8. Combining with Text Effects
# ============================================================================
echo -e "${HEADER}8. Gradients with Text Effects${RESET}"
echo "Combine gradients with transformation effects:"
echo "$DIVIDER"

cat << 'EOF'
# Apply effect then gradient
moji effect bubble "BUBBLES" | moji gradient - --theme rainbow

# Gradient then effect
moji gradient "EFFECT" --theme fire | moji effect zalgo -

# Complex pipeline
echo "COOL" | moji gradient - --theme neon | moji effect sparkle -
EOF
echo ""

# ============================================================================
# 9. Gradient Export
# ============================================================================
echo -e "${HEADER}9. Saving Gradients to Files${RESET}"
echo "Export gradient output:"
echo "$DIVIDER"

cat << 'EOF'
# Save to text file
moji gradient "TEXT" --theme rainbow > output.txt

# Save banner with gradient
moji banner "HEADER" --font shadow --gradient neon > header.txt

# Export as PNG (with colors)
moji banner "TEXT" --gradient rainbow --output gradient.png

# Copy gradient to clipboard
moji gradient "TEXT" --theme fire --copy
EOF
echo ""

# ============================================================================
# 10. Theme Showcase Script
# ============================================================================
echo -e "${HEADER}10. Showcase All Themes${RESET}"
echo "Script to preview all gradient themes:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Preview all gradient themes

themes=(
    "rainbow"
    "neon"
    "fire"
    "ice"
    "matrix"
    "sunset"
    "ocean"
    "c64"
    "dracula"
    "vaporwave"
    "retro"
    "pastel"
)

text="GRADIENT THEME"

for theme in "${themes[@]}"; do
    echo "=== $theme ==="
    moji gradient "$text" --theme "$theme"
    echo ""
done
EOF
echo ""

# ============================================================================
# 11. Performance and Size Tips
# ============================================================================
echo -e "${HEADER}11. Tips & Optimization${RESET}"
echo "$DIVIDER"

cat << 'EOF'
1. Gradient Mode Selection:
   - Horizontal: Smooth left-to-right transition
   - Vertical: Applies per-line (good for multi-line)
   - Diagonal: Smooth corner-to-corner transition

2. Theme Selection Guide:
   - rainbow: Classic, friendly, visible
   - neon: Modern, high contrast, glowing feel
   - fire: Warm, energetic, warning colors
   - ice: Cool, calming, tech-forward
   - matrix: Geeky, technical, Matrix movie theme
   - sunset: Warm, aesthetic, pleasant
   - ocean: Cool, calming, water theme
   - dracula: Dark, professional, popular theme
   - vaporwave: Retro, pastel, aesthetic
   - c64: Retro computing, limited palette
   - pastel: Soft, friendly, approachable

3. Combination Tips:
   - Gradients + shadow font = professional
   - Gradients + bold border = eye-catching
   - Gradients + effects = dramatic
   - Gradients + alignment = structured

4. Performance:
   - Large text + complex gradient: slower rendering
   - Simpler themes render faster
   - Piping between commands is efficient
EOF
echo ""

# ============================================================================
# 12. Quick Command Reference
# ============================================================================
echo -e "${HEADER}12. Quick Command Reference${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# Apply gradient to text
moji gradient "TEXT" --theme rainbow

# Specify gradient mode
moji gradient "TEXT" --theme fire --mode horizontal

# Gradient on banner
moji banner "TEXT" --font shadow --gradient neon

# Gradient with border
moji banner "TEXT" --gradient sunset --border bold

# Gradient piped input
echo "text" | moji gradient - --theme ocean

# List all themes
moji list-themes

# Interactive preview
moji interactive

# Save to file
moji gradient "TEXT" --theme dracula > output.txt

# Save as PNG
moji banner "TEXT" --gradient rainbow --output banner.png

# Copy to clipboard
moji gradient "TEXT" --theme vaporwave --copy
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Experiment with different theme combinations and modes!"
