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
	errs := make(chan error)
	go func() {
		errs <- http.ListenAndServe(":26692", nil)
	}()
	if (*flagCertFile != "" && *flagKeyFile != "") {
		go func() {
			errs <- http.ListenAndServeTLS(":26693", *flagCertFile, *flagKeyFile, nil)
		}()
	}
	for err := range errs {
		log.Print(err)
	}
	log.Panic(1)
}
