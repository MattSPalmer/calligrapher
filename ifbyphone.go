package main

import (
	"bytes"
  "flag"
	"fmt"
	"github.com/MattSPalmer/objcsv"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
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

type CallGraph struct {
	Type    int
	Records []CallRecord
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

func (cg CallGraph) Draw() {
	return
}

func (cg CallGraph) distribution() (map[int64]int, error) {
	dist := make(map[int64]int)
	var key int64

	for _, call := range cg.Records {
		switch cg.Type {
		case ByAgent:
			key = call.AgentNumber
		case ByDuration:
			key = call.Duration
		case ByHour:
			callTime, err := time.Parse("2006-01-02 15:04:05", call.Created_at)
			if err != nil {
				return nil, err
			}
			key = int64(callTime.Hour())
		default:
			return nil, fmt.Errorf("Error: invalid CallType %v", cg.Type)
		}
		dist[key] = dist[key] + 1
	}
	return dist, nil
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
  var start, end string

  flag.StringVar(&start, "s", "", "start date (\"YYYYMMDD\")")
  flag.StringVar(&end, "e", "", "end date (\"YYYYMMDD\")")

  flag.Parse()
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

	duration, err := CallGraph{ByDuration, calls}.distribution()
	hour, err := CallGraph{ByHour, calls}.distribution()
	agent, err := CallGraph{ByAgent, calls}.distribution()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Printf("Durations: %v\n", duration)
	fmt.Printf("Hours: %v\n", hour)
	fmt.Printf("Agents: %v\n", agent)
}
