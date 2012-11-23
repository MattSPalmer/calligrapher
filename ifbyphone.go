package main

import (
	"bytes"
	"fmt"
	"github.com/MattSPalmer/objcsv"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	APIUrl = "https://secure.ifbyphone.com/ibp_api.php?"
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

	if *toFile {
		ds := time.Now().Format("01-02-06_15:04:05")
		filePath := fmt.Sprintf("call_graph_%v.csv", ds)
		switch {
		case *showDuration:
			err := WriteToCSV(GraphByDuration(calls), filePath)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			fallthrough
		case *showAgent:
			err := WriteToCSV(GraphByAgent(calls), filePath)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
			fallthrough
		case *showHour:
			err := WriteToCSV(GraphByHour(calls), filePath)
			if err != nil {
				fmt.Printf("%v\n", err)
				return
			}
		}
		fmt.Printf("Wrote results to file %v\n", filePath)
	}

	switch {
	case *showDuration:
		durGraph, err := Draw(GraphByDuration(calls))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("Durations:\n%v\n\n", durGraph)
		fallthrough
	case *showHour:
		hourGraph, err := Draw(GraphByHour(calls))
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Printf("Hours:\n%v\n\n", hourGraph)
		fallthrough
	case *showAgent:
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
