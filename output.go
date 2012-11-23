package main

import (
	"fmt"
	"github.com/MattSPalmer/objcsv"
	"time"
)

func getCallsByDate(start, end string) ([]CallRecord, error) {
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

func graphOutput(data CallGraph) error {
	if *toFile {
		ds := time.Now().Format("01-02-06_15:04:05")
		filePath := fmt.Sprintf("call_graph_%v.%v", ds, extension)
		err := WriteToCSV(data, filePath)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote results to file %v\n\n", filePath)
	}

	graph, err := Draw(data)
	if err != nil {
		return err
	}
	fmt.Println(graph)
	return nil
}
