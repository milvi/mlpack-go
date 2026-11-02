package main

import (
	"fmt"
	"math"
	"os"

	"github.com/milvi/mlpack-go"
	"gonum.org/v1/gonum/mat"
)

func main() {
	if err := verifier(); err != nil {
		fmt.Fprintln(os.Stderr, "ÉCHEC :", err)
		os.Exit(1)
	}
}

func verifier() error {
	// Chaque ligne est un point de ℝ³ : neuf points, répartis en trois groupes.
	valeurs := []float64{
		0, 0, 0, 0, 1, 0, 1, 0, 0,
		10, 10, 10, 10, 11, 10, 11, 10, 10,
		-10, -10, -10, -10, -11, -10, -11, -10, -10,
	}
	attendus := [][3]float64{
		{1.0 / 3, 1.0 / 3, 0},
		{10 + 1.0/3, 10 + 1.0/3, 10},
		{-10 - 1.0/3, -10 - 1.0/3, -10},
	}
	for _, algorithme := range []string{"naive", "dualtree"} {
		entree := mat.NewDense(9, 3, append([]float64(nil), valeurs...))
		options := mlpack.KmeansOptions()
		options.Algorithm = algorithme
		options.LabelsOnly = true // Évite de recopier les coordonnées dans la sortie.
		options.KillEmptyClusters = true
		options.MaxIterations = 100
		options.Seed = 42
		options.InitialCentroids = mat.NewDense(3, 3, []float64{
			0, 0, 0, 10, 10, 10, -10, -10, -10,
		})

		// La signature effective place le nombre de classes AVANT les données.
		centres, etiquettes := mlpack.Kmeans(3, entree, options)
		nc, dc := centres.Dims()
		ne, de := etiquettes.Dims()
		if nc != 3 || dc != 3 || ne != 9 || de != 1 {
			return fmt.Errorf("%s : dimensions inattendues, centres %d×%d, étiquettes %d×%d", algorithme, nc, dc, ne, de)
		}
		comptes := make([]int, 3)
		vus := make(map[int]bool)
		labels := make([]int, 9)
		for groupe := range 3 {
			premier := -1
			for i := groupe * 3; i < (groupe+1)*3; i++ {
				v := etiquettes.At(i, 0)
				if math.IsNaN(v) || math.IsInf(v, 0) || v < 0 || v >= 3 || v != math.Trunc(v) {
					return fmt.Errorf("%s : étiquette invalide au point %d : %v", algorithme, i, v)
				}
				label := int(v)
				labels[i] = label
				comptes[label]++
				if premier == -1 {
					premier = label
				} else if premier != label {
					return fmt.Errorf("%s : le groupe %d est scindé", algorithme, groupe)
				}
			}
			if vus[premier] {
				return fmt.Errorf("%s : deux groupes ont fusionné", algorithme)
			}
			vus[premier] = true
			for d := range 3 {
				v := centres.At(premier, d)
				if math.IsNaN(v) || math.IsInf(v, 0) || math.Abs(v-attendus[groupe][d]) > 1e-12 {
					return fmt.Errorf("%s : centre erroné, groupe %d, dimension %d : %.17g", algorithme, groupe, d, v)
				}
			}
		}
		for i, v := range valeurs {
			if entree.At(i/3, i%3) != v {
				return fmt.Errorf("%s : l'entrée a été modifiée à l'indice %d", algorithme, i)
			}
		}
		fmt.Printf("%s : OK — 9 points, 3 dimensions, 3 classes non vides\n", algorithme)
		fmt.Printf("étiquettes : %v ; effectifs : %v\n", labels, comptes)
	}
	return nil
}
