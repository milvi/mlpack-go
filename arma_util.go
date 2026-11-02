package mlpack

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L${SRCDIR} -lmlpack_go_util
#include <capi/arma_util.h>
*/
import "C"

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"unsafe"
)

type mlpackArma struct{}

// Chaque ligne est un point ; Categoricals décrit les dimensions.
type matrixWithInfo struct {
	Categoricals []bool
	Data         *mat.Dense
}

func DataAndInfo() *matrixWithInfo { return &matrixWithInfo{} }

// Une vue Gonum avec un pas différent de sa largeur est compactée une fois.
// Le chemin contigu conserve son stockage initial, épinglé jusqu'à cleanParams.
func compactData(m *mat.Dense) (r, c int, data []float64) {
	if m == nil {
		panic("mlpack : matrice d'entrée absente")
	}
	r, c = m.Dims()
	if r == 0 || c == 0 {
		panic("mlpack : matrice d'entrée vide")
	}
	raw := m.RawMatrix()
	if raw.Stride == c || r == 1 {
		return r, c, raw.Data[:r*c]
	}
	data = make([]float64, r*c)
	for i := 0; i < r; i++ {
		copy(data[i*c:(i+1)*c], raw.Data[i*raw.Stride:i*raw.Stride+c])
	}
	return r, c, data
}

func gonumToArmaMat(p *params, identifier string, m *mat.Dense, trans bool) {
	r, c, data := compactData(m)
	p.pins.Pin(&data[0])
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackToArmaMat(p.mem, id, (*C.double)(unsafe.Pointer(&data[0])), C.size_t(c), C.size_t(r), C.bool(trans))
}
func gonumToArmaUmat(p *params, identifier string, m *mat.Dense) {
	r, c, data := compactData(m)
	id := C.CString(identifier)
	defer freeCString(id)
	// La conversion native vers size_t copie les données pendant cet appel.
	C.mlpackToArmaUmat(p.mem, id, (*C.double)(unsafe.Pointer(&data[0])), C.size_t(c), C.size_t(r))
}
func gonumToArmaVector(p *params, identifier string, m *mat.Dense, kind int) {
	r, c, data := compactData(m)
	if r != 1 && c != 1 {
		panic("mlpack : une seule ligne ou colonne est requise")
	}
	p.pins.Pin(&data[0])
	id := C.CString(identifier)
	defer freeCString(id)
	ptr := (*C.double)(unsafe.Pointer(&data[0]))
	n := C.size_t(len(data))
	switch kind {
	case 2:
		C.mlpackToArmaRow(p.mem, id, ptr, n)
	case 3:
		C.mlpackToArmaUrow(p.mem, id, ptr, n)
	case 4:
		C.mlpackToArmaCol(p.mem, id, ptr, n)
	case 5:
		C.mlpackToArmaUcol(p.mem, id, ptr, n)
	}
}
func gonumToArmaRow(p *params, id string, m *mat.Dense)  { gonumToArmaVector(p, id, m, 2) }
func gonumToArmaUrow(p *params, id string, m *mat.Dense) { gonumToArmaVector(p, id, m, 3) }
func gonumToArmaCol(p *params, id string, m *mat.Dense)  { gonumToArmaVector(p, id, m, 4) }
func gonumToArmaUcol(p *params, id string, m *mat.Dense) { gonumToArmaVector(p, id, m, 5) }

func gonumToArmaMatWithInfo(p *params, identifier string, m *matrixWithInfo) {
	if m == nil {
		panic("mlpack : données et informations absentes")
	}
	r, c, data := compactData(m.Data)
	if len(m.Categoricals) != c {
		panic("mlpack : une information catégorielle est requise par dimension")
	}
	p.pins.Pin(&data[0])
	id := C.CString(identifier)
	defer freeCString(id)
	C.mlpackToArmaMatWithInfo(p.mem, id, (*C.bool)(unsafe.Pointer(&m.Categoricals[0])),
		(*C.double)(unsafe.Pointer(&data[0])), C.size_t(c), C.size_t(r))
}

// La sortie appartient à Go ; le tampon natif reste la propriété de Params.
func copyArma(p *params, identifier string, kind int) (r, c int, data []float64) {
	id := C.CString(identifier)
	defer freeCString(id)
	var n int
	switch kind {
	case 0:
		r, c = int(C.mlpackNumColMat(p.mem, id)), int(C.mlpackNumRowMat(p.mem, id))
		n = int(C.mlpackNumElemMat(p.mem, id))
	case 1:
		r, c = int(C.mlpackNumColUmat(p.mem, id)), int(C.mlpackNumRowUmat(p.mem, id))
		n = int(C.mlpackNumElemUmat(p.mem, id))
	case 2:
		n = int(C.mlpackNumElemRow(p.mem, id))
		r, c = n, 1
	case 3:
		n = int(C.mlpackNumElemUrow(p.mem, id))
		r, c = n, 1
	case 4:
		n = int(C.mlpackNumElemCol(p.mem, id))
		r, c = 1, n
	case 5:
		n = int(C.mlpackNumElemUcol(p.mem, id))
		r, c = 1, n
	case 6:
		r, c = int(C.mlpackArmaMatWithInfoCols(p.mem, id)), int(C.mlpackArmaMatWithInfoRows(p.mem, id))
		n = int(C.mlpackArmaMatWithInfoElements(p.mem, id))
	default:
		panic("mlpack : type de conversion inconnu")
	}
	if n < 0 || r < 0 || c < 0 || r*c != n {
		panic("mlpack : dimensions natives incohérentes ou trop grandes")
	}
	data = make([]float64, n)
	if n != 0 && !bool(C.mlpackCopyParamMat(p.mem, id,
		(*C.double)(unsafe.Pointer(&data[0])), C.size_t(n), C.int(kind))) {
		panic(fmt.Sprintf("mlpack : échec de copie de %q", identifier))
	}
	return r, c, data
}
func denseArma(p *params, id string, kind int) *mat.Dense {
	r, c, data := copyArma(p, id, kind)
	// Compatibilité avec les anciennes sorties optionnelles absentes.
	if len(data) == 0 {
		return mat.NewDense(1, 1, nil)
	}
	return mat.NewDense(r, c, data)
}
func (*mlpackArma) armaToGonumMat(p *params, id string) *mat.Dense  { return denseArma(p, id, 0) }
func (*mlpackArma) armaToGonumUmat(p *params, id string) *mat.Dense { return denseArma(p, id, 1) }
func (*mlpackArma) armaToGonumRow(p *params, id string) *mat.Dense  { return denseArma(p, id, 2) }
func (*mlpackArma) armaToGonumUrow(p *params, id string) *mat.Dense { return denseArma(p, id, 3) }
func (*mlpackArma) armaToGonumCol(p *params, id string) *mat.Dense  { return denseArma(p, id, 4) }
func (*mlpackArma) armaToGonumUcol(p *params, id string) *mat.Dense { return denseArma(p, id, 5) }
func (*mlpackArma) armaToGonumMatWithInfo(p *params, id string) *mat.Dense {
	return denseArma(p, id, 6)
}
func (*mlpackArma) armaToGonumArray(p *params, id string) (int, int, []float64) {
	return copyArma(p, id, 0)
}
