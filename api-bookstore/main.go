package main

// Exemplo baseado no video https://www.youtube.com/watch?v=9sovjfz_loA
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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
func getFunction(w http.ResponseWriter, r *http.Request) {
	confHeader(&w)
	splitURL := strings.Split(r.URL.Path, "/")
	// URL for "/books" ou "/books/"
	if len(splitURL) == 2 || len(splitURL) == 3 && splitURL[2] == "" {
		if r.Method == "GET" {
			showBooks(w, r)
		} else if r.Method == "POST" {
			insertBook(w, r)
		} else {
			fmt.Fprintf(w, "Método não é GET e nem POST.")
			return
		}
	} else if len(splitURL) == 3 || len(splitURL) == 4 && splitURL[3] == "" {
		if r.Method == "GET" {
			findBook(w, r)
		} else if r.Method == "DELETE" {
			deleteBook(w, r)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}

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
	splitURL := strings.Split(r.URL.Path, "/")
	//index 2 do array é o ID passado na URL
	id, _ := strconv.Atoi(splitURL[2])
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
	splitURL := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(splitURL[2])
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

func confHeader(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Content-Type", "application/json")
}

func confHandler() {
	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/books", getFunction)
	http.HandleFunc("/books/", getFunction)
}

func confServer() {
	confHandler()
	fmt.Println("Rodando na porta 3000.")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func main() {
	confServer()
}
