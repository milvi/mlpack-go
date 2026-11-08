//go:build linux && mlpack_memtest

package mlpack

import (
	"gonum.org/v1/gonum/mat"
	"runtime"
	"testing"
)

func TestNativeMemoryCluster(t *testing.T) {
	data := make([]float64, 100000*2)
	for i := range 100000 {
		data[2*i] = float64(i%3) * 10
		data[2*i+1] = float64(i%3) * 10
	}
	o := DefaultOptions()
	o.InitialCentroids = []float64{0, 0, 10, 10, 20, 20}
	run := func() {
		if _, err := Cluster(data, 2, 3, &o); err != nil {
			t.Fatal(err)
		}
	}
	for range 4 {
		run()
	}
	runtime.GC()
	before := nativeBytes()
	for range 50 {
		run()
	}
	runtime.GC()
	delta := nativeBytes() - before
	t.Logf("Cluster : 50 appels × 100000 points, variation native %d octets", delta)
	if delta > 2<<20 {
		t.Fatalf("croissance native excessive : %d", delta)
	}
	runtime.KeepAlive(data)
}

func TestNativeMemoryModelClose(t *testing.T) {
	run := func() {
		o := LinearRegressionOptions()
		o.Training = mat.NewDense(5, 1, []float64{0, 1, 2, 3, 4})
		o.TrainingResponses = mat.NewDense(5, 1, []float64{1, 3, 5, 7, 9})
		m, _ := LinearRegression(o)
		if err := m.Close(); err != nil {
			t.Fatal(err)
		}
	}
	for range 5 {
		run()
	}
	runtime.GC()
	before := nativeBytes()
	for range 1000 {
		run()
	}
	runtime.GC()
	delta := nativeBytes() - before
	t.Logf("1000 modèles créés et fermés : variation native %d octets", delta)
	if delta > 2<<20 {
		t.Fatalf("modèles non restitués : %d", delta)
	}
}

func TestNativeMemoryFailedModel(t *testing.T) {
	const n = 10000
	data, labels := make([]float64, 2*n), make([]float64, n)
	for i := range n {
		data[2*i], data[2*i+1], labels[i] = float64(i), float64(i%31), float64(i%1000)
	}
	o := NbcOptions()
	o.Training, o.Labels = mat.NewDense(n, 2, data), mat.NewDense(n, 1, labels)
	o.Test = mat.NewDense(1, 3, []float64{0, 0, 0}) // erreur après construction du modèle
	run := func() {
		defer func() {
			if recover() == nil {
				t.Fatal("erreur native attendue")
			}
		}()
		Nbc(o)
	}
	for range 3 {
		run()
	}
	runtime.GC()
	before := nativeBytes()
	for range 100 {
		run()
	}
	runtime.GC()
	delta := nativeBytes() - before
	t.Logf("100 modèles dont l’appel échoue avant export : variation native %d octets", delta)
	if delta > 2<<20 {
		t.Fatalf("modèles perdus sur exception : %d octets", delta)
	}
	runtime.KeepAlive(o)
}
