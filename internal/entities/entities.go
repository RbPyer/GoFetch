package entities

const (
	Delimeter="\n______________________________________________________________________\n"
	Prefix="PRETTY_NAME="
	CPU_PATH="/proc/cpuinfo"
	RAM_PATH="/proc/meminfo"
)

type Response struct {
	Info []string
}

type RAM struct {
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

func NewResponse() *Response {
	return &Response{}
}