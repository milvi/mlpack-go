// Commande de vérification et de mesure : un appel k-means sur des points de [0,1]².
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand/v2"
	"os"
	"runtime"
	"syscall"
	"time"

	mlpack "github.com/milvi/mlpack-go"
)

type Measurement struct {
	N                    int     `json:"n"`
	K                    int     `json:"k_requested"`
	KActual              int     `json:"k_actual"`
	Seed                 uint64  `json:"seed"`
	Algorithm            string  `json:"algorithm"`
	Threads              int     `json:"threads"`
	IterationLimit       int     `json:"iteration_limit"`
	Iterations           int     `json:"iterations"`
	Converged            bool    `json:"converged"`
	Seconds              float64 `json:"seconds"`
	PeakRSSBytes         int64   `json:"peak_rss_bytes"`
	HeapAllocBytes       uint64  `json:"go_heap_alloc_bytes"`
	Inertia              float64 `json:"inertia"`
	DistanceCalculations uint64  `json:"distance_calculations"`
}

func main() {
	n := flag.Int("n", 100000, "nombre de points")
	k := flag.Int("k", 1000, "nombre maximal de classes")
	seed := flag.Uint64("seed", 42, "graine reproductible")
	threads := flag.Int("threads", 1, "threads OpenMP")
	iterations := flag.Int("iterations", 1000, "limite d’itérations")
	algorithm := flag.String("algorithm", "dualtree", "dualtree ou naive")
	flag.Parse()
	if *n < 1 || *n > int(^uint(0)>>1)/2 {
		fmt.Fprintln(os.Stderr, "n invalide")
		os.Exit(2)
	}
	data := make([]float64, *n*2)
	rng := rand.New(rand.NewPCG(*seed, 17))
	for i := range data {
		data[i] = rng.Float64()
	}
	options := mlpack.DefaultOptions()
	options.Seed, options.Threads, options.MaxIterations, options.Algorithm = *seed, *threads, *iterations, *algorithm
	runtime.GC()
	start := time.Now()
	result, err := mlpack.Cluster(data, 2, *k, &options)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	counts := make([]int, result.Clusters)
	for _, label := range result.Labels {
		if label < 0 || label >= len(counts) {
			panic("étiquette invalide")
		}
		counts[label]++
	}
	for _, count := range counts {
		if count == 0 {
			panic("classe vide")
		}
	}
	runtime.GC()
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		panic(err)
	}
	measurement := Measurement{*n, *k, result.Clusters, *seed, *algorithm, *threads, *iterations, result.Iterations, result.Converged, elapsed.Seconds(), usage.Maxrss * 1024, memory.HeapAlloc, result.Inertia, result.DistanceCalculations}
	if err := json.NewEncoder(os.Stdout).Encode(measurement); err != nil {
		panic(err)
	}
	runtime.KeepAlive(data)
	runtime.KeepAlive(result)
}
