package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

const GoodreadURL = "https://www.goodreads.com/search/index.xml"

func get(request Request) (*string, int, error) {
	// Set up URL
	u, err := url.Parse(GoodreadURL)
	if err != nil {
		log.Fatalf("Server Error: Cannot parse URL %s", GoodreadURL)
	}

	u.Scheme = "https"
	q := u.Query()
	q.Set("q", request.searchString)
	q.Set("p", string(rune(request.page)))

	// In production, I would read this from environment variables.
	q.Set("key", "RDfV4oPehM6jNhxfNQzzQ")
	u.RawQuery = q.Encode()

	// Query the GoodReads API
	resp, err := http.Get(u.String())

	if err != nil {
		return nil, resp.StatusCode, err
	}

	body, _ := io.ReadAll(resp.Body)
	out := string(body)
	return &out, 200, nil
}
