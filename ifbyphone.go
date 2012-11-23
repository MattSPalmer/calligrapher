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
	APIUrl    = "https://secure.ifbyphone.com/ibp_api.php?"
	extension = "csv"
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

	if *toFile {
		ds := time.Now().Format("01-02-06_15:04:05")
		filePath := fmt.Sprintf("call_graph_%v.%v", ds, extension)
		err := WriteToCSV(data, filePath)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("Wrote results to file %v\n\n", filePath)
	}

	graph, err := Draw(data)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Println(graph)
}
