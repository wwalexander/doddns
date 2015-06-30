package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "unable to get IP address", 500)
		return
	}
	fmt.Fprintln(w, host)
}

func main() {
	logFile, err := os.OpenFile("doduc-server.log", os.O_APPEND|os.O_CREATE, 0200)
	if err != nil {
		log.Fatal("unable to open log file")
	}
	log.SetOutput(logFile)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":18768", nil))
}
