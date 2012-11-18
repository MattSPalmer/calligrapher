package main

import "flag"

var (
	start = flag.String("s", "", "start date (\"YYYYMMDD\")")
	end   = flag.String("e", "", "end date (\"YYYYMMDD\")")

	all = flag.String("", "", "")
)
