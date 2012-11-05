package main

import (
	"reflect"
	"testing"
)

var testRecords = []CallRecord{
	CallRecord{"2012-10-01 12:00:00", "Customer Care", "5559426524", 5, 5550001111},
	CallRecord{"2012-10-03 16:21:00", "Customer Care", "5559192426", 5, 5550001112},
	CallRecord{"2012-10-03 04:24:00", "Customer Care", "5556842575", 0, 0},
	CallRecord{"2012-10-03 12:00:00", "POS Specialist", "5324205192", 5, 5550001111},
	CallRecord{"2012-10-04 23:03:00", "Customer Care", "8344988928", 5, 5550001113},
	CallRecord{"2012-10-04 13:02:00", "POS Specialist", "4549091415", 0, 0},
	CallRecord{"2012-10-04 13:02:00", "I don't even know", "4549091415", 7, 5550001111},
}

func TestIsCustomerCare(t *testing.T) {
	expected := []bool{true, true, true, false, true, false, false}
	for i, call := range testRecords {
		if call.isCustomerCare() != expected[i] {
			t.Errorf("expected %v, got %v for testRecords[%v] from isCustomerCare()", expected[i], call.isCustomerCare(), i)
		}
	}
}

func TestIsMissed(t *testing.T) {
	expected := []bool{false, false, true, false, false, true, false}
	for i, call := range testRecords {
		if call.isMissed() != expected[i] {
			t.Errorf("expected %v, got %v for testRecords[%v] from isMissed()", expected[i], call.isMissed(), i)
		}
	}
}

func TestFilter(t *testing.T) {
	graph := CallGraph{ByAgent, testRecords}
	graph.filter()
	for _, call := range graph.Records {
		if !call.isCustomerCare() {
			t.Errorf("expected to find no non-Customer Care calls, but found\n%v", call)
		}
	}
}

func TestDistribution(t *testing.T) {
	distByHour, err := CallGraph{ByHour, testRecords}.distribution()
	if err != nil {
		t.Errorf("distribution(): %v\n", err)
	}
	distByDuration, err := CallGraph{ByDuration, testRecords}.distribution()
	if err != nil {
		t.Errorf("distribution(): %v\n", err)
	}
	distByAgent, err := CallGraph{ByAgent, testRecords}.distribution()
	if err != nil {
		t.Errorf("distribution(): %v\n", err)
	}

	expectedByHour := map[int64]int{12: 2, 13: 2, 16: 1, 23: 1, 4: 1}
	expectedByDuration := map[int64]int{5: 4, 0: 2, 7: 1}
	expectedByAgent := map[int64]int{5550001111: 3, 5550001112: 1, 5550001113: 1, 0: 2}

	if !reflect.DeepEqual(distByHour, expectedByHour) {
		t.Errorf("distribution by hour: expected %v, got %v", expectedByHour, distByHour)
	}
	if !reflect.DeepEqual(distByDuration, expectedByDuration) {
		t.Errorf("distribution by duration: expected %v, got %v", expectedByDuration, distByDuration)
	}
	if !reflect.DeepEqual(distByAgent, expectedByAgent) {
		t.Errorf("distribution by agent: expected %v, got %v", expectedByAgent, distByAgent)
	}
}
