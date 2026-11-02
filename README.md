# mlpack pour Go — fork corrigé, Ubuntu 26.04 amd64

Ce fork fournit les wrappers Go de **mlpack 4.8.0**, des corrections de propriété
mémoire et une installation locale vérifiée avec **Go 1.27.3**. Le module est
`github.com/milvi/mlpack-go`. Les bibliothèques C++ correspondantes sont
construites par le script ; elles ne sont pas incluses dans le module Go.

Les corrections du générateur et de l'interface native sont conservées dans
[milvi/mlpack](https://github.com/milvi/mlpack). Le présent dépôt contient aussi
leur [patch pour la version 4.8.0](patches/mlpack-4.8.0-go-memory.patch), dont
l'empreinte est contrôlée par l'installateur. **Un seul checkout de ce dépôt
suffit pour installer et lancer l'exemple.**

Lire [l'audit amont](docs/audit-amont.md) pour les défauts identifiés et
[les corrections et validations](docs/corrections-et-validation.md) pour leur
portée exacte. Les modèles opaques, les exceptions des algorithmes et les appels
concurrents ne bénéficient pas d'une garantie générale dans ce fork.

## Installation

Ubuntu **26.04**, architecture **amd64 / x86-64**, dépôts comprenant `universe`.
Examiner d'abord la transaction APT, puis installer les prérequis :

```sh
sudo apt-get update
apt-get --simulate install --no-install-recommends \
  build-essential cmake pkg-config curl ca-certificates patch \
  libarmadillo-dev libensmallen-dev libcereal-dev

sudo apt-get install --no-install-recommends \
  build-essential cmake pkg-config curl ca-certificates patch \
  libarmadillo-dev libensmallen-dev libcereal-dev
```

Récupérer ensuite ce dépôt, par Git ou avec l'archive GitHub de la révision
corrigée. Avec Git déjà installé :

```sh
git clone https://github.com/milvi/mlpack-go.git
cd mlpack-go
git checkout a4cedd823d9c44d970a1330fe6ef8a2d4e27149c
```

Cette révision contient le code validé ; les commits de documentation ultérieurs
peuvent compléter les preuves de test. Git n'est pas un prérequis du script
lui-même. Depuis ce checkout, lancer **sans sudo** :

```sh
set -o pipefail
bash installer.sh 2>&1 | tee installation-mlpack.log
```

Deux compilations C++ sont lancées simultanément par défaut. Pour limiter leur
consommation mémoire : `JOBS=1 bash installer.sh`. La première construction
compile les **47 bibliothèques natives** liées par le package Go complet ; elle
peut prendre plusieurs minutes. Le téléchargement des dépendances par CMake est
désactivé. Go et mlpack sont téléchargés en HTTPS et vérifiés par SHA-256.

Le script installe un SDK Go privé. Il ne remplace pas une installation de Go
existante. Après succès :

```sh
source "$HOME/.local/opt/mlpack-go-milvi/4.8.0-go1.27.3/env.sh"
"$HOME/.local/opt/mlpack-go-milvi/4.8.0-go1.27.3/bin/kmeans-demo"

# Depuis le checkout :
go run -buildvcs=false ./examples/kmeans
go test -buildvcs=false ./...
```

L'exemple vérifie neuf points de ℝ³, trois groupes, les étiquettes, les centres
avec une tolérance de 10⁻¹² et l'absence de modification des entrées, pour les
algorithmes `naive` et `dualtree`. `INSTALLATION-VERIFIEE` est écrit dans le
préfixe uniquement après succès de la compilation et de cet exemple. Une reprise
commence par retirer l'ancien marqueur de succès et conserve les fichiers
intermédiaires pour diagnostic.

## Utiliser le fork dans son projet

Charger `env.sh`, puis, depuis le module de son application :

```sh
go get github.com/milvi/mlpack-go@v0.0.0-20261102004546-a4cedd823d9c
go mod tidy
```

Cette pseudo-version correspond à la révision corrigée ci-dessus et a été
téléchargée et exécutée depuis un projet séparé, sans directive `replace`.
La conserver dans le `go.mod` de l'application. Ne pas utiliser une révision
amont non corrigée avec ces instructions.

```go
import (
    mlpack "github.com/milvi/mlpack-go"
    "gonum.org/v1/gonum/mat"
)

// Une ligne par point : n × d côté Go.
donnees := mat.NewDense(4, 2, []float64{0, 0, 0, 1, 10, 10, 10, 11})
options := mlpack.KmeansOptions()
options.Algorithm = "dualtree"
options.LabelsOnly = true
options.Seed = 42
centres, etiquettes := mlpack.Kmeans(2, donnees, options)
```

`Kmeans` reçoit **le nombre de classes avant les données**. Les étiquettes sont
des valeurs entières stockées dans une matrice Gonum de `float64`. Sans
`LabelsOnly`, la sortie contient aussi les coordonnées des points. Les matrices
retournées appartiennent à Go et restent valides après libération des paramètres
natifs. Gonum **0.7.0** est conservé pour éviter une migration d'API dans ce
correctif. Un projet sélectionnant une autre version de Gonum doit être revalidé.

Lancer Helix et ses outils Go depuis le terminal ayant chargé `env.sh` pour
qu'ils héritent de la même chaîne de compilation et des chemins de liaison.

## Ce qui est écrit sur la machine

Notons P le préfixe `$HOME/.local/opt/mlpack-go-milvi/4.8.0-go1.27.3` et C le
cache `$HOME/.cache/mlpack-go-milvi/4.8.0-go1.27.3` avec les réglages habituels.

| Emplacement | Contenu | Gestionnaire |
|---|---|---|
| `/usr/bin`, `/usr/include`, `/usr/lib`, etc. | Compilateurs, outils et dépendances mathématiques, uniquement lors de la commande APT | APT / dpkg |
| Checkout courant | Wrappers Go, déclarations C et exemple ; le script les lit sans copier une distribution Go amont | Git / utilisateur |
| `P/go/` | SDK officiel Go 1.27.3 | Installateur |
| `P/lib/libmlpack_go_*.so` | 46 bibliothèques de méthodes et une bibliothèque utilitaire native | CMake |
| `P/include/mlpack.hpp`, `P/include/mlpack/` | En-têtes C++ mlpack 4.8.0 | CMake |
| `P/lib/pkgconfig/mlpack.pc` et autres fichiers du manifeste CMake | Métadonnées de l'installation C++ | CMake |
| `P/bin/kmeans-demo` | Exemple Go compilé depuis ce checkout | Installateur |
| `P/env.sh` | Variables à charger explicitement dans une session | Installateur |
| `P/share/audit/` | Manifeste CMake, empreintes des fichiers Go, versions et résultat k-means | Installateur |
| `C/` | Archives vérifiées, sources C++ patchées, objets et configuration de construction | Installateur |
| `go env GOMODCACHE`, `go env GOCACHE` | Modules téléchargés et cache de compilation | Go |

Les caches Go utilisent les réglages existants, habituellement `~/go/pkg/mod`
et `~/.cache/go-build`. Go peut aussi écrire les données de vérification sumdb
et ses compteurs locaux de télémétrie selon son mode préexistant. Le script ne
modifie pas ce mode. La transaction APT dépend des paquets déjà installés : elle
peut ajouter BLAS, LAPACK, ARPACK, SuperLU, HDF5 et leurs dépendances, ou mettre à
jour des paquets existants. Le script n'installe pas `libmlpack-dev` d'Ubuntu.

Le script ne fait aucun `sudo`, `go env -w`, `ldconfig`, changement de `.bashrc`
ou de `/etc/ld.so.conf.d`. Les bibliothèques restent dans P. Les paquets APT
peuvent enregistrer leurs alternatives BLAS/LAPACK ; l'installateur ne les
reconfigure pas et n'ajoute pas OpenBLAS silencieusement.

`env.sh` préfixe `PATH` avec `P/go/bin`, fixe `GOROOT`, `CGO_ENABLED=1` et
`GOTOOLCHAIN=local`, remplace `CGO_LDFLAGS` par `-LP/lib -Wl,-rpath,P/lib` et
ajoute `P/lib/pkgconfig` à `PKG_CONFIG_PATH`, uniquement dans la session qui le
charge. Il ne fixe pas `LD_LIBRARY_PATH`. Le chemin absolu de P est inscrit
dans les programmes Go compilés ainsi ; les bibliothèques natives retrouvent
leurs voisines grâce à `$ORIGIN`. Déplacer P demande de reconstruire ces
programmes.

Pour choisir les deux répertoires avant la première installation :

```sh
MLPACK_PREFIX="$HOME/.local/opt/mlpack-personnel" \
MLPACK_BUILD_CACHE="$HOME/.cache/mlpack-personnel" \
JOBS=2 bash installer.sh
```

Les chemins doivent être absolus, sans espaces, deux-points ni caractères de
citation. Éviter toute construction concurrente sur le même préfixe. Un verrou
du cache est pris lorsque `flock` est disponible.

## Installation antérieure et désinstallation

Une copie manuelle d'Armadillo, ensmallen ou cereal dans `/usr/local` peut masquer
APT. Le script refuse les emplacements usuels détectés et affiche le fichier
concerné. Il ne supprime rien. Examiner l'ancienne installation, ses utilisateurs
et son éventuel `install_manifest.txt` avant de déplacer ses fichiers. `dpkg -S
/chemin/du/fichier` permet de chercher son propriétaire APT.

La désinstallation locale consiste à supprimer P et C **après vérification de
leur contenu**. Les programmes liés à P perdent alors leurs bibliothèques. Les
paquets APT, le checkout et les caches partagés de Go restent en place. Aucune
configuration du shell n'est à réparer.

## Validation reproductible

Depuis le checkout et après chargement de `env.sh` :

```sh
# Contrôles supplémentaires des pointeurs Go conservés côté C.
GOEXPERIMENT=cgocheck2 go test -buildvcs=false ./...

# Diagnostic glibc : allocations natives, distinctes du tas Go.
GOMAXPROCS=1 OMP_NUM_THREADS=1 OPENBLAS_NUM_THREADS=1 \
  go test -buildvcs=false -tags mlpack_memtest -run TestNative -v .
```

Une recette Docker et un workflow GitHub Actions exécutent l'installation dans
Ubuntu 26.04 sous un utilisateur ordinaire, puis les tests et l'exemple :

```sh
docker build --progress=plain --build-arg JOBS=2 -t mlpack-go-verification .
docker run --rm mlpack-go-verification
```

Cette construction intégrale, les tests ordinaires, les contrôles cgo renforcés,
les tests de mémoire native et l'exécution finale de k-means ont
[réussi sur la révision de code indiquée plus haut](https://github.com/milvi/mlpack-go/actions/runs/41368925523).

Les résultats réellement obtenus et les limites du laboratoire sont détaillés
dans [corrections-et-validation.md](docs/corrections-et-validation.md). La
reproductibilité des sources mlpack/Go et du patch est fixée par leurs empreintes ;
les dépôts APT et le tag d'image Ubuntu peuvent recevoir des mises à jour.
