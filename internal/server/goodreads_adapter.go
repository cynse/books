package server

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
)

const goodreadsURL = "https://www.goodreads.com/search/index.xml"
const searchQueryKey = "q"
const pageSearchKey = "page"
const apiKeyQuery = "key"

// WorksList represents the way works are structured from Goodreads response
type WorksList []struct {
	BestBook struct {
		Title  string `xml:"title"`
		Author struct {
			Name string `xml:"name"`
		} `xml:"author"`
		Image string `xml:"image_url"`
	} `xml:"best_book"`
}

// GoodreadsResponse represents the XML structure response from Goodreads
type GoodreadsResponse struct {
	XMLName xml.Name `xml:"GoodreadsResponse"`
	Search  struct {
		Source  string `xml:"source"`
		Results struct {
			WorksList WorksList `xml:"work"`
		} `xml:"results"`
	} `xml:"search"`
}

// GetFromGoodreads queries Goodreads using parameters from the request and returns a list of books
func GetFromGoodreads(request Request) ([]Book, *BooksError) {
	// Set up URL
	u, err := url.Parse(goodreadsURL)
	if err != nil {
		message := fmt.Sprintf("Server Error: Cannot parse URL %s", goodreadsURL)
		log.Print(message)
		return nil, NewBooksError(http.StatusInternalServerError, message)
	}

	// Set query params
	queryParams := u.Query()
	queryParams.Set(searchQueryKey, request.searchString)
	if request.page != nil {
		queryParams.Set(pageSearchKey, fmt.Sprintf("%d", *request.page))
	}
	// In production, I would read this from environment variables.
	queryParams.Set(apiKeyQuery, "RDfV4oPehM6jNhxfNQzzQ")
	u.RawQuery = queryParams.Encode()

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
