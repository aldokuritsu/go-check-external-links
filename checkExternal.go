package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {

	// Demander à l'utilisateur de saisir l'URL à scanner
	fmt.Println("Entrer l'URL que vous souhaitez scanner")
	var url string
	fmt.Scanf("%s", &url)

	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	// Envoi d'une requête HTTP GET à l'URL et récupération de la réponse
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur liée à la requête :", err)
		return
	}
	defer response.Body.Close()

	// Création d'un parseur HTML à partir de la réponse
	doc, err := html.Parse(response.Body)
	if err != nil {
		fmt.Println("Erreur lors du parsing html :", err)
		return
	}

	/******* Fonction pour parcourir l'arbre DOM et récupérer tous les éléments <a> *******/

	// Déclaration d'un tableau de chaînes de caractères qui va contenir tous les liens
	var links []string

	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && a.Val != "/" && !strings.HasPrefix(a.Val, "/") && !strings.HasPrefix(a.Val, url) && !strings.HasPrefix(a.Val, "#") && !strings.HasPrefix(a.Val, "tel:") && !strings.HasPrefix(a.Val, "mailto:") {
					links = append(links, a.Val)
					break
				}
			}
		}
	}
	forEachNode(doc, visitNode, nil)

	// Ouverture d'un fichier CSV en écriture
	file, err := os.Create("links.csv")
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	// Création d'un writer CSV à partir du fichier
	writer := csv.NewWriter(file)

	// Si des liens externes ont été détectés, on les écrit dans un fichier csv
	if len(links) == 0 {
		fmt.Println("Pas de lien externe détecté !")
	} else {
		for _, link := range links {
			err := writer.Write([]string{link})
			if err != nil {
				fmt.Println("Error writing to CSV file:", err)
				return
			}
		}
	}
	// Vide le tampon et s'assure que toutes les données ont été écrites sur le disque
	writer.Flush()
}

// forEachNode appelle la fonction passée en premier argument pour chaque noeud de l'arbre DOM,
// en commençant par le noeud racine passé en deuxième argument.
// Si la fonction passée en troisième argument renvoie true, les enfants du noeud courant ne sont pas visités.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil {
		post(n)
	}
}
