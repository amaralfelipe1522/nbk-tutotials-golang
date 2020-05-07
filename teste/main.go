package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

//Book armazena os dados do livros
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var b Book

func main() {
	db, err := sql.Open("mysql", "root:Project@1522@/nbktutorial")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, _ := db.Begin()

	rows, err1 := tx.Query("select * from books;")
	if err1 != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var bookList []Book

	for rows.Next() {
		rows.Scan(&b.ID, &b.Title, &b.Author)
		bookList = append(bookList, b)
	}
	fmt.Println("Result:", bookList)
	json, _ := json.Marshal(bookList)
	//fmt.Println(string(json))
	fmt.Println("Result:", string(json))
}
