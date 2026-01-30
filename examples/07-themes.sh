#!/usr/bin/env bash

##############################################################################
# Example 7: Custom Themes
#
# This example demonstrates theme management, theme customization, and
# how to apply consistent themes across moji commands.
#
# Key Topics:
#   - Interactive TUI mode for theme exploration
#   - Theme management and selection
#   - Custom color definitions
#   - Theme persistence in configuration
#   - Applying themes across commands
#
##############################################################################

set -e

# Color codes for section headers
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 7: Custom Themes${RESET}\n"

# ============================================================================
# 1. Exploring Themes with Interactive Mode
# ============================================================================
echo -e "${HEADER}1. Interactive Theme Explorer${RESET}"
echo "The easiest way to explore and use themes:"
echo "$DIVIDER"

cat << 'EOF'
# Launch the interactive TUI
moji interactive

# Or use these shortcuts
moji i
moji ui
moji studio

# In the interactive mode:
# 1. Navigate to the "Gradients" or "Themes" tab
# 2. Use arrow keys to browse themes
# 3. See live preview of each theme
# 4. Press Enter to copy selection
# 5. Adjust settings in real-time
# 6. Explore all 15+ built-in themes

# Other tabs available:
#   - Kaomoji: Browse 170+ kaomoji
#   - Banners: Test banner fonts and styles
#   - Effects: Apply text transformations
#   - Filters: Chain color filters
#   - ASCII Art: Browse ASCII art database
#   - QR Codes: Generate QR codes
#   - Animations: See animation effects
#   - Patterns: Create decorative patterns
EOF
echo ""

# ============================================================================
# 2. Available Built-in Themes
# ============================================================================
echo -e "${HEADER}2. Built-in Gradient Themes${RESET}"
echo "List and explore all themes:"
echo "$DIVIDER"

echo "Listing all available themes:"
moji list-themes
echo ""

cat << 'EOF'
# Save the list
moji list-themes > themes.txt

# Grep for specific themes
moji list-themes | grep -i rainbow

# Count available themes
moji list-themes | wc -l
EOF
echo ""

# ============================================================================
# 3. Using Themes with Banners
# ============================================================================
echo -e "${HEADER}3. Applying Themes to Banners${RESET}"
echo "Use gradient themes with banner command:"
echo "$DIVIDER"

cat << 'EOF'
# Apply a theme to a banner
moji banner "TEXT" --gradient rainbow

# With font and theme
moji banner "TEXT" --font shadow --gradient neon

# With border and theme
moji banner "TEXT" --gradient fire --border bold

# With alignment and theme
moji banner "TEXT" --gradient ocean --align center

# With all options
moji banner "STYLED" \
    --font shadow \
    --gradient sunset \
    --border double \
    --align center
EOF
echo ""

# ============================================================================
# 4. Using Themes with Gradients
# ============================================================================
echo -e "${HEADER}4. Applying Themes with Gradient Command${RESET}"
echo "Use themes on plain text:"
echo "$DIVIDER"

cat << 'EOF'
# Simple gradient with theme
moji gradient "TEXT" --theme rainbow

# With specific mode
moji gradient "TEXT" --theme fire --mode horizontal

# Multiple lines
echo -e "Line1\nLine2" | moji gradient - --theme neon

# In a pipeline
echo "Piped" | moji gradient - --theme dracula
EOF
echo ""

# ============================================================================
# 5. Configuration File
# ============================================================================
echo -e "${HEADER}5. Theme Configuration${RESET}"
echo "Manage themes through configuration:"
echo "$DIVIDER"

cat << 'EOF'
# Show config location
moji config --path

# Display current config
moji config show

# Initialize/reset config
moji config init

# Typical config file location:
# ~/.config/moji/config.yaml

# Example config structure:
---
theme: "rainbow"           # Default theme
gradient_mode: "horizontal"  # Default gradient mode
color_enabled: true        # Enable colors
palette:                   # Custom color definitions
  primary: "#FF00FF"
  secondary: "#00FFFF"
  accent: "#FFFF00"

# After editing config:
# 1. Restart moji
# 2. New defaults apply to all commands
# 3. Command-line args override config
EOF
echo ""

# ============================================================================
# 6. Custom Color Palettes
# ============================================================================
echo -e "${HEADER}6. Custom Color Palettes${RESET}"
echo "Create and apply custom colors:"
echo "$DIVIDER"

cat << 'EOF'
# Define colors in config file
# ~/.config/moji/config.yaml

palettes:
  cyberpunk:
    primary: "#FF006E"
    secondary: "#00F0FF"
    accent: "#00FF00"

  sunset:
    primary: "#FF6B35"
    secondary: "#F7931E"
    accent: "#FDB833"

  ocean:
    primary: "#0A2463"
    secondary: "#3E92CC"
    accent: "#A23B72"

# Usage (once configured):
# moji banner "TEXT" --palette cyberpunk
# moji gradient "TEXT" --palette sunset

# Reference popular color schemes:
# - Dracula: #282a36, #f8f8f2, #ff79c6
# - Nord: #2e3440, #88c0d0, #81a1c1
# - Solarized: #002b36, #268bd2, #2aa198
# - Material: #263238, #64b5f6, #ff5252
EOF
echo ""

# ============================================================================
# 7. Environment Variables for Theming
# ============================================================================
echo -e "${HEADER}7. Theme Environment Variables${RESET}"
echo "Control themes via environment:"
echo "$DIVIDER"

cat << 'EOF'
# Set default theme
export MOJI_THEME="rainbow"

# Set gradient mode
export MOJI_GRADIENT_MODE="vertical"

# Disable colors globally
export MOJI_COLOR_DISABLED="true"

# Set in .bashrc or .zshrc for persistence
echo 'export MOJI_THEME="neon"' >> ~/.bashrc

# Override for single command
MOJI_THEME=fire moji banner "HOT"

# Use in scripts
#!/bin/bash
export MOJI_THEME="dracula"
moji banner "MY APP" --font shadow
EOF
echo ""

# ============================================================================
# 8. Theme-Based Script Templates
# ============================================================================
echo -e "${HEADER}8. Reusable Theme Templates${RESET}"
echo "Create consistent styled output with functions:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Define theme-based helper functions

# Professional theme
professional() {
    moji banner "$1" --font shadow --gradient dracula --border bold
}

# Friendly theme
friendly() {
    moji banner "$1" --font slant --gradient rainbow --border round
}

# Tech/Hacker theme
techno() {
    moji banner "$1" --font shadow --gradient matrix --border ascii
}

# Retro theme
retro() {
    moji banner "$1" --font block --gradient c64 --border hash
}

# Modern theme
modern() {
    moji banner "$1" --font shadow --gradient neon --border double
}

# Usage:
professional "MY APP"
friendly "Welcome"
techno "HACKING"
retro "1985"
modern "2024"

# Extend for other commands:
professional_text() {
    echo "$1" | moji gradient - --theme dracula | moji effect bold -
}
EOF
echo ""

# ============================================================================
# 9. Theme Selection in Scripts
# ============================================================================
echo -e "${HEADER}9. Dynamic Theme Selection${RESET}"
echo "Choose themes based on conditions:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Select theme based on time of day or mood

get_theme() {
    local hour=$(date +%H)

    # Morning: cheerful
    if [ $hour -ge 6 ] && [ $hour -lt 12 ]; then
        echo "rainbow"
    # Afternoon: focused
    elif [ $hour -ge 12 ] && [ $hour -lt 17 ]; then
        echo "neon"
    # Evening: calm
    elif [ $hour -ge 17 ] && [ $hour -lt 21 ]; then
        echo "sunset"
    # Night: dark
    else
        echo "dracula"
    fi
}

# Usage:
theme=$(get_theme)
moji banner "HELLO" --gradient "$theme"

# Or based on mood
case "${1:-normal}" in
    happy)
        theme="rainbow"
        ;;
    work)
        theme="dracula"
        ;;
    fun)
        theme="vaporwave"
        ;;
    *)
        theme="neon"
        ;;
esac

echo "Running in $theme mode"
moji banner "MOOD: ${1:-NORMAL}" --gradient "$theme"
EOF
echo ""

# ============================================================================
# 10. Theme Cycling Script
# ============================================================================
echo -e "${HEADER}10. Theme Demo Script${RESET}"
echo "Preview all themes in sequence:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Demo all themes

text="${1:-HELLO WORLD}"

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

echo "Theme Preview: '$text'"
echo "================================"
echo ""

for theme in "${themes[@]}"; do
    echo "Theme: $theme"
    moji gradient "$text" --theme "$theme"
    echo ""
done

echo "================================"
echo "Theme preview complete!"

# Optional: Save each theme to file
for theme in "${themes[@]}"; do
    moji gradient "$text" --theme "$theme" > "theme_$theme.txt"
done
echo "Saved individual theme files"
EOF
echo ""

# ============================================================================
# 11. Theme Consistency Across Project
# ============================================================================
echo -e "${HEADER}11. Project-Wide Theme Consistency${RESET}"
echo "Apply the same theme throughout your project:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Consistent branding across all output

BRAND_THEME="neon"
BRAND_FONT="shadow"
BRAND_BORDER="double"

# Helper functions using brand theme
brand_header() {
    moji banner "$1" \
        --font "$BRAND_FONT" \
        --gradient "$BRAND_THEME" \
        --border "$BRAND_BORDER" \
        --align center
}

brand_text() {
    echo "$1" | moji gradient - --theme "$BRAND_THEME"
}

# Usage throughout project
brand_header "PROJECT ALPHA"
brand_text "Welcome to our tool"

# In main script:
brand_header "Application Starting"
echo "Running..." | moji filter neon -
brand_header "Complete"

# Result: Consistent visual theme
# Easy to rebrand: just change BRAND_THEME variable
EOF
echo ""

# ============================================================================
# 12. Theme Selection Helper
# ============================================================================
echo -e "${HEADER}12. Interactive Theme Selector${RESET}"
echo "Let users choose their theme:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Interactive theme selector

select_theme() {
    echo "Choose a theme:"
    echo "1) Rainbow"
    echo "2) Neon"
    echo "3) Fire"
    echo "4) Ice"
    echo "5) Matrix"
    echo "6) Dracula"
    echo "7) Vaporwave"
    read -p "Enter choice [1-7]: " choice

    case $choice in
        1) echo "rainbow" ;;
        2) echo "neon" ;;
        3) echo "fire" ;;
        4) echo "ice" ;;
        5) echo "matrix" ;;
        6) echo "dracula" ;;
        7) echo "vaporwave" ;;
        *) echo "neon" ;;  # default
    esac
}

# Usage:
theme=$(select_theme)
moji banner "YOUR CHOICE" --gradient "$theme"
EOF
echo ""

# ============================================================================
# 13. Documentation with Themes
# ============================================================================
echo -e "${HEADER}13. Creating Themed Documentation${RESET}"
echo "Use themes in README and docs:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Generate themed README

{
    # Title
    moji banner "MY PROJECT" --font shadow --gradient neon --border bold
    echo ""

    # Description
    echo "A beautiful command-line tool"
    echo ""

    # Section headers (with theme)
    moji banner "FEATURES" --font slant --gradient rainbow
    echo "  • Feature 1"
    echo "  • Feature 2"
    echo ""

    moji banner "INSTALLATION" --font slant --gradient rainbow
    echo "  $ npm install my-project"
    echo ""

    moji banner "USAGE" --font slant --gradient rainbow
    echo "  $ my-project --help"
    echo ""

    # Footer
    moji banner "BUILT WITH MOJI" --font slant --gradient neon

} > README.md

echo "README.md generated with theme!"
EOF
echo ""

# ============================================================================
# 14. Quick Command Reference
# ============================================================================
echo -e "${HEADER}14. Quick Command Reference${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# List all themes
moji list-themes

# Launch interactive explorer
moji interactive

# Apply theme to banner
moji banner "TEXT" --gradient rainbow

# Apply theme to text
moji gradient "TEXT" --theme fire

# Show config
moji config show

# Config path
moji config --path

# Initialize config
moji config init

# Set theme via env var
MOJI_THEME=dracula moji banner "TEXT"

# With export (persistent)
export MOJI_THEME="neon"
moji banner "TEXT"

# Theme with other options
moji banner "TEXT" --gradient sunset --border bold --font shadow
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Create your own theme-based styling system for consistent, beautiful output!"
