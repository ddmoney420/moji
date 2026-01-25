package tree

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Options for tree generation
type Options struct {
	MaxDepth   int
	ShowHidden bool
	DirsOnly   bool
	FilesOnly  bool
	Pattern    string
	MaxItems   int
	ShowSize   bool
	SortBySize bool
	Indent     string
}

// DefaultOptions returns default tree options
func DefaultOptions() Options {
	return Options{
		MaxDepth:   -1, // unlimited
		ShowHidden: false,
		Indent:     "│   ",
	}
}

// Entry represents a file/directory entry
type Entry struct {
	Name    string
	Path    string
	IsDir   bool
	Size    int64
	Mode    os.FileMode
	Entries []Entry
}

// Generate creates a tree structure from a path
func Generate(root string, opts Options) (Entry, error) {
	info, err := os.Stat(root)
	if err != nil {
		return Entry{}, err
	}

	entry := Entry{
		Name:  filepath.Base(root),
		Path:  root,
		IsDir: info.IsDir(),
		Size:  info.Size(),
		Mode:  info.Mode(),
	}

	if info.IsDir() {
		entries, err := readDir(root, opts, 0)
		if err != nil {
			return entry, err
		}
		entry.Entries = entries
	}

	return entry, nil
}

func readDir(path string, opts Options, depth int) ([]Entry, error) {
	if opts.MaxDepth >= 0 && depth >= opts.MaxDepth {
		return nil, nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for _, f := range files {
		name := f.Name()

		// Skip hidden files
		if !opts.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// Pattern matching
		if opts.Pattern != "" {
			matched, _ := filepath.Match(opts.Pattern, name)
			if !matched {
				continue
			}
		}

		info, err := f.Info()
		if err != nil {
			continue
		}

		isDir := f.IsDir()

		// Filter dirs/files only
		if opts.DirsOnly && !isDir {
			continue
		}
		if opts.FilesOnly && isDir {
			continue
		}

		entry := Entry{
			Name:  name,
			Path:  filepath.Join(path, name),
			IsDir: isDir,
			Size:  info.Size(),
			Mode:  info.Mode(),
		}

		if isDir {
			subEntries, _ := readDir(entry.Path, opts, depth+1)
			entry.Entries = subEntries
		}

		entries = append(entries, entry)

		if opts.MaxItems > 0 && len(entries) >= opts.MaxItems {
			break
		}
	}

	// Sort
	if opts.SortBySize {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Size > entries[j].Size
		})
	} else {
		sort.Slice(entries, func(i, j int) bool {
			// Dirs first, then alphabetical
			if entries[i].IsDir != entries[j].IsDir {
				return entries[i].IsDir
			}
			return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
		})
	}

	return entries, nil
}

// Format returns ASCII tree representation
func Format(entry Entry, opts Options) string {
	var sb strings.Builder
	sb.WriteString(entry.Name)
	sb.WriteString("\n")

	formatEntries(&sb, entry.Entries, "", opts)

	return sb.String()
}

func formatEntries(sb *strings.Builder, entries []Entry, prefix string, opts Options) {
	for i, entry := range entries {
		isLast := i == len(entries)-1

		// Choose connector
		connector := "├── "
		if isLast {
			connector = "└── "
		}

		sb.WriteString(prefix)
		sb.WriteString(connector)

		// Add icon
		if entry.IsDir {
			sb.WriteString("📁 ")
		} else {
			sb.WriteString(getFileIcon(entry.Name))
			sb.WriteString(" ")
		}

		sb.WriteString(entry.Name)

		if opts.ShowSize && !entry.IsDir {
			sb.WriteString(" (")
			sb.WriteString(formatSize(entry.Size))
			sb.WriteString(")")
		}

		sb.WriteString("\n")

		// Recurse for directories
		if entry.IsDir && len(entry.Entries) > 0 {
			newPrefix := prefix
			if isLast {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			formatEntries(sb, entry.Entries, newPrefix, opts)
		}
	}
}

func getFileIcon(name string) string {
	ext := strings.ToLower(filepath.Ext(name))

	icons := map[string]string{
		".go":   "🔵",
		".py":   "🐍",
		".js":   "🟨",
		".ts":   "🔷",
		".rs":   "🦀",
		".rb":   "💎",
		".java": "☕",
		".c":    "🔧",
		".cpp":  "🔧",
		".h":    "📋",
		".md":   "📝",
		".txt":  "📄",
		".json": "📋",
		".yaml": "📋",
		".yml":  "📋",
		".xml":  "📋",
		".html": "🌐",
		".css":  "🎨",
		".jpg":  "🖼️",
		".jpeg": "🖼️",
		".png":  "🖼️",
		".gif":  "🖼️",
		".svg":  "🖼️",
		".mp3":  "🎵",
		".wav":  "🎵",
		".mp4":  "🎬",
		".mov":  "🎬",
		".pdf":  "📕",
		".zip":  "📦",
		".tar":  "📦",
		".gz":   "📦",
		".sh":   "⚙️",
		".bash": "⚙️",
		".zsh":  "⚙️",
	}

	if icon, ok := icons[ext]; ok {
		return icon
	}
	return "📄"
}

func formatSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return strings.TrimRight(strings.TrimRight(
			strings.Replace(string(rune(bytes/GB)), ".", ",", 1), "0"), ",") + "G"
	case bytes >= MB:
		return strings.TrimRight(strings.TrimRight(
			strings.Replace(string(rune(bytes/MB)), ".", ",", 1), "0"), ",") + "M"
	case bytes >= KB:
		return strings.TrimRight(strings.TrimRight(
			strings.Replace(string(rune(bytes/KB)), ".", ",", 1), "0"), ",") + "K"
	default:
		return string(rune(bytes)) + "B"
	}
}

// Simple tree without options
func Simple(root string, maxDepth int) string {
	opts := DefaultOptions()
	opts.MaxDepth = maxDepth
	entry, err := Generate(root, opts)
	if err != nil {
		return err.Error()
	}
	return Format(entry, opts)
}
