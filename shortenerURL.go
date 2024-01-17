package main

import (
	"fmt"
	"net/http"
)

type URLShortener struct {
	urls map[string]string
}

func (us *URLShortener) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	originalURL := r.FormValue("url")

	if originalURL == "" {
		http.Error(w, "URL parameter is missing", http.StatusBadRequest)
		return
	}

	// generate a unique shortened key for the original URL
	shortKey := generateShortKey()
	us.urls[shortKey] = originalURL

	// construct the full shortened URL
	shortenedURL := fmt.Sprintf("http://localhost:8080:/short/%s", shortKey)

	// render the HTML response with the shortened URL
	w.Header().Set("Content-Type", "text/html")
	responseHTML := fmt.Sprintf(`
		<h2>URL Shortener</h2>
		<p>Original URL : %s</p>
		<p>Shortened URL : <a href="%s">%s</a></p>
		<form method="post" action="/shorten">
			<input type="text" name="url" placeholder="Enter a URL">
			<input type="submit" value="Shorten">
		</form>
	`, originalURL, shortenedURL, shortenedURL)

	_, err := fmt.Fprintf(w, responseHTML)
	if err != nil {
		panic(err)
	}
}

func (us *URLShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := r.URL.Path[len("/short/"):]

	if shortKey == "" {
		http.Error(w, "Shortened key is missing", http.StatusBadRequest)
		return
	}

	// retrieve the original URL from the `urls` map using the shortened key
	originalURL, found := us.urls[shortKey]

	if !found {
		http.Error(w, "Shortened key not found", http.StatusNotFound)
		return
	}

	// redirect the user to the original URL
	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
