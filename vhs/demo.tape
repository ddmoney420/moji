# Demo recording for moji - terminal art toolkit
# Generate with: vhs vhs/demo.tape

Output vhs/demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 900
Set Height 600
Set Theme "Dracula"
Set Padding 20

# Title banner
Type "moji banner 'MOJI' --font slant --style rainbow"
Enter
Sleep 2s

# Kaomoji lookup
Type "moji shrug"
Enter
Sleep 1.5s

Type "moji tableflip"
Enter
Sleep 1.5s

# Color filter
Type "moji filter neon 'Hello World'"
Enter
Sleep 1.5s

# QR code
Type "moji qr 'https://github.com/ddmoney420/moji' --compact"
Enter
Sleep 2s

# Gradient text
Type "moji gradient 'Terminal Art Toolkit' --theme fire"
Enter
Sleep 1.5s

# Art database
Type "moji artdb dragon"
Enter
Sleep 2s

# Speech bubble with character
Type "moji say 'ASCII art is awesome!' --character owl"
Enter
Sleep 2s

# Calendar
Type "moji cal"
Enter
Sleep 1.5s

# Text effects
Type "moji effect bubble 'moji'"
Enter
Sleep 1.5s

# Pattern border
Type "moji pattern --border double --text 'Thanks for watching!'"
Enter
Sleep 2s
