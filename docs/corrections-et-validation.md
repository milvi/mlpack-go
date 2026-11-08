# Réduction à k-means et validation du 8 novembre 2026

Branche `codex/kmeans-only` des deux dépôts ; préversion `v0.2.0-kmeans.1`.
Le module C++ désigne le commit
[`fa1b738da2211f05215e7668c56877b20190d939`](https://github.com/milvi/mlpack/commit/fa1b738da2211f05215e7668c56877b20190d939).
Les sources conservées dérivent de mlpack C++ 4.8.0.

## Périmètre du source

| Élément | Avant | Sur cette branche |
|---|---|---|
| Unités C++ compilées | 49 | 1 : `native_kmeans_api.cpp` |
| Méthodes exposées | Ensemble des bindings mlpack | K-means : `Cluster`, adaptateur `ClusterDense` |
| En-têtes internes conservés | Ensemble de mlpack | 55 fichiers, 349 333 octets |
| Dépendances C++ de développement | Armadillo, ensmallen, cereal | Armadillo |
| Exécutable de mesure local | 25 997 776 octets | 4 129 464 octets |

Les autres méthodes, autres langages, modèles persistants, générateurs et
infrastructures CMake sont retirés. La recherche de voisins conservée est une
dépendance interne du dual-tree et de l’affectation finale exacte. Le mode
k-means naïf reste disponible pour la comparaison numérique, dans la même fonction.
La sérialisation inutilisée est retirée des classes conservées, ce qui supprime
cereal. La réduction des inclusions supprime aussi ensmallen.

Le contrôle [`tools/check-source.py`](https://github.com/milvi/mlpack/blob/codex/kmeans-only/tools/check-source.py)
vérifie la fermeture exacte des inclusions en modes normal et debug : aucun
fichier de `src` ne reste hors des dépendances du calcul et du pilote mémoire.
Les changements de contenu de ces en-têtes invalident le cache Go.

## Corrections mémoire conservées

Le C++ copie les entrées avant tout réordonnancement, écrit dans les slices Go
et détruit ses ressources avant le retour. Il ne conserve aucun pointeur Go.
Les sorties sont indépendantes et ne demandent ni `Close` ni libération manuelle.
L’interface simplifiée n’alloue plus de chaînes C ni de paramètres/modèles opaques.

Les propriétaires uniques des arbres, le nettoyage des constructeurs interrompus,
les copies et déplacements sûrs des arbres et des bornes et la restauration des
enfants masqués sur exception restent présents. Le k-means naïf conserve la
préallocation de ses accumulateurs OpenMP et la collecte des exceptions dans la
région parallèle. La frontière C intercepte les exceptions ; Go rend `error`.
Les réglages partagés sont protégés et restaurés au retour.

Le [pilote natif](https://github.com/milvi/mlpack/blob/codex/kmeans-only/tests/native_memory.cpp)
intercepte **390 échecs d’allocation k-means et 639 échecs d’arbres**, avec retour
exact au nombre initial de `new` actifs après chaque scénario. Le [journal local](validation/asan.log)
est fourni. AddressSanitizer contrôle les accès mémoire ; UndefinedBehaviorSanitizer
contrôle les comportements indéfinis. Leur exécution locale a réussi. LeakSanitizer,
qui détecte les allocations perdues, nécessite les permissions du conteneur Ubuntu.

Les [tests de restitution](validation/memoire.log), après échauffement,
ramasse-miettes Go et nettoyage de l’allocateur natif, donnent :

| Cas | Variation native |
|---|---:|
| 50 appels de 100 000 points | +6 816 octets |
| 200 calculs interrompus par une erreur C++ | 0 octet |

Le seuil de non-régression est de 2 Mio. Ces variations incluent les caches de
l’allocateur et ne constituent pas une estimation d’un débit de fuite.

## Calcul et cache vérifiés

Les [tests Go](validation/tests-go.log) et les [contrôles cgo renforcés](validation/tests-cgocheck2.log)
passent : parité dual-tree/naïf, dimensions 2 et 3, entrées inchangées, sorties
indépendantes, doublons, un centre, k = n, suppression des centres vides, vues
Gonum, concurrence, OpenMP et récupération après exception. L’[exemple](validation/exemple.log)
retrouve trois groupes de trois points en deux itérations.

Une [comparaison avant/après extraction](validation/comparaison.json) sur
100 000 points, 1 000 centres, graine 42 et budget de 100 itérations donne
exactement les mêmes nombres de classes et d’itérations, indicateur de convergence,
inertie et compteur de distances. Elle compare les deux exécutables sur les mêmes données.

| Action sur le cache | Compilations C++ | Temps local |
|---|---:|---:|
| Sans modification | 0 | 0,075 s |
| Horodatage seul | 0 | 0,096 s |
| Commentaire dans un en-tête imbriqué | 1 | 6,36 s |
| Contenu restauré | 0 | 0,52 s |

Le [rapport cache](validation/cache.json) et l’[outil](../tools/verify-cache.py)
vérifient aussi que le témoin de source n’est pas présent dans l’exécutable.
Les bibliothèques et en-têtes APT restent extérieurs à cette clé : après leur
modification, forcer `go build -a ./...`.

## Mesures locales de contrôle

Debian 13.6, Go 1.27.3, GCC 14.2, Armadillo 14.2.3 ; paquets extraits dans un
sysroot sans installation sur l’hôte. Un thread, points uniformes dans le carré
unité, graine 42, **un seul essai par taille**, budget de 100 itérations.
Le temps couvre l’appel entier ; la génération des points est exclue. Le pic
RSS (*Resident Set Size*, mémoire résidente du processus) inclut Go et C++.

| n | k | Temps (s) | Pic RSS (Mio) |
|---:|---:|---:|---:|
| 100 000 | 1 000 | 1,639 | 20,89 |
| 1 000 000 | 1 000 | 8,824 | 160,92 |
| 1 000 000 | 10 000 | 20,448 | 161,87 |

Pour n multiplié par 10 : temps ×5,38 à k fixé, ×12,47 avec k = n/100 ;
mémoire ×7,70 et ×7,75. Les trois essais atteignent le budget sans convergence.
Ces mesures comparent un coût à budget fixé et ne prouvent pas une complexité
asymptotique. Les [données brutes](validation/mesures.json) et les
[conditions locales](validation/local.json) sont fournies.

## Validation Ubuntu et distribution

La [validation native Ubuntu 26.04](https://github.com/milvi/mlpack/actions/runs/42468656957)
a réussi : 55 en-têtes, une seule unité, tests Go, 390 et 639 échecs d’allocation
interceptés et contrôles mémoire, y compris LeakSanitizer. Le tag natif a été
créé après succès. Go a ensuite téléchargé le [module publié](validation/module-natif.json)
et passé [tous les tests du frontend](validation/tests-module-publie.log) avec
`GOWORK=off`, sans remplacement local.

Les workflows [natif](https://github.com/milvi/mlpack/actions/runs/42468656957) et
[Go](../.github/workflows/ubuntu2604.yml) construisent dans Ubuntu 26.04 avec
uniquement `build-essential` et `libarmadillo-dev` comme dépendances C++ de
développement. Les tests s’exécutent sous un utilisateur ordinaire ; le contrôle
LeakSanitizer du pilote natif ajoute la capacité `SYS_PTRACE` dans son conteneur.
Les préversions ne sont créées qu’après succès de l’intégration de leur dépôt.

La publication du frontend est suivie d’un contrôle depuis un projet neuf
pour vérifier les deux modules téléchargés, leurs empreintes et l’absence de `replace`.
Les contrôles portent sur les chemins k-means conservés et leurs utilitaires ;
ils ne constituent pas une preuve pour toute entrée numérique imaginable.
