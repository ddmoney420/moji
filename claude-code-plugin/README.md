```
 __  __  ___      _ ___
|  \/  |/ _ \  _ | |_ _|
| |\/| | | | || || || |
| |  | | |_| || || || |
|_|  |_|\___(_)___/|___|
```

# Moji - Claude Code Plugin

> A Claude Code plugin wrapper for the [moji CLI](https://github.com/ddmoney420/moji) terminal art toolkit.

This plugin lets you use all moji features directly in Claude Code via the `/moji` slash command.

---

## Prerequisites

You must have the `moji` CLI installed:

```bash
# macOS (Homebrew)
brew install ddmoney420/tap/moji

# Or install from source
go install github.com/ddmoney420/moji@latest
```

Verify installation:
```bash
moji --version
moji doctor
```

---

## Installation

```bash
# Install the Claude Code plugin
claude plugin install moji

# Or use directly from directory
claude --plugin-dir /path/to/this/folder
```

---

## Usage

Once installed, use `/moji` in Claude Code:

### Kaomoji
```
/moji shrug
/moji tableflip
/moji random
/moji list
```

### ASCII Banners
```
/moji banner "HELLO"
/moji banner "DEPLOY" --font=doom
/moji preview "TEXT"
/moji list-fonts
```

### Text Effects
```
/moji effect flip "TEXT"
/moji effect zalgo "CHAOS"
/moji filter glitch "ERROR"
/moji gradient "Hello" --theme=fire
/moji lolcat "Rainbow text"
```

### ASCII Art & Fortune
```
/moji fortune
/moji say "Ship it!" --character=cow
/moji art cat
/moji artdb search
```

### Image Conversion
```
/moji convert image.png
/moji convert photo.jpg --width=60
```

### QR Codes
```
/moji qr "https://github.com"
```

### Patterns & Utilities
```
/moji pattern divider
/moji cal
/moji tree
/moji sysinfo
```

---

## Full Feature List

| Category | Features |
|----------|----------|
| Kaomoji | 170+ emoticons, categories, random |
| Banners | 47 FIGlet fonts |
| Effects | flip, reverse, mirror, wave, zalgo |
| Filters | rainbow, metal, glitch, neon, 23 total |
| Gradients | 12 color themes |
| ASCII Art | Database browser, speech bubbles |
| Images | Multiple charsets, width control |
| QR Codes | Multiple styles |
| Patterns | Borders, dividers, decorations |
| Utilities | Calendar, tree, sysinfo |

---

## Examples

**Banner with doom font:**
```
/moji banner "SUCCESS" --font=doom
```
```
_______ _     _ _______ _______ _______ _______ _______
|______ |     | |       |       |______ |______ |______
______| |_____| |_____  |_____  |______ ______| ______|
```

**Fortune with character:**
```
/moji fortune
```
```
 _______________________________________
/ Why do programmers prefer dark mode?  \
\ Because light attracts bugs.          /
 ---------------------------------------
        \   ^__^
         \  (oo)\_______
            (__)\       )\/\
                ||----w |
                ||     ||
```

**Glitch effect:**
```
/moji filter glitch "SYSTEM ERROR"
```

---

## Links

- **moji CLI:** [github.com/ddmoney420/moji](https://github.com/ddmoney420/moji)
- **Claude Code Plugins:** [docs.anthropic.com](https://docs.anthropic.com/claude-code/plugins)

---

## License

MIT

---

```
 _______________
< Stay radical. >
 ---------------
   \
    \
      ¯\_(ツ)_/¯
```
