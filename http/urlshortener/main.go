//go:build !solution

package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"sync"
)

type URLData struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

var (
	urlToKeyMap = make(map[string]string)
	keyToUrlMap = make(map[string]string)
	mutex       = &sync.Mutex{}
)

func generateKey() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var data URLData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var key string
	var ok bool

	mutex.Lock()
	if key, ok = urlToKeyMap[data.URL]; !ok {
		key = generateKey()

		urlToKeyMap[data.URL] = key
		keyToUrlMap[key] = data.URL
	}
	mutex.Unlock()

	response := URLData{URL: data.URL, Key: key}
	jsonResponse, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonResponse)
}

func goHandler(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Path[len("/go/"):]

	mutex.Lock()
	url, ok := keyToUrlMap[key]
	mutex.Unlock()

	if !ok {
		http.Error(w, "key not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shortenHandler)
	http.HandleFunc("/go/", goHandler)

	port := os.Args[2]
	_ = http.ListenAndServe(":"+port, nil)
}
