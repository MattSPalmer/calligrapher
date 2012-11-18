package main

import (
	"reflect"
	"testing"
	"time"
)

// Declarations
var testRecordsFromFile = []callRecordFromFile{
	callRecordFromFile{"2012-10-01 12:00:00", "Customer Care", "5559426524", 5, 5550001111},
	callRecordFromFile{"2012-10-03 16:21:00", "Customer Care", "5559192426", 5, 5550001112},
	callRecordFromFile{"2012-10-03 04:24:00", "Customer Care", "5556842575", 0, 0},
	callRecordFromFile{"2012-10-03 12:00:00", "POS Specialist", "5324205192", 5, 5550001111},
	callRecordFromFile{"2012-10-04 23:03:00", "Customer Care", "8344988928", 5, 5550001113},
	callRecordFromFile{"2012-10-04 13:02:00", "POS Specialist", "4549091415", 0, 0},
	callRecordFromFile{"2012-10-04 13:02:00", "I don't even know", "4549091415", 7, 5550001111},
}

var expectedRecords = []CallRecord{
	// not missed, is Customer Care, during hours
	CallRecord{time.Date(2012, time.October, 1, 12, 0, 0, 0, time.UTC), false, true, true, "5559426524", 5, 20},
	// not missed, is Customer Care, during hours
	CallRecord{time.Date(2012, time.October, 3, 16, 21, 0, 0, time.UTC), false, true, true, "5559192426", 5, 21},
	// missed, is Customer Care, not during hours
	CallRecord{time.Date(2012, time.October, 3, 4, 24, 0, 0, time.UTC), true, true, false, "5556842575", 0, 0},
	// not missed, not Customer Care, during hours
	CallRecord{time.Date(2012, time.October, 3, 12, 0, 0, 0, time.UTC), false, false, true, "5324205192", 5, 20},
	// not missed, is Customer Care, not during hours
	CallRecord{time.Date(2012, time.October, 4, 23, 3, 0, 0, time.UTC), false, true, false, "8344988928", 5, 22},
	// missed, not Customer Care, during hours
	CallRecord{time.Date(2012, time.October, 4, 13, 2, 0, 0, time.UTC), true, false, true, "4549091415", 0, 0},
	// not missed, not Customer Care, during hours
	CallRecord{time.Date(2012, time.October, 4, 13, 2, 0, 0, time.UTC), false, false, true, "4549091415", 7, 20},
}

var testRecords, _ = batchConvert(testRecordsFromFile)

func TestIsCustomerCare(t *testing.T) {
	expected := []bool{true, true, true, false, true, false, false}
	for i, call := range testRecordsFromFile {
		if call.isCustomerCare() != expected[i] {
			t.Errorf("expected %v, got %v for testRecords[%v] from isCustomerCare()", expected[i], call.isCustomerCare(), i)
		}
	}
}

func TestIsMissed(t *testing.T) {
	expected := []bool{false, false, true, false, false, true, false}
	for i, call := range testRecordsFromFile {
		if call.isMissed() != expected[i] {
			t.Errorf("expected %v, got %v for testRecords[%v] from isMissed()", expected[i], call.isMissed(), i)
		}
	}
}

func TestDuring(t *testing.T) {
	testBlock := timeBlock{
		time.Date(0, time.January, 1, 12, 0, 0, 0, time.UTC),
		time.Date(0, time.January, 1, 14, 0, 0, 0, time.UTC),
	}
	testMoment1 := time.Date(2012, time.March, 14, 13, 0, 0, 0, time.UTC)
	testMoment2 := time.Date(2012, time.February, 14, 16, 0, 0, 0, time.UTC)
	if !testBlock.during(testMoment1) {
		t.Errorf("unexpected value from during(): expected %v but got %v.", true, testBlock.during(testMoment1))
	}
	if testBlock.during(testMoment2) {
		t.Errorf("unexpected value from during(): expected %v but got %v.", false, testBlock.during(testMoment2))
	}
}

func TestConvert(t *testing.T) {
}

func TestFilter(t *testing.T) {
	testSet := expectedRecords
	actual := make([][]CallRecord, 3)
	expectedSets := [][]CallRecord{
		// Were support calls
		[]CallRecord{expectedRecords[0], expectedRecords[1], expectedRecords[2], expectedRecords[4]},
		// Were not missed
		[]CallRecord{expectedRecords[0], expectedRecords[1], expectedRecords[3], expectedRecords[4], expectedRecords[6]},
		// Occurred during business hours
		[]CallRecord{expectedRecords[0], expectedRecords[1], expectedRecords[3], expectedRecords[5], expectedRecords[6]},
	}

	actual[0] = Filter(testSet, func(cr CallRecord) bool {
		return cr.IsCustomerCare
	})

	actual[1] = Filter(testSet, func(cr CallRecord) bool {
		return !cr.IsMissed
	})

	actual[2] = Filter(testSet, func(cr CallRecord) bool {
		day := cr.Created_at.Weekday()
		return schedule[day].during(cr.Created_at)
	})

	for i := 0; i < 3; i++ {
		if !reflect.DeepEqual(actual[i], expectedSets[i]) {
			t.Errorf("unexpected value for filter level %d: expected\n", i+1)
			for i, elem := range expectedSets[i] {
				t.Logf("%02d: %v\n", i, elem)
			}
			t.Logf("but got\n")
			for i, elem := range actual[i] {
				t.Errorf("%02d: %v\n", i, elem)
			}
			t.Logf("\n")
		}
	}
}
