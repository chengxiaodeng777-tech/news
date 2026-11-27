package system

import (
	"runtime"
	"syscall"
)

func readHostInfo() hostMetrics {
	var info hostMetrics
	var memory syscall.Sysinfo_t
	if err := syscall.Sysinfo(&memory); err == nil {
		unit := uint64(memory.Unit)
		info.TotalMemoryMB = uint64(memory.Totalram) * unit / 1024 / 1024
		info.FreeMemoryMB = uint64(memory.Freeram) * unit / 1024 / 1024
		// Linux reports the one-minute load average as a fixed-point value.
		info.CPULoad = float64(memory.Loads[0]) / 65536 / float64(runtime.NumCPU()) * 100
	}
	var disk syscall.Statfs_t
	if err := syscall.Statfs(".", &disk); err == nil {
		info.TotalDiskGB = disk.Blocks * uint64(disk.Bsize) / 1024 / 1024 / 1024
		info.FreeDiskGB = disk.Bavail * uint64(disk.Bsize) / 1024 / 1024 / 1024
	}
	return info
}
