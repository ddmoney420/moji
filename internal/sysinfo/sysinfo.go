package sysinfo

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Info holds system information
type Info struct {
	OS           string
	Hostname     string
	Kernel       string
	Uptime       string
	Shell        string
	Terminal     string
	CPU          string
	Memory       string
	User         string
	GoVersion    string
	Architecture string
	HomeDir      string
}

// Collect gathers system information
func Collect() Info {
	info := Info{
		OS:           runtime.GOOS,
		Architecture: runtime.GOARCH,
		GoVersion:    runtime.Version(),
	}

	// Hostname
	if h, err := os.Hostname(); err == nil {
		info.Hostname = h
	}

	// User
	if u := os.Getenv("USER"); u != "" {
		info.User = u
	} else if u := os.Getenv("USERNAME"); u != "" {
		info.User = u
	}

	// Home directory
	if h, err := os.UserHomeDir(); err == nil {
		info.HomeDir = h
	}

	// Shell
	if s := os.Getenv("SHELL"); s != "" {
		info.Shell = s
	}

	// Terminal
	if t := os.Getenv("TERM"); t != "" {
		info.Terminal = t
	}
	if tp := os.Getenv("TERM_PROGRAM"); tp != "" {
		info.Terminal = tp
	}

	// OS-specific info
	switch runtime.GOOS {
	case "darwin":
		info.collectMacOS()
	case "linux":
		info.collectLinux()
	case "windows":
		info.collectWindows()
	}

	return info
}

func (i *Info) collectMacOS() {
	// Get macOS version
	if out, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
		i.OS = "macOS " + strings.TrimSpace(string(out))
	}

	// Get kernel
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		i.Kernel = strings.TrimSpace(string(out))
	}

	// Get CPU
	if out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
		i.CPU = strings.TrimSpace(string(out))
	}

	// Get memory
	if out, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
		var bytes int64
		fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &bytes)
		i.Memory = formatBytes(bytes)
	}

	// Get uptime
	if out, err := exec.Command("sysctl", "-n", "kern.boottime").Output(); err == nil {
		// Parse boot time
		bootStr := strings.TrimSpace(string(out))
		// Format: { sec = 1234567890, usec = 0 }
		var sec int64
		fmt.Sscanf(bootStr, "{ sec = %d", &sec)
		if sec > 0 {
			bootTime := time.Unix(sec, 0)
			i.Uptime = formatDuration(time.Since(bootTime))
		}
	}
}

func (i *Info) collectLinux() {
	// Get distro info
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				i.OS = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				break
			}
		}
	}

	// Get kernel
	if out, err := exec.Command("uname", "-r").Output(); err == nil {
		i.Kernel = strings.TrimSpace(string(out))
	}

	// Get CPU
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "model name") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					i.CPU = strings.TrimSpace(parts[1])
					break
				}
			}
		}
	}

	// Get memory
	if data, err := os.ReadFile("/proc/meminfo"); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					var kb int64
					fmt.Sscanf(parts[1], "%d", &kb)
					i.Memory = formatBytes(kb * 1024)
				}
				break
			}
		}
	}

	// Get uptime
	if data, err := os.ReadFile("/proc/uptime"); err == nil {
		var seconds float64
		fmt.Sscanf(string(data), "%f", &seconds)
		i.Uptime = formatDuration(time.Duration(seconds) * time.Second)
	}
}

func (i *Info) collectWindows() {
	// Basic Windows info
	i.OS = "Windows " + runtime.GOARCH

	// Get computer name
	if out, err := exec.Command("hostname").Output(); err == nil {
		i.Hostname = strings.TrimSpace(string(out))
	}
}

func formatBytes(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	var parts []string
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d days", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d hours", hours))
	}
	if minutes > 0 && days == 0 {
		parts = append(parts, fmt.Sprintf("%d mins", minutes))
	}

	if len(parts) == 0 {
		return "just now"
	}
	return strings.Join(parts, ", ")
}

// Format returns formatted system info
func (i Info) Format() string {
	var lines []string

	add := func(label, value string) {
		if value != "" {
			lines = append(lines, fmt.Sprintf("\033[1;36m%s:\033[0m %s", label, value))
		}
	}

	add("User", i.User+"@"+i.Hostname)
	add("OS", i.OS)
	add("Kernel", i.Kernel)
	add("Uptime", i.Uptime)
	add("Shell", i.Shell)
	add("Terminal", i.Terminal)
	add("CPU", i.CPU)
	add("Memory", i.Memory)
	add("Arch", i.Architecture)
	add("Go", i.GoVersion)

	return strings.Join(lines, "\n")
}

// FormatWithArt returns system info alongside ASCII art
func (i Info) FormatWithArt(art string) string {
	infoLines := strings.Split(i.Format(), "\n")
	artLines := strings.Split(art, "\n")

	// Find max art width
	maxArtWidth := 0
	for _, line := range artLines {
		w := len([]rune(line))
		if w > maxArtWidth {
			maxArtWidth = w
		}
	}

	// Pad and combine
	var result strings.Builder
	maxLines := len(artLines)
	if len(infoLines) > maxLines {
		maxLines = len(infoLines)
	}

	for i := 0; i < maxLines; i++ {
		artLine := ""
		if i < len(artLines) {
			artLine = artLines[i]
		}

		infoLine := ""
		if i < len(infoLines) {
			infoLine = infoLines[i]
		}

		// Pad art line
		artRunes := []rune(artLine)
		padding := maxArtWidth - len(artRunes) + 4

		result.WriteString(artLine)
		result.WriteString(strings.Repeat(" ", padding))
		result.WriteString(infoLine)
		result.WriteString("\n")
	}

	return result.String()
}

// Logos for different operating systems
var OSLogos = map[string]string{
	"darwin": `
       .:'
    _ :'_
 .'` + "`" + `_` + "`" + `-'` + "`" + `_` + "`" + `'.
:________.-'
:_______:
 :_______:
  ` + "`" + `:_____:'
`,
	"linux": `
    .--.
   |o_o |
   |:_/ |
  //   \ \
 (|     | )
/'\_   _/` + "`" + `\
\___)=(___/
`,
	"windows": `
 _____
|     |
|_____|
|     |
|_____|
`,
	"default": `
  ___
 /   \
|  O  |
 \___/
`,
}

// GetOSLogo returns ASCII logo for current OS
func GetOSLogo() string {
	if logo, ok := OSLogos[runtime.GOOS]; ok {
		return logo
	}
	return OSLogos["default"]
}
