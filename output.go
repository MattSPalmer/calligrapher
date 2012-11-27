package main

import (
	"fmt"
	"github.com/MattSPalmer/objcsv"
	"time"
)

var pathIncrement int

// GetCallsByDate takes two date strings of the format "YYYYMMDD" and returns,
// if successful, a slice of CallRecords.
func GetCallsByDate(start, end string) ([]CallRecord, error) {
	r, err := callReader(start, end)
	if err != nil {
		return nil, err
	}

	callsFromFile := make([]callRecordFromFile, 0)
	err = objcsv.ReadCSV(r, &callsFromFile)
	if err != nil {
		return nil, err
	}

	calls := make([]CallRecord, 0)
	calls, err = batchConvert(callsFromFile)
	if err != nil {
		return nil, err
	}

	return calls, nil
}

func rangeIntoDays(start, end string) ([]timeBlock, error) {
	sT, err := time.Parse("20060102", start)
	if err != nil {
		return nil, err
	}

	eT, err := time.Parse("20060102", end)
	if err != nil {
		return nil, err
	}

	days := int(eT.Sub(sT).Hours())/24 + 1
	dates := make([]timeBlock, days)
	for i := 0; i < days; i++ {
		durString := fmt.Sprintf("%dh", 24*i)
		d, err := time.ParseDuration(durString)
		if err != nil {
			return nil, err
		}

		dayBegin := sT.Add(d)
		dayEnd := dayBegin.Add(time.Duration(86399000000000))
		dates[i] = timeBlock{dayBegin, dayEnd}
	}
	return dates, nil
}

func graphOutput(calls []CallRecord, graphType string, toCSV, toSVG bool) error {
	var data CallGraph

	switch graphType {
	case "duration":
		data = GraphByDuration(calls)
	case "agent":
		data = GraphByAgent(calls)
	case "hour":
		data = GraphByHour(calls)
	default:
		return fmt.Errorf("invalid graphType specifed: %v", graphType)
	}

	if toCSV {
		ds := time.Now().Format("01-02-06_15:04:05")
        filePath := fmt.Sprintf("call_graph_%v.%v", ds, "csv")
		err := WriteToCSV(data, filePath)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote results to file %v\n\n", filePath)
	}
	if toSVG {
        filePath := fmt.Sprintf("calls_%v_%v_%v.%v", graphType, start, pathIncrement, "svg")
		err := WriteToSVG(data, filePath)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote results to file %v\n\n", filePath)
        pathIncrement++
	}

	graph, err := Draw(data)
	if err != nil {
		return err
	}
	fmt.Println(graph)
	return nil
}
