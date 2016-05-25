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

type tokenSource struct {
	accessToken string
}

func (ts tokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{AccessToken: ts.accessToken}, nil
}

// Update checks the client's current IP address against the DigitalOcean DNS
// record of the given type, and updates the record if necessary.
func Update(domain string, subdomain string, recordType string, ipServer string, client *godo.Client) error {
	drs, _, err := client.Domains.Records(domain, nil)
	if err != nil {
		return err
	}
	var id *int
	var addr string
	for _, dr := range drs {
		if dr.Type == recordType && dr.Name == subdomain {
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

const usage = `usage: doddns-client [-type type] domain subdomain URI token

doddns-client periodically updates the DNS record for subdomain.domain using
the IP address returned by the named HTTP/HTTPS URI. To authenticate,
doddns-client uses the DigitalOcean API token saved at the named path. The
record is updated to match the TTL of the domain, in order to avoid doing
useless DNS updates between TTL timeouts. By default, doddns-client will update
the A record; arbitrary record types can be specified with the -type flag.`

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, usage)
	}
	var recordType string
	flag.StringVar(&recordType, "type", "A", "the type of record to update")
	flag.Parse()
	args := flag.Args()
	if len(args) != 4 {
		flag.Usage()
		os.Exit(1)
	}
	domain := args[0]
	subdomain := args[1]
	uri := args[2]
	tokenPath := args[3]
	token, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Fatalf("unable to read token file '%s': %v", tokenPath, err)
	}
	ts := tokenSource{accessToken: string(token)}
	client := godo.NewClient(oauth2.NewClient(oauth2.NoContext, ts))
	dom, _, err := client.Domains.Get(domain)
	if err != nil {
		log.Fatal(err)
	}
	logFile, err := os.OpenFile("doddns-client.log", os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
	for {
		if err := Update(domain, subdomain, recordType, uri, client); err != nil {
			log.Println(err)
		}
		time.Sleep(time.Duration(dom.TTL) + time.Second)
	}
}
