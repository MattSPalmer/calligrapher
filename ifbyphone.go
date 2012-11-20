package main

import (
	"bytes"
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
	err := handleArgs()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	r, err := callReader(start, end)
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

	if showDuration {
		durGraph, err := Draw(GraphByDuration(calls))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("Durations:\n%v\n\n", durGraph)
	}

	if showHour {
		hourGraph, err := Draw(GraphByHour(calls))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("Hours:\n%v\n\n", hourGraph)
	}

	if showAgent {
		calls = Filter(calls, func(cr CallRecord) bool {
			return !cr.IsMissed
		})
		agentGraph, err := Draw(GraphByAgent(calls))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("Agents:\n%v\n\n", agentGraph)
	}
}
