package main

import "flag"

var start, end string

func flagInit() {
	flag.StringVar(&start, "s", "", "start date (\"YYYYMMDD\")")
	flag.StringVar(&end, "e", "", "end date (\"YYYYMMDD\")")

	flag.Parse()
}
