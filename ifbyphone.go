package main

import (
	"bytes"
	"fmt"
	"github.com/MattSPalmer/objcsv"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	APIUrl = "https://secure.ifbyphone.com/ibp_api.php?"

	ByAgent int = iota
	ByHour
	ByDuration
)

type CallRecord struct {
	Created_at   string
	ActivityInfo string
	CallerID     string
	Duration     int64
	AgentNumber  int64
}

func (cr CallRecord) isCustomerCare() bool {
	query, err := regexp.Compile("Customer Care")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return query.MatchString(cr.ActivityInfo)
}

func (cr CallRecord) isMissed() bool {
	return cr.AgentNumber == 0
}

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
	flagInit()

	r, err := callReader(start, end)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	calls := make([]CallRecord, 0)
	err = objcsv.ReadCSV(r, &calls)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	duration := CallGraph{ByDuration, calls}
	hour := CallGraph{ByHour, calls}
	agent := CallGraph{ByAgent, calls}

	duration.filter()
	hour.filter()
	agent.filter()

	durGraph, err := duration.distribution()
	hourGraph, err := hour.distribution()
	agentGraph, err := agent.distribution()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("Durations: %v\n\n", durGraph)
	fmt.Printf("Hours: %v\n\n", hourGraph)
	fmt.Printf("Agents: %v\n\n", agentGraph)
}
