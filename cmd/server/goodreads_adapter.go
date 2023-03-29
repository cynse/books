package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const GoodreadURL = "https://www.goodreads.com/search/index.xml"

func get(request Request) (*string, *BooksError) {
	// Set up URL
	u, err := url.Parse(GoodreadURL)
	if err != nil {
		message := fmt.Sprintf("Server Error: Cannot parse URL %s", GoodreadURL)
		log.Print(message)
		return nil, &BooksError{
			statusCode: 500,
			message:    message,
		}
	}

	q := u.Query()
	q.Set("q", request.searchString)

	if request.page != nil {
		q.Set("p", fmt.Sprintf("%d", *request.page))
	}

	// In production, I would read this from environment variables.
	q.Set("key", "RDfV4oPehM6jNhxfNQzzQ")
	u.RawQuery = q.Encode()

	// Query the GoodReads API
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, &BooksError{
			statusCode: resp.StatusCode,
			message:    err.Error(),
		}
	}

	body, _ := io.ReadAll(resp.Body)
	out := string(body)
	return &out, nil
}
