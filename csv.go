package main

import (
	"encoding/csv"
	"os"
	"sort"
	"strconv"
)

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
