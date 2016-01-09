// models for bookstore app

package models

import (
    "fmt"
)

type Book struct {
    Isbn   string
    Title  string
    Author string
    Price  float32
}

// get an array of Book with optional filter: isbn
func GetBooks(isbn ...string) ([]*Book, error) {

    var query string
    // apply isbn filter if available
    if len(isbn) == 1 {
        query = fmt.Sprintf("SELECT * FROM books WHERE isbn = '%s'", isbn[0])
    } else {
        query = "SELECT * FROM books"
    }

    rows, err := db.Query(query)
    if err != nil {
        return nil, err
    }

    defer rows.Close()

    bks := make([]*Book, 0)
    for rows.Next() {
        bk := new(Book)
        err := rows.Scan(&bk.Isbn, &bk.Title, &bk.Author, &bk.Price)
        if err != nil {
            return nil, err
        }
        bks = append(bks, bk)
    }
    if err = rows.Err(); err != nil {
        return nil, err
    }
    return bks, nil
}

func CreateBook(isbn string, title string, author string, price float64) (int64, error) {

    result, err := db.Exec("INSERT INTO books VALUES($1, $2, $3, $4)", isbn, title, author, price)

    if err != nil {
        return 0, err
    }

    rowsAffected, err := result.RowsAffected()

    if err != nil {
        return 0, err
    }

    return rowsAffected, nil
}

func DeleteBook(isbn string) (int64, error) {

    result, err := db.Exec("DELETE FROM books WHERE isbn=$1", isbn)

    if err != nil {
        return 0, err
    }

    rowsAffected, err := result.RowsAffected()

    if err != nil {
        return 0, err
    }

    return rowsAffected, nil
}


