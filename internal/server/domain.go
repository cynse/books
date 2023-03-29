package server

// Request is the type that we are sending to our goodreads adapter
type Request struct {
	searchString string
	sortBy       string
	page         *int
}

// Book will represent one entry in our list
type Book struct {
	Author string
	Title  string
	Image  string
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
