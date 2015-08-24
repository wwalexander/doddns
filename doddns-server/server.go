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

const (
	fportName = "port"
	fcertName = "cert"
	fkeyName  = "key"
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

func flagIsSet(name string) bool {
	set := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			set = true
		}
	})
	return set
}

func main() {
	fcert := flag.String(fcertName, "", "the TLS certificate to use")
	fkey := flag.String(fkeyName, "", "the TLS key to use")
	fport := flag.Uint("port", 18768, "the port to listen on")
	flag.Parse()
	fcertSet := flagIsSet(fcertName)
	fkeySet := flagIsSet(fkeyName)
	if fcertSet && !fkeySet {
		log.Fatalf("%s set without %s", fcertName, fkeyName)
	}
	if !fcertSet && fkeySet {
		log.Fatalf("%s set without %s", fkeyName, fcertName)
	}
	tls := fcertSet && fkeySet
	logFile, err := os.OpenFile(filepath.Base(strings.TrimSuffix(os.Args[0],
				filepath.Ext(os.Args[0])))+".log",
		os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("unable to open log file")
	}
	log.SetOutput(logFile)
	http.HandleFunc("/", Handler)
	addr := fmt.Sprintf(":%d", *fport)
	if tls {
		err = http.ListenAndServeTLS(addr, *fcert, *fkey, nil)
	} else {
		err = http.ListenAndServe(addr, nil)
	}
	log.Fatal(err)
}
