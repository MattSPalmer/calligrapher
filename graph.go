package main

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	frameRow = strings.Repeat("=", 80)
)

type CallGraph interface {
	Distribution() (map[int64][]CallRecord, error)
	DrawRows() (string, error)
	Labels(int64) string
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
		key = call.Duration - call.Duration%5
		if key > 60 {
			key = 60
		}
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
		row := fmt.Sprintf("%4v|", bh.Labels(i))

		// Iterate over calls that occurred in hour i
		for _, call := range dist[i] {
			if !call.IsMissed {
				var char, val string

				// Get first character of agent name
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
	rows := make([]string, 13)
	for i := 0; i < 13; i++ {
		rows[i] = fmt.Sprintf("%v|", bd.Labels(int64(5*i)))
	}
	counts := make(map[int64]int)
	for k, v := range dist {
		counts[k/5] = len(v)
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
		agentName := ba.Labels(k)
		rows = append(rows, fmt.Sprintf("%9v|", agentName))
		rows[len(rows)-1] += fmt.Sprintf(" %d", len(v))
	}
	return strings.Join(rows, "\n"), err
}

func (bh GraphByHour) Labels(v int64) string {
	return fmt.Sprintf("%02v00", v)
}

func (bd GraphByDuration) Labels(v int64) string {
	if v == 60 {
		return "60+"
	}
	return fmt.Sprintf("%02d-%02d", v, v+5)
}

func (ba GraphByAgent) Labels(v int64) string {
	if val, ok := agentsByID[v]; ok {
		return val
	}
	return strconv.FormatInt(v, 10)
}

// Implement sort for Table so we can output nicely sorted CSV files.
type Table [][]string

func (t Table) Less(i, j int) bool { return t[i][0] < t[j][0] }
func (t Table) Len() int           { return len(t) }
func (t Table) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
