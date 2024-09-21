package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type URL struct { // url-structure in database
	ID           string    `json:"id"`
	OriginalURL  string    `json:"original_url"`
	ShortURL     string    `json:"short_url"`
	CreationDate time.Time `json:"creation_date"`
}

// in memory database.., mapping each shortURL to its whole URL structure
var urlDB = make(map[string]URL)

func generateShortURL(OriginalURL string) string {
	hasher := md5.New()               //a checksum to convert string into hash
	hasher.Write([]byte(OriginalURL)) //converts org.URL str to a byte slice
	fmt.Println("hasher: ", hasher)
	data := hasher.Sum(nil)
	fmt.Println("Hasher data:", data)
	hash := hex.EncodeToString(data)
	fmt.Println("EncodeTOString :", hash)
	fmt.Println("final String :", hash[:8])
	return hash[:8]
}

func createURL(originalURL string) string {
	shortURL := generateShortURL(originalURL)
	id := shortURL
	urlDB[id] = URL{
		ID:           id,
		OriginalURL:  originalURL,
		ShortURL:     shortURL,
		CreationDate: time.Now(),
	}
	return shortURL
}

func getURL(id string) (URL, error) { //function to get our original url from the shorturl sstring
	url, ok := urlDB[id]
	if !ok {
		return URL{}, errors.New("URL not found")
	}
	return url, nil
}

func RootPageurlhandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hiey there")
}

func ShorturlHandler(w http.ResponseWriter, r *http.Request) {

	var data struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortURL := createURL(data.URL)
	//fmt.Fprintf(w, shortURL)

	response := struct {
		ShortURL string `json:"short_url`
	}{ShortURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func redirectURLHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/redirect/"):]
	url, err := getURL(id)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusNotFound)
	}
	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

func main() {
	// fmt.Println("..Starting the URL shortner")
	// OriginalURL := "https://github.com/JainnPranjal"
	// generateShortURL(OriginalURL)

	//register the handler fn to handle all the reqs to the rootUrl"/....as handler fn for handling the page at sm address
	http.HandleFunc("/", RootPageurlhandler)
	http.HandleFunc("/shorten", ShorturlHandler)
	http.HandleFunc("/redirect/", redirectURLHandler)

	//starting the server at port 4003
	fmt.Println("starting server on 4003")
	err := http.ListenAndServe(":4003", nil)
	if err != nil {
		fmt.Println("errror on starting the server:", err)
	}
}
