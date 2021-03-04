package main

// Import necessary packages
import (
	"log"
	"net/http"
	"net/url"
	"encoding/json"
	"math/rand"
	"strings"
)

// LongUrl formats long URLs for easy conversion to and from JSON
type LongUrl struct {
	Long_url string `json:"url"`
}

// ShortUrl formats short URLs for easy conversion to and from JSON
type ShortUrl struct {
	Short_url string `json:"short_url"`
}

// Map of shortened URLs keyed by their ID
var urlMap map[string]LongUrl
// Constant string containing possible characters for shortened URL IDs
const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// The main function creates the urlMap, handles HTTP requests, and starts a web server on port 8080
func main() {
	urlMap = make(map[string]LongUrl)
	http.HandleFunc("/shorten", handleShorten)		// Go to handleShorten if URL has endpoint "/shorten"
	http.HandleFunc("/", handleRedirect)		// Go to handleRedirect if URL has endpoint "/%ID"
	log.Fatal(http.ListenAndServe(":8080", nil))		// Host web server on port 8080
}

// The handleShorten function takes in a long URL, shortens it into a ShortUrl,
// tracks the long and short URL pairs in a map, and responds with the shortened URL in JSON
func handleShorten(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")		// Set the return Content-Type as JSON
	w.WriteHeader(http.StatusCreated)

	var longUrl LongUrl
	json.NewDecoder(r.Body).Decode(&longUrl)        // Decode the POST request body from JSON into a LongUrl

	_, err := url.ParseRequestURI(longUrl.Long_url)
	if err != nil {
		w.Write([]byte(longUrl.Long_url + " is not a valid URL."))      // Respond with error message if longUrl is not a valid URL
		return
	}

	var id string
	var exists bool = true
	// While id exists in urlMap, keep generating random IDs
	for exists {
		id = generateID()
		if _, key := urlMap[id]; !key {
			exists = false
		}
	}
	urlMap[id] = longUrl        // Add id: longUrl to urlMap

	var urlStr string = "http://127.0.0.1:8080/" + id       // Create the shortened url
	shortUrl := ShortUrl{urlStr}
	json.NewEncoder(w).Encode(shortUrl)     // Encode the response from a ShortUrl into JSON
}

// The handleRedirect function takes in a short URL's ID, finds its corresponding LongUrl
// using the urlMap, and redirects to the long form of the URL
func handleRedirect(w http.ResponseWriter, r *http.Request) {
	var requestURI string = strings.TrimPrefix(r.RequestURI, "/")
	if _, key := urlMap[requestURI]; key {
		http.Redirect(w, r, urlMap[requestURI].Long_url, 302)       // 302 redirect to the long URL
	} else {
		w.Write([]byte(requestURI + " is not linked to a long URL."))       // Respond with error message if id is not in urlMap
		return
	}
}

// The generateID function generates and returns a random string of length 8 to use as the ID
// for a shortened URL
func generateID() string {
	var id string = ""
	for i := 0; i < 8; i++ {
		id += string(chars[rand.Intn(len(chars))])
	}
	return id
}

/*
For testing purposes:
curl -v -X POST http://127.0.0.1:8080/shorten -H "Content-Type: application/json" -d '{"url": "http://google.com"}'
curl -v -X GET http://127.0.0.1:8080/$ID
*/
