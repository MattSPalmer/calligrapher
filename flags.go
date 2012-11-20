package main

import (
	"flag"
	"fmt"
	"regexp"
	"strings"
)

var (
	start, end                        string
	showAgent, showDuration, showHour bool
)

func handleArgs() error {
	var e error

	flag.BoolVar(&showAgent, "a", false, "show a breakdown of calls by agent")
	flag.BoolVar(&showDuration, "d", false, "show a breakdown of calls by duration")
	flag.BoolVar(&showHour, "c", false, "show a breakdown of calls by hour")

	flag.Parse()

	if !showHour && !showAgent && !showDuration {
		showHour = true
	}

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
	pattern, err := regexp.Compile("(20[0-9]{2})?[/-]?([0-3][0-9])[/-]?([0-9]{2})")
	if err != nil {
		return "", err
	}
	if e := pattern.FindStringSubmatch(s); len(e) > 2 && len(e) < 5 {
		s = strings.Join(e[1:len(e)], "")
		if len(s) == 8 {
			return s, nil
		} else {
			return "2012" + s, nil
		}
	}
	return "", fmt.Errorf("formatDateString: with string %v, something happened that should never happen.", s)
}
