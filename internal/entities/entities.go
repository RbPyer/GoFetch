package entities

const (
	Delimeter="\n______________________________________________________________________\n"
	Prefix="PRETTY_NAME="
	CPU_PATH="/proc/cpuinfo"
	RAM_PATH="/proc/meminfo"
)

// var (
// 	TotalMemErr = errors.New("")
// )

type Response struct {
	Info []string
}

type RAM struct {
	TrueFree uint64
	Total uint64
	Free uint64
	Shared uint64
	SReclaimable uint64
	Buffers uint64
	Cached uint64
}

type CPU struct {
	ModelName string
	Cores int
	Siblings int
}

type MemoryInfo struct {
	Total uint64
	Free uint64
	Buffers uint64
	Cache uint64
	Shared uint64
	SReclaimable uint64	
}


func NewResponse() *Response {
	return &Response{}
}