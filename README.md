# Book Listing Exercise #

The purpose of this exercise is to test your familiarity with Go development. You'll be building a small book listing app using the Goodreads' public API.

## Functional Requirements ##
* The client app will accept the following command line arguments:
    - --help Output a usage message and exit
    - -s, --search _TERMS_ Search the Goodreads' API and display the results on screen.
        + Results must include author, title, and a link or display of the image of the book
    - --sort _FIELD_ where field is one of "author" or "title"
        + Sorts the results by the specified field, if no sort is specified, title is the default
    - -p _NUMBER_ if you choose to implement pagination, display the _NUMBER_ page of results
    - -h, --host _HOSTNAME_ the hostname or ip address where the server can be found, should default to 127.0.0.1

* There should be a server component as well. The server component should provide REST endpoints that the client
  communicates with. The client should not directly contact the Goodreads API.
* The server should listen on a non-restricted port and the client should connect to that port.

## System Requirements ##

* The application must be written in Go
* Errors that occur during processing should be logged and the user should be presented with a message asking them to retry.

## Non-Requirements ##

* Security measures, including user authentication / authorization
* Unit testing
* UX, as long as the application is usable. As this is just an exercise the UX can be command-line only
* Logging, with the exception of errors

## Misc Notes ##

* https://www.goodreads.com/api/index#search.books
* The Goodreads search API returns XML. Transform the XML into JSON and only send what your app will need
* Be sure to document your code, _especially cases where you might have made a different choice in a 'real' application_
* Upon completion, be sure that your code is accessible through a git repo, and provide the link to that repo to Dotdash Meredith

## Bonus Points ##

* Include pagination in the UI.


## Chris' instructions ##
To run server, simply run:
```shell
$ make run-server
```

To run client, run
```shell
go run cmd/client/main.go [args]
```