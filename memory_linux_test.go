//go:build linux && mlpack_memtest

package mlpack

import (
	"runtime"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestNativeMemoryAfterKmeans(t *testing.T) {
	const n = 100_000
	data := make([]float64, n*3)
	for i := range n {
		for j := range 3 {
			data[3*i+j] = float64(i%3) * 10
		}
	}
	input := mat.NewDense(n, 3, data)
	options := KmeansOptions()
	options.LabelsOnly, options.KillEmptyClusters = true, true
	options.MaxIterations, options.Seed = 2, 42
	options.InitialCentroids = mat.NewDense(3, 3, []float64{0, 0, 0, 10, 10, 10, 20, 20, 20})
	run := func() { _, _ = Kmeans(3, input, options) }
	for range 5 {
		run()
	}
	runtime.GC()
	before := nativeBytes()
	for range 50 {
		run()
	}
	runtime.GC()
	after := nativeBytes()
	t.Logf("50 appels × %d points : variation native %d octets", n, after-before)
	// Le défaut historique laisse plus de 40 Mo d'étiquettes natives.
	if delta := after - before; delta > 2<<20 {
		t.Fatalf("croissance native excessive : %d octets", delta)
	}
	runtime.KeepAlive(input)
	runtime.KeepAlive(options)
}

func TestNativeCStringLifetime(t *testing.T) {
	run := func() {
		p := getParams("kmeans")
		defer cleanParams(p)
		setParamString(p, "algorithm", "naive")
		setPassed(p, "algorithm")
		_ = getParamString(p, "algorithm")
		_ = hasParam(p, "algorithm")
	}
	for range 20 {
		run()
	}
	runtime.GC()
	before := nativeBytes()
	for range 50_000 {
		run()
	}
	runtime.GC()
	delta := nativeBytes() - before
	t.Logf("50 000 cycles de paramètres : variation native %d octets", delta)
	if delta > 2<<20 {
		t.Fatalf("fuite de chaînes C : %d octets", delta)
	}
}
