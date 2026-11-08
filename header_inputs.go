package mlpack

import "embed"

// Les déclarations C imbriquées participent également à la clé du cache Go.
//
//go:embed capi/*.h
var _ embed.FS
