package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
)

const (
	// cookieName is, surprise, the name of a session cookie that this service populates
	cookieName = "authtest-request"
)

var (
	// requestCount tracks the number of requests received
	requestCount int
)

// main is the command line entry point to the application.
func main() {

	// Route the annoying requests for favicon.ico
	http.HandleFunc("/favicon.ico", handleFavicon)

	// Route request for the root URL path and everything else under it
	http.HandleFunc("/", handleRoot)

	// Start the server
	const serverAddress = ":8080"
	log.Printf("starting server at %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

// handleRoot is the primary HTTP request handler for the root path and everything under it for this very crude website.
func handleRoot(w http.ResponseWriter, r *http.Request) {

	// Bump the count of requests seen
	requestCount++

	// Log the URL that we have been asked for etc
	log.Printf("handling request = %s, request count = %d\n", r.URL.Path, requestCount)

	// Add a cookie containing the path and request count
	setCookie(w, r.URL.Path, requestCount)

	// Dump some text onto the response
	_, _ = fmt.Fprintf(w, "Host:\t\t%q\n", r.Host)
	_, _ = fmt.Fprintf(w, "Path:\t\t%q\n", r.URL.Path)
	_, _ = fmt.Fprintf(w, "Count:\t\t%d\n", requestCount)

	// Dump the headers, alphanumerically sorted, onto the response
	_, _ = fmt.Fprintf(w, "\nHEADERS (%d)\n=======\n\n", len(r.Header))

	// ... populate a list of the header names
	names := make([]string, 0, len(r.Header))
	for name := range r.Header {
		names = append(names, name)
	}

	// ... sort the list of header names
	sort.Strings(names)

	// ... output the headers in name order
	for _, name := range names {
		_, _ = fmt.Fprintf(w, "%v: %v\n", name, r.Header[name])
	}
}

// setCookie adds an
func setCookie(w http.ResponseWriter, path string, count int) {

	// Build the cookie (its a session cookie with no expiration or max age
	cookie := &http.Cookie{
		Value:    fmt.Sprintf("last-path=%s count=%d", path, requestCount),
		Name:     cookieName,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	// Set the cookie in the response
	http.SetCookie(w, cookie)
}

// handleFavicon is the HTTP request handler for "/favicon.ico" requests: 404 NOT FOUND it!
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request = %s\n", r.URL.Path)
	http.Error(w, "Not Found", http.StatusNotFound)
}
