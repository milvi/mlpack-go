# Budget temporel de k-means — 9 novembre 2026

La préversion `v0.2.0-kmeans.2` de la branche `codex/kmeans-only` remplace le
plafond d’itérations par `Options.MaxDuration`, de type `time.Duration`.
La valeur par défaut et la valeur zéro choisissent **20 secondes** ; une durée
négative est refusée. La commande utilise `-max-duration 20s`.
`MaxIterations` et le drapeau `-iterations` sont retirés. `Result.Iterations`
compte uniquement les itérations terminées et ne constitue plus un critère d’arrêt.

## Contrat d’expiration

Le budget couvre la préparation des données dans l’API Go, l’attente du verrou
natif et le calcul C++, y compris construction des arbres et affectation finale.
Pour `ClusterDense`, la copie éventuelle d’une vue Gonum fait partie du budget.
La génération aléatoire des points par le CLI et la compilation par `go run`
sont extérieures au calcul chronométré.

Le C++ emploie une horloge monotone et des contrôles amortis dans les boucles,
les distances et les parcours. Les workers OpenMP reçoivent le même budget.
Les exceptions restent interceptées à l’intérieur de la région parallèle,
puis relancées sur le thread appelant après la jonction des workers.

À expiration : `errors.Is(err, mlpack.ErrTimeLimit)` est vrai, les statistiques
diagnostiques sont rendues, `Converged` est faux et les slices `Centers` et
`Labels` sont nil. Une itération interrompue n’est pas comptée. Aucun résultat
partiel n’est présenté comme une partition complète et aucun calcul natif ne
continue après le retour. Le CLI émet son JSON avec `timed_out: true` et le
binaire termine avec le code 124.

L’arrêt est **coopératif**, pas une garantie temps réel : une opération mémoire
indivisible, l’ordonnancement du système et la destruction des ressources peuvent
faire légèrement dépasser le budget. Le nettoyage n’est jamais interrompu,
notamment la restauration des enfants masqués des arbres dual-tree.

## Contrôles exécutés

Le [rapport local](validation/duree.json) enregistre :

- les suites Go, les contrôles de pointeurs renforcés `cgocheck2`, la parité
  dual-tree/naïf et l’exécution après une erreur ;
- des expirations réelles sur dual-tree et naïf avec un et deux threads :
  un budget de 20 ms revient en environ 20–22 ms sur le laboratoire ;
- une attente de verrou bornée et le refus de durées négatives ;
- 52 expirations injectées à des contrôles successifs, y compris après des
  itérations terminées, avec retour exact au nombre initial de `new` actifs ;
- 390 pannes d’allocation k-means et 639 pannes de construction/affectation
  d’arbres, contrôlées avec AddressSanitizer et UndefinedBehaviorSanitizer ;
- 50 expirations répétées : variation native de +6 896 octets, sous le seuil
  de non-régression de 2 Mio ; les caches de l’allocateur expliquent que cette
  mesure ne soit pas une preuve de fuite strictement nulle.

Le test réel du budget par défaut sur 1 000 000 de points et 10 000 centres
termine avec `ErrTimeLimit` en **20,054 s**, après **50 itérations terminées**.
Les données sont uniformes dans le carré unité, graine 42 et un thread.
Il ne s’agit pas d’un résultat convergé.

## Code et distribution

Le [module natif](https://github.com/milvi/mlpack/tree/codex/kmeans-only) conserve
une unité C++ et 56 en-têtes nécessaires, dont le contrôle de durée.
Les copies des entrées, les résultats appartenant à Go, les propriétaires uniques
des arbres et le nettoyage des constructeurs interrompus restent présents.
Aucune nouvelle dépendance APT n’est introduite : `build-essential` et
`libarmadillo-dev` suffisent.

Le tag natif est créé après [validation Ubuntu](https://github.com/milvi/mlpack/actions/runs/42656799239).
Le frontend dépend de cette version exacte et [son workflow](../.github/workflows/ubuntu2604.yml)
vérifie les sommes des modules, les tests ordinaires, `cgocheck2`, les restitutions,
l’exemple et les mesures. Les opérations de test et de compilation sont bornées
à deux minutes. Les résultats de la préversion précédente restent consultables
[dans son rapport versionné](https://github.com/milvi/mlpack-go/blob/v0.2.0-kmeans.1/docs/corrections-et-validation.md) ;
leurs budgets de 100 itérations ne décrivent pas le présent contrat temporel.
