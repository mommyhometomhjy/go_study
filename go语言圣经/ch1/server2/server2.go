package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		count++
		mu.Unlock()

		fmt.Fprintf(w, "url path = %q\n", r.URL.Path)
	})
	http.HandleFunc("/count", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		fmt.Fprintf(w, "count = %d\n", count)
		mu.Unlock()

	})
	log.Fatal(http.ListenAndServe(":8000", nil))
}
