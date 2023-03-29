package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
)

const GoodreadURL = "https://www.goodreads.com/search/index.xml"

type WorksList []struct {
	BestBook struct {
		Title  string `xml:"title"`
		Author struct {
			Name string `xml:"name"`
		} `xml:"author"`
		Image string `xml:"image_url"`
	} `xml:"best_book"`
}

type GoodreadsResponse struct {
	XMLName xml.Name `xml:"GoodreadsResponse"`
	Search  struct {
		Source  string `xml:"source"`
		Results struct {
			WorksList WorksList `xml:"work"`
		} `xml:"results"`
	} `xml:"search"`
}

// Book will represent one entry in our list
type Book struct {
	Author string
	Title  string
	Image  string
}

func get(request Request) ([]Book, *BooksError) {
	// Set up URL
	u, err := url.Parse(GoodreadURL)
	if err != nil {
		message := fmt.Sprintf("Server Error: Cannot parse URL %s", GoodreadURL)
		log.Print(message)
		return nil, NewBooksError(http.StatusInternalServerError, message)
	}

	// Set query params
	q := u.Query()
	q.Set("q", request.searchString)
	if request.page != nil {
		q.Set("p", fmt.Sprintf("%d", *request.page))
	}
	// In production, I would read this from environment variables.
	q.Set("key", "RDfV4oPehM6jNhxfNQzzQ")
	u.RawQuery = q.Encode()

	// Get from the GoodReads API
	resp, err := http.Get(u.String())
	if err != nil {
		// If there are "client errors" reading to Goodreads, it is actually a server error on our end since that would be a bug.
		code := resp.StatusCode
		if code >= 400 {
			code = http.StatusInternalServerError
		}
		return nil, NewBooksError(code, err.Error())
	}

	body, _ := io.ReadAll(resp.Body)
	bookList, booksErr := convertGoodreadsResponseToBookList(body)
	if booksErr != nil {
		return nil, booksErr
	}

	// Sort the result
	if request.sortBy == "title" {
		sort.Slice(bookList, func(i, j int) bool {
			return bookList[i].Title < bookList[j].Title
		})
	} else {
		sort.Slice(bookList, func(i, j int) bool {
			return bookList[i].Author < bookList[j].Author
		})
	}

	return bookList, nil
}

// Take the goodreads XML response and convert it into Book list
func convertGoodreadsResponseToBookList(body []byte) ([]Book, *BooksError) {
	responseBody := &GoodreadsResponse{}
	if err := xml.Unmarshal(body, responseBody); err != nil {
		return nil, NewBooksError(http.StatusInternalServerError, "Issue unmarshalling Goodreads XML Response. Please try again")
	}
	works := responseBody.Search.Results.WorksList

	books := convertWorksListToBookList(works)
	return books, nil
}

// This function converts the Goodreads formatted object to a list of our Book object
func convertWorksListToBookList(works WorksList) []Book {
	var bookList []Book
	for _, work := range works {
		book := Book{
			Author: work.BestBook.Author.Name,
			Title:  work.BestBook.Title,
			Image:  work.BestBook.Image,
		}
		bookList = append(bookList, book)
	}

	return bookList
}
