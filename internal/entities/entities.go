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


type CPU struct {
	ModelName string
	Cores int
	Siblings int
}

func NewResponse() *Response {
	return &Response{}
}