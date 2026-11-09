# Maintenir les deux modules k-means

Les deux branches sont `codex/kmeans-only`. Les modules restent
`github.com/milvi/mlpack` pour le C++ et `github.com/milvi/mlpack-go` pour l’interface
utilisateur. Le module natif est `v0.2.0-kmeans.2` et le module utilisateur
est `v0.2.0-kmeans.3`, distincts de mlpack amont 4.8.0.
Les tags publiés sont immuables ; le développement avance par nouveaux commits.

1. Modifier directement le C++ dans `milvi/mlpack`, puis lancer
   `python3 tools/check-source.py`. Tout en-tête ajouté doit être réellement
   nécessaire. Les anciens générateurs et modèles persistants sont retirés.
2. Pour travailler sur les deux checkouts voisins, créer localement un `go.work`
   avec `go work init ./mlpack ./mlpack-go`. Ce fichier ne fait pas partie de la
   distribution. Le module utilisateur publié doit fonctionner avec `GOWORK=off`.
3. Lancer les tests Go ordinaires, les contrôles `cgocheck2`, les restitutions
   natives et le pilote C++ avec injection de pannes et détection des fuites.
   Les commandes de compilation et de test sont bornées à 120 secondes dans
   les Dockerfiles et les outils ; les tests Go ont aussi leur propre plafond.
4. Publier une nouvelle version sémantique immuable du module C++, puis exécuter
   `GOWORK=off go get github.com/milvi/mlpack@VERSION` et `go mod tidy` dans le
   module utilisateur. Les empreintes de `go.sum` doivent être produites par Go.
5. Valider Ubuntu 26.04 et un projet neuf, sans `replace` ni environnement privé,
   avant de publier la version correspondante du module utilisateur.

Le job de publication crée la version indiquée dans son workflow après succès
de l’intégration Ubuntu. Les versions des deux modules sont indépendantes. Un tag
déjà présent est conservé ; les publications suivantes choisissent une nouvelle
version. `@latest` ne sélectionne pas cette préversion tant qu’une version stable
existe. Les branches `master` et les versions stables existantes restent disponibles.
