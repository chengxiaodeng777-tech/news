package system

import (
	"syscall"
	"unsafe"
)

func readHostInfo() hostMetrics {
	var info hostMetrics

	kernel := syscall.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx := kernel.NewProc("GlobalMemoryStatusEx")

	type memoryStatusEx struct {
		Length               uint32
		MemoryLoad           uint32
		TotalPhys            uint64
		AvailPhys            uint64
		TotalPageFile        uint64
		AvailPageFile        uint64
		TotalVirtual         uint64
		AvailVirtual         uint64
		AvailExtendedVirtual uint64
	}

	mem := memoryStatusEx{Length: uint32(unsafe.Sizeof(memoryStatusEx{}))}
	ret, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mem)))
	if ret != 0 {
		info.CPULoad = float64(mem.MemoryLoad)
		info.TotalMemoryMB = mem.TotalPhys / 1024 / 1024
		info.FreeMemoryMB = mem.AvailPhys / 1024 / 1024
	}

	var freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64
	diskRet, _, _ := syscall.NewLazyDLL("kernel32.dll").NewProc("GetDiskFreeSpaceExW").Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:\\"))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)),
	)
	if diskRet != 0 {
		info.TotalDiskGB = totalNumberOfBytes / 1024 / 1024 / 1024
		info.FreeDiskGB = totalNumberOfFreeBytes / 1024 / 1024 / 1024
	}

	return info
}
