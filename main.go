package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		// if we want to print a %v we cannot use log.Fatal, instead print and then os.Exit(1)
		log.Printf("error: could not read from input: %v\n", err)
		os.Exit(1)
	}
}

func checkDomain(domain string) {
	// set variables for all of the data to verify
	var hasMX, hasSPF, hasDMARC bool
	var spfRecord, dmarcRecord string

// set a variable mxRecords with the result of the LookupMX from net
	mxRecords, err := net.LookupMX(domain)
// if there's an error, print it
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	// otherwise if there's anything in mxRecords > hasMX=true
	if len(mxRecords) > 0 {
		hasMX = true
	}

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("error: %v\n", err)
	}

	// loop over txtRecords
	// if we don't want an iterator variable, use "_"
	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			hasSPF = true
			// set spfRecord
			spfRecord = record
			// break out
			break
		}
	}

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	// loop over dmarcRecords
	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			hasDMARC = true
			// set dmarcRecord
			dmarcRecord = record
			break
		}
	}

	fmt.Printf("%v, %v, %v, %v, %v, %v", domain, hasMX, hasSPF, spfRecord, hasDMARC, dmarcRecord)
}