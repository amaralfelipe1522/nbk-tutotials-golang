package main

// Exemplo baseado no video https://www.youtube.com/watch?v=9sovjfz_loA
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

//Book armazena os dados do livros
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

//Books simula um banco de dados de livros
var Books []Book = []Book{
	Book{
		ID:     1,
		Title:  "Guia do Mochileiro das Galáxias",
		Author: "Douglas Adams",
	},
	Book{
		ID:     2,
		Title:  "Guerra do Apocalipse",
		Author: "Eduardo Spohr",
	},
	Book{
		ID:     3,
		Title:  "Deuses Americanos",
		Author: "Neil Gaiman",
	},
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

//Função que define qual ação realizar (showBooks ou insertBooks) a partir do método
func getMethod(w http.ResponseWriter, r *http.Request) {
	confHeader(&w)
	if r.Method == "GET" {
		showBooks(w, r)
	} else if r.Method == "POST" {
		insertBook(w, r)
	} else {
		fmt.Fprintf(w, "Método não é GET e nem POST.")
	}
}

func showBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Books)
}

func insertBook(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Insertando")
}

func confHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Content-Type", "application/json")
}

func confHandler() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/books", getMethod)
	//http.HandleFunc("/books", insertBook)
}

func confServer() {
	confHandler()
	fmt.Println("Rodando na porta 3000.")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	confServer()
}
