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

  switch cg.Type {
  case ByAgent:
    for _, call := range cg.Records {
      dist[call.AgentNumber]++
    }
  case ByDuration:
    for _, call := range cg.Records {
      dist[call.Duration]++
    }
  case ByHour:
    for _, call := range cg.Records {
      callTime, err := time.Parse("2006-01-02 15:04:05", call.Created_at)
      if err != nil {
        return nil, err
      }
      dist[int64(callTime.Hour())]++
    }
  default:
    return nil, fmt.Errorf("Error: invalid CallType %v", cg.Type)
}
	return dist, nil
}

func (cg *CallGraph) filter() {
  newRecords := make([]CallRecord, 0)
  for _, call := range cg.Records {
    if call.isCustomerCare() {
      newRecords = append(newRecords, call)
    }
  }
  cg.Records = newRecords
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
