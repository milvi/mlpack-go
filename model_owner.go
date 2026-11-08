package mlpack

import (
	"fmt"
	"runtime"
	"unsafe"
	"weak"

	native "github.com/milvi/mlpack"
)

type modelState struct {
	ptr     unsafe.Pointer
	destroy func(unsafe.Pointer)
}
type modelOwner struct {
	state   *modelState
	cleanup runtime.Cleanup
}

// Ce registre faible ne retient pas les modèles après disparition de leurs utilisateurs.
// Tous ses accès se font sous le verrou commun native.Acquire.
var modelOwners = make(map[unsafe.Pointer]weak.Pointer[modelOwner])

func destroyModelLocked(state *modelState) {
	if state.ptr == nil {
		return
	}
	ptr := state.ptr
	state.ptr = nil
	delete(modelOwners, ptr)
	state.destroy(ptr)
}

func cleanupModel(state *modelState) {
	unlock := native.Acquire()
	defer unlock()
	destroyModelLocked(state)
}

func adoptModel(p *params, ptr unsafe.Pointer, destroy func(unsafe.Pointer)) *modelOwner {
	if ptr == nil {
		return nil
	}
	if existing := modelOwners[ptr].Value(); existing != nil {
		return existing
	}
	owner := &modelOwner{state: &modelState{ptr: ptr, destroy: destroy}}
	owner.cleanup = runtime.AddCleanup(owner, cleanupModel, owner.state)
	modelOwners[ptr] = weak.Make(owner)
	p.created = append(p.created, owner)
	return owner
}

func (owner *modelOwner) address() unsafe.Pointer {
	if owner == nil || owner.state.ptr == nil {
		panic(fmt.Errorf("mlpack : modèle absent ou fermé"))
	}
	return owner.state.ptr
}
func (owner *modelOwner) closeLocked() {
	if owner == nil {
		return
	}
	owner.cleanup.Stop()
	destroyModelLocked(owner.state)
}
func (owner *modelOwner) Close() error {
	unlock := native.Acquire()
	defer unlock()
	owner.closeLocked()
	return native.Error()
}
