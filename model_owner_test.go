package mlpack

import (
	"gonum.org/v1/gonum/mat"
	"math"
	"runtime"
	"testing"
)

func TestModelCloseAndBorrowedAlias(t *testing.T) {
	options := LinearRegressionOptions()
	options.Training = mat.NewDense(5, 1, []float64{0, 1, 2, 3, 4})
	options.TrainingResponses = mat.NewDense(5, 1, []float64{1, 3, 5, 7, 9})
	model, _ := LinearRegression(options)
	defer model.Close()
	prediction := LinearRegressionOptions()
	prediction.InputModel = &model
	prediction.Test = mat.NewDense(1, 1, []float64{5})
	alias, values := LinearRegression(prediction)
	runtime.GC()
	if values.At(0, 0) < 10.999 || values.At(0, 0) > 11.001 {
		t.Fatal("modèle emprunté libéré prématurément")
	}
	if alias.owner != model.owner {
		t.Fatal("deux propriétaires pour un même modèle")
	}
	if err := alias.Close(); err != nil {
		t.Fatal(err)
	}
	if err := model.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if recover() == nil {
			t.Fatal("modèle fermé réutilisable")
		}
	}()
	LinearRegression(prediction)
}

func TestLegacyCPPExceptionBecomesGoPanic(t *testing.T) {
	opts := KmeansOptions()
	opts.Algorithm = "algorithme-invalide"
	defer func() {
		if recover() == nil {
			t.Fatal("exception C++ non signalée")
		}
	}()
	Kmeans(2, mat.NewDense(3, 2, []float64{0, 0, 1, 1, 2, 2}), opts)
}

func TestBorrowedModelSurvivesCPPError(t *testing.T) {
	training := ApproxKfnOptions()
	training.NumTables, training.NumProjections = 1, 1
	training.Reference = mat.NewDense(12, 2, []float64{0, 0, 1, 2, 2, 1, 3, 4, 4, 3, 5, 6, 6, 5, 7, 8, 8, 7, 9, 10, 10, 9, 11, 12})
	_, _, model := ApproxKfn(training)
	defer model.Close()
	bad := ApproxKfnOptions()
	bad.InputModel, bad.K = &model, 1
	bad.Query = mat.NewDense(2, 2, []float64{0, 0, 1, 1})
	bad.CalculateError = true
	bad.ExactDistances = mat.NewDense(2, 2, []float64{1, 1, 1, 1})
	func() {
		defer func() {
			if recover() == nil {
				t.Fatal("erreur native attendue")
			}
		}()
		ApproxKfn(bad)
	}()
	good := ApproxKfnOptions()
	good.InputModel, good.K, good.Query = &model, 1, bad.Query
	_, _, alias := ApproxKfn(good)
	if alias.owner != model.owner {
		t.Fatal("modèle emprunté perdu après erreur")
	}
}

func TestModelOwnsTrainingData(t *testing.T) {
	values := []float64{0, 0, 0, 1, 10, 10, 10, 11, 20, 20, 20, 21}
	original := append([]float64(nil), values...)
	training := KnnOptions()
	training.Reference = mat.NewDense(6, 2, values)
	_, _, model := Knn(training)
	defer model.Close()
	for i := range values {
		if values[i] != original[i] {
			t.Fatal("construction du modèle modifiant l’entrée Go")
		}
		values[i] = 1e9
	}
	training.Reference = nil
	values = nil
	runtime.GC()
	query := KnnOptions()
	query.InputModel, query.K = &model, 1
	query.Query = mat.NewDense(1, 2, []float64{0.1, 0})
	distances, neighbors, alias := Knn(query)
	if math.Abs(distances.At(0, 0)-0.1) > 1e-12 || neighbors.At(0, 0) != 0 || alias.owner != model.owner {
		t.Fatal("modèle dépendant du stockage Go après retour")
	}
}
