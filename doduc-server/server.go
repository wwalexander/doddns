package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintln(w, host)
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":18768", nil))
}
