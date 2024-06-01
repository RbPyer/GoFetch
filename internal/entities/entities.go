package entities

const (
	Delimeter="\n____________________________________________________________\n"
	Prefix="PRETTY_NAME="
)

type Response struct {
	Info []string
}

func NewResponse() *Response {
	return &Response{}
}