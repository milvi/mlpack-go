#!/usr/bin/env bash
# Installation locale des bindings corrigés de mlpack 4.8.0 pour Go 1.27.3.
# Aucun sudo, aucune modification de /usr/local, de .bashrc ou de go env -w.
set -Eeuo pipefail
umask 022
export LC_ALL=C.UTF-8

readonly VERSION_MLPACK=4.8.0
readonly VERSION_GO=1.27.3
readonly SHA_MLPACK=0ab06e5c506c7ed5f072faa4c04477c65624e4a1ff62eee8c0d996ef850ec51c
readonly SHA_GO=65efbc80fdc1fa2cd2631cad5b481777a68e8725f34e89f50fc6ac714e379b40
readonly MODULE_GO=github.com/milvi/mlpack-go
readonly SHA_CORRECTIF=e2f1a1111a05d6313e47921ed5015ec29e8fc3a5a09cb4e2fbc7d38991c99a9d
readonly DOSSIER_SCRIPT="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd -P)"
PREFIXE="${MLPACK_PREFIX:-${HOME}/.local/opt/mlpack-go-milvi/4.8.0-go1.27.3}"
CACHE="${MLPACK_BUILD_CACHE:-${XDG_CACHE_HOME:-${HOME}/.cache}/mlpack-go-milvi/4.8.0-go1.27.3}"
JOBS="${JOBS:-2}"

erreur() { printf '\nERREUR : %s\n' "$*" >&2; exit 1; }
trap 'printf "\nÉchec ligne %s. Les fichiers partiels sont conservés pour diagnostic et reprise.\n" "$LINENO" >&2' ERR

case "${1:-}" in
  --help|-h)
    cat <<'AIDE'
Utilisation : bash installer.sh
Variables facultatives :
  MLPACK_PREFIX       répertoire d'installation absolu
  MLPACK_BUILD_CACHE  répertoire de construction absolu
  JOBS               nombre de compilations simultanées (2 par défaut)
Prérequis Ubuntu 26.04 amd64 :
  sudo apt-get update
  sudo apt-get install --no-install-recommends build-essential cmake pkg-config \
    curl ca-certificates patch libarmadillo-dev libensmallen-dev libcereal-dev
Aucun paquet APT n'est installé par ce script.
AIDE
    exit 0 ;;
  '') ;;
  *) erreur "Argument inconnu : $1. Utiliser --help." ;;
esac

[[ "$JOBS" =~ ^[1-9][0-9]*$ ]] || erreur "JOBS doit être un entier strictement positif."
for chemin in "$PREFIXE" "$CACHE"; do
  [[ "$chemin" == /* && "$chemin" != / ]] || erreur "Chemin absolu requis : $chemin"
  [[ "$chemin" != *[[:space:]\'\"\$\`\\:]* ]] || erreur "Chemin non pris en charge par les outils amont : $chemin"
done
[[ -f /etc/os-release ]] || erreur "/etc/os-release absent."
# shellcheck source=/dev/null
. /etc/os-release
[[ "${ID:-}" == ubuntu && "${VERSION_ID:-}" == 26.04 ]] || erreur "Cette procédure vise Ubuntu 26.04."
[[ "$(uname -m)" == x86_64 ]] || erreur "Cette archive de Go vise amd64/x86_64."

for commande in cmake make curl tar sha256sum pkg-config readelf dpkg-query sed patch; do
  command -v "$commande" >/dev/null || erreur "Commande absente : $commande. Voir --help pour les paquets APT."
done
LIB_ARMADILLO="$(dpkg-query -L libarmadillo-dev | sed -n '\|/libarmadillo\.so$|p')"
[[ -n "$LIB_ARMADILLO" && -f "$LIB_ARMADILLO" ]] || erreur "La bibliothèque du paquet libarmadillo-dev est absente. Voir --help."
command -v "${CC:-cc}" >/dev/null || erreur "Compilateur C absent."
command -v "${CXX:-c++}" >/dev/null || erreur "Compilateur C++ absent."
[[ -f "$DOSSIER_SCRIPT/examples/kmeans/main.go" && -f "$DOSSIER_SCRIPT/go.mod" ]] || erreur "Le dépôt complet doit accompagner installer.sh."

# Une ancienne copie manuelle peut masquer les paquets APT à la compilation.
# On la signale sans supprimer de fichiers qui n'appartiennent pas à ce kit.
for ancien in /usr/local/include/armadillo /usr/local/include/armadillo_bits \
  /usr/local/include/ensmallen.hpp /usr/local/include/ensmallen_bits \
  /usr/local/include/cereal /usr/local/lib/libarmadillo.so \
  /usr/local/lib64/libarmadillo.so /usr/local/lib/x86_64-linux-gnu/libarmadillo.so; do
  [[ ! -e "$ancien" ]] || erreur "Une installation manuelle peut masquer APT : $ancien. Voir README.md avant de la déplacer."
done

if [[ -e "$PREFIXE" && ! -d "$PREFIXE" ]]; then
  erreur "Le préfixe existe et n'est pas un répertoire : $PREFIXE"
fi
if [[ -d "$PREFIXE" && ! -f "$PREFIXE/.installation-mlpack-go" ]] &&
   [[ -n "$(find "$PREFIXE" -mindepth 1 -maxdepth 1 -print -quit)" ]]; then
  erreur "Le préfixe est déjà occupé par une autre installation : $PREFIXE"
fi
if [[ -f "$PREFIXE/.installation-mlpack-go" ]]; then
  [[ "$(cat "$PREFIXE/.installation-mlpack-go")" == 'mlpack-go-milvi 4.8.0 go1.27.3' ]] \
    || erreur "Le préfixe appartient à une autre version : $PREFIXE"
fi
mkdir -p -- "$PREFIXE" "$CACHE" "$PREFIXE/share/audit" "$PREFIXE/bin"
printf 'mlpack-go-milvi 4.8.0 go1.27.3\n' > "$PREFIXE/.installation-mlpack-go"
exec 9>"$CACHE/installation.lock"
if command -v flock >/dev/null; then
  flock -n 9 || erreur "Une installation utilise déjà ce cache."
fi
rm -f -- "$PREFIXE/INSTALLATION-VERIFIEE"
printf 'Installation : %s\nConstruction : %s\nCompilations simultanées : %s\n' "$PREFIXE" "$CACHE" "$JOBS"
printf '#include <armadillo>\n#include <ensmallen.hpp>\n#include <cereal/cereal.hpp>\nint main() { return 0; }\n' > "$CACHE/verification-dependances.cpp"
"${CXX:-c++}" -std=c++17 -c "$CACHE/verification-dependances.cpp" -o "$CACHE/verification-dependances.o"
pkg-config --modversion armadillo

telecharger() {
  local url="$1" fichier="$2" somme="$3"
  if [[ -f "$fichier" ]] && printf '%s  %s\n' "$somme" "$fichier" | sha256sum --check --status; then
    printf 'Archive vérifiée en cache : %s\n' "$fichier"
    return
  fi
  curl --fail --show-error --location --retry 3 --connect-timeout 30 \
    --proto '=https' --tlsv1.2 --output "$fichier.part" "$url"
  printf '%s  %s\n' "$somme" "$fichier.part" | sha256sum --check --status \
    || erreur "SHA-256 incorrect : $fichier.part"
  mv -- "$fichier.part" "$fichier"
}

telecharger "https://go.dev/dl/go${VERSION_GO}.linux-amd64.tar.gz" "$CACHE/go${VERSION_GO}.linux-amd64.tar.gz" "$SHA_GO"
if [[ ! -x "$PREFIXE/go/bin/go" ]]; then
  tar -xzf "$CACHE/go${VERSION_GO}.linux-amd64.tar.gz" -C "$PREFIXE"
fi

# Le fichier d'environnement est limité à la session qui le charge.
{
  printf '# Environnement local mlpack 4.8.0 / Go 1.27.3.\n'
  printf 'export PATH=%q:"$PATH"\n' "$PREFIXE/go/bin"
  printf 'export GOROOT=%q\n' "$PREFIXE/go"
  printf 'export CGO_ENABLED=1\nexport GOTOOLCHAIN=local\n'
  printf 'export CGO_LDFLAGS=%q\n' "-L$PREFIXE/lib -Wl,-rpath,$PREFIXE/lib"
  printf 'export PKG_CONFIG_PATH=%q${PKG_CONFIG_PATH:+:"$PKG_CONFIG_PATH"}\n' "$PREFIXE/lib/pkgconfig"
} > "$PREFIXE/env.sh"
# shellcheck source=/dev/null
. "$PREFIXE/env.sh"
[[ "$(go version)" == "go version go${VERSION_GO} linux/amd64" ]] || erreur "La version effective de Go est inattendue."

telecharger "https://www.mlpack.org/files/mlpack-${VERSION_MLPACK}.tar.gz" "$CACHE/mlpack-${VERSION_MLPACK}.tar.gz" "$SHA_MLPACK"
if [[ ! -f "$CACHE/mlpack-${VERSION_MLPACK}/CMakeLists.txt" ]]; then
  tar -xzf "$CACHE/mlpack-${VERSION_MLPACK}.tar.gz" -C "$CACHE"
fi

# Vérifier et appliquer les correctifs correspondant aux wrappers du dépôt.
CORRECTIF="$DOSSIER_SCRIPT/patches/mlpack-4.8.0-go-memory.patch"
printf '%s  %s\n' "$SHA_CORRECTIF" "$CORRECTIF" | sha256sum --check --status \
  || erreur "Le correctif C++ ne correspond pas à cette version du script."
if patch -d "$CACHE/mlpack-${VERSION_MLPACK}" -p1 --batch --forward --dry-run < "$CORRECTIF" > "$CACHE/verification-correctif.log" 2>&1; then
  patch -d "$CACHE/mlpack-${VERSION_MLPACK}" -p1 --batch --forward < "$CORRECTIF"
elif patch -d "$CACHE/mlpack-${VERSION_MLPACK}" -p1 --batch --reverse --dry-run < "$CORRECTIF" >> "$CACHE/verification-correctif.log" 2>&1; then
  printf 'Correctif C++ déjà appliqué et vérifié.\n'
else
  erreur "Sources C++ incompatibles avec le correctif. Voir $CACHE/verification-correctif.log."
fi

# Le CMake amont réécrit ce fichier même quand son contenu reste identique.
# Conserver sa date dans ce cas évite de recompiler toutes les méthodes à la reprise.
CONFIG_LOCALE="$CACHE/build/include/mlpack/config-local.hpp"
CONFIG_PRECEDENTE="$CACHE/config-local.precedent.hpp"
if [[ -f "$CONFIG_LOCALE" ]]; then
  cp -p -- "$CONFIG_LOCALE" "$CONFIG_PRECEDENTE"
else
  rm -f -- "$CONFIG_PRECEDENTE"
fi
cmake -S "$CACHE/mlpack-${VERSION_MLPACK}" -B "$CACHE/build" \
  -DCMAKE_BUILD_TYPE=Release \
  -DCMAKE_INSTALL_PREFIX="$PREFIXE" \
  -DCMAKE_INSTALL_LIBDIR=lib \
  '-DCMAKE_INSTALL_RPATH=$ORIGIN' \
  -DARMADILLO_INCLUDE_DIR=/usr/include \
  -DARMADILLO_LIBRARY="$LIB_ARMADILLO" \
  -DENSMALLEN_INCLUDE_DIR=/usr/include \
  -DCEREAL_INCLUDE_DIR=/usr/include \
  -DBUILD_SHARED_LIBS=ON \
  -DDOWNLOAD_DEPENDENCIES=OFF \
  -DBUILD_TESTS=OFF \
  -DBUILD_CLI_EXECUTABLES=OFF \
  -DBUILD_PYTHON_BINDINGS=OFF \
  -DBUILD_JULIA_BINDINGS=OFF \
  -DBUILD_R_BINDINGS=OFF \
  -DBUILD_MARKDOWN_BINDINGS=OFF \
  -DBUILD_GO_BINDINGS=OFF \
  -DBUILD_GO_SHLIB=ON \
  -DUSE_OPENMP=ON
if [[ -f "$CONFIG_PRECEDENTE" ]] && cmake -E compare_files "$CONFIG_PRECEDENTE" "$CONFIG_LOCALE"; then
  touch -r "$CONFIG_PRECEDENTE" "$CONFIG_LOCALE"
fi
rm -f -- "$CONFIG_PRECEDENTE"
cmake --build "$CACHE/build" --parallel "$JOBS"
cmake --install "$CACHE/build"

# L'exemple utilise le module corrigé du checkout courant, sans renvoi amont.
(
  cd -- "$DOSSIER_SCRIPT"
  export GOWORK=off
  export GOFLAGS=
  go mod download
  go mod verify
  [[ "$(go list -m -f '{{.Path}}')" == "$MODULE_GO" ]] || erreur "Module Go inattendu."
  go build -mod=readonly -buildvcs=false -trimpath -o "$PREFIXE/bin/kmeans-demo" ./examples/kmeans
  OMP_NUM_THREADS=2 OPENBLAS_NUM_THREADS=1 "$PREFIXE/bin/kmeans-demo" | tee "$PREFIXE/share/audit/kmeans.txt"
  sha256sum ./*.go go.mod go.sum capi/*.h examples/kmeans/main.go > "$PREFIXE/share/audit/sources-go.sha256"
)

cp -- "$CACHE/build/install_manifest.txt" "$PREFIXE/share/audit/install_manifest.txt"
{
  date -u '+Date UTC : %Y-%m-%dT%H:%M:%SZ'
  printf 'mlpack : %s\nModule Go : %s\nCorrectif SHA-256 : %s\n' "$VERSION_MLPACK" "$MODULE_GO" "$SHA_CORRECTIF"
  go version
  "${CXX:-c++}" --version | head -n 1
  cmake --version | head -n 1
  dpkg-query -W -f='${binary:Package}\t${Version}\n' \
    libarmadillo-dev libensmallen-dev libcereal-dev
  readelf -d "$PREFIXE/lib/libmlpack_go_kmeans.so"
} > "$PREFIXE/share/audit/versions-et-liaison.txt"
printf 'mlpack-go-milvi 4.8.0 go1.27.3\n' > "$PREFIXE/INSTALLATION-VERIFIEE"
printf '\nInstallation et vérification k-means réussies.\nPour la session courante :\n  source %q\nPour relancer la vérification :\n  %q\n' "$PREFIXE/env.sh" "$PREFIXE/bin/kmeans-demo"
