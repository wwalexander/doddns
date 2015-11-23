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
	"path/filepath"
	"strings"
	"time"
)

type tokenSource struct {
	accessToken string
}

func (ts tokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: ts.accessToken}, nil
}

// Update checks the client's current IP address against the DigitalOcean DNS A
// record, and updates the record if necessary.
func Update(domain string, subdomain string, ipServer string, client *godo.Client) error {
	drs, _, err := client.Domains.Records(domain, nil)
	if err != nil {
		return err
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
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	newAddr := strings.TrimRight(string(body), "\n")
	drer := &godo.DomainRecordEditRequest{
		Type: "A",
		Name: subdomain,
		Data: newAddr,
	}
	if id == nil {
		if _, _, err := client.Domains.CreateRecord(domain, drer); err != nil {
			return err
		}
	} else if newAddr != addr {
		if _, _, err := client.Domains.EditRecord(domain, *id, drer); err != nil {
			return err
		}
	}
	return nil
}

const usage = `usage: doddns-client [-interval seconds] [domain] [subdomain] [server] [token]

doddns-client periodically updates the DNS A record for subdomain.domain using
the IP address returned by the named server via HTTP. To authenticate,
doddns-client uses the DigitalOcean API token saved at the named path. By
default, doddns-client updates the record every 5 minutes; this interval can be
set to a give number of second using the -interval flag.

doddns-server provides an implementation of a server that can be used by
doddns-client.`

func main() {
	finterval := flag.Uint("interval", 300, "the interval between updates in seconds")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
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
	ts := tokenSource{accessToken: string(token)}
	client := godo.NewClient(oauth2.NewClient(oauth2.NoContext, ts))
	logFile, err := os.OpenFile(filepath.Base(strings.TrimSuffix(os.Args[0],
		filepath.Ext(os.Args[0])))+".log",
		os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal("unable to open log file")
	}
	log.SetOutput(logFile)
	for {
		if err := Update(domain, subdomain, ipServer, client); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Duration(*finterval) * time.Second)
	}
}
