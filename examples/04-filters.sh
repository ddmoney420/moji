#!/usr/bin/env bash

##############################################################################
# Example 4: Text Filters
#
# This example demonstrates how to chain powerful text transformation
# filters together to create complex visual effects.
#
# Key Topics:
#   - Filter chaining with comma-separated syntax
#   - Individual filters (rainbow, fire, ice, neon, matrix, glitch, etc.)
#   - Piped input support
#   - Filter combinations and stacking
#   - Color and style transformations
#
##############################################################################

set -e

# Color codes for section headers
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 4: Text Filters${RESET}\n"

# ============================================================================
# 1. Individual Filters - Basic Usage
# ============================================================================
echo -e "${HEADER}1. Individual Filters${RESET}"
echo "Apply single filters to text:"
echo "$DIVIDER"

echo "Rainbow Filter:"
moji filter rainbow "Rainbow Text"
echo ""

echo "Fire Filter:"
moji filter fire "Burning Hot"
echo ""

echo "Ice Filter:"
moji filter ice "Frozen Cold"
echo ""

echo "Neon Filter:"
moji filter neon "Glowing Neon"
echo ""

echo "Matrix Filter:"
moji filter matrix "Matrix Code"
echo ""

echo "Glitch Filter:"
moji filter glitch "Glitchy Text"
echo ""

echo "Metal Filter:"
moji filter metal "Heavy Metal"
echo ""

echo "Retro Filter:"
moji filter retro "Retro Vibes"
echo ""

echo "3D Filter:"
moji filter 3d "Three Dee"
echo ""

echo "Shadow Filter:"
moji filter shadow "Shadowy"
echo ""

echo "Border Filter:"
moji filter border "Bordered"
echo ""

# ============================================================================
# 2. Combining Multiple Filters
# ============================================================================
echo -e "${HEADER}2. Filter Chaining - Combining Multiple Filters${RESET}"
echo "Use comma-separated syntax to chain filters:"
echo "$DIVIDER"

echo "Rainbow + Border:"
moji filter rainbow,border "Rainbow Box"
echo ""

echo "Fire + Shadow:"
moji filter fire,shadow "Fiery Shadow"
echo ""

echo "Neon + Border:"
moji filter neon,border "Neon Box"
echo ""

echo "Matrix + 3D:"
moji filter matrix,3d "Matrix 3D"
echo ""

echo "Ice + Shadow:"
moji filter ice,shadow "Cool Shadow"
echo ""

echo "Glitch + Border:"
moji filter glitch,border "Glitch Box"
echo ""

# ============================================================================
# 3. Three-Way Filter Combinations
# ============================================================================
echo -e "${HEADER}3. Complex Filter Chains (3+ Filters)${RESET}"
echo "Combine three or more filters for dramatic effects:"
echo "$DIVIDER"

echo "Rainbow + Border + Shadow:"
moji filter rainbow,border,shadow "Triple Effect"
echo ""

echo "Fire + Shadow + Border:"
moji filter fire,shadow,border "Fire Box Shadow"
echo ""

echo "Neon + 3D + Border:"
moji filter neon,3d,border "Neon 3D Box"
echo ""

echo "Matrix + Glitch + Shadow:"
moji filter matrix,glitch,shadow "Chaotic"
echo ""

# ============================================================================
# 4. Filter Order Matters
# ============================================================================
echo -e "${HEADER}4. Filter Application Order${RESET}"
echo "The order of filters affects the final result:"
echo "$DIVIDER"

echo "Rainbow THEN Border:"
moji filter rainbow,border "Order1"
echo ""

echo "Border THEN Rainbow (different order):"
moji filter border,rainbow "Order2"
echo ""

echo "Neon THEN Shadow:"
moji filter neon,shadow "Order3"
echo ""

echo "Shadow THEN Neon (reversed):"
moji filter shadow,neon "Order4"
echo ""

# ============================================================================
# 5. Using Piped Input
# ============================================================================
echo -e "${HEADER}5. Filtering Piped Input${RESET}"
echo "Apply filters to text from stdin:"
echo "$DIVIDER"

cat << 'EOF'
# Simple piped filter
echo "Hello World" | moji filter rainbow -

# Multiple lines with filter
echo -e "Line One\nLine Two\nLine Three" | moji filter fire -

# From file
cat myfile.txt | moji filter neon -

# From command output
whoami | moji filter matrix -

# Pipe with filter chaining
echo "Combined" | moji filter rainbow,border -

# From another moji command
moji banner "BANNER" | moji filter fire -

# Complex pipeline
moji effect bubble "EFFECT" | moji filter neon -
EOF
echo ""

# ============================================================================
# 6. Combining with Banners and Effects
# ============================================================================
echo -e "${HEADER}6. Filters with Banners and Effects${RESET}"
echo "Combine filters with other moji features:"
echo "$DIVIDER"

echo "Banner then filtered:"
moji banner "BANNER" --font slant | moji filter fire -
echo ""

echo "Effect then filtered:"
moji effect bubble "BUBBLE" | moji filter rainbow -
echo ""

echo "Banner with built-in style AND filter:"
moji banner "STYLED" --style rainbow | moji filter border -
echo ""

# ============================================================================
# 7. Color Manipulation Filters
# ============================================================================
echo -e "${HEADER}7. Color and Style Manipulation${RESET}"
echo "Filters that transform colors and styles:"
echo "$DIVIDER"

cat << 'EOF'
# Available Color Filters:
# rainbow - Rainbow spectrum colors
# fire - Red to yellow, warm tones
# ice - Blue to cyan, cool tones
# neon - Bright, glowing neon colors
# matrix - Green Matrix-style colors
# glitch - Chaotic color shifts
# metal - Metallic gray to silver
# retro - Retro 8-bit colors
# 3d - Depth-based coloring
# shadow - Dark shadow effect
# border - Box border drawing
# bold - Bold/bright text
# italic - Italic text style
# underline - Underlined text

# Experiment with combinations
echo "Text" | moji filter bold,rainbow -
echo "Text" | moji filter italic,neon -
echo "Text" | moji filter underline,fire -
EOF
echo ""

# ============================================================================
# 8. Advanced Piping with Filters
# ============================================================================
echo -e "${HEADER}8. Advanced Filter Pipelines${RESET}"
echo "Complex multi-step transformations:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Example: Create a fancy section header

# Simple approach
section="My Section"
moji banner "$section" --font shadow --style rainbow | moji filter border -

# More elaborate
echo "PROJECT ALPHA" | moji filter neon,border -
echo "========================================" | moji filter neon -
echo "Some content here" | moji filter rainbow -

# With effects
echo "Important" | moji effect bold - | moji filter fire,border -

# Multi-line with filters
echo -e "Header\nContent\nFooter" | moji filter matrix,shadow -

# Generate ASCII art, apply filters
moji banner "ART" --font graffiti | moji filter fire,shadow -
EOF
echo ""

# ============================================================================
# 9. Using lolcat for Rainbow Effects
# ============================================================================
echo -e "${HEADER}9. Rainbow Animation with lolcat${RESET}"
echo "Special rainbow command with animation:"
echo "$DIVIDER"

cat << 'EOF'
# Rainbow text with animation
moji lolcat "Rainbow Animated Text"

# With animation loop (cycles colors)
moji lolcat "Loop" --animate

# Specific speed (milliseconds between frames)
moji lolcat "Text" --speed 50

# From piped input
echo "Piped Rainbow" | moji lolcat -

# Combined with other features
moji banner "ANIMATED" --font slant | moji lolcat -
EOF
echo ""

# ============================================================================
# 10. Saving Filtered Output
# ============================================================================
echo -e "${HEADER}10. Saving Filtered Output${RESET}"
echo "Export filtered text to files:"
echo "$DIVIDER"

cat << 'EOF'
# Save filtered text to file
moji filter rainbow "TEXT" > rainbow.txt

# Save piped result
echo "Text" | moji filter fire - > fire.txt

# Save complex pipeline
moji banner "HEADER" --font shadow | moji filter neon,border - > styled.txt

# Save from stdin
cat input.txt | moji filter matrix - > matrix_output.txt

# Export as PNG (if supported)
moji banner "TEXT" --style fire --output banner.png

# Copy to clipboard
moji filter rainbow "Text" --copy

# Pipe to another command
echo "Data" | moji filter neon - | less

# Pipe to grep or other tools
moji banner "TITLE" --font slant | moji filter rainbow - | grep "pattern"
EOF
echo ""

# ============================================================================
# 11. Performance Considerations
# ============================================================================
echo -e "${HEADER}11. Performance Tips${RESET}"
echo "$DIVIDER"

cat << 'EOF'
1. Filter complexity:
   - Single filter: fast
   - Two filters: slightly slower
   - Three+ filters: noticeable slowdown
   - More complex filters (glitch): slower

2. Text length:
   - Short text (< 50 chars): instant
   - Medium text (50-500 chars): fast
   - Long text (> 500 chars): observable delay

3. Optimization tips:
   - Use simpler filter chains for real-time output
   - Save complex filtering for batch jobs
   - Test filter combinations before large batches
   - Use piping efficiently (avoid re-parsing)

4. Measurement:
   time moji filter rainbow,border "Text"
   # Shows processing time
EOF
echo ""

# ============================================================================
# 12. Creative Combinations Guide
# ============================================================================
echo -e "${HEADER}12. Creative Filter Combinations${RESET}"
echo "$DIVIDER"

cat << 'EOF'
Professional:
  - neon,border: Modern, technical look
  - matrix,shadow: Geeky, sophisticated

Friendly:
  - rainbow,border: Colorful, approachable
  - ice,shadow: Cool, friendly

Dramatic:
  - fire,shadow: Hot and dramatic
  - glitch,border: Chaotic, edgy

Retro:
  - retro,border: 8-bit, nostalgic
  - metal,shadow: Metal, bold

Playful:
  - rainbow,3d: Fun, dimensional
  - bubble effect + rainbow filter: Very playful
EOF
echo ""

# ============================================================================
# 13. Quick Command Reference
# ============================================================================
echo -e "${HEADER}13. Quick Command Reference${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# Single filter
moji filter rainbow "TEXT"

# Chained filters
moji filter rainbow,border "TEXT"

# Piped input
echo "text" | moji filter fire -

# From file
cat file.txt | moji filter neon -

# Complex pipeline
moji banner "BANNER" --font shadow | moji filter neon,border -

# Rainbow animation
moji lolcat "TEXT"

# List available effects
moji list-effects

# Save to file
moji filter rainbow "TEXT" > output.txt

# Copy to clipboard
moji filter neon "TEXT" --copy

# Show help
moji filter --help
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Experiment with different filter combinations to find your style!"
