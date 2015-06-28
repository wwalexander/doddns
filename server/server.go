package main

import (
	"flag"
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
	flagCertFile := flag.String("cert", "", "the SSL certificate")
	flagKeyFile := flag.String("key", "", "the  SSL private key")
	flag.Parse()
	http.HandleFunc("/", handler)
	port := ":26992"
	var err error
	if (*flagCertFile != "" && *flagKeyFile != "") {
		err = http.ListenAndServeTLS(port, *flagCertFile, *flagKeyFile, nil)
	} else {
		err = http.ListenAndServe(port, nil)
	}
	log.Fatal(err)
}
