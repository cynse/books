package main

import (
	"log"
	"net/http"
	"strconv"
)

type Request struct {
	searchString string
	page         int
}

// In production, this handler would be registered in a separate file to keep the main file clean. However, since this
// is such a small application, I've decided to leave it here for simplicity
func handler(w http.ResponseWriter, r *http.Request) {
	qParams := r.URL.Query()
	searchString := qParams.Get("search")
	page, err := strconv.Atoi(qParams.Get("page"))
	if err != nil {
		http.Error(w, "Format of page must be an integer!", 400)
		return
	}

	r := Request{
		searchString: searchString,
		page:         page,
	}

	s, err := get(r)

	// Respond to client

}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
