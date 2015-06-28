package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.RemoteAddr)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":26992", nil))
}
