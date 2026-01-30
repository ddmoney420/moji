// Package sysinfo provides system information collection and formatting.
//
// It gathers system details including OS, hostname, kernel, uptime, CPU, memory, shell, and terminal.
// Information can be formatted as plain text or displayed with OS logos.
//
// Example usage:
//
//	info := sysinfo.Collect()
//	output := sysinfo.Format(info)
//	output := sysinfo.FormatWithArt(info)
//	sysinfo.Print()
package sysinfo
