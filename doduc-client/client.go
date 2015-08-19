package main

import (
	"flag"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type tokenSource struct {
	accessToken string
}

func (ts tokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: ts.accessToken}, nil
}

func update(domain string, subdomain string, ipServer string, client *godo.Client) {
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
	flagInterval := flag.Uint("interval", 300, "the interval between updates in seconds")
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
	ts := tokenSource{accessToken: string(token)}
	client := godo.NewClient(oauth2.NewClient(oauth2.NoContext, ts))
	log.SetOutput(logFile)
	for {
		update(domain, subdomain, ipServer, client)
		time.Sleep(time.Duration(*flagInterval) * time.Second)
	}
}
