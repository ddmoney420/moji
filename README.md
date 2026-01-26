# moji

[![Go Version](https://img.shields.io/github/go-mod/go-version/ddmoney420/moji?style=flat-square)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ddmoney420/moji?style=flat-square)](https://goreportcard.com/report/github.com/ddmoney420/moji)
[![License](https://img.shields.io/github/license/ddmoney420/moji?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/ddmoney420/moji?style=flat-square&cacheSeconds=60)](https://github.com/ddmoney420/moji/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/ddmoney420/moji.svg)](https://pkg.go.dev/github.com/ddmoney420/moji)

A terminal art toolkit with 13 feature categories. 170+ kaomoji, 47 FIGlet fonts, 20 text effects, 23 color filters, 12 gradient themes, 24 image charsets, and more.

```
  __  __  ___      _ ___
 |  \/  |/ _ \  _ | |_ _|
 | |\/| | (_) || || || |
 |_|  |_|\___/  \__/|___|
```

## Demo

![moji demo](vhs/demo.gif)

<details>
<summary>Regenerate demo GIF</summary>

Requires [VHS](https://github.com/charmbracelet/vhs):
```bash
vhs vhs/demo.tape
```
</details>

## Install

**Homebrew:**
```bash
brew install ddmoney420/tap/moji
```

**Go:**
```bash
go install github.com/ddmoney420/moji@latest
```

**Binary releases:**
Download from [Releases](https://github.com/ddmoney420/moji/releases).

**From source:**
```bash
git clone https://github.com/ddmoney420/moji.git
cd moji
make install
```

## Quick Start

```bash
# Kaomoji
moji shrug                        # => ¯\_(ツ)_/¯
moji tableflip                    # => (╯°□°）╯︵ ┻━┻

# ASCII banners
moji banner "HELLO" --font slant --style rainbow

# Color filters
moji filter neon "Glow"
moji filter fire "Burn"

# QR codes
moji qr "https://github.com/ddmoney420/moji"

# Image to ASCII
moji convert photo.jpg --color --charset blocks

# Image with Sixel/Kitty protocol (true image rendering)
moji convert photo.jpg --protocol auto

# Text effects
moji effect zalgo "Cursed"
moji effect bubble "Bubbles"

# Interactive TUI
moji interactive
```

## Features

### Kaomoji & Emoji
170+ built-in kaomoji with search and categories.

```bash
moji list                         # List all
moji list --category happy        # Filter by category
moji random                       # Random kaomoji
moji shrug                        # Direct lookup
moji emoji ":)"                   # ASCII to emoji
```

### ASCII Banners
FIGlet-powered text banners with 47 fonts and 18 color styles.

```bash
moji banner "TEXT" --font standard
moji banner "FIRE" --font slant --style fire
moji banner "COOL" --font big --style neon
moji list-fonts                   # See all fonts
moji preview "Hi"                 # Preview all fonts
```

### Color Filters
Apply color effects to any text or piped input.

```bash
moji filter rainbow "Text"
moji filter matrix "Hack"
moji filter neon "Glow"
echo "Pipe me" | moji filter fire
moji banner "WOW" | moji filter ice
```

Available: `rainbow`, `fire`, `ice`, `neon`, `matrix`, `glitch`, `metal`, `retro`, `3d`, `shadow`, `border`, `bold`, `italic`, `underline`, `invert`

### Text Effects
Unicode text transformations.

```bash
moji effect flip "Hello"          # Upside down
moji effect bubble "Hello"        # Circled letters
moji effect fraktur "Hello"       # Gothic style
moji effect zalgo "Hello"         # Corrupted text
moji list-effects                 # See all effects
```

### Image Conversion
Convert images to ASCII art or render with terminal graphics protocols.

```bash
# ASCII conversion
moji convert image.png
moji convert image.jpg --width 120 --color
moji convert image.png --charset braille --color
moji convert image.png --edge            # Line art style

# Terminal graphics protocols (true image rendering)
moji convert image.png --protocol sixel
moji convert image.png --protocol kitty
moji convert image.png --protocol iterm2
moji convert image.png --protocol auto   # Auto-detect best
```

Supported protocols:
- **Sixel** - xterm, mlterm, foot, mintty, contour
- **Kitty** - Kitty terminal
- **iTerm2** - iTerm2 on macOS

### QR Codes
Generate ASCII QR codes.

```bash
moji qr "Hello World"
moji qr "https://example.com" --charset blocks
moji qr "text" --invert --compact
```

### Patterns & Borders
Decorative borders and dividers.

```bash
moji pattern border "Framed" --style double
moji pattern divider --style stars --width 40
moji list-patterns
```

### Calendar
ASCII calendar views.

```bash
moji cal                          # Current month
moji cal --week                   # Week view
moji cal --year                   # Full year
```

### Color Gradients
Apply color gradients to text.

```bash
moji gradient "Rainbow Text" --theme rainbow
moji gradient "Sunset Vibes" --theme sunset --mode diagonal
moji list-themes
```

### System Info
Neofetch-style system information display.

```bash
moji sysinfo
moji sysinfo --theme neon
```

### Interactive TUI
Full-featured terminal UI with live preview, 9 feature tabs, and export.

```bash
moji interactive    # or: moji i, moji ui, moji studio
```

Controls: `1-9` tabs, `arrows` navigate, `Enter` copy, `e` export, `?` help

### Export
Export any output to PNG, SVG, HTML, or TXT.

```bash
moji banner "Hi" -o banner.png
moji banner "Hi" -o banner.svg
```

The TUI export modal (`e` key) includes a file browser for path selection.

## Configuration

```bash
moji config init                  # Create config file
moji config show                  # Show current config
```

Config file: `~/.config/moji/config.yaml`

## Shell Completions

```bash
moji completions bash >> ~/.bashrc
moji completions zsh >> ~/.zshrc
moji completions fish > ~/.config/fish/completions/moji.fish
```

## Diagnostics

```bash
moji doctor                       # Check terminal capabilities
moji term                         # Show terminal info
```

## Web Playground

Try moji in your browser at [ddmoney420.github.io/moji](https://ddmoney420.github.io/moji). The playground runs via WebAssembly and includes live preview, copy-to-clipboard, and save-as-image.

To run locally:
```bash
make web-serve    # Builds WASM and serves at http://localhost:8080
```

## Development

```bash
make help          # Show all targets
make build         # Build binary
make test          # Run tests
make lint          # Run linter
make ci            # Full CI pipeline (fmt, vet, lint, test, build)
make ci-quick      # Quick CI (skip lint)
make install       # Build and install to /usr/local/bin
```

**Generate terminal recordings:**
```bash
# Requires: https://github.com/charmbracelet/vhs
vhs vhs/demo.tape
vhs vhs/banner.tape
vhs vhs/tui.tape
vhs vhs/filters.tape
```

**Release:**
```bash
# Requires: https://goreleaser.com
goreleaser release --clean
```

## License

MIT
