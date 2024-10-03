package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	// sessionCookieName is, surprise, the name of the session cookie that this service creates and destroys
	sessionCookieName = "session"
)

// main is the command line entry point to the application.
func main() {

	// Route the annoying requests for favicon.ico
	http.HandleFunc("/favicon.ico", handleFavicon)

	// Route request for the login URL path and everything else under it
	http.HandleFunc("/login", handleLogin)

	// Route request for the logout URL path and everything else under it
	http.HandleFunc("/logout", handleLogout)

	// Start the server
	const serverAddress = ":8080"
	log.Printf("starting server at %s\n", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

// handleLogin is the HTTP request handler for login URLs. It crudely "logs in" the
// username passed as a query parameter, creating a "session" cookie containing that
// value.
//
// Obviously, this is not in any way a realistic implementation of user authentication. It's
// sole purpose is to create a cookie value whose presence that affects the behavior of the
// Envoy external authorization filter.
func handleLogin(w http.ResponseWriter, r *http.Request) {

	// See if we have a user query parameter
	user := r.URL.Query().Get("user")
	log.Printf("handling request = %s, user = %s\n", r.URL.Path, user)

	// If we don't' have a username in the query parameters, tell the user how to give us one
	if len(user) == 0 {

		// Respond that a user= parameter is required
		_, _ = fmt.Fprintln(w, "To login, add a user=username query parameter to this ULR path")

	} else {

		// We have a username - redirect to /dashboard setting a session cookie
		log.Printf("setting session cookie for user = %s\n", user)
		cookie := buildSessionCookie(user)
		http.SetCookie(w, cookie)

		// Build the redirect URL back to the same host
		redirectUrl := "http://" + r.Host + "/dashboard"

		// Redirect to the dashboard URL without making assumptions about our domain name
		log.Printf("redirecting to %s\n", redirectUrl)
		http.Redirect(w, r, redirectUrl, http.StatusFound)
	}
}

// handleLogin is the HTTP request handler for logout URLs. It simply deletes any session cookie.
func handleLogout(w http.ResponseWriter, r *http.Request) {

	// See if we have a session cookie?
	cookie, err := r.Cookie(sessionCookieName)

	// If we don't' have session cookie, simply inform the user that they are not logged in
	if err != nil {

		// Respond that a user= parameter is required
		_, _ = fmt.Fprintln(w, "not currently logged in")

	} else {

		// We have a session cookie, instruct the browser to remove it by setting its maximum age to 0
		user := cookie.Value
		cookie := buildSessionCookie("")
		cookie.MaxAge = -1 // Destroy immediately
		http.SetCookie(w, cookie)
		msg := fmt.Sprintf("user %s has been logged out\n", user)
		log.Print(msg)
		_, _ = fmt.Fprintln(w, msg)
	}
}

// buildSessionCookie returns a cookie structure describing a non-secure session cookie containing the
// given string value.
func buildSessionCookie(value string) *http.Cookie {

	// This cookies is implicitly a session cookie because no expiration time or max age is set.
	return &http.Cookie{
		Value:    value,
		Name:     sessionCookieName,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

// handleFavicon is the HTTP request handler for "/favicon.ico" requests: 404 NOT FOUND it!
func handleFavicon(w http.ResponseWriter, r *http.Request) {
	log.Printf("handling request = %s\n", r.URL.Path)
	http.Error(w, "Not Found", http.StatusNotFound)
}
