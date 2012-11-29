package main

import (
	"fmt"
	"time"
)

var (
	// filterCriteria stores filter functions that return a boolean based on
	// certain CallRecord values.
	filterCriteria = map[string](func(CallRecord) bool){
		"customer care":   func(cr CallRecord) bool { return cr.IsCustomerCare },
		"answered":        func(cr CallRecord) bool { return !cr.IsMissed },
		"during business": func(cr CallRecord) bool { return cr.DuringHours },
	}
)

func main() {
	err := handleArgs()
	if err != nil {
		fmt.Printf("argument error: %v\n", err)
		return
	}

	calls, err := GetCallsByDate(start, end)
	if err != nil {
		fmt.Printf("getCallsByDate error: %v\n", err)
		return
	}

	calls = Filter(calls, filterCriteria["customer care"])

	if *test {
		t := time.Now()
		m, err := mapCallTime(calls)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		durationStackGraph(3000, 500, m, "moo.svg")
		fmt.Printf("%v\n", time.Since(t))
		return
	}

	if *byDate {
		days, err := rangeIntoDays(start, end)
		if err != nil {
			fmt.Println(err)
			return
		}

		for i, day := range days {
			dayData := Filter(calls, func(cr CallRecord) bool {
				return cr.Created_at.After(day.start) && cr.Created_at.Before(day.end)
			})
			fmt.Println(day.start.Format("Monday, Jan 2 2006"))
			if err := graphOutput(dayData, *graphType, *toCSV, *toSVG); err != nil {
				fmt.Printf("graphOutput error: %v\n", err)
				return
			}
			if i != len(days)-1 {
				fmt.Scanln()
			}
		}
	} else {
		if err := graphOutput(calls, *graphType, *toCSV, *toSVG); err != nil {
			fmt.Printf("graphOutput error: %v\n", err)
			return
		}
	}
}
