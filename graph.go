package main

import (
	"fmt"
	"regexp"
	"strings"
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

func filter(g *[]CallRecord) {
	newRecords := make([]CallRecord, 0)
	for _, call := range *g {
		if call.isCustomerCare() {
			newRecords = append(newRecords, call)
		}
	}
	g = &newRecords
}

type CallGraph interface {
	DrawRows(int64, int) string
	Distribution() (map[int64]int, error)
}

func Draw(g CallGraph) (string, error) {
	strSlice := make([]string, 0)
	strSlice = append(strSlice, frameRow)
	dist, err := g.Distribution()
	if err != nil {
		return "", err
	}
	for k, v := range dist {
		strSlice = append(strSlice, g.DrawRows(k, v))
	}
	strSlice = append(strSlice, frameRow)
	return strings.Join(strSlice, "\n"), nil
}

type GraphByHour []CallRecord
type GraphByDuration []CallRecord
type GraphByAgent []CallRecord

func (bh GraphByHour) Distribution() (map[int64]int, error) {
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

func (bd GraphByDuration) Distribution() (map[int64]int, error) {
	dist := make(map[int64]int)
	for _, call := range bd {
		dist[call.Duration]++
	}
	return dist, nil
}

func (ba GraphByAgent) Distribution() (map[int64]int, error) {
	dist := make(map[int64]int)
	for _, call := range ba {
		dist[call.AgentNumber]++
	}
	return dist, nil
}

func (bh GraphByHour) DrawRows(k int64, v int) string {
	return ""
}

func (bd GraphByDuration) DrawRows(k int64, v int) string {
	return ""
}

func (ba GraphByAgent) DrawRows(k int64, v int) string {
	return ""
}
