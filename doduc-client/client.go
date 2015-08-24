package main

import (
	"flag"
	"fmt"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// TokenSource returns a token.
type TokenSource struct {
	AccessToken string
}

// Token returns a token.
func (ts TokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: ts.AccessToken}, nil
}

// Update checks the client's current IP address against the DigitalOcean DNS
// record, and updates the record if necessary.
func Update(domain string, subdomain string, ipServer string, client *godo.Client) {
	drs, _, err := client.Domains.Records(domain, nil)
	if err != nil {
		log.Printf("unable to fetch domain records: %v", err)
		return
	}
	var id *int
	var addr string
	for _, dr := range drs {
		if dr.Type == "A" && dr.Name == subdomain {
			id = &dr.ID
			addr = dr.Data
		}
	}
	resp, err := http.Get(ipServer)
	if err != nil {
		log.Printf("unable to open request to IP server: %v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("unable to read response from IP server: %v", err)
	}
	newAddr := strings.TrimRight(string(body), "\n")
	drer := &godo.DomainRecordEditRequest{
		Type: "A",
		Name: subdomain,
		Data: newAddr,
	}
	if id == nil {
		log.Println("record not found; creating new record")
		_, _, err := client.Domains.CreateRecord(domain, drer)
		if err != nil {
			log.Printf("unable to create new record: %v", err)
			return
		}
	} else if newAddr != addr {
		log.Println("address changed; updating record")
		_, _, err := client.Domains.EditRecord(domain, *id, drer)
		if err != nil {
			log.Printf("unable to update record: %v", err)
			return
		}
	} else {
		log.Println("record is up to date")
	}
}

func main() {
	logFile, err := os.OpenFile("doduc-client.log", os.O_APPEND|os.O_CREATE, 0200)
	if err != nil {
		log.Fatal("unable to open log file")
	}
	finterval := flag.Uint("interval", 300, "the interval between updates in seconds")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [OPTIONS] [DOMAIN] [SUBDOMAIN] [SERVER] [TOKEN]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) != 4 {
		flag.Usage()
		os.Exit(1)
	}
	domain := args[0]
	subdomain := args[1]
	ipServer := args[2]
	tokenPath := args[3]
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Fatalf("unable to read token file '%s': %v", tokenPath, err)
	}
	ts := TokenSource{AccessToken: string(token)}
	client := godo.NewClient(oauth2.NewClient(oauth2.NoContext, ts))
	log.SetOutput(logFile)
	for {
		Update(domain, subdomain, ipServer, client)
		time.Sleep(time.Duration(*finterval) * time.Second)
	}
}
