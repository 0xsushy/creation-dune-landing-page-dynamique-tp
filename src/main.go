package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

type Product struct {
	Id          int
	Name        string
	Price       float64
	PathImage   string
	ReducePrice float64
	HasReduce   bool
}

var listeProducts = []Product{
	{1, "Product 1", 19.99, "/static/img/products/image1.webp", 14.99, false},
	{2, "Product 2", 19.99, "/static/img/products/image2.webp", 140, true},
	{3, "Product 3", 19.99, "/static/img/products/image3.webp", 122.3, false},
	{4, "Product 4", 19.99, "/static/img/products/image4.webp", 19.99, false}}

func main() {
	listTemplates, errTemplates := template.ParseGlob("./templates/*.html")
	if errTemplates != nil {
		fmt.Println("Erreur chargement Template:", errTemplates)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		listTemplates.ExecuteTemplate(w, "index", listeProducts)
	})

	http.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		idProduit := r.FormValue("id")
		produitId, produiterr := strconv.Atoi(idProduit)
		if produiterr != nil {
			http.Error(w, "Erreur l'id du produit n'est pas bon", http.StatusBadRequest)
			return
		}

		for _, product := range listeProducts {
			if product.Id == produitId {
				listTemplates.ExecuteTemplate(w, "produit", product)
				break
			}
		}

		http.Error(w, "Produit non trouvé", http.StatusNotFound)
	})

	http.HandleFunc("/liste", func(w http.ResponseWriter, r *http.Request) {
		listTemplates.ExecuteTemplate(w, "liste", nil)
	})

	fileserver := http.FileServer(http.Dir("./assets"))
	http.Handle("/static/", http.StripPrefix("/static/", fileserver))

	errServer := http.ListenAndServe(":8000", nil)
	if errServer != nil {
		fmt.Printf("Erreur lancement de serveur")
		os.Exit(1)
	}
}
