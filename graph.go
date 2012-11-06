package main

import (
	"fmt"
	"regexp"
	"time"
)

const (
	frameRow = "================================================================================"
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

type CallGraph interface {
	DrawRow() string
	Distribution() map[int64]int
}

func Draw(g CallGraph) (string, error) {
	return "", nil
}

type GraphByHour []CallRecord
type GraphByDuration []CallRecord
type GraphByAgent []CallRecord

func (bh GraphByHour) distribution() (map[int64]int, error) {
	dist := make(map[int64]int)
	for _, call := range bh {
		callTime, err := time.Parse("2006-01-02 15:04:05", call.Created_at)
		if err != nil {
			return nil, err
		}
		dist[int64(callTime.Hour())]++
	}
	return dist, nil
}

func (bd GraphByDuration) distribution() (map[int64]int, error) {
	dist := make(map[int64]int)
	for _, call := range bd {
		dist[call.Duration]++
	}
	return dist, nil
}

func (ba GraphByAgent) distribution() (map[int64]int, error) {
	dist := make(map[int64]int)
	for _, call := range ba {
		dist[call.AgentNumber]++
	}
	return dist, nil
}

func filter(g *[]CallRecord) {
	newRecords := make([]CallRecord, 0)
	for _, call := range *g {
		if call.isCustomerCare() {
			newRecords = append(newRecords, call)
		}
	}
	g = &newRecords
}
