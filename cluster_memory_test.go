//go:build linux && mlpack_memtest

package mlpack

import (
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

// Le calcul lève après les allocations natives ; l’appel suivant doit rester utilisable.
func TestNativeMemoryFailure(t *testing.T) {
	options := DefaultOptions()
	options.Algorithm = "naive"
	options.InitialCentroids = []float64{0, 0, 0, 0}
	run := func() {
		if _, err := Cluster([]float64{1e200, 0, -1e200, 0}, 2, 2, &options); err == nil {
			t.Fatal("débordement natif attendu")
		}
	}
	for range 4 {
		run()
	}
	runtime.GC()
	before := nativeBytes()
	for range 200 {
		run()
	}
	runtime.GC()
	delta := nativeBytes() - before
	t.Logf("200 calculs interrompus : variation native %d octets", delta)
	if delta > 2<<20 {
		t.Fatalf("croissance native après erreur : %d", delta)
	}
	if _, err := Cluster([]float64{0, 0, 1, 1, 10, 10}, 2, 2, nil); err != nil {
		t.Fatalf("appel suivant l’erreur impossible : %v", err)
	}
}
