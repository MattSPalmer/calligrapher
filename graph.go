package main

import (
	"fmt"
	"time"
)

type CallGraph struct {
	Type    int
	Records []CallRecord
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
