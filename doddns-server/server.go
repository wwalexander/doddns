package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Handler responds a request with its IP address.
func Handler(w http.ResponseWriter, r *http.Request) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "unable to get IP address", 500)
		return
	}
	fmt.Fprintln(w, host)
}

func main() {
	fport := flag.Uint("port", 18768, "the port to listen on")
	flag.Parse()
	logFile, err := os.OpenFile(filepath.Base(strings.TrimSuffix(os.Args[0],
				filepath.Ext(os.Args[0])))+".log",
		os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("unable to open log file")
	}
	log.SetOutput(logFile)
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *fport), nil))
}
