package internal

import (
	"reflect"
	"testing"
	"time"
)

func TestAggregateByTest(t *testing.T) {
	cl := ChangeList{
		Id:               "master@1",
		PrevId:           "",
		DependentTestIds: []string{"TEST1", "TEST2"},
		Timestamp:        time.Now(),
		TestRuns: []TestRun{{
			BuildId: "build-1",
			TestResults: map[string]TestResult{
				"test1": {Runs: 1, Fails: 0},
				"test2": {Runs: 1, Fails: 1}},
			OutputUrl: "build-1.build",
			Timestamp: time.Now(),
		}, {
			BuildId: "build-2",
			TestResults: map[string]TestResult{
				"test1": {Runs: 1, Fails: 0},
				"test2": {Runs: 1, Fails: 0},
				"test3": {Runs: 1, Fails: 0}},
			OutputUrl: "build-2.build",
			Timestamp: time.Now(),
		}},
	}

	aggregateMap := cl.AggregateByTest()

	expectedAggregateMap := map[string]TestResult{
		"test1": {
			Runs:  2,
			Fails: 0,
		},
		"test2": {
			Runs:  2,
			Fails: 1,
		},
		"test3": {
			Runs:  1,
			Fails: 0,
		},
	}

	if !reflect.DeepEqual(aggregateMap, expectedAggregateMap) {
		t.Fatalf("not equal expected=%s actual=%s", aggregateMap, expectedAggregateMap)
	}
}
