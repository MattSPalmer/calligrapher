package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/MattSPalmer/objcsv"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	APIUrl = "https://secure.ifbyphone.com/ibp_api.php?"

	ByAgent int = iota
	ByHour
	ByDuration
)

func callReader(start, end string) (io.Reader, error) {
	ibpParams.Add("start_date", start)
	ibpParams.Add("end_date", end)
	theURL := APIUrl + ibpParams.Encode()

	resp, err := http.Get(theURL)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(body), nil
}

func main() {
	flag.Parse()

	r, err := callReader(*start, *end)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	callsFromFile := make([]callRecordFromFile, 0)
	err = objcsv.ReadCSV(r, &callsFromFile)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	calls := make([]CallRecord, 0)
	calls, err = batchConvert(callsFromFile)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	calls = Filter(calls, func(cr CallRecord) bool {
		return cr.IsCustomerCare
	})

	duration := GraphByDuration(calls)
	hour := GraphByHour(calls)

	calls = Filter(calls, func(cr CallRecord) bool {
		return !cr.IsMissed
	})

	agent := GraphByAgent(calls)

	durGraph, err := Draw(duration)
	hourGraph, err := Draw(hour)
	agentGraph, err := Draw(agent)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("Hours:\n%v\n\n", hourGraph)
	fmt.Printf("Durations:\n%v\n\n", durGraph)
	fmt.Printf("Agents:\n%v\n\n", agentGraph)
}
