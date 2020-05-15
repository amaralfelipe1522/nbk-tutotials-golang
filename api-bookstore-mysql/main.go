package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Book armazena os dados do livros
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func showBooks(w http.ResponseWriter, r *http.Request) {
	db := confDB()
	defer db.Close()

	op, _ := db.Begin()

	rows, err := op.Query("select * from books;")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookList []Book

	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			log.Fatal(err)
		}
		bookList = append(bookList, book)
	}
	json.NewEncoder(w).Encode(bookList)
}

func insertBook(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	db := confDB()
	defer db.Close()

	op, _ := db.Begin()

	////	Inserindo livro a partir de uma nova variável
	// var newBook = Book{
	// 	Title:  "Laranja Mecanica",
	// 	Author: "Anthony Burgess",
	// }
	// stmt, _ := op.Prepare("insert into books (title, author) values (?, ?)")
	// _, err := stmt.Exec(newBook.Title, newBook.Author)
	// if err != nil {
	// 	op.Rollback()
	// 	log.Fatal(err)
	// }
	// op.Commit()

	/* Inserindo livro passado no Body da requisição:
	{
		"title":  "Dragões de Ether",
		"author": "Raphael Dracoon"
	}
	*/
	var newBook Book
	rBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(rBody, &newBook)
	stmt, _ := op.Prepare("insert into books (title, author) values (?, ?)")
	_, err := stmt.Exec(newBook.Title, newBook.Author)
	if err != nil {
		op.Rollback()
		log.Fatal(err)
	}

	op.Commit()

	fmt.Fprintf(w, "Livro inserido com sucesso:")
	json.NewEncoder(w).Encode(newBook.Title)
}

func findBook(w http.ResponseWriter, r *http.Request) {
	// Utiliza o bookID presente na URL
	vars := mux.Vars(r)["bookID"]
	id, _ := strconv.Atoi(vars)

	if id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db := confDB()
	defer db.Close()

	op, _ := db.Begin()

	var b Book

	err := op.QueryRow("select * from books where id = ?", id).Scan(&b.ID, &b.Title, &b.Author)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(b)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["bookID"]
	id, _ := strconv.Atoi(vars)

	if id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db := confDB()
	defer db.Close()

	op, _ := db.Begin()

	stmt, _ := op.Prepare("delete from books where id = ?")
	result, err := stmt.Exec(id)
	if err != nil {
		fmt.Fprintf(w, "Erro na execução do script.")
		op.Rollback()
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		fmt.Fprintf(w, "ID não encontrado.")
		op.Rollback()
		return
	}

	op.Commit()
	fmt.Fprintf(w, "O livro de código %d foi deletado", id)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)["bookID"]
	id, _ := strconv.Atoi(vars)
	if id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var newBook Book
	// armazena o JSON vindo do body da requisição no newBook
	rBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(rBody, &newBook)

	db := confDB()
	defer db.Close()

	op, _ := db.Begin()

	stmt, _ := op.Prepare("update books set title = ?, author = ? where id = ?")
	result, err := stmt.Exec(newBook.Title, newBook.Author, id)
	if err != nil {
		fmt.Fprintf(w, "Erro na execução do script.")
		op.Rollback()
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		fmt.Fprintf(w, "ID não encontrado.")
		op.Rollback()
		return
	}

	op.Commit()
	fmt.Fprintf(w, "O livro de código %d foi atualizado", id)
}

func confHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
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

func confDB() *sql.DB {
	db, err := sql.Open("mysql", "root:Project@1522@/nbktutorial")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

func main() {
	confServer()
}
