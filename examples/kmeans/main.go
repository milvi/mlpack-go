// Exemple autonome : un appel et des résultats possédés par Go.
package main

import (
	"fmt"
	"os"

	mlpack "github.com/milvi/mlpack-go"
)

func main() {
	// Neuf points de ℝ³, répartis en trois groupes.
	points := []float64{
		0, 0, 0, 0, 1, 0, 1, 0, 0,
		10, 10, 10, 10, 11, 10, 11, 10, 10,
		-10, -10, -10, -10, -11, -10, -11, -10, -10,
	}
	options := mlpack.DefaultOptions()
	options.InitialCentroids = []float64{0, 0, 0, 10, 10, 10, -10, -10, -10}
	result, err := mlpack.Cluster(points, 3, 3, &options)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Classes : %d ; itérations : %d ; convergence : %t\n", result.Clusters, result.Iterations, result.Converged)
	fmt.Printf("Centres : %v\nÉtiquettes : %v\n", result.Centers, result.Labels)
}
