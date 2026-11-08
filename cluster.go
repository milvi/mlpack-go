package mlpack

import (
	"fmt"
	"math/rand/v2"

	native "github.com/milvi/mlpack"
	"gonum.org/v1/gonum/mat"
)

// Options contrôle un appel indépendant de k-means. Une ligne représente un point.
type Options struct {
	Algorithm         string
	Seed              uint64
	MaxIterations     int
	Tolerance         float64
	Threads           int
	InitialCentroids  []float64
	KeepEmptyClusters bool
}

// Result contient exclusivement des slices Go ; aucun Close n’est nécessaire.
type Result = native.Result

// DefaultOptions choisit dual-tree avec KD-tree, une graine fixe et un thread.
func DefaultOptions() Options {
	return Options{Algorithm: "dualtree", Seed: 42, MaxIterations: 1000, Tolerance: 1e-5, Threads: 1}
}

// Cluster traite des points contigus n × dimensions. Les entrées restent inchangées.
// Les centres vides sont conservés pendant les itérations, puis supprimés à la fin.
func Cluster(points []float64, dimensions, clusters int, options *Options) (Result, error) {
	if dimensions <= 0 || len(points) == 0 || len(points)%dimensions != 0 ||
		clusters <= 0 || clusters > len(points)/dimensions {
		return Result{}, fmt.Errorf("mlpack : dimensions ou nombre de classes invalides")
	}
	o := DefaultOptions()
	if options != nil {
		o = *options
	}
	if o.Algorithm == "" {
		o.Algorithm = "dualtree"
	}
	if o.MaxIterations == 0 {
		o.MaxIterations = 1000
	}
	if o.Tolerance == 0 {
		o.Tolerance = 1e-5
	}
	if o.Threads == 0 {
		o.Threads = 1
	}
	algorithm := 0
	switch o.Algorithm {
	case "dualtree":
	case "naive":
		algorithm = 1
	default:
		return Result{}, fmt.Errorf("mlpack : algorithme inconnu %q", o.Algorithm)
	}
	initial := o.InitialCentroids
	if initial == nil {
		// Échantillonnage sans remise, local à cet appel : aucun générateur global.
		random := rand.New(rand.NewPCG(o.Seed, o.Seed^0x9e3779b97f4a7c15))
		initial = make([]float64, clusters*dimensions)
		selected := make(map[int]int, clusters)
		n := len(points) / dimensions
		for i := 0; i < clusters; i++ {
			j := i + random.IntN(n-i)
			value, ok := selected[j]
			if !ok {
				value = j
			}
			replacement, ok := selected[i]
			if !ok {
				replacement = i
			}
			selected[j] = replacement
			copy(initial[i*dimensions:(i+1)*dimensions], points[value*dimensions:(value+1)*dimensions])
		}
	}
	return native.KMeans(points, dimensions, clusters, initial, o.MaxIterations, o.Tolerance, algorithm, o.Threads, o.KeepEmptyClusters)
}

// ClusterDense accepte aussi les vues Gonum non contiguës.
func ClusterDense(points *mat.Dense, clusters int, options *Options) (Result, error) {
	if points == nil {
		return Result{}, fmt.Errorf("mlpack : matrice absente")
	}
	r, c := points.Dims()
	if r == 0 || c == 0 {
		return Result{}, fmt.Errorf("mlpack : matrice vide")
	}
	raw := points.RawMatrix()
	data := raw.Data
	if raw.Stride != c && r != 1 {
		data = make([]float64, r*c)
		for i := 0; i < r; i++ {
			copy(data[i*c:(i+1)*c], raw.Data[i*raw.Stride:i*raw.Stride+c])
		}
	} else {
		data = data[:r*c]
	}
	return Cluster(data, c, clusters, options)
}
