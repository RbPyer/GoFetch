package models

const (
	Template = "\033[32m%s@%s\nOS: %s\nKernel: %s\nUptime: %s\nRAM: %d/%d MiB [%.2f%%]\n\nCPU Info: %s; %d cores / %d threads\nTemperature zones: %v\n\nDisk Info: %.2f/%.2f GiB [%.2f%%]\nGPU: %s\nShell: %s"
	Prefix    = "PRETTY_NAME="
	CpuPath   = "/proc/cpuinfo"
	RamPath   = "/proc/meminfo"
	GpuPath   = "/sys/bus/pci/devices"
	B         = 1
	KB        = 1024 * B
	MB        = 1024 * KB
	GB        = 1024 * MB
)

func New() Response {
	return Response{CPU: CPU{Temperatures: make([]uint64, 0, 5)}}
}

type Response struct {
	Hostname string
	Username string
	OSRelease string
	KernelVersion string
	Uptime string
	GPUModel string
	Shell string
	RAM
	DiskInfo
	CPU
}

type RAM struct {
	TrueFree     uint64
	Total        uint64
	Free         uint64
	Shared       uint64
	SReclaimable uint64
	Buffers      uint64
	Cached       uint64
}

type DiskInfo struct {
	All  uint64
	Used uint64
}

type CPU struct {
	ModelName    string
	Cores        int
	Siblings     int
	Temperatures []uint64
}
