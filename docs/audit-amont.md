> Audit historique du 2 novembre 2026 : décrit la distribution antérieure à la compilation automatique. Voir le rapport actuel pour les corrections et limites restantes.

# Audit amont de mlpack / Go — avant corrections

> **État amont avant correction.** Ce rapport est conservé comme historique de
> l'audit et de la première installation sans corrections. Les références au
> « kit » désignent ce premier kit d'audit, pas l'installateur actuel du fork.
> Pour installer le fork et connaître ses résultats, lire le [README](../README.md)
> et [corrections-et-validation.md](corrections-et-validation.md).

Audit réalisé le 2 novembre 2026. Architecture visée : Ubuntu 26.04 amd64.
APT signifie Advanced Package Tool, le gestionnaire de paquets Ubuntu utilisé
ici. GCC signifie GNU Compiler Collection. Les empreintes SHA-256 utilisent
Secure Hash Algorithm sur 256 bits.
Le périmètre couvre l'installation, la génération et la distribution des
bindings, leur liaison native et le chemin k-means. Il ne prétend pas certifier
les implémentations mathématiques de tous les algorithmes de mlpack.

## Versions et sources

| Élément | Version / référence fixée |
|---|---|
| mlpack C++ testé | 4.8.0, tag au commit `cef4041029a62be209d9a5e4f8f8bc06036230f1` |
| Branche principale examinée | `9767815cb73a924300535b09c5d08a3a8cc33f57` |
| Distribution des wrappers Go | `3505bb06392e345e1a15bb641cfd1d20974cb172` |
| Version de module Go correspondante | `v0.0.0-20260617213524-3505bb06392e` |
| Site mlpack.org examiné | `f733af5f11fb938866561636a9535e8b325e10b5` |
| Go | 1.27.3, dernière version stable renvoyée par go.dev lors de l'audit |
| Gonum | 0.7.0, dépendance déclarée par le module officiel |
| Armadillo Ubuntu | `1:15.2.1+dfsg-2` |
| ensmallen Ubuntu | `3.11.0-1` |
| cereal Ubuntu | `1.3.2+dfsg-6` |

Empreinte SHA-256 de l'archive officielle mlpack 4.8.0 utilisée :

```text
0ab06e5c506c7ed5f072faa4c04477c65624e4a1ff62eee8c0d996ef850ec51c
```

Empreinte publiée par go.dev pour l'archive Linux amd64 Go 1.27.3 :

```text
65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40
```

Le téléchargement mlpack est fixé par l'empreinte de l'archive récupérée en HTTPS
pendant cet audit ; il ne s'agit pas d'une signature indépendante de l'éditeur.
Les modules Go sont contrôlés par leurs entrées `go.sum` et par `go mod verify`.
Les paquets Ubuntu sont pris dans les dépôts APT signés de la distribution.
Les 1 985 fichiers de `src/`, `CMake/`, `doc/` et le CMake principal présents
dans l'archive ont également été comparés aux empreintes des blobs du tag GitHub
4.8.0 : aucune différence ni absence. Les six fichiers Go utilitaires distribués
(`arma_util.go`, `io_util.go`, `numcsv.go`, `doc.go`, `go.mod`, `go.sum`) sont
identiques à ceux des sources C++ 4.8.0.

## Il existe plusieurs voies, avec des objets différents

| Voie | Ce qui est compilé | Ce qui doit être lié |
|---|---|---|
| Programme C++ direct | Les modèles C++ de mlpack inclus par les en-têtes | Armadillo et ses dépendances ; OpenMP si activé |
| Module Go officiel + bibliothèques natives | Les wrappers Go par `go build` ; séparément, les instanciations C++ par CMake | `libmlpack_go_util.so`, les `libmlpack_go_<méthode>.so` et leurs dépendances |
| Génération intégrale des bindings depuis le dépôt principal | Les générateurs, les fichiers Go produits et les bibliothèques natives | Les mêmes bibliothèques natives |
| Binding cgo personnalisé avec un `.cpp` dans son propre paquet | cgo fait compiler ce `.cpp` pendant `go build` | Armadillo ; les modèles mlpack sont instanciés dans ce `.cpp` |

La quatrième voie est possible, mais **ce n'est pas le contenu du module Go
officiel actuel**. Dans ce dernier, `capi/` distribue les déclarations C et les
fichiers Go demandent explicitement les bibliothèques natives par `#cgo LDFLAGS`.
La procédure de publication retire volontairement les fichiers `.so` avant de
committer le dépôt de distribution.

Le cœur mlpack 4 est constitué d'en-têtes, d'où l'absence de nécessité d'une
`libmlpack.so` centrale pour un programme C++ direct. Le code machine existe
bien après compilation : il est instancié dans le programme ou dans les
bibliothèques des bindings. « Header-only » ne signifie pas « aucun code natif ».

L'appel officiel traverse successivement le code Go, cgo, une fonction C à ABI
stable (`extern "C"`), puis le C++ mlpack. ABI signifie interface binaire
d'application. Armadillo est l'interface d'algèbre linéaire ; BLAS signifie
Basic Linear Algebra Subprograms et LAPACK signifie Linear Algebra Package.
OpenMP désigne Open Multi-Processing, utilisé pour le parallélisme CPU.

`BUILD_GO_BINDINGS=ON` sert à générer les fichiers Go, en plus des bibliothèques.
`BUILD_GO_SHLIB=ON` suffit pour construire les bibliothèques si l'on utilise les
fichiers Go déjà publiés. SHLIB désigne une bibliothèque partagée.
Le message CMake « Not building Go bindings » est donc attendu avec les options
du kit : les fichiers Go ne sont pas régénérés, les bibliothèques le sont.

Importer tout le paquet officiel agrège les directives cgo de ses fichiers.
Construire seulement `libmlpack_go_kmeans.so` ne suffit pas à satisfaire tous
les noms `-l` demandés par ce paquet. Un paquet Go réduit à k-means aurait un
autre périmètre et doit être présenté comme tel.

## Rôle des trois dépôts et apparente circularité

`mlpack/mlpack` est la source des algorithmes, générateurs, utilitaires du binding,
documents Markdown et règles CMake. C'est la source de référence de construction.

`mlpack/mlpack-go` est le dépôt de distribution des fichiers générés. Son README
renvoie au dépôt principal pour les rapports de problèmes. Sa procédure de
publication demeure manuelle. Le commit de juin 2026 met bien mlpack à 4.8.0,
mais ne modernise pas l'ensemble du processus d'installation.

`mlpack/mlpack.org` fournit le site et l'import Go de présentation
`mlpack.org/v1/mlpack`. Le fichier `v1/mlpack` contient une balise `go-import`
qui donne à Go l'URL Git réelle : `https://github.com/mlpack/mlpack-go`.
Son script `scripts/build-website.sh` extrait les documents depuis l'archive
mlpack choisie et construit les pages HTML.

Il y a des renvois documentaires qui tournent en rond entre installation,
quickstart et dépôt de distribution. Il n'y a pas de dépendance de compilation
intrinsèquement circulaire : les sources C++ ne nécessitent pas l'installation
préalable du module Go pour construire les bibliothèques natives.

Le binding Go n'est pas entièrement abandonné : le dépôt principal conserve un
travail Go dans `.github/workflows/ci.yml`, avec génération et exécution de tests.
CI signifie intégration continue. Cela ne remet pas à jour les anciens fichiers
d'installation du dépôt de distribution ni ne certifie leur fonctionnement.

## Défauts concrets du parcours publié

| Fichier / mécanisme | Constat | Conséquence |
|---|---|---|
| Quickstart Go | `go get -u -d`, recherche d'un chemin de module non versionné, puis `make` et `sudo make install` dans le cache | Mélange de résolution de modules et d'installation système ; cache normalement en lecture seule ; `-d` est obsolète avec les modules modernes |
| Makefile, cible par défaut | Le premier objectif ordinaire est `deps` | Un simple `make` prépare des paquets, il ne constitue pas une compilation complète du binding |
| Makefile, `setup_armadillo` | Reconstruit Armadillo 10.8.2 et fait `sudo make install` | Duplique la bibliothèque déjà fournie par Ubuntu dans `/usr/local` |
| Makefile, `setup_cereal` | URL `cereal-1.3.0.tar.xz`, décompression forcée avec `tar -z`, copie directe dans `/usr/local/include` | Incohérence de format ; téléchargement ancien fragile ; installation hors APT |
| Makefile, `build` | `BUILD_TESTS=ON`, `DOWNLOAD_DEPENDENCIES=ON`, parallélisme sur tous les cœurs | Travail inutile pour l'installation ; dépendances potentiellement non fixées ; consommation mémoire non bornée par le script |
| Makefile, `install` | Relance préparation et construction depuis l'objectif d'installation | Avec `sudo make install`, une grande partie du processus s'exécute aussi comme administrateur |
| Makefile, `clean` | `go clean --cache` et répertoires fixes dans `/tmp` | Efface le cache Go partagé ; risque de collisions entre constructions |
| Dockerfile | Debian Stretch, Go 1.13.5 par défaut, Armadillo 8.400.0 en dur | Recette historique inadaptée à Ubuntu 26.04 / mlpack 4.8.0 |
| `.travis.yml` du dépôt Go | Xenial, Go 1.14, Xcode 9.4 | Ne constitue pas une validation actuelle de cette installation |
| Documentation générée de `Kmeans` | Exemples avec données avant le nombre de classes | Contredit la signature effective `Kmeans(clusters int, input *mat.Dense, ...)` |
| Texte `KillEmptyClusters` | Décrit encore le remplissage par DBL_MAX | Le corps C++ 4.8.0 supprime désormais une colonne avec `shed_col` et le compteur correspondant |

L'URL ancienne de cereal renvoyait HTTP 404 lors du contrôle de la ressource
normalisée. L'archive mlpack 4.8.0 reste téléchargeable et son empreinte est fixée.
Le kit ne réutilise aucun des téléchargements d'Armadillo/cereal du Makefile.

Ubuntu 26.04 propose également `libmlpack-dev` 4.7.0-2. Ce paquet installe
l'interface C++ de cette version ; il n'est pas l'installation des bibliothèques
Go 4.8.0 attendues ici. Le kit utilise donc directement les sources C++ 4.8.0
et les dépendances mathématiques Ubuntu, sans ajouter ce paquet mlpack 4.7.0.

## Choix de la procédure nettoyée

Les dépendances natives restent gérées par APT. Les fonctionnalités image, audio
et HTTP conservent les options du CMake officiel : les bibliothèques d'en-têtes
qui font partie de l'archive mlpack ne sont pas remplacées par des copies
supplémentaires. Aucun Boost, Python, Julia, R, outil de génération de pages
de manuel ni installation globale de Gonum n'est nécessaire pour cette voie.

La compilation désactive les tests C++ généraux, les autres langages, les
programmes en ligne de commande et le téléchargement automatique de dépendances.
Elle conserve OpenMP. Le test numérique Go du kit est exécuté après installation.
Il valide le résultat demandé ; il ne remplace pas la suite de tests amont.

Les fichiers C++ et les bibliothèques sont installés dans un préfixe utilisateur.
`CMAKE_INSTALL_LIBDIR=lib` évite une variation de chemin liée au multiarch.
Le RPATH des bibliothèques vaut `$ORIGIN`, c'est-à-dire leur propre répertoire.
Les exécutables Go reçoivent le chemin de ce préfixe via les options de cgo.
RPATH / RUNPATH désignent les chemins de recherche inscrits dans les objets ELF
(Executable and Linkable Format), le format binaire utilisé ici par Linux.

Le préfixe local ne devient pas un second dossier système d'en-têtes ou de
bibliothèques. Son utilisation est explicite. Aucun `sudo make install`, aucun
ajout à la configuration globale du chargeur et aucun nettoyage du cache Go
partagé n'est effectué.

APT peut installer des dépendances transitives, notamment HDF5, SuperLU,
ARPACK, BLAS, LAPACK et des bibliothèques de développement associées.
Leur nombre dépend de la machine. `apt-get --simulate` montre la transaction
réelle à effectuer, y compris d'éventuelles mises à jour de dépendances.
Les bibliothèques installées par APT peuvent enregistrer des alternatives BLAS
et LAPACK ; le script utilisateur ne les reconfigure pas.
OpenBLAS est une optimisation facultative distincte, pas un prérequis ajouté
silencieusement par ce kit. Installer une variante d'OpenBLAS peut changer les
alternatives utilisées par d'autres applications ; ce choix doit être explicite.

## Défauts du binding indépendants de son installation

La directive `go 1.13` du module indique la version minimale et le niveau de
langage de ce module ; elle ne demande pas d'installer l'ancien compilateur Go
1.13. Le cas testé utilise Go 1.27.3 avec le graphe de modules du kit. Dans un
projet Go préexistant, la sélection minimale des versions peut retenir une
version de Gonum plus récente à cause d'autres dépendances : ce serait un graphe
différent de celui validé ici.

### Allocations C de chaînes non libérées

`getParams`, `setPassed`, `setParamInt`, `setParamString`, les conversions de
matrices et d'autres fonctions appellent `C.CString`. Ces chaînes sont allouées
sur le tas C et doivent être libérées par l'appelant. Le chemin audité ne contient
pas les appels de libération correspondants. Les conversions vers `std::string`
côté C++ ne libèrent pas l'original alloué par `C.CString`.
Il s'agit d'un défaut de gestion de mémoire, pas d'un manque de bibliothèque APT.

### Propriété des matrices de sortie

`GetMemory` dans `capi/arma_util.hpp` copie les petites matrices vers une nouvelle
allocation native ou détache la propriété du tampon d'une matrice plus grande.
`armaToGonumMat` transforme ensuite ce pointeur en slice Go et en `mat.Dense`.
Il n'enregistre ni propriétaire avec méthode de restitution ni finaliseur qui
libère ce tampon. Les appels `runtime.KeepAlive` visibles ne sont pas des
finaliseurs et ne libèrent pas la mémoire native.
La destruction de l'objet `Params` a précisément été empêchée de libérer ce
tampon ; le ramasse-miettes Go ne peut pas le faire à sa place.

Pour un appel k-means avec n points de dimension d et k centres, le stockage
des sorties allouées côté C++, avec `LabelsOnly=true`, représente environ
8(n + kd) octets. À n = 10⁷, d = 3, k = 10⁶, cela représente environ 104 Mo
décimaux par appel, dont 80 Mo pour les seules étiquettes, hors données d'entrée,
structures d'arbres et autres allocations. Les options et le chemin d'itération
peuvent changer le stockage des centres ; le calcul de 104 Mo décrit le cas où
les deux sorties sont allouées côté C++. Le fait de
ne plus référencer les matrices Go ne constitue pas un mécanisme de restitution
de ces allocations natives dans cette version.

Le programme `diagnostic-memoire/main.go` mesure les allocations glibc avec
`mallinfo2` après une chauffe, puis après 50 appels sur 100 000 points de dimension
3, avec trois centres et les sorties immédiatement abandonnées. Il force
`runtime.GC()` et `malloc_trim(0)` avant chaque mesure. Résultat observé :

```text
avant : 1318816 octets natifs
après : 41400704
différence : 40081888
étiquettes seules, 50 allocations de 100 000 float64 : 40000000 octets
```

Ce diagnostic mesure les octets considérés alloués par glibc et les blocs qu'elle
a obtenus par mappage mémoire ; il ne s'agit pas simplement de la mémoire résidente
du processus. La croissance observée concorde avec la perte des tampons
d'étiquettes identifiée dans le code. Ce contrôle a utilisé les 47 bibliothèques
de la première compilation, copiées pour éviter une modification des objets
pendant l'exécution, avec un chemin de recherche adapté au répertoire de test.
Les sources C++ et Go n'ont pas été corrigées pour ce diagnostic.

### Autres limites statiques

- `gonumToArmaMat` transmet une vue sur le tampon Go, retenue dans l'objet de
  paramètres C++ entre plusieurs appels cgo. Le chemin ne montre pas de
  `runtime.Pinner` pour prolonger l'épinglage au-delà du premier appel. Une
  conservation d'un pointeur Go par C doit respecter les règles d'épinglage cgo.
- Le champ `Stride` d'une vue Gonum n'est pas transmis : une sous-matrice non
  contiguë ne doit pas être utilisée comme si son stockage était compact.
- Les convertisseurs prennent `&data[0]` : les matrices vides ne sont pas une
  entrée neutre prise en charge par ces fonctions.
- Le wrapper C++ généré appelle directement l'algorithme, sans transformer une
  exception C++ en valeur `error` Go. Le binding n'offre pas de contrat général
  d'erreurs Go pour tous les paramètres invalides.
- Les réglages de verbosité et d'autres utilitaires sont globaux côté C++ :
  aucun contrat de sûreté des appels concurrents n'a été établi par cet audit.

Ces constats viennent de la lecture du code. Ils ne doivent pas être transformés
en affirmation selon laquelle chaque appel planterait : un exemple peut réussir
tout en laissant subsister des problèmes lors d'appels répétés ou d'entrées
différentes. Un petit binding cgo spécifique à k-means, avec propriété mémoire
explicite et interception des exceptions, est une autre voie raisonnable pour
une utilisation intensive ; il constitue du code à développer et à valider.

## Validation effectivement réalisée

| Vérification | Résultat |
|---|---|
| Sources C++ comparées au tag 4.8.0 | 1 985 fichiers identiques |
| Dépendances APT de la racine de test | Transaction de 151 ajouts et 9 mises à jour ; configuration achevée |
| Construction native | 47 bibliothèques Go compilées avec GCC 15.2.0, C++17, `-O3`, OpenMP |
| Installation dans le préfixe utilisateur | Réussie ; manifeste CMake conservé |
| Modules Go | `go mod verify` réussi ; versions fixées |
| Programme d'exemple Go 1.27.3 | Compilation puis exécution réussies |
| `go run -buildvcs=false .` depuis l'exemple installé | Réussi |
| k-means `naive` | 9 étiquettes correctes, 3 classes de 3 points, centres exacts à 10⁻¹² près |
| k-means `dualtree` | Même validation numérique réussie |
| Données d'entrée | Inchangées après les deux appels |
| Dépendances ELF de l'exécutable | 47 dépendances directes `libmlpack_go_*`, toutes résolues |
| Reprise du script sans modification de configuration | Réussie ; 0 objet C++ recompilé |
| Diagnostic des allocations natives | Croissance de 40 081 888 octets pour 50 appels de 100 000 points |

Le préfixe complet mesure environ 296 Mo décimaux de fichiers, dont 37,5 Mo pour
les 47 bibliothèques natives. Cette mesure comprend le SDK Go privé et les
en-têtes ; elle exclut les paquets APT, le cache de construction et les caches Go.

La sortie numérique vérifiée comprend :

```text
naive : OK — 9 points, 3 dimensions, 3 classes non vides
étiquettes : [0 0 0 1 1 1 2 2 2] ; effectifs : [3 3 3]
dualtree : OK — 9 points, 3 dimensions, 3 classes non vides
étiquettes : [0 0 0 1 1 1 2 2 2] ; effectifs : [3 3 3]
```

Le C++ n'a pas été rendu artificiellement exempt d'avertissements. GCC 15 émet
notamment `-Wreturn-local-addr` sur les convertisseurs `mlpackArmaPtrUrow` et
`mlpackArmaPtrUcol`. Les journaux les conservent. La présente validation n'exerce
pas tous ces convertisseurs et ne conclut pas à l'absence générale de défauts.

### Portée et limites du laboratoire

L'hôte de l'outil est Ubuntu 24.04, avec un noyau Linux 6.18.57. Il ne permet ni
conteneur Docker/Podman, ni `chroot`, ni PRoot avec `ptrace`. Le test a donc
utilisé une **racine distincte Ubuntu Base 26.04.1 amd64**, avec les vrais paquets,
en-têtes, compilateur et bibliothèques d'Ubuntu 26.04, exécutés via fakechroot et
le chargeur glibc 2.43 de cette racine. Ce n'est ni une machine virtuelle bootée
sous Ubuntu 26.04, ni une exécution native dans un conteneur.

Image de départ :
`ubuntu-base-26.04.1-base-amd64.tar.gz`, archive officielle dont le SHA-256
`a496a960472ce474a59590b8987d3a1135d3cbef1991f3b1abe8cacfea8bf85a`
correspond au fichier SHA256SUMS du serveur Ubuntu.

Les adaptations du laboratoire ont porté sur la virtualisation des chemins,
pas sur les algorithmes de mlpack :

- outils GNU coreutils de la distribution utilisés dans la racine de test,
  l'invocation par un chargeur explicite perturbant l'identification des commandes
  par les coreutils Rust de l'image de base ;
- adaptation de `update-alternatives` via son option `--root`, et normalisation
  des liens symboliques pour préserver leur sens dans une racine déplacée ;
- lancement du GCC Ubuntu 26.04 dans cette racine, avec préfixe d'outils explicite ;
- interpréteur ELF et chemins des bibliothèques système adaptés à cette racine
  pour les programmes Go liés avec cgo ;
- certificat du proxy de l'environnement de test ajouté à la confiance de cette
  racine, sans désactiver la vérification TLS.

Le binaire officiel Go 1.27.3 a été utilisé. La construction et l'exécution font
bien appel au C++ et aux bibliothèques Ubuntu 26.04 ; les adaptateurs évitent
d'utiliser les bibliothèques C++ et glibc 24.04 de l'hôte pour ce programme.
Le script livré ne contient pas ces adaptateurs de laboratoire.

Le compte du laboratoire a l'identifiant 0 sans capacités Linux. Le chemin
d'installation utilisateur et l'absence d'écritures système dans le script ont
été vérifiés, mais le scénario sous UID 1000 n'a pas été exécuté ici.
`Dockerfile.validation` fournit un scénario Ubuntu 26.04 avec un utilisateur
ordinaire pour revalidation sur une machine disposant de Docker/Podman ; ce
Dockerfile est explicitement fourni **sans affirmation de test dans ce laboratoire**.

Les 47 bibliothèques ont d'abord été construites depuis les sources. Le test
final du script a réutilisé ces objets vérifiés pour sa reprise après les
ajustements de chemins ; les téléchargements HTTPS de Go et mlpack, leur contrôle
SHA-256 et l'extraction du SDK dans un nouveau préfixe ont aussi été exercés.
Le succès de ce test est une preuve d'intégration sur le cas décrit ; il ne
constitue pas une garantie universelle sur une machine dont la configuration,
les bibliothèques manuelles et le graphe de modules Go sont inconnus.

La dernière version du script évite de recompiler inutilement l'ensemble à chaque
reprise : si CMake réécrit `config-local.hpp` avec un contenu identique, sa date
antérieure est conservée. Une modification du contenu conserve naturellement sa
nouvelle date et les dépendances de compilation usuelles restent actives.
L'exemple est compilé sans métadonnées d'un éventuel dépôt Git parent.

Les preuves du diagnostic mémoire initial et les journaux des nouveaux tests
sont conservés dans [validation/](validation/historique-2026-11-02/). Le premier kit d’audit conservait
aussi les journaux complets de construction, la transaction APT et les manifestes
de cette première installation. Le script actuel produit ses propres manifestes
dans le préfixe d’installation.

## Références primaires

- [Installation C++ et bindings](https://github.com/mlpack/mlpack/blob/4.8.0/doc/user/install.md)
- [Quickstart Go](https://github.com/mlpack/mlpack/blob/4.8.0/doc/quickstart/go.md)
- [CMake principal](https://github.com/mlpack/mlpack/blob/4.8.0/CMakeLists.txt)
- [CMake du binding Go](https://github.com/mlpack/mlpack/blob/4.8.0/src/mlpack/bindings/go/CMakeLists.txt)
- [Recherche et téléchargement des dépendances](https://github.com/mlpack/mlpack/blob/4.8.0/CMake/mlpack.cmake)
- [Utilitaire C++ de propriété mémoire](https://github.com/mlpack/mlpack/blob/4.8.0/src/mlpack/bindings/go/mlpack/capi/arma_util.hpp)
- [Utilitaire C++ de conversion des matrices](https://github.com/mlpack/mlpack/blob/4.8.0/src/mlpack/bindings/go/mlpack/capi/arma_util.cpp)
- [Générateur de wrapper C++](https://github.com/mlpack/mlpack/blob/4.8.0/src/mlpack/bindings/go/go_method.cpp.in)
- [Programme k-means appelé par les bindings](https://github.com/mlpack/mlpack/blob/4.8.0/src/mlpack/methods/kmeans/kmeans_main.cpp)
- [Politique de suppression des classes vides](https://github.com/mlpack/mlpack/blob/4.8.0/src/mlpack/methods/kmeans/kill_empty_clusters.hpp)
- [Makefile distribué](https://github.com/mlpack/mlpack-go/blob/3505bb06392e345e1a15bb641cfd1d20974cb172/Makefile)
- [Dockerfile distribué](https://github.com/mlpack/mlpack-go/blob/3505bb06392e345e1a15bb641cfd1d20974cb172/Dockerfile)
- [Publication du module Go](https://github.com/mlpack/mlpack-go/blob/3505bb06392e345e1a15bb641cfd1d20974cb172/rel/deployment.md)
- [Conversions Go des matrices](https://github.com/mlpack/mlpack-go/blob/3505bb06392e345e1a15bb641cfd1d20974cb172/arma_util.go)
- [Utilitaires Go de paramètres](https://github.com/mlpack/mlpack-go/blob/3505bb06392e345e1a15bb641cfd1d20974cb172/io_util.go)
- [Fonction Go Kmeans](https://github.com/mlpack/mlpack-go/blob/3505bb06392e345e1a15bb641cfd1d20974cb172/kmeans.go)
- [Import Go de présentation](https://github.com/mlpack/mlpack.org/blob/f733af5f11fb938866561636a9535e8b325e10b5/v1/mlpack)
- [Construction du site](https://github.com/mlpack/mlpack.org/blob/f733af5f11fb938866561636a9535e8b325e10b5/scripts/build-website.sh)
- [Intégration continue actuelle du dépôt principal](https://github.com/mlpack/mlpack/blob/9767815cb73a924300535b09c5d08a3a8cc33f57/.github/workflows/ci.yml)
- [Allocations C et règles des pointeurs cgo](https://go.dev/src/cmd/cgo/doc.go)
- [Distribution officielle Go](https://go.dev/dl/?mode=json)
- [Paquets Armadillo d'Ubuntu 26.04](https://packages.ubuntu.com/resolute/libarmadillo-dev)
