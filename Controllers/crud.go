package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"goTutorial/src/Models"
	"log"
	"net/http"
	"strconv"
)

// Init books var as a slice Book struct
var books []Models.Book

// Get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConnect()
	json.NewEncoder(w).Encode(db.Find(&books))
}
func dbConnect() *gorm.DB {
	db, err := gorm.Open("mysql", "root:@(localhost:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		//defer db.Close()
		fmt.Println(err)
	}
	return db

}
func migration() {
	db := dbConnect()
	db.AutoMigrate(&Models.Book{})
	db.AutoMigrate(&Models.Author{})
}

// Get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// Gets params
	// Loop through books and find one with the id from the params
	for _, item := range books {
		id, err1 := strconv.ParseUint(r.FormValue("id"), 10, 32)
		if err1 != nil {
			fmt.Println(err1)
		}
		if uint(id) == item.ID {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Models.Book{})
}

// Add new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var book Models.Book
	err := decoder.Decode(&book)
	if err != nil {
		panic(err)
	}
	log.Println("new book :" , book)


	db := dbConnect()
	db.NewRecord(book)
	db.Create(&book)
	_ = json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

// Update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		id, err1 := strconv.ParseUint(params["id"], 10, 32)
		if err1 != nil {
			fmt.Println(err1)
		}
		if item.ID == uint(id) {
			books = append(books[:index], books[index+1:]...)
			var book Models.Book
			_ = json.NewDecoder(r.Body).Decode(&book)
			book.ID = uint(id)
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

// Delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range books {
		id, err1 := strconv.ParseUint(params["id"], 10, 32)
		if err1 != nil {
			fmt.Println(err1)
		}
		if item.ID == uint(id) {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

// Main function
func main() {
	// Init router
	migration()
	r := mux.NewRouter()

	// Hardcoded data - @todo: add database
	books = append(books, Models.Book{ Isbn: "438227", Title: "Book One"})
	books = append(books, Models.Book{ Isbn: "454555", Title: "Book Two"})
	books = append(books, Models.Book{ Isbn: "454555", Title: "Book Two"})
	// Route handles & endpoints
	r.HandleFunc("/books", getBooks).Methods("GET")
	r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", createBook).Methods("POST")
	r.HandleFunc("/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", deleteBook).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":1230", r))
}
