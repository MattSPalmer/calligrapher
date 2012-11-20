package main

import (
	"fmt"
	"regexp"
	"time"
)

var (
	weekendHours = timeBlock{
		time.Date(0, time.January, 1, 10, 0, 0, 0, time.UTC),
		time.Date(0, time.January, 1, 18, 0, 0, 0, time.UTC),
	}
	weekdayHours = timeBlock{
		time.Date(0, time.January, 1, 9, 0, 0, 0, time.UTC),
		time.Date(0, time.January, 1, 20, 0, 0, 0, time.UTC),
	}

	schedule = map[time.Weekday]timeBlock{
		1: weekendHours,
		2: weekdayHours,
		3: weekdayHours,
		4: weekdayHours,
		5: weekdayHours,
		6: weekdayHours,
		7: weekendHours,
	}
)

type timeBlock struct {
	start time.Time
	end   time.Time
}

type callRecordFromFile struct {
	Created_at   string
	ActivityInfo string
	CallerID     string
	Duration     int64
	AgentNumber  int64
}

type CallRecord struct {
	Created_at     time.Time
	IsMissed       bool
	IsCustomerCare bool
	DuringHours    bool
	IncomingNumber string
	Duration       int64
	AgentID        int64
}

func (tb timeBlock) during(test time.Time) bool {
	test = time.Date(0, time.January, 1, test.Hour(), test.Minute(), test.Second(), 0, time.UTC)
	return test.After(tb.start) && test.Before(tb.end)
}

func (cr callRecordFromFile) isCustomerCare() bool {
	query, err := regexp.Compile("Customer Care")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return query.MatchString(cr.ActivityInfo)
}

func (cr callRecordFromFile) isMissed() bool {
	return cr.AgentNumber == 0
}

func (cr callRecordFromFile) convert() (out CallRecord, err error) {
	out = CallRecord{
		IsCustomerCare: cr.isCustomerCare(),
		IsMissed:       cr.isMissed(),
		IncomingNumber: cr.CallerID,
		Duration:       cr.Duration,
		AgentID:        agentsByNumber[cr.AgentNumber],
	}
	if val, ok := agentsByNumber[cr.AgentNumber]; ok {
		out.AgentID = val
	} else {
		out.AgentID = -1
	}
	out.Created_at, err = time.Parse("2006-01-02 15:04:05", cr.Created_at)
	if err != nil {
		return
	}
	out.DuringHours = schedule[out.Created_at.Weekday()].during(out.Created_at)
	return
}

func batchConvert(fromFile []callRecordFromFile) ([]CallRecord, error) {
	out := make([]CallRecord, 0)

	for _, call := range fromFile {
		convertedCall, err := call.convert()
		if err != nil {
			return nil, err
		}
		out = append(out, convertedCall)
	}
	return out, nil
}

func Filter(calls []CallRecord, f func(CallRecord) bool) []CallRecord {
	newCalls := make([]CallRecord, 0)
	for _, call := range calls {
		if f(call) {
			newCalls = append(newCalls, call)
		}
	}
	return newCalls
}
