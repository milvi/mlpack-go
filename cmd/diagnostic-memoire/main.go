package main

/*
#include <malloc.h>
#include <stddef.h>
static size_t octets_natifs(void) {
    struct mallinfo2 m = mallinfo2();
    return m.uordblks + m.hblkhd;
}
static void purger_arene(void) { malloc_trim(0); }
*/
import "C"

import (
	"fmt"
	"runtime"

	"github.com/milvi/mlpack-go"
	"gonum.org/v1/gonum/mat"
)

// Diagnostic spécifique à glibc : observation des allocations natives après
// abandon des sorties Go. Ce fichier ne remplace pas un analyseur de mémoire.
func main() {
	const n = 100_000
	const repetitions = 50
	valeurs := make([]float64, n*3)
	for i := range n {
		g := float64(i % 3)
		valeurs[i*3], valeurs[i*3+1], valeurs[i*3+2] = g*10, g*10, g*10
	}
	donnees := mat.NewDense(n, 3, valeurs)
	options := mlpack.KmeansOptions()
	options.Algorithm = "naive"
	options.LabelsOnly = true
	options.KillEmptyClusters = true
	options.MaxIterations = 2
	options.Seed = 42
	options.InitialCentroids = mat.NewDense(3, 3, []float64{0, 0, 0, 10, 10, 10, 20, 20, 20})
	appel := func() { _, _ = mlpack.Kmeans(3, donnees, options) }
	appel() // Chauffe les bibliothèques avant de mesurer.
	runtime.GC()
	C.purger_arene()
	avant := uint64(C.octets_natifs())
	for range repetitions {
		appel()
	}
	runtime.GC()
	C.purger_arene()
	apres := uint64(C.octets_natifs())
	runtime.KeepAlive(donnees)
	runtime.KeepAlive(options)
	fmt.Printf("avant : %d octets natifs\naprès : %d\ndifférence : %d\n", avant, apres, int64(apres)-int64(avant))
	fmt.Printf("étiquettes seules, 50 allocations de 100 000 float64 : %d octets\n", repetitions*n*8)
}
