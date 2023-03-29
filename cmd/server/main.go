package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// Request is the type that we are sending to our goodreads adapter
type Request struct {
	searchString string
	sortBy       string
	page         *int
}

// BooksError would be defined in a separate file, but it is a small application, so I keep it here
type BooksError struct {
	statusCode int
	message    string
}

func (b *BooksError) Error() string {
	return b.message
}

// NewBooksError is a helper function to create a booksError
func NewBooksError(statusCode int, message string) *BooksError {
	return &BooksError{
		statusCode: statusCode,
		message:    message,
	}
}

// In production, this handler would be registered in a separate file to keep the main file clean. However, since this
// is such a small application, I've decided to leave it here for simplicity
//
// This handler validates a request, calls the goodreads adapter and responds
func handler(w http.ResponseWriter, r *http.Request) {
	req, booksErr := validateAndWrapRequest(r)
	if booksErr != nil {
		http.Error(w, booksErr.message, booksErr.statusCode)
		return
	}

	bookList, booksErr := get(*req)
	if booksErr != nil {
		http.Error(w, booksErr.Error(), booksErr.statusCode)
		log.Printf("Internal error: %s", booksErr.Error())
		return
	}

	s, err := json.Marshal(bookList)
	if err != nil {
		http.Error(w, "Error marshalling bookList", http.StatusInternalServerError)
		log.Printf("Internal error: Error marshalling booklist, %v", err)
		return
	}

	// Respond to client
	w.WriteHeader(http.StatusOK)
	w.Write(s)
}

// This function will wrap the request from client and validate query params
func validateAndWrapRequest(r *http.Request) (*Request, *BooksError) {
	requestParams := r.URL.Query()

	// Validate searchString
	searchString := requestParams.Get("q")
	if searchString == "" {
		return nil, NewBooksError(http.StatusBadRequest, "Search string must not be empty! Please try again.")
	}

	// Validate page number
	page := requestParams.Get("p")
	var requestPage *int
	if page != "" {
		if intPage, err := strconv.Atoi(page); err != nil {
			return nil, NewBooksError(http.StatusBadRequest, "Page number must be provided as an integer.")
		} else {
			requestPage = new(int)
			*requestPage = intPage
		}
	}

	// Validate sortBy
	sortBy := requestParams.Get("s")
	if sortBy == "" {
		sortBy = "title"
	} else if sortBy != "title" && sortBy != "author" {
		return nil, NewBooksError(http.StatusBadRequest, "SortBy must be either 'title' or 'author'")
	}

	return &Request{
		searchString: searchString,
		page:         requestPage,
		sortBy:       sortBy,
	}, nil

}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
