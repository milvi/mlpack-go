# Corrections du binding Go et validation

Travail du 2 novembre 2026. Ce document décrit le **fork corrigé**.
[L'audit amont](audit-amont.md) conserve les constats et mesures avant correction.

## Organisation des deux dépôts

| Dépôt | Responsabilité |
|---|---|
| `milvi/mlpack` | Interface C native, utilitaires Go communs, générateur de wrappers et de modèles |
| `milvi/mlpack-go` | Distribution Go corrigée, module `github.com/milvi/mlpack-go`, installateur, exemple et tests |

La base C++ du fork est `9767815cb73a924300535b09c5d08a3a8cc33f57`.
La base de la distribution Go est `3505bb06392e345e1a15bb641cfd1d20974cb172`.
La correction C++ correspond au commit
[`d391b0bc6701aa1a5ebdfb057e00375a30526b2c`](https://github.com/milvi/mlpack/commit/d391b0bc6701aa1a5ebdfb057e00375a30526b2c),
proposé dans [milvi/mlpack#1](https://github.com/milvi/mlpack/pull/1).
Pour l'installation, la source native est l'archive officielle **4.8.0**, dont
les fichiers concernés acceptent le même correctif. Cela évite de construire
une branche C++ mouvante à chaque installation.

Le patch `patches/mlpack-4.8.0-go-memory.patch` contient les dix fichiers de code
modifiés dans le dépôt C++ ; il ne contient pas cette documentation. Son SHA-256 :

```text
e2f1a1111a05d6313e47921ed5015ec29e8fc3a5a09cb4e2fbc7d38991c99a9d
```

Le script contrôle cette empreinte et vérifie l'application complète du patch.
Une source déjà entièrement corrigée est reconnue par une vérification inverse
sans écriture. Une source partiellement corrigée ou incompatible est refusée.
Une empreinte fixe garantit le contenu sélectionné ; elle ne constitue pas une
signature indépendante de l'éditeur.

## Comment le binding fonctionne

L'application Go appelle les wrappers `.go` du module, qui utilisent **cgo**.
Les déclarations de `capi/*.h` exposent des fonctions à ABI C. Les implémentations
de ces fonctions, compilées en C++, appellent les algorithmes mlpack et Armadillo.
L'exécutable Go est lié aux `libmlpack_go_*.so` ; celles-ci restent nécessaires
à l'exécution et utilisent les bibliothèques mathématiques d'Ubuntu.

Le module Go contient donc les wrappers et les déclarations C, **pas les sources
complètes C++ ni des bibliothèques précompilées**. `go get` seul n'effectue pas
l'installation native. Un programme C++ direct inclut les en-têtes mlpack et lie
ses dépendances mathématiques ; ce chemin ne produit pas automatiquement les
wrappers Go.

`BUILD_GO_SHLIB=ON` construit les bibliothèques natives. `BUILD_GO_BINDINGS=ON`
active aussi la génération des fichiers Go. L'installateur conserve cette
dernière option à `OFF`, car il utilise les wrappers déjà présents dans ce
checkout. Le message CMake « Not building Go bindings » désigne la génération,
pas l'absence de construction des `.so`.

Toutes les directives cgo du package Go sont agrégées, même si l'application ne
fait appel qu'à k-means. La distribution actuelle demande donc 47 bibliothèques
(46 méthodes et l'utilitaire commun), pas seulement celle de k-means.

### Entrées

Une matrice Gonum utilise une ligne par observation et un stockage par lignes.
Le binding présente ce même tampon à Armadillo avec les dimensions inversées :
un point devient une colonne du côté C++. Une vue Gonum dont `Stride` dépasse
la largeur doit d'abord être compactée ; le chemin corrigé effectue cette copie.

Les tampons Go que C++ conserve entre plusieurs appels cgo sont épinglés avec
`runtime.Pinner`. Leur épinglage dure jusqu'à la destruction de `Params` natif.
Cette durée ne peut pas être remplacée par un simple `runtime.KeepAlive`.
Les conversions vers des matrices d'entiers effectuent leur conversion native
pendant l'appel. Les entrées vides ou absentes sont refusées explicitement par
les convertisseurs Go au lieu de déréférencer `&data[0]`.

### Sorties et nettoyage

L'ancien chemin utilisait `GetMemory` pour détacher le stockage C++ et le montrer
à Go comme une slice, sans mécanisme de restitution. Le nouveau chemin alloue
la sortie en Go et appelle `mlpackCopyParamMat` pour y copier les valeurs.
Les matrices natives restent la propriété de `Params` et sont détruites avec lui.
Le vecteur d'entiers utilise de même `mlpackCopyParamVectorInt`.

Chaque wrapper généré programme `defer cleanParams(params)` et
`defer cleanTimers(timers)` dès leur création. Toutes les chaînes temporaires
`C.CString` des utilitaires et des accesseurs de modèles sont libérées.
Le nettoyage des paramètres natifs précède le désépinglage de leurs entrées Go.
Il reste une copie de sortie de coût O(nd) ; il n'y a pas de promesse de coût nul
ni de mémoire maximale indépendante de la taille des résultats.

Les anciens symboles C de transfert de pointeurs sont conservés pour compatibilité
ABI. Les nouveaux wrappers ne les appellent plus. Installer uniquement les
nouvelles `.so` ne corrige donc pas un ancien wrapper Go. Inversement, utiliser
les nouveaux wrappers avec des `.so` amont ne fonctionne pas : les deux nouveaux
symboles de copie y manquent. Le script installe un ensemble cohérent.

## Vérifications réalisées dans le laboratoire

Les commandes ordinaires sont dans le README. Les journaux de tests sont
conservés dans [validation/](validation/). Les résultats suivants ne sont pas
des résultats supposés de GitHub Actions.

| Vérification | Résultat observé |
|---|---|
| Application du patch à une extraction neuve de 4.8.0 | Vérification puis application réussies |
| Utilitaires natifs corrigés | Compilation réussie avec le GCC d'Ubuntu 26.04 |
| Générateur réel `generate_go_kmeans` | Compilation et exécution réussies ; `kmeans.go` strictement identique au wrapper corrigé distribué |
| Génération des accesseurs de modèles | 52 occurrences de `C.CString`, chacune associée à `defer freeCString(id)` |
| Copies de matrices | Sortie indépendante après destruction de `Params` et collecte Go ; stockage supérieur au petit tampon interne d'Armadillo testé |
| Vues Gonum | Sous-matrice avec `Stride > nombre de colonnes` conservée correctement |
| Types de conversion exercés | Matrice double, matrice entière, ligne double, ligne entière, colonne double, matrice avec informations catégorielles |
| Vecteurs | Entiers vides et non vides, chaînes UTF-8 |
| k-means | `naive` et `dualtree` ; résultats indépendants du stockage natif |
| `go test -buildvcs=false -tags mlpack_memtest -v ./...` | Réussi |
| `GOEXPERIMENT=cgocheck2 go test -buildvcs=false ./...` | Réussi |
| Script complet du fork | Installation dans un nouveau préfixe, construction Go et validation numérique réussies |
| Reprise du script | Réussie ; correctif reconnu, zéro objet C++ recompilé |
| `go run -buildvcs=false ./examples/kmeans` avec l'installation du fork | Réussi pour les deux algorithmes |
| Projet externe : `go get`, `go mod verify`, `go mod tidy`, exemple | Réussis, sans remplacement local du module |

Le préfixe de validation contient 1 196 fichiers inscrits au manifeste CMake,
en plus du SDK et des fichiers produits par le script. L'ensemble mesure environ
296 Mo de fichiers, dont 37,5 Mo pour les bibliothèques natives. Ces chiffres
excluent les paquets APT, les sources et objets en cache, et les caches Go.

Le projet externe a téléchargé
`github.com/milvi/mlpack-go@v0.0.0-20261102004546-a4cedd823d9c`, correspondant au
commit publié `a4cedd823d9c44d970a1330fe6ef8a2d4e27149c`. Son exemple utilise les
bibliothèques natives du préfixe installé par le script corrigé. Le journal et
son `go.mod` sont conservés avec les autres preuves.

Pour le diagnostic natif : glibc `mallinfo2`, somme des octets alloués dans les
arènes et par mappage, après chauffe, `runtime.GC()` et `malloc_trim(0)`.
Ce n'est ni le tas Go ni simplement la mémoire résidente du processus.

| Mesure | Avant, laboratoire | Après, laboratoire | Après, Docker |
|---|---:|---:|---:|
| 50 appels k-means × 100 000 points, dimension 3, trois centres ; sorties abandonnées | +40 081 888 octets | +2 464 octets | +3 920 octets |
| 50 000 cycles création / paramètres chaînes / destruction | Non mesuré | +25 568 octets | +26 848 octets |

Le seuil de régression est fixé à 2 Mio pour tolérer les allocations internes
stabilisées de l'allocateur et des bibliothèques. Les deux variations corrigées
sont largement inférieures ; ce test n'affirme pas une variation strictement nulle
sur toutes les machines.

### Environnement réellement exécuté

L'hôte est Ubuntu 24.04 avec noyau Linux 6.18.57. Il ne permet pas Docker/Podman,
un vrai `chroot` ou PRoot avec `ptrace`. Une **racine distincte Ubuntu Base
26.04.1 amd64** a été installée avec les vrais paquets Ubuntu 26.04 : GCC 15.2.0,
glibc 2.43, Armadillo `1:15.2.1+dfsg-2`, ensmallen `3.11.0-1`, cereal
`1.3.2+dfsg-6`. Les outils et bibliothèques de cette racine sont exécutés via
fakechroot et son chargeur glibc. Le SDK officiel Go 1.27.3 est utilisé.

Les adaptateurs du laboratoire virtualisent les chemins et fixent l'interpréteur
ELF ; ils ne sont pas dans le script livré. Le compte a l'identifiant 0 sans
capacités Linux : le scénario sous utilisateur ordinaire est fourni dans le
Dockerfile, pas attesté par ce laboratoire. La méthode exacte et le SHA-256 de
l'image de départ figurent dans l'audit amont.

Les 47 bibliothèques ont d'abord été compilées depuis la source 4.8.0. Le test
du script corrigé réutilise les objets de méthodes déjà construits, recompile les
utilitaires modifiés, réinstalle dans un nouveau préfixe et reconstruit l'exemple
contre le module de ce fork. Les contrôles d'archives, l'extraction du SDK et
l'application du patch à des sources propres ont été exercés. Cette reprise
ne vaut pas construction native intégralement neuve sous UID ordinaire.

### Construction neuve en conteneur : réussie

Le [workflow GitHub Actions du 2 novembre 2026](https://github.com/milvi/mlpack-go/actions/runs/41368925523)
est **réussi**, sur le commit `a4cedd823d9c44d970a1330fe6ef8a2d4e27149c`.
Le Dockerfile y a réellement été exécuté dans Ubuntu 26.04, sans objets C++
préchargés. APT s'exécute d'abord comme administrateur ; le script et les tests
s'exécutent ensuite sous l'utilisateur ordinaire `mlpack`.

Le journal atteste 48 unités de compilation C++ (deux pour l'utilitaire commun,
46 pour les méthodes), les 47 bibliothèques construites, puis :

- installation et validation numérique de k-means par l'installateur ;
- suite `go test ./...` réussie ;
- suite avec `GOEXPERIMENT=cgocheck2` réussie ;
- tests des allocations natives réussis, avec les mesures Docker du tableau ;
- `docker run` de l'image produite, avec k-means `naive` et `dualtree` réussis.

Des [extraits du journal](validation/ci-ubuntu2604.log) sont conservés dans le
dépôt. Cette exécution complète la validation via fakechroot ; elle ne repose
pas sur les adaptateurs du laboratoire. Les commits de documentation ultérieurs
peuvent utiliser `[skip ci]` : ils ne modifient pas le code testé et référencent
explicitement le commit de cette exécution.

## Limites conservées

- Les objets modèles opaques gardent le contrat de durée de vie historique ;
  aucune API générale `Close` ni restitution complète des modèles n'est ajoutée.
  La correction des chaînes de leurs accesseurs ne résout pas cette question.
- Le wrapper d'algorithme ne convertit toujours pas toutes les exceptions C++
  en erreurs Go. Les deux nouvelles fonctions de copie interceptent leurs propres
  exceptions, ce qui ne rend pas tous les algorithmes résistants à toute entrée.
- Les utilitaires de journalisation et d'autres réglages C++ sont globaux.
  La sûreté de plusieurs appels concurrents n'a pas été établie.
- La conversion de colonne entière (`ucol`) reste dans l'interface commune ; elle
  n'a pas été exercée par les tests dynamiques décrits ici. Les anciens accesseurs
  C de pointeurs peuvent encore produire des avertissements GCC ; ce fork Go
  utilise le nouveau chemin de copie.
- Les limites de tailles des accesseurs de dimensions historiques restent en
  place. Cet audit ne certifie pas toutes les tailles de matrices ni tous les
  algorithmes de la bibliothèque.
- Les performances, le pic mémoire et le résultat pour des volumes beaucoup plus
  grands doivent être mesurés sur le jeu de données réel. Le cas de validation
  n'est pas un benchmark de cette charge.

## Références de l'interface Go

- [Règles cgo et pointeurs conservés](https://go.dev/src/cmd/cgo/doc.go)
- [runtime.Pinner](https://pkg.go.dev/runtime#Pinner)
- [Source amont de GetMemory](https://github.com/mlpack/mlpack/blob/4.8.0/src/mlpack/bindings/go/mlpack/capi/arma_util.hpp)
- [Générateur amont des wrappers](https://github.com/mlpack/mlpack/blob/4.8.0/src/mlpack/bindings/go/print_go.cpp)
