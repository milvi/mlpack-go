# Maintenir deux modules avec un historique linéaire

Le module C++ est `github.com/milvi/mlpack`, le module utilisateur est
`github.com/milvi/mlpack-go`. Leurs tags commencent à `v0.1.0` ; mlpack amont reste
4.8.0. Aucun patch n’est appliqué lors de l’installation. Les commits publiés et
les tags existants ne sont pas réécrits.

1. Modifier le C++ dans `milvi/mlpack`. Si l’interface C change, lancer
   `python3 tools/generate-native.py` et conserver les fichiers générés dans Git.
   Les méthodes partagent le même exécutable : les fonctions communes ne doivent
   pas dépendre d’un nom d’algorithme figé par macro. Conserver les contrôles
   d’erreur, de propriété et de cache des en-têtes.
2. Pour travailler simultanément sur les deux dépôts : créer un `go.work` dans
   leur répertoire parent avec `go work init ./mlpack ./mlpack-go`. Ce fichier est
   local au développement ; il n’est pas publié et ne sert pas au déploiement.
3. Vérifier les deux modules, les tests cgo renforcés, les tests de restitution,
   le pilote ASan/UBSan et une construction Docker Ubuntu 26.04. Mettre à jour le
   rapport avec les mesures réelles. Les générateurs CMake conservés doivent
   rendre les mêmes aides Go et les mêmes déclarations C que la distribution.
4. Faire des commits linéaires. Rebaser les commits de travail non publiés si
   nécessaire, puis avancer `master` sans commit de fusion. Ne pas rebaser un tag
   ou l’historique déjà consommé par Go.
5. Publier un nouveau tag sémantique immuable du module C++, puis exécuter dans
   le module Go `GOWORK=off go get github.com/milvi/mlpack@v0.x.y`, `go mod tidy`
   et les tests. `go.sum` doit contenir les empreintes réelles fournies par Go.
6. Publier le tag du module utilisateur après ses tests. Vérifier depuis un
   nouveau projet, avec `GOWORK=off`, sans `replace` et sans environnement mlpack.

Pour la première publication, une branche `codex/publication-v0.1.0` déclenche
le job GitHub Actions qui crée le tag seulement après validation Ubuntu.
Le jeton standard `GITHUB_TOKEN` ne dispose de l’écriture que dans ce job de
publication. Un tag préexistant doit désigner exactement le même commit ; aucun
`push --force` n’est utilisé. Les futures publications changent explicitement
la version et leur condition de déclenchement.

Les wrappers C++ sont générés sans CMake par `tools/generate-native.py` et sont
committés. Les aides Go de conversion/propriété et les générateurs historiques
se trouvent dans `src/mlpack/bindings/go/`. Le module utilisateur conserve ses
fichiers modernes `cluster*.go`, ses tests, sa documentation et ses outils ;
ne jamais remplacer tout son répertoire par une génération historique.
