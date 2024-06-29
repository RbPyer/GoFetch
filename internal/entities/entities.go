package entities

const (
	Delimeter = "\n______________________________________________________________________\n"
	Prefix    = "PRETTY_NAME="
	CpuPath   = "/proc/cpuinfo"
	RamPath   = "/proc/meminfo"
	B         = 1
	KB        = 1024 * B
	MB        = 1024 * KB
	GB        = 1024 * MB
)

//var (
//	NoTempInfo = errors.New("no information about CPU temperature")
//)

type Response struct {
	Info []string
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

type MemoryInfo struct {
	Total        uint64
	Free         uint64
	Buffers      uint64
	Cache        uint64
	Shared       uint64
	SReclaimable uint64
}

func NewResponse() *Response {
	return &Response{}
}
