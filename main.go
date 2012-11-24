package main

import (
	"fmt"
)

var (
	filterCriteria = map[string](func(CallRecord) bool){
		"customer care":   func(cr CallRecord) bool { return cr.IsCustomerCare },
		"answered":        func(cr CallRecord) bool { return !cr.IsMissed },
		"during business": func(cr CallRecord) bool { return cr.DuringHours },
	}
)

func main() {
	err := handleArgs()
	if err != nil {
		fmt.Printf("handleArgs error: %v\n", err)
		return
	}

	// The output of this function depends on the toFile flag.
	calls, err := getCallsByDate(start, end)
	if err != nil {
		fmt.Printf("getCallsByDate error: %v\n", err)
		return
	}

	calls = Filter(calls, filterCriteria["customer care"])

	var data CallGraph

	switch *graphType {
	case "duration":
		data = GraphByDuration(calls)
	case "agent":
		data = GraphByAgent(calls)
	case "hour":
		data = GraphByHour(calls)
	default:
		fmt.Printf("invalid graphType specifed: %v", *graphType)
		return
	}

	if err := graphOutput(data); err != nil {
		fmt.Printf("graphOutput error: %v\n", err)
		return
	}
}
