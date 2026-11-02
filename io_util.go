package mlpack

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L${SRCDIR} -Wl,-rpath,${SRCDIR} -lmlpack_go_util
#include <capi/io_util.h>
#include <stdlib.h>
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type params struct {
	mem  unsafe.Pointer
	pins runtime.Pinner
}
type timers struct{ mem unsafe.Pointer }

// Chaque CString a un propriétaire local et une libération explicite.
func freeCString(s *C.char) { C.free(unsafe.Pointer(s)) }
func getParams(binding string) *params {
	name := C.CString(binding)
	defer freeCString(name)
	return &params{mem: C.mlpackGetParams(name)}
}
func getTimers() *timers { return &timers{mem: C.mlpackGetTimers()} }
func cleanParams(p *params) {
	if p == nil || p.mem == nil {
		return
	}
	// Les vues natives doivent disparaître avant de désépingler leurs données Go.
	C.mlpackCleanParams(p.mem)
	p.mem = nil
	p.pins.Unpin()
}
func cleanTimers(t *timers) {
	if t == nil || t.mem == nil {
		return
	}
	C.mlpackCleanTimers(t.mem)
	t.mem = nil
}
func setParamDouble(p *params, identifier string, value float64) {
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamDouble(p.mem, id, C.double(value))
}
func setParamInt(p *params, identifier string, value int) {
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamInt(p.mem, id, C.int(value))
}
func setParamFloat(p *params, identifier string, value float64) {
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamFloat(p.mem, id, C.float(value))
}
func setParamBool(p *params, identifier string, value bool) {
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamBool(p.mem, id, C.bool(value))
}
func hasParam(p *params, identifier string) bool {
	id := C.CString(identifier)
	defer freeCString(id)
	return bool(C.mlpackHasParam(p.mem, id))
}
func setPassed(p *params, identifier string) {
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetPassed(p.mem, id)
}
func setParamString(p *params, identifier, value string) {
	id := C.CString(identifier)
	defer freeCString(id)
	val := C.CString(value)
	defer freeCString(val)
	C.mlpackSetParamString(p.mem, id, val)
}
func setParamPtr(p *params, identifier string, ptr unsafe.Pointer) {
	if ptr != nil {
		p.pins.Pin(ptr)
	}
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamPtr(p.mem, id, (*C.double)(ptr))
}
func enableTimers()     { C.mlpackEnableTimers() }
func disableBacktrace() { C.mlpackDisableBacktrace() }
func disableVerbose()   { C.mlpackDisableVerbose() }
func enableVerbose()    { C.mlpackEnableVerbose() }
func getParamString(p *params, identifier string) string {
	id := C.CString(identifier)
	defer freeCString(id)
	return C.GoString(C.mlpackGetParamString(p.mem, id))
}
func getParamBool(p *params, identifier string) bool {
	id := C.CString(identifier)
	defer freeCString(id)
	return bool(C.mlpackGetParamBool(p.mem, id))
}
func getParamInt(p *params, identifier string) int {
	id := C.CString(identifier)
	defer freeCString(id)
	return int(C.mlpackGetParamInt(p.mem, id))
}
func getParamDouble(p *params, identifier string) float64 {
	id := C.CString(identifier)
	defer freeCString(id)
	return float64(C.mlpackGetParamDouble(p.mem, id))
}
func setParamVecInt(p *params, identifier string, values []int) {
	data := make([]C.longlong, len(values))
	for i, v := range values {
		data[i] = C.longlong(v)
	}
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamVectorInt(p.mem, id, unsafe.SliceData(data), C.size_t(len(data)))
}
func setParamVecString(p *params, identifier string, values []string) {
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamVectorStrLen(p.mem, id, C.size_t(len(values)))
	for i, value := range values {
		s := C.CString(value)
		C.mlpackSetParamVectorStr(p.mem, id, s, C.size_t(i))
		freeCString(s)
	}
}
func getParamVecInt(p *params, identifier string) []int {
	id := C.CString(identifier)
	defer freeCString(id)
	n := int(C.mlpackVecIntSize(p.mem, id))
	if n < 0 {
		panic("mlpack : taille de vecteur invalide")
	}
	data := make([]C.longlong, n)
	if !bool(C.mlpackCopyParamVectorInt(p.mem, id, unsafe.SliceData(data), C.size_t(n))) {
		panic("mlpack : échec de copie du vecteur entier")
	}
	out := make([]int, n)
	for i, v := range data {
		out[i] = int(v)
	}
	return out
}
func getParamVecString(p *params, identifier string) []string {
	id := C.CString(identifier)
	defer freeCString(id)
	n := int(C.mlpackVecStringSize(p.mem, id))
	if n < 0 {
		panic("mlpack : taille de vecteur invalide")
	}
	out := make([]string, n)
	for i := range out {
		out[i] = C.GoString(C.mlpackGetVecStringPtr(p.mem, id, C.size_t(i)))
	}
	return out
}
