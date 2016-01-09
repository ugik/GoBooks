// Web server for bookstore app

package main

import (
	"./models"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// dataserver url
var dataserver = "http://127.0.0.1:4000"

func main() {

	http.HandleFunc("/books", bookIndex)
	http.HandleFunc("/books/show", bookShow)
	http.HandleFunc("/books/create_book", create_book)
	http.HandleFunc("/books/delete_book", delete_book)
	http.HandleFunc("/error", booksError)
	http.HandleFunc("/", booksLanding)

	if ping() {
		fmt.Println("Bookstore: webserver (port:3000)")
		http.ListenAndServe(":3000", nil)
	} else {
		fmt.Printf("Dataserver not available %s \n", dataserver)
	}
}

// ping data server
func ping() bool {

	// url with endpoint
	url := fmt.Sprintf("%s", dataserver)

	_, err := http.Get(url)
	if err != nil {
		return false
	}
	return true
}

// get Books from data server
func getBooks(w http.ResponseWriter, r *http.Request, endpoint string) []models.Book {

	// url with endpoint
	url := fmt.Sprintf("%s/%s", dataserver, endpoint)

	// GET JSON data for books
	resp, err := http.Get(url)

	if err != nil {
		http.Redirect(w, r, "/error", http.StatusFound)
	} else {
		defer resp.Body.Close()

		// read json http response
		jsonDataFromHttp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var bks []models.Book

		// unmarshal json into our struct
		err = json.Unmarshal([]byte(jsonDataFromHttp), &bks)
		if err != nil {
			panic(err)
		}
		return bks
	}
	return nil
}

func booksLanding(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("welcome.html")
	t.Execute(w, nil)
}

func booksError(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("error.html")
	t.Execute(w, nil)
}

func bookIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	// render html template
	t, _ := template.ParseFiles("index.html", "base.html")
	t.Execute(w, getBooks(w, r, "books"))

}

func bookShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), 405)
		return
	}

	isbn := r.FormValue("isbn")
	if isbn == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	// render html template
	t, _ := template.ParseFiles("book.html", "base.html")
	t.Execute(w, getBooks(w, r, fmt.Sprintf("books/show?isbn=%s", isbn)))

}

// handle HTTP post to dataserver endpoint
func http_post(endpoint string, data url.Values) {

	req, err := http.NewRequest("POST", endpoint, bytes.NewBufferString(data.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.Status != "200 OK" {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("response Body:", string(body))
	}

}

func create_book(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("new_book.html", "base.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		isbn := r.FormValue("isbn")
		title := r.FormValue("title")
		author := r.FormValue("author")

		if isbn == "" || title == "" || author == "" {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		// confirm price is a float
		_, err := strconv.ParseFloat(r.FormValue("price"), 32)
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		// take the form submitted value for price
		price := r.FormValue("price")

		// package the data for HTTP POST
		data := url.Values{}
		data.Set("isbn", isbn)
		data.Add("title", title)
		data.Add("author", author)
		data.Add("price", price)

		url := fmt.Sprintf("%s/books/create", dataserver)
		http_post(url, data)

		http.Redirect(w, r, "/books", http.StatusFound)
	}

}

func delete_book(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	isbn := r.FormValue("isbn")

	if isbn == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// package the data for HTTP POST
	data := url.Values{}
	data.Set("isbn", isbn)

	url := fmt.Sprintf("%s/books/delete", dataserver)
	http_post(url, data)

	http.Redirect(w, r, "/books", http.StatusFound)

}
