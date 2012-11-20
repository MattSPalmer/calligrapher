package main

import (
	"fmt"
	"strings"
)

var (
	frameRow = strings.Repeat("=", 80)
)

type CallGraph interface {
	Distribution() (map[int64][]CallRecord, error)
	DrawRows() (s string, err error)
}

func Draw(g CallGraph) (string, error) {
	s, err := g.DrawRows()
	if err != nil {
		return "", err
	}

	strSlice := make([]string, 0)
	strSlice = append(strSlice, frameRow, s, frameRow)

	return strings.Join(strSlice, "\n"), nil
}

type GraphByHour []CallRecord
type GraphByDuration []CallRecord
type GraphByAgent []CallRecord

func (bh GraphByHour) Distribution() (map[int64][]CallRecord, error) {
	dist := make(map[int64][]CallRecord)
	var key int64
	for _, call := range bh {
		key = int64(call.Created_at.Hour())
		dist[key] = append(dist[key], call)
	}
	return dist, nil
}

func (bd GraphByDuration) Distribution() (map[int64][]CallRecord, error) {
	dist := make(map[int64][]CallRecord)
	var key int64
	for _, call := range bd {
		key = call.Duration
		dist[key] = append(dist[key], call)
	}
	return dist, nil
}

func (ba GraphByAgent) Distribution() (map[int64][]CallRecord, error) {
	dist := make(map[int64][]CallRecord)
	var key int64
	for _, call := range ba {
		key = call.AgentID
		dist[key] = append(dist[key], call)
	}
	return dist, nil
}

func (bh GraphByHour) DrawRows() (s string, err error) {
	dist, err := bh.Distribution()
	if err != nil {
		return
	}
	rows := make([]string, 0)
	for i := int64(8); i < 21; i++ {
		// Begin row with the row name (hour)
		row := fmt.Sprintf("%02v|", i)

		// Iterate over calls that occurred in hour i
		for _, call := range dist[i] {
			if !call.IsMissed {
				// Get first character of agent name
				var (
					char string
					val  string
				)
				if val = agentsByID[call.AgentID]; val != "" {
					char = string(val[0])
				} else {
					char = "?"
				}
				row += fmt.Sprintf(" %v", char)
			} else {
				row += " -"
			}
		}
		rows = append(rows, row)
	}
	return strings.Join(rows, "\n"), err
}

func (bd GraphByDuration) DrawRows() (s string, err error) {
	dist, err := bd.Distribution()
	if err != nil {
		return
	}
	rows := make([]string, 12)
	for i := 0; i < 11; i++ {
		rows[i] = fmt.Sprintf("%02d-%02d|", 5*i, 5*(i+1))
	}
	rows[11] = "  60+|"
	var (
		counts = make(map[int64]int)
	)
	for k, v := range dist {
		for _, call := range v {
			if !call.IsMissed && call.IsCustomerCare {
				if k/5 > 10 {
					counts[11] += 1
				} else {
					counts[k/5] += 1
				}
			}
		}
	}
	for k, v := range counts {
		rows[k] += fmt.Sprintf(" %d", v)
	}
	return strings.Join(rows, "\n"), err
}

func (ba GraphByAgent) DrawRows() (s string, err error) {
	dist, err := ba.Distribution()
	if err != nil {
		return
	}

	rows := make([]string, 0)

	for k, v := range dist {
		if k == -1 {
			continue
		}
		rows[len(rows)-1] += fmt.Sprintf(" %d", len(v))
	}
	return strings.Join(rows, "\n"), err
}
