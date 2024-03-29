package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
)

var (
	start, end string
	graphType  = flag.String("g", "hour", "graph by agent, duration or call time")
	byDate     = flag.Bool("d", false, "split data into separate days")
	toCSV      = flag.Bool("f", false, "write results to a CSV file")
	toSVG      = flag.Bool("s", false, "create a graph in SVG format")

	testStack     = flag.Bool("t", false, "")
	testDensity   = flag.Bool("n", false, "")
	testCallbacks = flag.Bool("c", false, "")
)

func handleArgs() error {
	var e error

	flag.Parse()

	if start, e = formatDateString(flag.Arg(0)); e != nil {
		return e
	}

	switch flag.NArg() {
	case 1:
		if end, e = formatDateString(flag.Arg(0)); e != nil {
			return e
		}
	case 2:
		if end, e = formatDateString(flag.Arg(1)); e != nil {
			return e
		}
	default:
		return fmt.Errorf("invalid number of args: expected 1 or 2, got %v", flag.NArg())
	}
	return nil
}

func formatDateString(s string) (string, error) {
	pattern, _ := regexp.Compile("(20[0-9]{2})?[/-]?([0-3][0-9])[/-]?([0-9]{2})")
	if e := pattern.FindStringSubmatch(s); len(e) > 2 && len(e) < 5 {
		s = strings.Join(e[1:len(e)], "")
		if len(s) == 8 {
			return s, nil
		} else {
			return "2012" + s, nil
		}
	}
	return "", fmt.Errorf("formatDateString: invalid date string %v", s)
}
