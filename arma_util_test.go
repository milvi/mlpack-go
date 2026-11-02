package mlpack

import (
	"runtime"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestMatrixOutputOwnsItsData(t *testing.T) {
	p := getParams("kmeans")
	defer cleanParams(p)
	values := make([]float64, 30)
	for i := range values {
		values[i] = float64(i) + 0.25
	}
	input := mat.NewDense(10, 3, values)
	gonumToArmaMat(p, "input", input, false)
	runtime.GC()
	var conv mlpackArma
	output := conv.armaToGonumMat(p, "input")
	cleanParams(p)
	input.Set(0, 0, -100)
	runtime.GC()
	if output.At(0, 0) != 0.25 || output.At(9, 2) != 29.25 {
		t.Fatal("la sortie ne possède pas une copie indépendante")
	}
}

func TestStridedMatrixInput(t *testing.T) {
	backing := mat.NewDense(6, 5, nil)
	for i := range 6 {
		for j := range 5 {
			backing.Set(i, j, float64(i*10+j))
		}
	}
	view := backing.Slice(0, 6, 0, 3).(*mat.Dense)
	p := getParams("kmeans")
	defer cleanParams(p)
	gonumToArmaMat(p, "input", view, false)
	runtime.GC()
	var conv mlpackArma
	output := conv.armaToGonumMat(p, "input")
	if !mat.Equal(output, view) {
		t.Fatal("le pas de la vue Gonum a été ignoré")
	}
}

func TestCommonConversions(t *testing.T) {
	tests := []struct {
		name, binding, id string
		set               func(*params, string, *mat.Dense)
		kind              int
	}{
		{"mat", "kmeans", "input", func(p *params, s string, m *mat.Dense) { gonumToArmaMat(p, s, m, false) }, 0},
		{"umat", "knn", "true_neighbors", gonumToArmaUmat, 1},
		{"row", "linear_regression", "training_responses", gonumToArmaRow, 2},
		{"urow", "nbc", "labels", gonumToArmaUrow, 3},
		{"col", "kde", "predictions", gonumToArmaCol, 4},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := getParams(tc.binding)
			defer cleanParams(p)
			input := mat.NewDense(6, 1, []float64{0, 1, 2, 3, 4, 5})
			tc.set(p, tc.id, input)
			_, _, output := copyArma(p, tc.id, tc.kind)
			cleanParams(p)
			for i, v := range output {
				if v != float64(i) {
					t.Fatalf("indice %d : %v", i, v)
				}
			}
		})
	}
}

func TestCategoricalMatrixCopy(t *testing.T) {
	p := getParams("preprocess_one_hot_encoding")
	defer cleanParams(p)
	input := mat.NewDense(3, 2, []float64{0, 1, 1, 2, 0, 3})
	gonumToArmaMatWithInfo(p, "input", &matrixWithInfo{Categoricals: []bool{true, false}, Data: input})
	out := denseArma(p, "input", 6)
	cleanParams(p)
	if !mat.Equal(input, out) {
		t.Fatal("la matrice catégorielle n'a pas été copiée")
	}
}

func TestVectorOutputs(t *testing.T) {
	p := getParams("preprocess_one_hot_encoding")
	defer cleanParams(p)
	for _, values := range [][]int{{}, {0, 2, 5}} {
		setParamVecInt(p, "dimensions", values)
		got := getParamVecInt(p, "dimensions")
		if len(got) != len(values) {
			t.Fatal("taille entière erronée")
		}
		for i, v := range values {
			if got[i] != v {
				t.Fatal("conversion entière erronée")
			}
		}
	}
	strings := getParams("image_converter")
	defer cleanParams(strings)
	values := []string{"été", "🙂"}
	setParamVecString(strings, "input", values)
	got := getParamVecString(strings, "input")
	cleanParams(strings)
	for i, v := range values {
		if got[i] != v {
			t.Fatal("chaîne incorrecte")
		}
	}
}

func TestKmeansReturnsIndependentLabels(t *testing.T) {
	input := mat.NewDense(4, 2, []float64{0, 0, 0, 1, 10, 10, 10, 11})
	for _, algorithm := range []string{"naive", "dualtree"} {
		t.Run(algorithm, func(t *testing.T) {
			opts := KmeansOptions()
			opts.Algorithm = algorithm
			opts.LabelsOnly = true
			opts.Seed = 42
			opts.InitialCentroids = mat.NewDense(2, 2, []float64{0, 0, 10, 10})
			centres, labels := Kmeans(2, input, opts)
			runtime.GC()
			if r, c := centres.Dims(); r != 2 || c != 2 {
				t.Fatal("dimensions des centres incorrectes")
			}
			if labels.At(0, 0) != labels.At(1, 0) || labels.At(2, 0) != labels.At(3, 0) || labels.At(0, 0) == labels.At(2, 0) {
				t.Fatal("partition incorrecte")
			}
		})
	}
}
