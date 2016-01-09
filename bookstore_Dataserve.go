// Data server for bookstore app
// borrows from and extends: http://www.alexedwards.net/blog/practical-persistence-sql

package main

import (
	"./models"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	models.InitDB("postgres://test:pass@localhost/Bookstore")

	http.HandleFunc("/books", booksIndex)
	http.HandleFunc("/books/show", booksShow)
	http.HandleFunc("/books/create", booksCreate)
	http.HandleFunc("/books/delete", booksDelete)

	fmt.Println("Bookstore: dataserver (port:4000)")
	http.ListenAndServe(":4000", nil)
}

// return an index of all Books
func booksIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	bks, err := models.GetBooks()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//  render json
	b, err := json.Marshal(bks)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))

}

// return a subset of Book records
func booksShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	// use isbn as a filter
	isbn := r.FormValue("isbn")
	if isbn == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	bks, err := models.GetBooks(isbn)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//  render json
	b, err := json.Marshal(bks)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))

}

// create a new Book record
func booksCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")
	title := r.FormValue("title")
	author := r.FormValue("author")
	if isbn == "" || title == "" || author == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	price, err := strconv.ParseFloat(r.FormValue("price"), 32)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	rowsAffected, err := models.CreateBook(isbn, title, author, price)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// output confirmation to console
	fmt.Printf("Book %s created successfully (%d row affected)\n", isbn, rowsAffected)
}

// deletes a Book record
func booksDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")

	if isbn == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	rowsAffected, err := models.DeleteBook(isbn)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// output confirmation to console
	fmt.Printf("Book %s deleted successfully (%d row affected)\n", isbn, rowsAffected)
}
