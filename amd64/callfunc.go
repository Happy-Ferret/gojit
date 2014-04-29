package amd64

import (
	"reflect"
	"unsafe"
)

func (a *Assembler) CallFunc(f interface{}) {
	if reflect.TypeOf(f).Kind() != reflect.Func {
		panic("CallFunc: Can't call non-func")
	}
	ival := *(*struct {
		typ uintptr
		fun uintptr
	})(unsafe.Pointer(&f))

	// runtime·cgocallback_gofunc(f, frame, framesize)
	// plan9 ABI

	framesize := calculateFramesize(f)

	a.Sub(Imm{24}, Rsp)
	a.Mov(Imm{int32(framesize)}, Indirect{Rsp, 16, 64})
	a.Lea(Indirect{Rsp, 24, 64}, Rax)
	a.Mov(Rax, Indirect{Rsp, 8, 64})
	a.MovAbs(uint64(ival.fun), Rax)
	a.Mov(Rax, Indirect{Rsp, 0, 64})
	a.MovAbs(uint64(get_runtime_cgocallback_gofunc()), Rax)
	a.Call(Rax)
	a.Add(Imm{24}, Rsp)
}

func calculateFramesize(f interface{}) uintptr {
	t := reflect.TypeOf(f)
	s := uintptr(0)

	for i := 0; i < t.NumIn(); i++ {
		s += t.In(i).Size()
	}
	for i := 0; i < t.NumOut(); i++ {
		s += t.Out(i).Size()
	}
	return s
}

func get_runtime_cgocallback_gofunc() uintptr
