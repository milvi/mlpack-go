// Package mlpack expose uniquement k-means, avec dual-tree et KD-tree par défaut.
// Cluster rend des slices Go indépendantes et une erreur. Aucun objet natif ne
// survit au calcul ; aucun Close n’est nécessaire. cgo compile le C++ automatiquement.
package mlpack // import "github.com/milvi/mlpack-go"
