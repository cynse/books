package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"text/tabwriter"
)

const port = 8080
const searchQueryKey = "q"
const sortQueryKey = "s"
const pageQueryKey = "p"

// Book represents the structure of our server response
type Book struct {
	Author string
	Title  string
	Image  string
}

func Execute(hostname string, searchTerm string, sortBy string, page int) {
	requestURL := fmt.Sprintf("http://%s:%d", hostname, port)
	u, err := url.Parse(requestURL)
	if err != nil {
		log.Printf("client: could not parse URL: %s\n. Please try again", err)
		return
	}

	queryParams := u.Query()
	queryParams.Set(searchQueryKey, searchTerm)
	queryParams.Set(sortQueryKey, sortBy)
	queryParams.Set(pageQueryKey, fmt.Sprintf("%d", page))
	u.RawQuery = queryParams.Encode()

	if resp, err := http.Get(u.String()); err == nil {
		if respBody, err := io.ReadAll(resp.Body); err == nil {
			// our response body contains the error message and a new line, we can trim that now.
			bodyStr := strings.TrimSuffix(string(respBody), "\n")
			if resp.StatusCode >= 500 {
				log.Printf("server error: %s. Please try again later", bodyStr)
			} else if resp.StatusCode == 400 {
				log.Printf("client error: %s. Please check the input and try again.", bodyStr)
			} else {
				displayBooks(respBody)
			}
		} else {
			log.Printf("client: error reading response body: %s", err.Error())
		}
	} else {
		log.Printf("Error: %s", err.Error())
	}
}

// Displays the json results in a readable format for the CLI
func displayBooks(responseBody []byte) {
	var books []Book
	if err := json.Unmarshal(responseBody, &books); err != nil {
		log.Printf("Unmarshalling error: %s", err.Error())
		return
	}
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "Author", "Title", "Image Link")
	for _, book := range books {
		fmt.Fprintf(w, "%s\t%s\t%s\n", book.Author, book.Title, book.Image)
	}
	_ = w.Flush()
}
