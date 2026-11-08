package mlpack

import native "github.com/milvi/mlpack"

// Les fonctions historiques gardent leur signature et signalent une erreur par panic.
// La nouvelle API Cluster renvoie directement une valeur error.
func checkNative() {
	if err := native.Error(); err != nil {
		panic(err)
	}
}
