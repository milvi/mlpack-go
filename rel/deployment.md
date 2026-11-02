# Régénérer la distribution Go corrigée

La source des wrappers est le dépôt [milvi/mlpack](https://github.com/milvi/mlpack),
avec les corrections Go de ce fork. Le dépôt présent ajoute le module sous
`github.com/milvi/mlpack-go`, l'installateur, ses tests et sa documentation.

Ne pas exécuter l'ancienne procédure qui effaçait tout le dépôt : elle supprimerait
ces ajouts. Ne pas régénérer depuis une source C++ non corrigée, qui réintroduirait
les défauts de propriété mémoire.

1. Configurer la source C++ corrigée avec `BUILD_GO_BINDINGS=ON`,
   `FORCE_BUILD_GO_BINDINGS=ON`, `BUILD_GO_SHLIB=ON`, une version de Go ≥ 1.21
   et les dépendances natives appropriées. Pour cette distribution, utiliser
   Go 1.27.3 et mlpack 4.8.0.
2. Construire la cible `go` de CMake. Elle génère les fichiers Go et les
   bibliothèques natives dans
   `build/src/mlpack/bindings/go/src/mlpack.org/v1/mlpack/`.
3. Copier uniquement les fichiers de wrappers `.go` et les déclarations
   `capi/*.h` utiles. Conserver les tests, exemples, documentation, installateur,
   workflow, Dockerfile, Makefile et le module du présent fork.
4. Vérifier que `arma_util.go`, `io_util.go` et les déclarations C correspondent
   à la source corrigée. Le module généré dans le dépôt C++ conserve le chemin
   amont ; le module distribué ici doit rester `github.com/milvi/mlpack-go`.
   Adapter aussi le commentaire d'import canonique de `doc.go` à ce chemin.
5. Si le correctif natif change, reconstruire le patch depuis la base 4.8.0,
   vérifier son application à une archive neuve et mettre à jour le SHA-256
   attendu dans `installer.sh`.
6. Exécuter les tests décrits dans le README, puis construire le Dockerfile
   depuis une base propre. Ne jamais committer les `.so`, archives ou caches.

La distribution de wrappers et l'installation native sont donc deux étapes.
L'installation normale n'active pas le générateur : les wrappers présents ici
sont déjà générés, tandis que les 47 bibliothèques natives sont reconstruites.
