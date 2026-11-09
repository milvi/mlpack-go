// Commande de vérification et de mesure : un appel k-means sur des points de [0,1]².
package main

import (
	"encoding/json"
	"errors"
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
	DurationLimitSeconds float64 `json:"duration_limit_seconds"`
	TimedOut             bool    `json:"timed_out"`
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
	maxDuration := flag.Duration("max-duration", 20*time.Second, "durée maximale du calcul (ex. 20s, 500ms)")
	algorithm := flag.String("algorithm", "dualtree", "dualtree ou naive")
	jsonOutput := flag.Bool("json", false, "sortie JSON pour les scripts (code 124 si le budget est atteint)")
	flag.Parse()
	if *maxDuration <= 0 || *n < 1 || *n > int(^uint(0)>>1)/2 {
		fmt.Fprintln(os.Stderr, "Le nombre de points et la durée maximale doivent être positifs et représentables.")
		os.Exit(2)
	}
	data := make([]float64, *n*2)
	rng := rand.New(rand.NewPCG(*seed, 17))
	for i := range data {
		data[i] = rng.Float64()
	}
	options := mlpack.DefaultOptions()
	options.Seed, options.Threads, options.MaxDuration, options.Algorithm = *seed, *threads, *maxDuration, *algorithm
	runtime.GC()
	start := time.Now()
	result, err := mlpack.Cluster(data, 2, *k, &options)
	elapsed := time.Since(start)
	timedOut := errors.Is(err, mlpack.ErrTimeLimit)
	if err != nil && !timedOut {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	counts := make([]int, result.Clusters)
	for i, label := range result.Labels {
		if label < 0 || label >= len(counts) {
			fmt.Fprintf(os.Stderr, "Résultat incohérent : classe %d invalide au point %d.\n", label, i)
			os.Exit(1)
		}
		counts[label]++
	}
	for _, count := range counts {
		if count == 0 {
			fmt.Fprintln(os.Stderr, "Résultat incohérent : une classe vide est présente.")
			os.Exit(1)
		}
	}
	runtime.GC()
	var memory runtime.MemStats
	runtime.ReadMemStats(&memory)
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		fmt.Fprintln(os.Stderr, "Impossible de mesurer la mémoire du processus :", err)
		os.Exit(1)
	}
	measurement := Measurement{
		N: *n, K: *k, KActual: result.Clusters, Seed: *seed,
		Algorithm: *algorithm, Threads: *threads, DurationLimitSeconds: maxDuration.Seconds(),
		TimedOut: timedOut, Iterations: result.Iterations, Converged: result.Converged,
		Seconds: elapsed.Seconds(), PeakRSSBytes: usage.Maxrss * 1024, HeapAllocBytes: memory.HeapAlloc,
		Inertia: result.Inertia, DistanceCalculations: result.DistanceCalculations,
	}
	var outputError error
	if *jsonOutput {
		outputError = json.NewEncoder(os.Stdout).Encode(measurement)
	} else {
		outputError = writeReport(os.Stdout, measurement)
	}
	if outputError != nil {
		fmt.Fprintln(os.Stderr, "Impossible d’afficher le résultat :", outputError)
		os.Exit(1)
	}
	runtime.KeepAlive(data)
	runtime.KeepAlive(result)
	if timedOut && *jsonOutput {
		os.Exit(124)
	}
}
