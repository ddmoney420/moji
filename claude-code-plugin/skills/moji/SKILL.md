---
name: moji
description: Terminal art toolkit - kaomoji, ASCII banners, text effects, image conversion, QR codes, and more. Use for terminal flair, markdown decoration, and creative output.
argument-hint: <command> [args]
user-invocable: true
disable-model-invocation: false
allowed-tools: Bash
---

# Moji - Terminal Art Toolkit

Execute the user's moji command using the installed `moji` CLI.

## Command: $ARGUMENTS

## Execution

Run the command directly:

```bash
moji $ARGUMENTS
```

## Available Commands

### Kaomoji & Emoji
- `moji <name>` - Get kaomoji by name (shrug, lenny, tableflip, etc.)
- `moji random` - Random kaomoji
- `moji list` - List all kaomoji
- `moji categories` - List categories
- `moji emoji :)` - Convert ASCII smiley to emoji

### ASCII Banners
- `moji banner "TEXT"` - Generate ASCII banner
- `moji banner "TEXT" --font=doom` - Use specific font
- `moji list-fonts` - List all 47 fonts
- `moji preview "TEXT"` - Preview in multiple fonts

### Text Effects
- `moji effect flip "TEXT"` - Flip text
- `moji effect reverse "TEXT"` - Reverse text
- `moji effect wave "TEXT"` - Wave effect
- `moji effect zalgo "TEXT"` - Zalgo corruption
- `moji list-effects` - List all effects

### Colors & Gradients
- `moji filter rainbow "TEXT"` - Rainbow filter
- `moji filter metal "TEXT"` - Metal filter
- `moji filter glitch "TEXT"` - Glitch effect
- `moji gradient "TEXT" --theme=fire` - Color gradient
- `moji lolcat "TEXT"` - Lolcat rainbow
- `moji list-themes` - List gradient themes

### ASCII Art
- `moji art <name>` - Display ASCII art
- `moji artdb <query>` - Browse art database
- `moji say "TEXT"` - Speech bubble with character
- `moji fortune` - Random fortune with character

### Image Conversion
- `moji convert image.png` - Convert image to ASCII
- `moji convert image.png --width=80` - Set width
- `moji list-charsets` - List character sets

### Patterns & Borders
- `moji pattern border` - Generate border
- `moji pattern divider` - Generate divider
- `moji list-patterns` - List all patterns

### QR Codes
- `moji qr "https://example.com"` - Generate QR code
- `moji list-qr-charsets` - List QR styles

### Utilities
- `moji cal` - ASCII calendar
- `moji tree` - Directory tree
- `moji sysinfo` - System info with ASCII art
- `moji animate <effect>` - Animated effects

### Interactive
- `moji interactive` - Launch interactive studio
- `moji demo` - Run feature demo

## Output Guidelines

1. Run the moji command and capture output
2. Wrap ASCII art output in markdown code blocks (```)
3. For long output (list-fonts, etc.), summarize or truncate
4. The `--copy` flag copies to clipboard
5. Use `--no-color` if colors cause issues in output

## Common Examples

```bash
# Quick kaomoji
moji shrug
moji tableflip

# Banners
moji banner "DEPLOY" --font=doom
moji banner "SUCCESS" --font=slant

# Fortune with creature
moji fortune

# Apply effects
moji gradient "Hello World" --theme=rainbow
moji filter glitch "ERROR"

# Speech bubble
moji say "Ship it!" --character=cow
```

## Help

Run `moji --help` or `moji <command> --help` for detailed usage.
