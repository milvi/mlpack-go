# Corrections et validation du 8 novembre 2026

Le module utilisateur est `github.com/milvi/mlpack-go@v0.1.0` ; il dépend de
`github.com/milvi/mlpack@v0.1.0`, basé sur mlpack C++ 4.8.0. Le tag C++ désigne
[`cb3d9db1842f20904d24cc01b71dae2ce9f21ded`](https://github.com/milvi/mlpack/commit/cb3d9db1842f20904d24cc01b71dae2ce9f21ded).
Les versions de module suivent les règles Go ; elles sont distinctes de la version amont C++.

L’installation utilise APT puis `go get` et `go build`. Le script d’installation,
le patch de déploiement et le SDK Go privé sont retirés. Les changements C++ sont
committés dans leur fork ; aucun patch n’est appliqué sur la machine qui compile.
Les anciennes mesures restent dans [l’historique](corrections-historiques-2026-11-02.md).

## Corrections effectives

| Problème | Correction et vérification |
|---|---|
| Sorties Armadillo détachées sans propriétaire | Copie vers Go ; matrices natives détruites avec les paramètres. Les anciens getters de pointeurs transférés sont supprimés. |
| Chaînes C non restituées | Chaque `C.CString` est libéré dans la même portée ; 50 000 cycles de paramètres contrôlés. |
| Matrices Gonum avec pas non contigu | Copie ligne par ligne lorsque nécessaire ; cas de vues testé. |
| Un modèle pouvait garder la mémoire Go après désépinglage | Les données réelles sont copiées dans une mémoire possédée en C++ ; les conversions entières possèdent leur résultat. Le pinning persistant et le setter générique de pointeur brut sont supprimés. |
| Modèle emprunté détruit sur erreur | Un propriétaire local protège les allocations des 26 bindings créateurs de modèles ; les entrées restent empruntées. Réutilisation d’un modèle après erreur testée. |
| Modèle créé puis perdu avant export | Propriétaire local jusqu’au transfert dans les paramètres ; contrôle des appels qui échouent après entraînement. |
| Modèles opaques sans durée de vie Go | Propriétaire partagé entre alias, `Close` idempotent et nettoyage automatique de secours avec registre faible. |
| Arbres et temporaires dual-tree possédés par pointeurs bruts | Propriétaires uniques, temporaires automatiques et restauration des enfants masqués si le parcours lève. |
| Constructeurs d’arbres interrompus | Nettoyage des enfants et de la matrice possédée ; copies et déplacements corrigés, parents rétablis. |
| Bornes rectangulaires | L’affectation par déplacement libère l’ancien tableau ; les réallocations précèdent la libération de l’ancien stockage. |
| Allocation ou exception dans une région OpenMP naïve | Accumulateurs alloués avant la région ; exceptions collectées puis relancées sur le thread appelant ; ordre de réduction fixe. |
| Exceptions traversant la frontière C | Interception des exceptions aux entrées C, tampon d’erreur borné. `Cluster` rend `error` ; signatures historiques conservées avec `panic(error)`. |
| Vérifications d’un autre algorithme dans les fonctions inline communes | Elles utilisent désormais les paramètres de l’appel. Les journaux sont partagés dans le module monolithique. |
| En-têtes C++ imbriqués ignorés par le cache cgo | Entrées explicites `go:embed`, sans conservation du système de fichiers dans l’exécutable ; essai de modification d’en-tête exécuté. |

`Cluster` fait un seul appel de calcul natif. Le C++ copie les entrées, construit
ses objets, écrit dans les slices Go et détruit ses ressources avant de revenir.
Les centres vides restent en place pendant les itérations, puis sont retirés à
la fin par défaut. Les étiquettes sont renumérotées. L’affectation finale utilise
une recherche exacte avec KD-tree et conserve l’ordre initial des points.
Elle évite le balayage final n × k de l’interface k-means historique.
Le cas d’un centre unique est aussi corrigé dans l’ancienne implémentation dual-tree.

## Preuves exécutées

Les [tests Go](validation/tests-go.log) couvrent les conversions communes,
les vues Gonum, l’indépendance des sorties, les cas dégénérés, la parité avec le
Lloyd naïf sur dix jeux de points en dimensions 2 et 3, les moyennes des centres,
les erreurs C++, l’appel suivant une erreur, la concurrence, OpenMP à deux threads,
les alias de modèles et la conservation d’un modèle après modification des
entrées Go et passage du ramasse-miettes.

Les [mesures de restitution](validation/memoire.log), après échauffement et
ramasse-miettes Go, donnent :

| Cas | Variation native mesurée |
|---|---:|
| 50 appels `Cluster` de 100 000 points | +6 384 octets |
| 50 appels k-means historiques de 100 000 points | +9 072 octets |
| 1 000 modèles créés puis fermés | +6 528 octets |
| 100 modèles dont l’appel échoue avant export | +5 040 octets |
| 50 000 cycles de paramètres et chaînes | +208 octets |

Le seuil de non-régression de ces tests est de 2 Mio. Ces variations reflètent
aussi les caches de l’allocateur ; elles ne sont pas des estimations d’un débit
de fuite. L’ancien défaut des étiquettes laissait plus de 40 Mo après 50 appels.

Le pilote C++ indépendant vérifie les restitutions après **chaque** scénario :
390 pannes d’allocation interceptées sur k-means, 639 sur la construction ou
l’affectation des arbres et 2 400 scénarios sur un modèle via l’interface C,
avec 1 378 erreurs attendues ou allocations interrompues. Le nombre de `new`
actifs revient exactement au niveau initial. Les journaux locaux sont
[conservés](validation/asan-local.log).

**Ubuntu 26.04 amd64** : [validation native et publication du tag](https://github.com/milvi/mlpack/actions/runs/42460225318)
réussies, avec compilation et test Go sous un utilisateur ordinaire, puis
AddressSanitizer (ASan, accès mémoire), UndefinedBehaviorSanitizer (UBSan,
comportements indéfinis) et LeakSanitizer (LSan, allocations perdues).
Le contrôle LSan utilise une capacité `SYS_PTRACE` dans son conteneur de diagnostic.
Le [journal](validation/ci-native-ubuntu2604.log) enregistre Go 1.27.3, GCC 15.2.0,
Armadillo 15.2.1, ensmallen 3.11.0 et cereal 1.3.2.

Le [workflow Go](../.github/workflows/ubuntu2604.yml) vérifie en plus les sommes
des modules, tous les tests, `GOEXPERIMENT=cgocheck2` (contrôles de pointeurs
renforcés), les restitutions, l’exemple et le benchmark. Son job de publication
crée le tag seulement après succès de ces étapes. La branche de publication
est `codex/publication-v0.1.0` ; les résultats de son exécution restent consultables
[dans GitHub Actions](https://github.com/milvi/mlpack-go/actions/workflows/ubuntu2604.yml).

Le laboratoire local utilise Debian 13.6 avec les paquets APT extraits dans un
sysroot ; aucune installation système n’y a été effectuée. Ses variables CGO
temporaires servent exclusivement à ce laboratoire. Le déploiement Ubuntu testé
utilise les chemins système habituels et ne charge aucun `env.sh`.
Les [versions et conditions du laboratoire](validation/laboratoire.json) sont archivées.
LSan y est indisponible à cause de la restriction ptrace : le contrôle de fuite
LSan réellement exécuté est celui d’Ubuntu ci-dessus.

## Cache réellement vérifié

| Action | Compilations C++ lancées | Durée |
|---|---:|---:|
| Build inchangé, cache chaud | 0 | 0.157 s |
| Seul l’horodatage d’un en-tête change | 0 | 0.170 s |
| Ajout d’un commentaire dans un en-tête imbriqué | 49 | 437.2 s |
| Contenu initial restauré | 0 | 1.077 s |

Le [rapport brut](validation/cache.json) et le [contrôle reproductible](../tools/verify-cache.py)
sont fournis. Le témoin de source recherché est absent de l’exécutable de
25997776 octets. Les en-têtes système APT restent extérieurs au cache Go ;
après leur mise à jour, forcer `go build -a ./...`.
`go.sum` vérifie les sources distribuées et les `go.mod` lors du téléchargement,
avec le journal public de vérification ; ces empreintes sont distinctes des clés
utilisées par le cache de compilation. Le [module C++ téléchargé](validation/module-cpp-publie.json)
porte bien le tag v0.1.0 et le commit annoncé.

## Temps et mémoire : points uniformes dans le carré unité

Mesure locale sur AMD EPYC 9V74, Go 1.27.3, GCC 14.2, Armadillo 14.2.3,
ensmallen 2.22.1 et cereal 1.3.2. Un thread OpenMP, trois graines (42, 43, 44),
un processus neuf par essai. Le temps couvre l’appel complet, y compris copies,
arbres, itérations et affectation finale ; la génération des points est exclue.
Le pic RSS (*Resident Set Size*, mémoire résidente du processus) inclut Go et C++.
Les statistiques du tas Go après l’appel ne mesurent pas la mémoire C++.

**Budget fixé à 100 itérations pour comparer le coût**, avec arrêt anticipé si
le seuil global de déplacement des centres 10⁻⁵ est atteint. Le tableau ne
présente donc pas tous les résultats comme convergés.

| n | k | Temps médian (s) | Pic RSS médian (Mio) | Itérations par graine | Convergence par graine |
|---:|---:|---:|---:|---|---|
| 100000 | 1000 | 1.536 | 29.01 | 100, 100, 95 | non, non, oui |
| 1000000 | 1000 | 8.683 | 169.30 | 100, 100, 100 | non, non, non |
| 1000000 | 10000 | 21.271 | 170.29 | 100, 100, 100 | non, non, non |

Pour n multiplié par 10 :

- k fixé à 1 000 : temps ×5.65, pic RSS ×5.84, temps moyen par itération ×5.59.
- k = n/100 : temps ×13.85, pic RSS ×5.87, temps moyen par itération ×13.70.

Ces quotients sont empiriques. Le second modifie aussi k ; aucun des deux ne
prouve une complexité asymptotique. Le nombre d’itérations I et les couples de
nœuds visités dépendent des données et de l’initialisation. Un parcours dual-tree
peut perdre son avantage de pruning ; une garantie globale O(n log n) ne découle
pas de ces mesures. La mémoire observée est compatible avec un coût linéaire,
plus une partie fixe, sur ces deux tailles ; deux tailles ne prouvent pas une borne.

Un [essai complémentaire à la limite par défaut de 1 000 itérations](validation/limite1000.json),
n = 1 000 000, k = 1 000, graine 42, rend encore `Converged: false`.
Il faut lire cet indicateur et choisir le budget ou le seuil selon l’usage.
Les points uniformes n’offrent pas de groupes nettement séparés ; les centres
peuvent continuer à se déplacer faiblement. Le résultat conserve des affectations
exactes vers ses centres, même si l’arrêt n’est pas encore un point fixe de Lloyd.

[Mesures JSON complètes](validation/mesures-debian.json) et
[outil de mesure](../tools/benchmark.py). Le CLI accepte n libre :

```sh
go run ./cmd/kmeans -n 1000000 -k 10000 -iterations 100
```

## Portée et maintenance

Le chemin k-means moderne et les utilitaires communs reçoivent les contrôles
ci-dessus. Les 46 méthodes historiques sont compilées et leurs allocations de
modèles principales sont protégées. Tous les chemins internes de tous les
algorithmes et toutes les entrées possibles n’ont pas subi une injection exhaustive.
Les modèles historiques peuvent rester mutables ; l’API `Cluster` conserve ses
entrées et ne garde pas d’état de calcul entre appels.

La première compilation reste coûteuse tant que les autres méthodes sont
conservées. L’extraction de k-means et une éventuelle réécriture en Go pur
constituent les étapes suivantes ; elles ne font pas partie de cette publication.
Les tags publiés ne sont jamais déplacés. Le travail continue par nouveaux
commits et nouveaux tags, avec rebase limité aux commits de travail non publiés.
