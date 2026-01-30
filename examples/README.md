# moji Examples

Welcome to the moji examples directory! These scripts demonstrate the breadth of moji's ASCII art, terminal graphics, and creative text manipulation capabilities.

## Quick Start

All examples are standalone shell scripts. Run any of them directly:

```bash
bash examples/01-basic-banner.sh
bash examples/02-image-convert.sh
# ... and so on
```

Or make them executable and run:

```bash
chmod +x examples/*.sh
./examples/01-basic-banner.sh
```

## Examples Overview

### 1. **Basic Banner** (`01-basic-banner.sh`)
Learn the fundamentals of ASCII art banner generation with different fonts and styles.

**Topics:**
- Simple banner creation
- Font selection (slant, shadow, block, 3d, graffiti)
- Basic color styles (rainbow, fire, ice, neon)
- Border styles (single, double, round, bold, ascii)
- Text alignment (left, center, right)

**Output:** Colorful ASCII banners with various visual styles

---

### 2. **Image Conversion** (`02-image-convert.sh`)
Transform images into ASCII art with different character sets and rendering modes.

**Topics:**
- PNG/JPEG to ASCII conversion
- Character set selection (blocks, shade, braille, boxes, decorative)
- Width/height customization
- Edge detection for line art
- Terminal graphics protocols (Sixel, Kitty, iTerm2)
- Watch mode for real-time updates

**Output:** ASCII representations of images, with optional color preservation

---

### 3. **Gradient Effects** (`03-gradients.sh`)
Apply beautiful color gradients to text and banners.

**Topics:**
- Built-in gradient themes (rainbow, neon, fire, ice, matrix, sunset, ocean, dracula, vaporwave)
- Gradient modes (horizontal, vertical, diagonal)
- Per-line gradients
- Combining gradients with banners
- Custom color palettes

**Output:** Smooth color transitions across text

---

### 4. **Text Filters** (`04-filters.sh`)
Chain powerful text transformation filters together.

**Topics:**
- Filter chaining with comma-separated syntax
- Individual filters: rainbow, fire, ice, neon, matrix, glitch, metal, retro, 3d, shadow, border
- Piped input support
- Filter combinations
- Inverting colors

**Output:** Variously filtered and styled text

---

### 5. **Text Effects** (`05-effects.sh`)
Apply Unicode text effects and mathematical transformations.

**Topics:**
- 20 different text effects
- Zalgo effects (mild, medium, intense)
- Mathematical fonts (bold, italic, script, fraktur, monospace)
- Fancy text (bubble, square, smallcaps, fullwidth)
- Formatting effects (strikethrough, underline)

**Output:** Creative text transformations

---

### 6. **Pipeline Chaining** (`06-chaining.sh`)
Combine multiple moji features in powerful workflows.

**Topics:**
- Piping output between moji commands
- Combining banners with effects and filters
- Using stdin/stdout effectively
- Complex multi-step transformations
- Shell redirection tricks

**Output:** Complex multi-feature compositions

---

### 7. **Custom Themes** (`07-themes.sh`)
Create and use custom color themes for consistent styling.

**Topics:**
- Interactive TUI mode
- Theme management
- Custom color definitions
- Theme persistence
- Applying themes across commands

**Output:** Consistent themed outputs

---

### 8. **Batch Processing** (`08-batch.sh`)
Process multiple images or files in parallel.

**Topics:**
- Batch image conversion
- Parallel processing with pattern matching
- Output organization
- Loop-based processing
- Handling multiple file types

**Output:** Multiple ASCII files from batch input

---

### 9. **Watch Mode** (`09-watch-mode.sh`)
Live development with real-time updates as files change.

**Topics:**
- Banner generation with watch mode
- Image conversion with file watching
- Auto-refresh for rapid iteration
- Development workflows
- Real-time testing

**Output:** Dynamic content that updates as source files change

---

### 10. **Advanced Workflows** (`10-advanced.sh`)
Complex, production-ready examples combining multiple techniques.

**Topics:**
- Creating formatted documentation
- Terminal UI demonstrations
- System information display (neofetch-style)
- Directory tree visualization
- Calendar generation
- Interactive features
- QR code generation with styling
- Fortune and speech bubbles
- Combined ASCII art compositions

**Output:** Professional-looking terminal applications and formatted content

---

## Common Patterns

### Running a Single Script

```bash
./examples/01-basic-banner.sh
```

### Running All Examples

```bash
for script in examples/*.sh; do
  echo "=== Running $script ==="
  bash "$script"
  echo ""
done
```

### Piping Between Tools

```bash
# Combine banner generation with effects
moji banner "Hello" --font slant | \
  moji filter neon -

# Convert image to ASCII and apply filter
moji convert photo.jpg --charset blocks | \
  moji filter rainbow -
```

### Using with Shell Variables

```bash
title="My Project"
moji banner "$title" --font shadow --style rainbow
```

### Saving Output

```bash
# Save to file
moji banner "Banner" > banner.txt

# Save as PNG
moji banner "Banner" --output banner.png

# Copy to clipboard (macOS/Linux)
moji banner "Banner" --copy
```

## Feature Matrix

| Feature | Example | Key Commands |
|---------|---------|--------------|
| Banners | 01 | `moji banner` with `--font`, `--style`, `--border` |
| Image Conversion | 02 | `moji convert` with `--charset`, `--protocol`, `--color` |
| Gradients | 03 | `moji gradient` or `--gradient` flag on banners |
| Filters | 04 | `moji filter` with chainable syntax |
| Effects | 05 | `moji effect` with various transformation types |
| Chaining | 06 | Piping and shell composition |
| Themes | 07 | `moji interactive` or gradient themes |
| Batch | 08 | `moji batch` with pattern matching |
| Watch | 09 | `--watch` flag on banners and image conversion |
| Advanced | 10 | TUI, system info, calendars, QR codes, fortunes |

## Installation

Make sure `moji` is installed and available in your PATH:

```bash
# Check installation
moji --version

# If not installed, build from source
make install

# Or use go install
go install github.com/ddmoney420/moji@latest
```

## Tips & Tricks

1. **Preview Fonts:** Use `moji list-fonts` to see all available fonts
2. **List Effects:** Use `moji list-effects` to discover text transformations
3. **List Character Sets:** Use `moji list-charsets` for image conversion options
4. **List Themes:** Use `moji list-themes` for gradient options
5. **Interactive Mode:** Run `moji interactive` or `moji i` for a feature-rich TUI
6. **Shell Completions:** Generate completions: `moji completions bash > ~/.moji-completions.bash`
7. **Clipboard:** Add `--copy` to any command to copy output to clipboard
8. **JSON Output:** Add `--json` to most commands for programmatic use

## Contributing Examples

Have a cool example? Consider contributing:

1. Create a new numbered script following the pattern
2. Add clear comments explaining what's happening
3. Include output examples in comments
4. Update this README with a brief description
5. Submit a pull request

## Troubleshooting

**Command not found:** Make sure `moji` is in your PATH and built

**Colors not showing:** Check terminal color support with `moji term` or `moji doctor`

**Image conversion fails:** Verify image format is supported (PNG, JPEG, GIF, BMP, WebP)

**Watch mode not updating:** Ensure the watched file actually changes

## Related

- **Main Repository:** https://github.com/ddmoney420/moji
- **ASCII Art Database:** Check `examples/10-advanced.sh` for artdb examples
- **Kaomoji:** Browse 170+ kaomoji with `moji list`
- **Interactive Demo:** Run `moji demo` for an interactive feature showcase

Happy creating! ✨
