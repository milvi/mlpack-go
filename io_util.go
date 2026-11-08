package mlpack

/*
#cgo CFLAGS: -I.
#include <capi/io_util.h>
#include <stdlib.h>
*/
import "C"

import (
	"unsafe"
)

type params struct {
	mem      unsafe.Pointer
	refs     []*modelOwner
	created  []*modelOwner
	complete bool
}
type timers struct{ mem unsafe.Pointer }

// Chaque CString a un propriétaire local et une libération explicite.
func freeCString(s *C.char) { C.free(unsafe.Pointer(s)) }
func getParams(binding string) *params {
	defer checkNative()
	name := C.CString(binding)
	defer freeCString(name)
	return &params{mem: C.mlpackGetParams(name)}
}
func getTimers() *timers {
	defer checkNative()
	return &timers{mem: C.mlpackGetTimers()}
}
func cleanParams(p *params) {
	defer checkNative()
	if p == nil || p.mem == nil {
		return
	}
	// Les paramètres natifs possèdent leurs données et modèles non exportés.
	C.mlpackCleanParams(p.mem)
	p.mem = nil
	if !p.complete {
		for _, owner := range p.created {
			owner.closeLocked()
		}
	}
	p.refs = nil
	p.created = nil
}
func cleanTimers(t *timers) {
	defer checkNative()
	if t == nil || t.mem == nil {
		return
	}
	C.mlpackCleanTimers(t.mem)
	t.mem = nil
}
func setParamDouble(p *params, identifier string, value float64) {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamDouble(p.mem, id, C.double(value))
}
func setParamInt(p *params, identifier string, value int) {
	if value < -2147483648 || value > 2147483647 {
		panic("mlpack : entier hors de la plage C int")
	}
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamInt(p.mem, id, C.int(value))
}
func setParamFloat(p *params, identifier string, value float64) {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamFloat(p.mem, id, C.float(value))
}
func setParamBool(p *params, identifier string, value bool) {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamBool(p.mem, id, C.bool(value))
}
func hasParam(p *params, identifier string) bool {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	return bool(C.mlpackHasParam(p.mem, id))
}
func setPassed(p *params, identifier string) {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetPassed(p.mem, id)
}
func setParamString(p *params, identifier, value string) {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	val := C.CString(value)
	defer freeCString(val)
	C.mlpackSetParamString(p.mem, id, val)
}
func enableTimers() {
	defer checkNative()
	C.mlpackEnableTimers()
}
func disableBacktrace() {
	defer checkNative()
	C.mlpackDisableBacktrace()
}
func disableVerbose() {
	defer checkNative()
	C.mlpackDisableVerbose()
}
func enableVerbose() {
	defer checkNative()
	C.mlpackEnableVerbose()
}
func getParamString(p *params, identifier string) string {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	return C.GoString(C.mlpackGetParamString(p.mem, id))
}
func getParamBool(p *params, identifier string) bool {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	return bool(C.mlpackGetParamBool(p.mem, id))
}
func getParamInt(p *params, identifier string) int {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	return int(C.mlpackGetParamInt(p.mem, id))
}
func getParamDouble(p *params, identifier string) float64 {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	return float64(C.mlpackGetParamDouble(p.mem, id))
}
func setParamVecInt(p *params, identifier string, values []int) {
	defer checkNative()
	data := make([]C.longlong, len(values))
	for i, v := range values {
		data[i] = C.longlong(v)
	}
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamVectorInt(p.mem, id, unsafe.SliceData(data), C.size_t(len(data)))
}
func setParamVecString(p *params, identifier string, values []string) {
	defer checkNative()
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackSetParamVectorStrLen(p.mem, id, C.size_t(len(values)))
	checkNative()
	for i, value := range values {
		s := C.CString(value)
		C.mlpackSetParamVectorStr(p.mem, id, s, C.size_t(i))
		freeCString(s)
	}
}
func getParamVecInt(p *params, identifier string) []int {
	defer checkNative()
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
	defer checkNative()
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
