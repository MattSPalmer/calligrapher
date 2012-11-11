package main

import (
	"io/ioutil"
	"os"
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

var graphByHour = GraphByHour(testRecords)
var graphByDuration = GraphByDuration(testRecords)
var graphByAgent = GraphByAgent(testRecords)

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
	testCopy := make([]CallRecord, 0)
	copy(testRecords, testCopy)

	filter(&testCopy)
	graph := GraphByAgent(testCopy)

	for _, call := range graph {
		if !call.isCustomerCare() {
			t.Errorf("expected to find no non-Customer Care calls, but found\n%v", call)
		}
	}
}

func TestDistribution(t *testing.T) {
	distByHour, err := graphByHour.Distribution()
	if err != nil {
		t.Errorf("distribution(): %v\n", err)
	}
	distByDuration, err := graphByDuration.Distribution()
	if err != nil {
		t.Errorf("distribution(): %v\n", err)
	}
	distByAgent, err := graphByAgent.Distribution()
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

func TestDraw(t *testing.T) {
	expectedHour := readTestFile("tests/byhour.txt")
	expectedDuration := readTestFile("tests/byduration.txt")
	expectedAgent := readTestFile("tests/byagent.txt")

	hourDraw, err := Draw(graphByHour)
	if err != nil {
		t.Errorf("draw error: %v\n", err)
	}
	durationDraw, err := Draw(graphByDuration)
	if err != nil {
		t.Errorf("draw error: %v\n", err)
	}
	agentDraw, err := Draw(graphByAgent)
	if err != nil {
		t.Errorf("draw error: %v\n", err)
	}

	if hourDraw != expectedHour {
		t.Errorf("graph mismatch: expected:\n%v\nbut got:\n%v\n", expectedHour, hourDraw)
	}
	if durationDraw != expectedDuration {
		t.Errorf("graph mismatch: expected:\n%v\nbut got:\n%v\n", expectedDuration, durationDraw)
	}
	if agentDraw != expectedAgent {
		t.Errorf("graph mismatch: expected:\n%v\nbut got:\n%v\n", expectedAgent, agentDraw)
	}
}

func readTestFile(inS string) string {
	f, _ := os.Open(inS)
	defer f.Close()
	outS, _ := ioutil.ReadAll(f)
    // Exclude last character (EOF)
    return string(outS)[:len(string(outS))-1]
}
