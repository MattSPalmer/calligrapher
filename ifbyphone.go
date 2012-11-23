package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	APIUrl    = "https://secure.ifbyphone.com/ibp_api.php?"
	extension = "csv"
)

var (
	filterCriteria = map[string](func(CallRecord) bool){
		"customer care":   func(cr CallRecord) bool { return cr.IsCustomerCare },
		"answered":        func(cr CallRecord) bool { return !cr.IsMissed },
		"during business": func(cr CallRecord) bool { return cr.IsMissed },
	}
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

	// The output of this function depends on the toFile flag.
	calls, err := getCallsByDate(start, end)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	calls = Filter(calls, filterCriteria["customer care"])

	var data CallGraph

	switch *graphType {
	case "duration":
		data = GraphByDuration(calls)
	case "agent":
		data = GraphByAgent(calls)
	case "hour":
		data = GraphByHour(calls)
	default:
		fmt.Printf("invalid graphType specifed: %v", *graphType)
		return
	}

	if err := graphOutput(data); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}
