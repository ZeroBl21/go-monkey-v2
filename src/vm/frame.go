package vm

import (
	"github.com/ZeroBl21/go-monkey/src/code"
	"github.com/ZeroBl21/go-monkey/src/object"
)

type Frame struct {
	fn *object.CompiledFuction
	ip int
}

func NewFrame(fn *object.CompiledFuction) *Frame {
	return &Frame{
		fn: fn,
		ip: -1,
	}
}

func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
