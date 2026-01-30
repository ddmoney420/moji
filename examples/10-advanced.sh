#!/usr/bin/env bash

##############################################################################
# Example 10: Advanced Workflows
#
# This example demonstrates complex, production-ready workflows that
# combine multiple moji features and techniques.
#
# Key Topics:
#   - Creating formatted documentation
#   - Terminal UI demonstrations
#   - System information display (neofetch-style)
#   - Directory tree visualization
#   - Calendar generation
#   - Interactive features
#   - QR code generation with styling
#   - Fortune and speech bubbles
#   - Combined ASCII art compositions
#
##############################################################################

set -e

# Color codes for section headers
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 10: Advanced Workflows${RESET}\n"

# ============================================================================
# 1. Interactive TUI Mode
# ============================================================================
echo -e "${HEADER}1. Interactive Terminal UI${RESET}"
echo "Explore moji's full-featured interactive mode:"
echo "$DIVIDER"

cat << 'EOF'
# Launch the interactive TUI
moji interactive

# Keyboard shortcuts:
# Tab / Shift+Tab     - Switch between tabs
# Arrow Keys          - Navigate options
# Enter               - Select/Copy
# Ctrl+C              - Exit
# Ctrl+S              - Save/Export
# q                   - Quick quit

# Available tabs:
#   1. Kaomoji         - Browse 170+ kaomoji
#   2. Banners         - Create ASCII banners with preview
#   3. Effects         - Apply text transformations
#   4. Filters         - Chain color filters
#   5. ASCII Art       - Browse ASCII art database
#   6. QR Codes        - Generate QR codes
#   7. Gradients       - Apply color gradients
#   8. Patterns        - Create decorative patterns
#   9. Animations      - See animation effects

# Workflow:
# 1. Launch: moji i
# 2. Browse available features
# 3. Test combinations
# 4. Copy results to clipboard
# 5. Export to files
# 6. Build documentation
EOF
echo ""

# ============================================================================
# 2. System Information Display
# ============================================================================
echo -e "${HEADER}2. System Information Display${RESET}"
echo "Display system info in ASCII art style:"
echo "$DIVIDER"

echo "System Info (neofetch-style):"
moji sysinfo
echo ""

cat << 'EOF'
# Alternative commands
moji info       # Same as sysinfo
moji neofetch   # Same as sysinfo

# Show system info as part of output
{
    moji banner "SYSTEM STATUS" --font shadow --gradient neon
    echo ""
    moji sysinfo
} | less

# Extract specific info
moji doctor    # Detailed diagnostics

# Show terminal capabilities
moji term
EOF
echo ""

# ============================================================================
# 3. Directory Tree Visualization
# ============================================================================
echo -e "${HEADER}3. Directory Tree Visualization${RESET}"
echo "Show project structure as ASCII tree:"
echo "$DIVIDER"

cat << 'EOF'
# Show directory tree
moji tree .

# With limited depth
moji tree . --depth 2

# With pattern filtering
moji tree src --depth 3

# Styled tree output
moji tree . --depth 2 | moji filter rainbow -

# Tree with different widths
moji tree . --depth 2 --width 100

# Complex tree usage
{
    moji banner "PROJECT STRUCTURE" --font slant --gradient neon
    echo ""
    moji tree . --depth 3 | moji filter matrix -
    echo ""
} > project_structure.txt
EOF
echo ""

# ============================================================================
# 4. Calendar and Date Displays
# ============================================================================
echo -e "${HEADER}4. Calendar Display${RESET}"
echo "Show calendars in ASCII format:"
echo "$DIVIDER"

cat << 'EOF'
# Current month calendar
moji cal

# Also works with:
moji calendar

# With effects
moji cal | moji filter rainbow -

# Styled calendar
{
    moji banner "JANUARY 2024" --font slant --gradient rainbow
    echo ""
    moji cal
} > calendar_display.txt

# Year view
moji cal --year 2024 | moji filter matrix -

# Week view
moji cal --week 15
EOF
echo ""

# ============================================================================
# 5. QR Code Generation
# ============================================================================
echo -e "${HEADER}5. QR Code Generation${RESET}"
echo "Generate ASCII QR codes with styling:"
echo "$DIVIDER"

cat << 'EOF'
# Simple QR code
moji qr "https://github.com"

# QR with specific charset
moji qr "https://github.com" --charset blocks

# QR with filter
moji qr "https://github.com" | moji filter neon -

# QR with border
moji qr "https://github.com" | moji filter border -

# Generate multiple QR codes
{
    moji banner "Quick Links" --font shadow --gradient rainbow
    echo ""

    echo "GitHub:"
    moji qr "https://github.com" --charset blocks

    echo ""
    echo "Documentation:"
    moji qr "https://docs.example.com" --charset blocks

} > qr_codes.txt

# QR code list
moji list-qr-charsets
EOF
echo ""

# ============================================================================
# 6. Kaomoji and Emoji Usage
# ============================================================================
echo -e "${HEADER}6. Kaomoji and Emoji Integration${RESET}"
echo "Use kaomoji in creative ways:"
echo "$DIVIDER"

cat << 'EOF'
# Display single kaomoji
moji shrug          # ¯\_(ツ)_/¯
moji tableflip      # (╯°□°）╯︵ ┻━┻
moji lenny          # ( ͡° ͜ʖ ͡°)

# Random kaomoji
moji random

# List all kaomoji
moji list

# Search kaomoji
moji list --search "happy"
moji list --category emotions

# Use in messages
echo "Success! $(moji cool)"
echo "Failed $(moji tableflip)"
echo "Thinking... $(moji lenny)"

# Combine with banners
{
    echo "$(moji bear) Bear necessities"
    echo "$(moji cat) Cat facts"
    echo "$(moji flex) Flex on them"
}

# Kaomoji categories
moji categories

# Create a kaomoji cheatsheet
{
    moji banner "KAOMOJI CHEATSHEET" --font slant --gradient rainbow
    echo ""

    for category in $(moji categories); do
        echo "=== $category ==="
        moji list --category "$category" | head -3
        echo ""
    done
} > kaomoji_guide.txt
EOF
echo ""

# ============================================================================
# 7. Fortune and Funny Messages
# ============================================================================
echo -e "${HEADER}7. Fortune and Messages${RESET}"
echo "Generate fortunes and funny output:"
echo "$DIVIDER"

cat << 'EOF'
# Random fortune
moji fortune

# Programming joke
moji fortune --joke

# Speech bubbles
moji say "Hello World"

# Different bubble styles
moji say "Thinking..." --bubble think

# Fortune with style
{
    moji banner "WISDOM" --font shadow --gradient rainbow
    echo ""
    moji fortune | moji filter neon -
    echo ""
    moji banner "END OF WISDOM" --font slant --gradient rainbow
} | less

# Combine fortune with effects
moji fortune | moji effect bold - | moji filter rainbow -

# Create motivational message
{
    moji say "You can do it!" --bubble speech
    echo ""
    moji fortune --joke | moji filter neon -
}

# Daily affirmation script
#!/bin/bash
echo "Good morning! Here's your inspiration:"
moji fortune | moji effect italic - | moji filter rainbow -
EOF
echo ""

# ============================================================================
# 8. ASCII Art Database
# ============================================================================
echo -e "${HEADER}8. ASCII Art Database${RESET}"
echo "Browse and use ASCII art from database:"
echo "$DIVIDER"

cat << 'EOF'
# Browse ASCII art
moji art

# List categories
moji artdb --list

# Get specific art
moji art dragon
moji art skeleton

# Art with filters
moji art dragon | moji filter fire -

# Art with gradients
moji art dragon | moji gradient - --theme rainbow

# Create a gallery
{
    moji banner "ASCII ART GALLERY" --font shadow --gradient rainbow
    echo ""

    arts=("dragon" "skeleton" "unicorn" "phoenix")
    for art in "${arts[@]}"; do
        echo "=== $art ==="
        moji art "$art" | head -10
        echo ""
    done

} > ascii_gallery.txt

# Combine art with text
{
    moji art skeleton
    echo ""
    echo "Spooky message:"
    echo "Something wicked this way comes..." | moji filter fire -
}
EOF
echo ""

# ============================================================================
# 9. Comprehensive Documentation Generation
# ============================================================================
echo -e "${HEADER}9. Generate Formatted Documentation${RESET}"
echo "Create professional-looking documentation:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Generate comprehensive README

{
    # Title with styling
    moji banner "MY AWESOME PROJECT" --font shadow --gradient neon --border double
    echo ""
    echo "A fantastic tool for amazing things"
    echo ""

    # Table of Contents
    echo "## Table of Contents"
    echo "- Installation"
    echo "- Usage"
    echo "- Features"
    echo "- Contributing"
    echo ""

    # Installation section
    moji banner "INSTALLATION" --font slant --gradient rainbow
    cat << 'INSTALL'
    \$ npm install my-project
    \$ npm run setup
    INSTALL
    echo ""

    # Features section
    moji banner "FEATURES" --font slant --gradient rainbow
    echo "$(moji star-power) Feature 1"
    echo "$(moji lightning) Feature 2"
    echo "$(moji magic) Feature 3"
    echo ""

    # Usage example
    moji banner "QUICK START" --font slant --gradient rainbow
    echo "Basic usage:"
    echo '$ my-project --help' | moji effect monospace -
    echo ""

    # System requirements
    echo "Requirements: (tested on)"
    moji sysinfo | head -5
    echo ""

    # Footer
    moji banner "HAPPY CODING" --font slant --gradient neon
    echo ""
    moji fortune --joke | moji filter neon -

} > README.md

echo "README.md generated!"
EOF
echo ""

# ============================================================================
# 10. Presentation Slides
# ============================================================================
echo -e "${HEADER}10. Create ASCII Presentation Slides${RESET}"
echo "Build a text-based presentation:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Create presentation slides

create_slide() {
    local title="$1"
    local content="$2"

    {
        moji banner "$title" --font shadow --gradient neon --border double
        echo ""
        echo "$content"
        echo ""
        echo "─────────────────────────────────────────"
        echo "[Press any key for next slide]"
    }
}

# Slide 1: Title
clear
create_slide "PRESENTATION TITLE" "Your Name - $(date +%Y-%m-%d)" | less -X

# Slide 2: Agenda
clear
create_slide "AGENDA" "$(cat << 'SLIDE'
1. Introduction
2. Main Content
3. Demo
4. Questions
SLIDE
)" | less -X

# Slide 3: Key Points
clear
create_slide "KEY POINTS" "$(echo "
$(moji star-power) Point 1: Important concept
$(moji lightning) Point 2: Critical insight
$(moji magic) Point 3: Technical detail
" | moji filter rainbow -)" | less -X

# Slide 4: Demo
clear
{
    moji banner "DEMO TIME" --font shadow --gradient fire --border bold
    echo ""
    echo "Live demonstration of feature X"
    echo ""
    moji art dragon | head -15
} | less -X

# Slide 5: Questions
clear
create_slide "QUESTIONS?" "$(moji fortune --joke | moji filter rainbow -)" | less -X
EOF
echo ""

# ============================================================================
# 11. Status Dashboard
# ============================================================================
echo -e "${HEADER}11. Status Dashboard${RESET}"
echo "Create a live status dashboard:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Live status dashboard

show_dashboard() {
    while true; do
        clear

        # Header
        moji banner "STATUS DASHBOARD" --font shadow --gradient rainbow
        echo "Last updated: $(date)"
        echo ""

        # System status
        echo "=== System Status ==="
        echo "Uptime: $(uptime -p 2>/dev/null || uptime)"
        echo ""

        # Disk usage
        echo "=== Disk Usage ==="
        df -h / | tail -1
        echo ""

        # Process count
        echo "=== Active Processes ==="
        echo "Total: $(ps aux | wc -l)"
        echo ""

        # ASCII art indicator
        if [ $(load_avg_good) -eq 1 ]; then
            echo "Status: $(moji cool) ALL GOOD"
        else
            echo "Status: $(moji tableflip) NEEDS ATTENTION"
        fi

        echo ""
        echo "Refreshing in 5 seconds... (Ctrl+C to stop)"
        sleep 5
    done
}

show_dashboard
EOF
echo ""

# ============================================================================
# 12. Color Palette Display
# ============================================================================
echo -e "${HEADER}12. Color Palette Visualization${RESET}"
echo "Display color palettes and themes:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Display available themes and palettes

{
    moji banner "COLOR THEMES" --font shadow --gradient rainbow
    echo ""

    echo "Available Gradient Themes:"
    moji list-themes
    echo ""

    # Show each theme with preview
    themes=(rainbow neon fire ice matrix)
    for theme in "${themes[@]}"; do
        echo "Theme: $theme"
        moji gradient "████████████████" --theme "$theme"
        echo ""
    done

    echo "Color Filters:"
    filters=(rainbow fire ice neon matrix glitch metal retro 3d shadow)
    for filter in "${filters[@]}"; do
        echo -n "$filter: "
        echo "███████" | moji filter "$filter" -
    done

} > palette_reference.txt

cat palette_reference.txt
EOF
echo ""

# ============================================================================
# 13. Project Initialization Script
# ============================================================================
echo -e "${HEADER}13. Project Setup Wizard${RESET}"
echo "Initialize a new project with styled output:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Project setup wizard

setup_project() {
    clear

    # Welcome
    moji banner "NEW PROJECT SETUP" --font shadow --gradient rainbow --border bold
    echo ""

    # Get project name
    read -p "Project name: " proj_name
    read -p "Author name: " author_name

    clear

    # Show progress
    echo "Creating project: $proj_name"
    echo ""

    # Create structure
    mkdir -p "$proj_name"/{src,test,docs,examples}

    echo "$(moji lightning) Creating directories..."
    sleep 1

    # Create files
    cat > "$proj_name/README.md" << README
# $proj_name

Created by: $author_name
Date: $(date)

## Description
Your project description here.
README

    echo "$(moji lightning) Creating README..."
    sleep 1

    # Final message
    clear
    moji banner "PROJECT READY" --font shadow --gradient rainbow
    echo ""
    echo "Project '$proj_name' initialized successfully! $(moji star-power)"
    echo ""
    echo "Next steps:"
    echo "  cd $proj_name"
    echo "  vim README.md"
    echo ""
}

setup_project
EOF
echo ""

# ============================================================================
# 14. Command Reference
# ============================================================================
echo -e "${HEADER}14. Advanced Commands Reference${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# Interactive features
moji interactive                          # Full-featured TUI
moji i / moji ui / moji studio           # Shortcuts

# System and info
moji sysinfo / moji info / moji neofetch # System display
moji term                                 # Terminal capabilities
moji doctor                               # System diagnostics

# Visual elements
moji tree [path]                         # Directory tree
moji cal / moji calendar                 # Calendar display
moji qr [text]                           # QR codes
moji art [name]                          # ASCII art

# Fun/Utility
moji fortune                              # Random fortune
moji fortune --joke                      # Programming joke
moji say [text]                          # Speech bubbles

# Database browsing
moji artdb --list                        # Browse art database
moji list                                # Browse kaomoji
moji list-themes                         # Show themes
moji list-effects                        # Show effects
moji list-fonts                          # Show fonts
moji categories                          # Kaomoji categories

# Configuration
moji config init                         # Initialize config
moji config show                         # Display config
moji config --path                       # Config location

# Help
moji --help                              # General help
moji [command] --help                    # Command specific help
moji demo                                # Interactive demo
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Combine these advanced techniques to create professional tools and documentation!"
