package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func showBooks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Books)
}

func insertBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated) //201
	//Inserindo livro a partir de uma nova variável
	var newBook = Book{
		ID:     len(Books) + 1,
		Title:  "Laranja Mecanica",
		Author: "Anthony Burgess",
	}
	Books = append(Books, newBook)

	/* Inserindo livro pelo Body da requisição:
	{
		"title":  "Dragões de Ether",
		"author": "Raphael Dracoon"
	}
	*/
	rBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(rBody, &newBook)
	//Gambiarra para autoincremento do ID
	newBook.ID = len(Books) + 1
	Books = append(Books, newBook)

	json.NewEncoder(w).Encode(Books)
}

func findBook(w http.ResponseWriter, r *http.Request) {
	// Utiliza o bookID presente na URL
	vars := mux.Vars(r)["bookID"]
	id, _ := strconv.Atoi(vars)

	if id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for _, book := range Books {
		if book.ID == id {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["bookID"]
	id, _ := strconv.Atoi(vars)

	if id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	for index, book := range Books {
		if book.ID == id {
			fmt.Fprintf(w, "O livro %s foi deletado", book.Title)
			Books = append(Books[:index], Books[index+1:]...)
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["bookID"]
	id, _ := strconv.Atoi(vars)
	if id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var newBook Book
	// armazena o JSON do body no newBook para atualizar o Books
	rBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(rBody, &newBook)

	for index, book := range Books {
		if book.ID == id {
			fmt.Fprintf(w, "O livro %s foi atualizado", book.Title)
			Books[index].Title = newBook.Title
			Books[index].Author = newBook.Author
			return
		}
	}
	w.WriteHeader(http.StatusNoContent)
}

func confHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func confHandler(router *mux.Router) {
	router.HandleFunc("/", mainHandler)
	router.HandleFunc("/books", showBooks).Methods("GET")
	router.HandleFunc("/books", insertBook).Methods("POST")
	router.HandleFunc("/books/{bookID}", findBook).Methods("GET")
	router.HandleFunc("/books/{bookID}", deleteBook).Methods("DELETE")
	router.HandleFunc("/books/{bookID}", updateBook).Methods("PUT")
}

func confServer() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(confHeader)
	confHandler(router)
	fmt.Println("Rodando na porta 3000.")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func main() {
	confServer()
}
