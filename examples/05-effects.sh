#!/usr/bin/env bash

##############################################################################
# Example 5: Text Effects
#
# This example demonstrates how to apply Unicode text transformations
# and mathematical font effects to create creative text variations.
#
# Key Topics:
#   - 20+ different text effects
#   - Zalgo effects (mild, medium, intense)
#   - Mathematical fonts (bold, italic, script, fraktur)
#   - Fancy text transformations (bubble, square, smallcaps)
#   - Formatting effects (strikethrough, underline)
#
##############################################################################

set -e

# Color codes for section headers
HEADER="\033[1;36m"  # Cyan bold
RESET="\033[0m"
DIVIDER="─────────────────────────────────────────────────────────────"

echo -e "${HEADER}Example 5: Text Effects${RESET}\n"

# ============================================================================
# 1. Basic Text Effects
# ============================================================================
echo -e "${HEADER}1. Basic Text Effects${RESET}"
echo "Fundamental text transformations:"
echo "$DIVIDER"

echo "Original: Hello World"
echo ""

echo "Reverse:"
moji effect reverse "Hello World"
echo ""

echo "Flip (upside-down):"
moji effect flip "Hello World"
echo ""

echo "Mirror (backwards):"
moji effect mirror "Hello World"
echo ""

echo "Wave:"
moji effect wave "Hello World"
echo ""

echo "Bounce:"
moji effect bounce "Hello World"
echo ""

# ============================================================================
# 2. Zalgo Effects - Chaotic Text
# ============================================================================
echo -e "${HEADER}2. Zalgo Effects - Chaotic Styling${RESET}"
echo "Add diacritics for a chaotic, creepy look:"
echo "$DIVIDER"

echo "Zalgo Mild:"
moji effect zalgo "ZALGO" --intensity mild
echo ""

echo "Zalgo Medium:"
moji effect zalgo "ZALGO" --intensity medium
echo ""

echo "Zalgo Intense:"
moji effect zalgo "ZALGO" --intensity intense
echo ""

# ============================================================================
# 3. Mathematical Fonts
# ============================================================================
echo -e "${HEADER}3. Mathematical Unicode Fonts${RESET}"
echo "Unicode mathematical text transformations:"
echo "$DIVIDER"

echo "Bold:"
moji effect bold "Bold Text"
echo ""

echo "Italic:"
moji effect italic "Italic Text"
echo ""

echo "Bold Italic:"
moji effect bolditalic "Bold Italic"
echo ""

echo "Script (cursive):"
moji effect script "Script Font"
echo ""

echo "Fraktur (Gothic):"
moji effect fraktur "Fraktur"
echo ""

echo "Double-struck (blackboard):"
moji effect doublestruck "Double Struck"
echo ""

echo "Monospace:"
moji effect monospace "Monospace Code"
echo ""

# ============================================================================
# 4. Fancy Text Boxes
# ============================================================================
echo -e "${HEADER}4. Fancy Text Transformations${RESET}"
echo "Encapsulate text in visual styles:"
echo "$DIVIDER"

echo "Bubble:"
moji effect bubble "BUBBLE"
echo ""

echo "Square:"
moji effect square "SQUARE"
echo ""

echo "Small Caps:"
moji effect smallcaps "Small Capitals"
echo ""

echo "Fullwidth (Wide characters):"
moji effect fullwidth "FULLWIDTH"
echo ""

# ============================================================================
# 5. Text Decorations
# ============================================================================
echo -e "${HEADER}5. Text Decorations${RESET}"
echo "Add decorative marks to text:"
echo "$DIVIDER"

echo "Strikethrough:"
moji effect strikethrough "Strikethrough"
echo ""

echo "Underline:"
moji effect underline "Underlined"
echo ""

echo "Overline:"
moji effect overline "Overlined"
echo ""

# ============================================================================
# 6. Combining Effects with Filters
# ============================================================================
echo -e "${HEADER}6. Effects Combined with Filters${RESET}"
echo "Apply filters to effect outputs:"
echo "$DIVIDER"

echo "Bubble + Rainbow Filter:"
moji effect bubble "COLORFUL" | moji filter rainbow -
echo ""

echo "Zalgo + Fire Filter:"
moji effect zalgo "CHAOS" --intensity medium | moji filter fire -
echo ""

echo "Bold + Neon Filter:"
moji effect bold "GLOWING" | moji filter neon -
echo ""

echo "Fraktur + Matrix Filter:"
moji effect fraktur "CODE" | moji filter matrix -
echo ""

# ============================================================================
# 7. Using Effects with Piped Input
# ============================================================================
echo -e "${HEADER}7. Effects with Piped Input${RESET}"
echo "Apply effects to stdin:"
echo "$DIVIDER"

cat << 'EOF'
# Simple effect
echo "Text" | moji effect bold -

# Zalgo from input
echo "Creepy" | moji effect zalgo - --intensity intense

# Multiple lines
echo -e "Line 1\nLine 2" | moji effect bubble -

# From command output
whoami | moji effect smallcaps -

# From file
cat mytext.txt | moji effect script -

# Pipeline combination
echo "Pipeline" | moji effect bubble - | moji filter neon -
EOF
echo ""

# ============================================================================
# 8. Effect Combinations
# ============================================================================
echo -e "${HEADER}8. Layering Effects and Filters${RESET}"
echo "Complex pipelines combining multiple transformations:"
echo "$DIVIDER"

echo "Italic Text with Rainbow Filter:"
moji effect italic "Styled" | moji filter rainbow -
echo ""

echo "Script Text with Neon Filter:"
moji effect script "Elegant" | moji filter neon -
echo ""

echo "Bubble Text with Fire Filter:"
moji effect bubble "HOT" | moji filter fire -
echo ""

echo "Zalgo with Matrix Filter:"
moji effect zalgo "GLITCH" --intensity intense | moji filter matrix -
echo ""

# ============================================================================
# 9. Creating Fancy Headers
# ============================================================================
echo -e "${HEADER}9. Creating Fancy Headers and Titles${RESET}"
echo "Combine effects for impressive headers:"
echo "$DIVIDER"

cat << 'EOF'
# Fancy section header
echo "SECTION TITLE" | moji effect bold - | moji filter neon -

# Script font header
echo "Elegant Header" | moji effect script - | moji filter rainbow -

# Bubble title
echo "PROJECT" | moji effect bubble - | moji filter fire -

# Mix techniques
moji banner "BANNER" --font shadow | moji effect bold - | moji filter neon -
EOF
echo ""

# ============================================================================
# 10. Effect Preview Script
# ============================================================================
echo -e "${HEADER}10. Preview All Effects${RESET}"
echo "Script to see all effects at once:"
echo "$DIVIDER"

cat << 'EOF'
#!/bin/bash
# Preview all text effects

text="Hello"

echo "Text Effects Preview:"
echo "===================="
echo ""

effects=(
    "reverse"
    "flip"
    "mirror"
    "wave"
    "bounce"
    "bold"
    "italic"
    "bolditalic"
    "script"
    "fraktur"
    "doublestruck"
    "monospace"
    "bubble"
    "square"
    "smallcaps"
    "fullwidth"
    "strikethrough"
    "underline"
    "overline"
)

for effect in "${effects[@]}"; do
    echo "$effect:"
    moji effect "$effect" "$text"
    echo ""
done

echo "Zalgo Effects:"
echo "Mild:"
moji effect zalgo "$text" --intensity mild
echo "Medium:"
moji effect zalgo "$text" --intensity medium
echo "Intense:"
moji effect zalgo "$text" --intensity intense
EOF
echo ""

# ============================================================================
# 11. Real-World Use Cases
# ============================================================================
echo -e "${HEADER}11. Real-World Use Cases${RESET}"
echo "$DIVIDER"

cat << 'EOF'
1. Section Headers:
   echo "INSTALLATION" | moji effect bold - | moji filter neon -

2. Error Messages:
   echo "ERROR" | moji effect bubble - | moji filter fire -

3. Success Messages:
   echo "SUCCESS" | moji effect bold - | moji filter rainbow -

4. Code Examples:
   echo "$ command" | moji effect monospace -

5. Quotes/Emphasis:
   echo "Important" | moji effect italic - | moji filter rainbow -

6. Warnings:
   echo "WARNING" | moji effect zalgo - --intensity medium | moji filter fire -

7. Fun Output:
   echo "Party" | moji effect bubble - | moji filter rainbow - | moji lolcat -

8. README Sections:
   - Use bold for important terms
   - Use script for elegant section headers
   - Use bubble for highlights
EOF
echo ""

# ============================================================================
# 12. Performance Notes
# ============================================================================
echo -e "${HEADER}12. Performance and Compatibility${RESET}"
echo "$DIVIDER"

cat << 'EOF'
1. Effect Rendering:
   - Most effects are instant
   - Zalgo with intense level: slightly slower
   - Piping to filters: minimal overhead

2. Terminal Compatibility:
   - Unicode effects require UTF-8 support
   - Check with: moji doctor
   - Some effects may not display on limited terminals
   - Mathematical fonts need good Unicode support

3. Copying to Clipboard:
   - Effects maintain Unicode: moji effect bubble "Text" --copy
   - Works on macOS, Linux (with xclip), Windows (WSL)

4. File Compatibility:
   - Save UTF-8 format: moji effect bubble "Text" > output.txt
   - Text will display in any UTF-8 compatible editor
   - Ensure terminal uses UTF-8 encoding

5. Size/Performance:
   - Text length doesn't significantly impact effects
   - Zalgo effects on very long text: might be slower
   - Efficient for most reasonable text lengths
EOF
echo ""

# ============================================================================
# 13. Creative Combinations Guide
# ============================================================================
echo -e "${HEADER}13. Creative Combination Ideas${RESET}"
echo "$DIVIDER"

cat << 'EOF'
Professional:
  echo "Title" | moji effect bold - | moji filter neon -

Elegant:
  echo "Header" | moji effect script - | moji filter rainbow -

Fun/Playful:
  echo "Fun" | moji effect bubble - | moji filter rainbow - | moji lolcat -

Technical:
  echo "Code" | moji effect monospace - | moji filter matrix -

Warning:
  echo "Alert" | moji effect zalgo - --intensity medium | moji filter fire -

Mysterious:
  echo "Secret" | moji effect fraktur - | moji filter matrix -

Retro:
  echo "Nostalgia" | moji effect fullwidth - | moji filter retro -

Elegant+Modern:
  moji banner "Title" --font shadow | moji effect bold - | moji filter neon -
EOF
echo ""

# ============================================================================
# 14. Quick Command Reference
# ============================================================================
echo -e "${HEADER}14. Quick Command Reference${RESET}"
echo "$DIVIDER"

cat << 'EOF'
# Apply effect to text
moji effect bold "TEXT"

# Effect with options
moji effect zalgo "TEXT" --intensity intense

# Piped input
echo "text" | moji effect bubble -

# Effect + filter
moji effect bold "TEXT" | moji filter neon -

# Save to file
moji effect script "TEXT" > fancy.txt

# Copy to clipboard
moji effect bubble "TEXT" --copy

# List all effects
moji list-effects

# Preview with interactive mode
moji interactive
# (then go to Effects tab)

# Help
moji effect --help
EOF
echo ""

echo -e "${HEADER}Example Complete!${RESET}"
echo "Mix and match effects to create unique text styles!"
