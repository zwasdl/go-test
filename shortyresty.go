package main

// Import necessary packages
import (
	"fmt"
	"io"
	"log"
	"net/http"
	"encoding/json"
	"math/rand"
	"github.com/gorilla/mux"
)

type Url struct {
	url string
}

type ShortUrl struct {
	short_url string
}

var urlMap map[string]Url		// Declare map of shortened URLs keyed by their ID
const chars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"		// Constant string containing possible characters for shortened URL IDs

func main() {
	urlMap = make(map[string]Url)
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", createUrl)
	myRouter.HandleFunc("/shorten", handleShorten)		// Go to handleShorten if URL has endpoint "/shorten"
	myRouter.HandleFunc("/$ID", handleRedirect)		// Go to handleRedirect if URL has endpoint "/%ID"
	log.Fatal(http.ListenAndServe(":8080", nil))		// Host web server on port 8080
}

func createUrl(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func handleShorten(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")		// Set the return Content-Type as JSON

	w.WriteHeader(http.StatusCreated)
	reqBody, _ := io.ReadAll(r.Body)
	var longUrl Url = Url{url: string(reqBody)}
	//fmt.Println(longUrl.url)

	var id string = ""
	for i := 0; i < 8; i++ {
		id += string(chars[rand.Intn(len(chars))])
	}

	var urlStr string = "http://127.0.0.1:8080/" + id
	shortUrl := ShortUrl{urlStr}
	//fmt.Println(shortUrl)

	urlMap[id] = longUrl
	//fmt.Println(urlMap[id])
	shortjson, err := json.Marshal(shortUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(shortjson)
	//w.Write([]byte(`{"short_url": "` + urlStr + `"}`))
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if (r.Method == "GET") {
		w.WriteHeader(http.StatusCreated)
	}
}

/*
curl -v -H "Content-Type: application/json" -X POST -d "{\"url\": \"http://www.abc.com/details\"}" http://127.0.0.1:8080/shorten
*/
