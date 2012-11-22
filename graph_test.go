package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var graphByHour = GraphByHour(testRecords)
var graphByDuration = GraphByDuration(testRecords)
var graphByAgent = GraphByAgent(testRecords)

// Tests
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

	expectedByHour := map[int64][]CallRecord{
		12: []CallRecord{expectedRecords[0], expectedRecords[3]},
		13: []CallRecord{expectedRecords[5], expectedRecords[6]},
		16: []CallRecord{expectedRecords[1]},
		23: []CallRecord{expectedRecords[4]},
		4:  []CallRecord{expectedRecords[2]},
	}
	expectedByDuration := map[int64][]CallRecord{
		5: []CallRecord{expectedRecords[0], expectedRecords[1], expectedRecords[3], expectedRecords[4]},
		0: []CallRecord{expectedRecords[2], expectedRecords[5]},
		7: []CallRecord{expectedRecords[6]},
	}
	expectedByAgent := map[int64][]CallRecord{
		20: []CallRecord{expectedRecords[0], expectedRecords[3], expectedRecords[6]},
		21: []CallRecord{expectedRecords[1]},
		22: []CallRecord{expectedRecords[4]},
		0:  []CallRecord{expectedRecords[2], expectedRecords[5]},
	}

	if !reflect.DeepEqual(distByHour, expectedByHour) {
		t.Errorf("distribution by hour: expected %v, got %v", expectedByHour, distByHour)
		for k, v := range distByHour {
			fmt.Printf("%v\n%v\n%v\n\n", k, expectedByHour[k], v)
		}
	}
	if !reflect.DeepEqual(distByDuration, expectedByDuration) {
		t.Errorf("distribution by duration: expected %v, got %v", expectedByDuration, distByDuration)
	}
	if !reflect.DeepEqual(distByAgent, expectedByAgent) {
		t.Errorf("distribution by agent: expected %v, got %v", expectedByAgent, distByAgent)
	}
}

func TestDraw(t *testing.T) {
    expected := make([]string, 3)
    actual := make([]string, 3)
    testRecs := []CallGraph{graphByHour, graphByDuration, graphByAgent}

	expected[0] = readTestFile("tests/byhour.txt")
	expected[1] = readTestFile("tests/byduration.txt")
	expected[2] = readTestFile("tests/byagent.txt")

    for i := 0; i < 3; i++ {
        var err error
        actual[i], err = Draw(testRecs[i])
        if err != nil {
            t.Errorf("draw error: %v\n", err)
        }
        if actual[i] != expected[i] {
            t.Errorf("graph mismatch: expected:\n%v\nbut got:\n%v\n", expected[i], actual[i])
        }
    }
}

func readTestFile(inS string) string {
	f, _ := os.Open(inS)
	defer f.Close()
	outS, _ := ioutil.ReadAll(f)
	// Exclude last character (EOF)
	return string(outS)[:len(string(outS))-1]
}
