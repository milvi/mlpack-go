package mlpack

import (
	"math"
	"math/rand/v2"
	"reflect"
	"runtime"
	"sync"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestClusterMatchesNaive(t *testing.T) {
	for _, dimension := range []int{2, 3} {
		for seed := uint64(1); seed <= 5; seed++ {
			rng := rand.New(rand.NewPCG(seed, 7))
			points := make([]float64, 600*dimension)
			for i := range points {
				points[i] = rng.Float64()
			}
			original := append([]float64(nil), points...)
			opts := DefaultOptions()
			opts.Seed = seed
			dual, err := Cluster(points, dimension, 12, &opts)
			if err != nil {
				t.Fatal(err)
			}
			opts.Algorithm = "naive"
			naive, err := Cluster(points, dimension, 12, &opts)
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(dual.Inertia-naive.Inertia) > 1e-8*math.Max(1, naive.Inertia) {
				t.Fatalf("coût dual-tree %g contre naïf %g", dual.Inertia, naive.Inertia)
			}
			if !reflect.DeepEqual(points, original) {
				t.Fatal("entrée modifiée")
			}
			validateResult(t, points, dual)
		}
	}
}

func validateResult(t *testing.T, points []float64, result Result) {
	t.Helper()
	counts := make([]int, result.Clusters)
	sums := make([]float64, len(result.Centers))
	for i, label := range result.Labels {
		if label < 0 || label >= result.Clusters {
			t.Fatal("étiquette hors limites")
		}
		counts[label]++
		chosen := 0.0
		for d := range result.Dimensions {
			delta := points[i*result.Dimensions+d] - result.Centers[label*result.Dimensions+d]
			chosen += delta * delta
			sums[label*result.Dimensions+d] += points[i*result.Dimensions+d]
		}
		for c := range result.Clusters {
			distance := 0.0
			for d := range result.Dimensions {
				delta := points[i*result.Dimensions+d] - result.Centers[c*result.Dimensions+d]
				distance += delta * delta
			}
			if distance+1e-10 < chosen {
				t.Fatal("étiquette différente du centre le plus proche")
			}
		}
	}
	for c, n := range counts {
		if n == 0 {
			t.Fatal("classe vide conservée")
		}
		if result.Converged {
			for d := range result.Dimensions {
				if math.Abs(result.Centers[c*result.Dimensions+d]-sums[c*result.Dimensions+d]/float64(n)) > 1e-5 {
					t.Fatal("centre convergé différent de la moyenne des points affectés")
				}
			}
		}
	}
}

func TestClusterOwnershipAndDegeneracies(t *testing.T) {
	for _, tc := range []struct {
		points []float64
		k      int
	}{
		{[]float64{1, 2, 1, 2, 1, 2}, 3},
		{[]float64{0, 0, 2, 2, 4, 4}, 1},
		{[]float64{0, 0, 1, 1, 2, 2}, 3},
	} {
		result, err := Cluster(tc.points, 2, tc.k, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !result.Converged {
			t.Fatal("cas dégénéré sans convergence")
		}
		before := append([]float64(nil), result.Centers...)
		validateResult(t, tc.points, result)
		tc.points[0] = 900
		runtime.GC()
		if !reflect.DeepEqual(before, result.Centers) {
			t.Fatal("sortie empruntée")
		}
	}
}

func TestClusterErrors(t *testing.T) {
	for _, tc := range []struct {
		data []float64
		d, k int
	}{
		{nil, 2, 1}, {[]float64{1, 2}, 0, 1}, {[]float64{1, 2, 3}, 2, 1},
		{[]float64{1, 2}, 2, 2}, {[]float64{math.NaN(), 2}, 2, 1},
	} {
		if _, err := Cluster(tc.data, tc.d, tc.k, nil); err == nil {
			t.Fatal("entrée invalide acceptée")
		}
	}
	opts := DefaultOptions()
	opts.MaxIterations = -1
	if _, err := Cluster([]float64{0, 0}, 2, 1, &opts); err == nil {
		t.Fatal("itérations invalides acceptées")
	}
}

func TestClusterConcurrency(t *testing.T) {
	points := []float64{0, 0, 0, 1, 10, 10, 10, 11}
	reference, err := Cluster(points, 2, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for range 6 {
		wg.Go(func() {
			result, err := Cluster(points, 2, 2, nil)
			if err != nil {
				t.Error(err)
				return
			}
			if !reflect.DeepEqual(reference, result) {
				t.Error("appel dépendant d’un état global")
			}
		})
	}
	wg.Wait()
}

func TestClusterDenseStride(t *testing.T) {
	storage := mat.NewDense(4, 3, []float64{0, 0, 99, 0, 1, 99, 10, 10, 99, 10, 11, 99})
	view := storage.Slice(0, 4, 0, 2).(*mat.Dense)
	result, err := ClusterDense(view, 2, nil)
	if err != nil {
		t.Fatal(err)
	}
	validateResult(t, []float64{0, 0, 0, 1, 10, 10, 10, 11}, result)
}

func TestClusterOpenMPParity(t *testing.T) {
	points := make([]float64, 5000*2)
	rng := rand.New(rand.NewPCG(42, 7))
	for i := range points {
		points[i] = rng.Float64()
	}
	options := DefaultOptions()
	options.Algorithm = "naive"
	single, err := Cluster(points, 2, 15, &options)
	if err != nil {
		t.Fatal(err)
	}
	options.Threads = 2
	parallel, err := Cluster(points, 2, 15, &options)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(single.Inertia-parallel.Inertia) > 1e-8 {
		t.Fatal("coût OpenMP différent")
	}
	validateResult(t, points, parallel)
}

func TestLegacyDualTreeSingleCluster(t *testing.T) {
	options := KmeansOptions()
	options.Algorithm = "dualtree"
	options.LabelsOnly = true
	options.InitialCentroids = mat.NewDense(1, 2, []float64{0, 0})
	centers, labels := Kmeans(1, mat.NewDense(3, 2, []float64{0, 0, 2, 2, 4, 4}), options)
	if centers.At(0, 0) != 2 || centers.At(0, 1) != 2 {
		t.Fatal("moyenne incorrecte")
	}
	for i := range 3 {
		if labels.At(i, 0) != 0 {
			t.Fatal("étiquette incorrecte")
		}
	}
}

func TestClusterCPPFailureAndRecovery(t *testing.T) {
	options := DefaultOptions()
	options.Algorithm = "naive"
	options.InitialCentroids = []float64{0, 0, 0, 0}
	if _, err := Cluster([]float64{1e200, 0, -1e200, 0}, 2, 2, &options); err == nil {
		t.Fatal("débordement C++ non transmis à Go")
	}
	if _, err := Cluster([]float64{0, 0, 1, 1, 10, 10}, 2, 2, nil); err != nil {
		t.Fatalf("appel suivant l’exception impossible : %v", err)
	}
}
