# Schéma quiver

[Ouvrir le schéma modifiable dans quiver](https://q.uiver.app/#q=WzAsNixbMCwwLCJcXHRleHR7UHJvamV0IEdvfSJdLFsxLDAsIlxcdGV4dHtNb2R1bGVzIEdvIGV0IEMrK30iXSxbMCwxLCJcXHRleHR7Q29tcGlsYXRldXJzIEdvIC8gQysrfSJdLFsxLDEsIlxcdGV4dHtBUFR9Il0sWzAsMiwiXFx0ZXh0e0NhY2hlIEdvfSJdLFsxLDIsIlxcdGV4dHtFeMOpY3V0YWJsZX0iXSxbMCwxLCJcXHRleHR7Z28gZ2V0fSJdLFswLDIsIlxcdGV4dHtnbyBidWlsZH0iXSxbMSwyLCJcXHRleHR7c291cmNlc30iXSxbMywyLCJcXHRleHR7ZW4tdMOqdGVzIC8gYmlibGlvdGjDqHF1ZXN9Il0sWzIsNCwiXFx0ZXh0e29iamV0c30iXSxbNCw1LCJcXHRleHR7bGllbn0iXSxbMyw1LCJcXHRleHR7YmlibGlvdGjDqHF1ZXMgZHluYW1pcXVlc30iLDFdXQ==).

Les deux modules sont téléchargés par Go. Les en-têtes APT servent à compiler ;
les bibliothèques dynamiques APT servent aussi à exécuter. Le cache Go conserve
les objets compilés. Ce dessin montre des flux de construction ; il ne s’agit
pas d’un diagramme commutatif au sens mathématique. Mermaid reste plus adapté
aux processus d’installation, quiver à leur mise en page manuelle.

Contrôle du 8 novembre 2026 : le navigateur de test a affiché « KaTeX failed to load » et n’a pas rendu le schéma, même après rechargement. Le rendu de cette proposition dans quiver reste donc non validé ; le schéma Mermaid du README est la référence.

Le format v0 a été vérifié dans le parseur officiel `src/quiver.mjs` ; son URL utilise la base64 brute, comme l’export de quiver.
