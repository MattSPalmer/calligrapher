package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

var (
	frameRow = strings.Repeat("=", 80)
)

type CallGraph interface {
	Distribution() (map[int64][]CallRecord, error)
	DrawRows() (string, error)
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

func WriteToCSV(cg CallGraph, fp string) error {
	f, err := os.Create(fp)
	if err != nil {
		return err
	}

	defer f.Close()

	dist, err := cg.Distribution()
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	t := make(Table, 0)
	for key, calls := range dist {
		var sum int64
		l := int64(len(calls))

		for _, call := range calls {
			sum += call.Duration
		}

		keyS := strconv.FormatInt(key, 10)
		sumS := strconv.FormatInt(sum, 10)
		avgS := strconv.FormatInt(sum/l, 10)
		lenS := strconv.FormatInt(l, 10)

		row := []string{keyS, lenS, sumS, avgS}
		t = append(t, row)
	}
	sort.Sort(t)
	if err = w.WriteAll(t); err != nil {
		return err
	}
	w.Flush()
	return nil
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
	for _, call := range ba {
		key := call.AgentID
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
		var agentName string
		if val, ok := agentsByID[k]; ok {
			agentName = val
		} else {
			agentName = strconv.FormatInt(k, 10)
		}
		rows = append(rows, fmt.Sprintf("%9v|", agentName))
		rows[len(rows)-1] += fmt.Sprintf(" %d", len(v))
	}
	return strings.Join(rows, "\n"), err
}

type Table [][]string

func (t Table) Less(i, j int) bool { return t[i][0] < t[j][0] }
func (t Table) Len() int           { return len(t) }
func (t Table) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
