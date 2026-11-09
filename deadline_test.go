package mlpack

import (
	"errors"
	"gonum.org/v1/gonum/mat"
	"math/rand/v2"
	"testing"
	"time"
)

func TestDefaultDuration(t *testing.T) {
	if DefaultOptions().MaxDuration != 20*time.Second {
		t.Fatal("budget par défaut différent de 20 s")
	}
}

func TestDurationStopsInsideNativeWork(t *testing.T) {
	points := make([]float64, 100000*2)
	rng := rand.New(rand.NewPCG(42, 17))
	for i := range points {
		points[i] = rng.Float64()
	}
	for _, algorithm := range []string{"dualtree", "naive"} {
		for _, threads := range []int{1, 2} {
			options := DefaultOptions()
			options.Algorithm, options.Threads = algorithm, threads
			options.MaxDuration = 20 * time.Millisecond
			started := time.Now()
			result, err := Cluster(points, 2, 1000, &options)
			elapsed := time.Since(started)
			if !errors.Is(err, ErrTimeLimit) {
				t.Fatalf("%s : expiration attendue, obtenu %v", algorithm, err)
			}
			if result.Centers != nil || result.Labels != nil || result.Converged {
				t.Fatal("sorties partielles exposées")
			}
			if elapsed > time.Second {
				t.Fatalf("%s : arrêt tardif : %s", algorithm, elapsed)
			}
			t.Logf("%s, %d threads : arrêt en %s, %d itérations terminées", algorithm, threads, elapsed, result.Iterations)
		}
	}
	if _, err := Cluster([]float64{0, 0, 1, 1, 10, 10}, 2, 2, nil); err != nil {
		t.Fatal(err)
	}
}

func TestTinyDurationAndDense(t *testing.T) {
	o := DefaultOptions()
	o.MaxDuration = time.Nanosecond
	if _, err := Cluster([]float64{0, 0, 1, 1}, 2, 1, &o); !errors.Is(err, ErrTimeLimit) {
		t.Fatal(err)
	}
	storage := mat.NewDense(4, 3, nil)
	view := storage.Slice(0, 4, 0, 2).(*mat.Dense)
	if _, err := ClusterDense(view, 1, &o); !errors.Is(err, ErrTimeLimit) {
		t.Fatal(err)
	}
}
