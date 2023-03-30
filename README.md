# EASYACCOUNTING

Projet pour m'aider à faire ma comptabilité et apprendre la programmation en Golang.
Il a pour but d'aider à la visualisation de ses dépenses, en reformatant un fichier tableur.

Le programme traite un fichier csv récupérable sur le site de sa banque. C'est un fichier contenant les transactions effectuées ce mois-ci.
Le programme édite les données et ajoute des colonnes puis sort le fichier sous forme de csv.

## Prérequis
Le programme n'est compatible qu'avec le format de fichier fourni par la Société Générale.
 
## Fonctionnement

Pour exécuter le programme,

- Créer un dossier "input" et "output" au même emplacement où se trouve votre binaire
- Placer le fichier de transactions dans le dossier input
- Lancer le programme,  `./easyaccounting` sous linux, double-click sous windows
- Entrer le nombre identifiant la catégorie de la transaction.
- Le programme se ferme et génère un fichier du même nom dans le dossier output

## Exemples de fichiers

Fichier attendu en entrée, récupéré d'un export de son compte sur le site de la Société Générale.
| **Date de l'opération** | **Libellé**        | **Détail de l'écriture**                                 | **Montant de l'opération** | **Devise** |
|-------------------------|--------------------|----------------------------------------------------------|----------------------------|------------|
| 29/12/2022              | CARTE XCCCC RETRAI | CARTE XCCCC RETRAIT DAB SG 28/12 16H44 PARIS ST AMBROISE | -10,00                     | EUR        |


Fichier généré en sortie
| **Date transaction** | **Date prelevement** | **Label**                                    | **Montant** | **Type**   |
|----------------------|----------------------|----------------------------------------------|-------------|------------|
| NULL                 | 29/12/2022           | RETRAIT DAB SG 28/12 16H44 PARIS ST AMBROISE | -10,00      | Nourriture |

## TODO
X Sauvegarde fichier en cours pour pouvoir s'arrêter et ne pas perdre le travail actuel
- Résoudre bug chargement fichier => 1ère ligne n'est pas un csv
- Choisir soi même le fichier d'input
- Choisir le lieu de sauvegarde
- Rendre le projet paramétrable, on donne un fichier de config avec les traitements à faire, colonnes à supprimer, nom des colonnes désirés...
- Une vrai interface pour faciliter l'utilisation