package main

import (
	"log"
	"net/http"
	"strconv"
)

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

	s, statusCode, err := get(*req)
	if err != nil {
		// If there are "client errors" reading to Goodreads, it is actually a server error on our end.
		if statusCode >= 400 {
			statusCode = 500
		}
		http.Error(w, err.Error(), statusCode)
	}

	// Respond to client
	w.WriteHeader(statusCode)
	w.Write([]byte(*s))
}

// This function will wrap
func validateAndWrapRequest(r *http.Request) (*Request, *BooksError) {
	requestParams := r.URL.Query()

	// Validate searchString
	searchString := requestParams.Get("q")
	if searchString == "" {
		return nil, &BooksError{
			statusCode: 400,
			message:    "Search string must not be empty! Please try again.",
		}
	}

	// Validate page number
	page := requestParams.Get("p")
	var requestPage *int
	if page != "" {
		if intPage, err := strconv.Atoi(page); err != nil {
			return nil, &BooksError{
				statusCode: 400,
				message:    "Page number must be provided as an integer.",
			}
		} else {
			requestPage = new(int)
			*requestPage = intPage
		}
	}

	// Validate sortBy
	sortBy := requestParams.Get("sortby")
	if sortBy == "" {
		sortBy = "title"
	} else if sortBy != "title" && sortBy != "author" {
		return nil, &BooksError{
			statusCode: 400,
			message:    "SortBy must be either 'title' or 'author'",
		}
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
