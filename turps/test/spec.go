package test

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
	pb "turps/api"
	"turps/internal"
)

type DSL interface {
	given_a_turps_server()
	given_a_change(ref string, changeList *pb.ChangeList)
	given_a_test_run(ref string, run *pb.TestRun)
	when_the_change_is_saved(ref string)
	when_the_test_run_is_saved(ref string)
	when_the_change_is_fetched(changeListId string)
	then_the_change_should_be(changeList *pb.ChangeList)
}

// TruncatedNow truncates the current time to microsecond resolution, so that it matches the PostgreSQL TIMESTAMP datatype.
func TruncatedNow() *timestamp.Timestamp {
	tz := time.Now().Truncate(time.Microsecond)
	pbTz, _ := ptypes.TimestampProto(tz)
	return pbTz
}

func ShouldCreateAndFetchChangeList(dsl DSL) {
	tz := TruncatedNow()
	clID := internal.RandomStringOfLength(10)
	dsl.given_a_turps_server()
	dsl.given_a_change("A", &pb.ChangeList{
		ChangeListId: clID,
		Tz:           tz,
	})
	dsl.when_the_change_is_saved("A")
	dsl.when_the_change_is_fetched(clID)
	dsl.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           tz,
	})
}

func ShouldUpdateOnDoubleUpsert(dsl DSL) {
	tz1 := TruncatedNow()
	tz2 := TruncatedNow()
	clID := internal.RandomStringOfLength(10)

	dsl.given_a_turps_server()
	dsl.given_a_change("A", &pb.ChangeList{
		ChangeListId: clID,
		Tz:           tz1,
		TestIds:      []string{"test-1"},
	})
	dsl.given_a_change("B", &pb.ChangeList{
		ChangeListId: clID,
		Tz:           tz2,
		TestIds:      []string{"test-2"},
	})
	dsl.when_the_change_is_saved("A")
	dsl.when_the_change_is_saved("B")
	dsl.when_the_change_is_fetched(clID)
	dsl.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           tz2,
		TestIds:      []string{"test-2"}, //changed
	})
}

func ShouldUpdateChangeListWithSingleTestRun(dsl DSL) {
	changeListTz := TruncatedNow()
	runTz := TruncatedNow()
	clID := internal.RandomStringOfLength(10)

	dsl.given_a_turps_server()
	dsl.given_a_change("change-A", &pb.ChangeList{
		ChangeListId: clID,
		Tz:           changeListTz,
	})
	dsl.given_a_test_run("run-B", &pb.TestRun{
		Id:           "run-1",
		ChangeListId: clID,
		OutputUrl:    "run-1",
		Tz:           runTz,
		TestResult: map[string]*pb.TestResult{
			"test-1": {NumFails: 1, NumRuns: 10},
		},
	})
	dsl.when_the_change_is_saved("change-A")
	dsl.when_the_test_run_is_saved("run-B")
	dsl.when_the_change_is_fetched(clID)
	dsl.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           changeListTz,
		TestRun: []*pb.TestRun{
			{
				Id:           "run-1",
				ChangeListId: clID,
				OutputUrl:    "run-1",
				Tz:           runTz,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 10},
				},
			},
		},
	})
}
func ShouldUpdateChangeListWithDoubleTestRun(dsl DSL) {
	changeListTz := TruncatedNow()
	clID := internal.RandomStringOfLength(10)

	dsl.given_a_turps_server()
	dsl.given_a_change("change-A", &pb.ChangeList{
		ChangeListId: clID,
		Tz:           changeListTz,
	})

	runTz1 := TruncatedNow()
	runTz2 := TruncatedNow()
	dsl.given_a_test_run("run-B", &pb.TestRun{
		Id:           "run-1",
		ChangeListId: clID,
		OutputUrl:    "run-1",
		Tz:           runTz1,
		TestResult: map[string]*pb.TestResult{
			"test-1": {NumFails: 1, NumRuns: 10},
		},
	})
	dsl.given_a_test_run("run-C",&pb.TestRun{
		Id:           "run-2",
		ChangeListId: clID,
		OutputUrl:    "run-2",
		Tz:           runTz2,
		TestResult: map[string]*pb.TestResult{
			"test-1": {NumFails: 1, NumRuns: 10},
			"test-2": {NumFails: 1, NumRuns: 1},
		},
	})
	dsl.when_the_change_is_saved("change-A")
	dsl.when_the_test_run_is_saved("run-B")
	dsl.when_the_test_run_is_saved("run-C")
	dsl.when_the_change_is_fetched(clID)
	dsl.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           changeListTz,
		TestRun: []*pb.TestRun{
			{
				Id:           "run-1",
				ChangeListId: clID,
				OutputUrl:    "run-1",
				Tz:           runTz1,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 10},
				},
			},
			{
				Id:           "run-2",
				ChangeListId: clID,
				OutputUrl:    "run-2",
				Tz:           runTz2,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 10},
					"test-2": {NumFails: 1, NumRuns: 1},
				},
			},
		},
	})
}

func ShouldUpsertTestRun(dsl DSL) {
	clID := internal.RandomStringOfLength(10)
	changeListTz := TruncatedNow()

	dsl.given_a_turps_server()
	dsl.given_a_change("change-A",&pb.ChangeList{
		ChangeListId: clID,
		Tz:           changeListTz,
	})
	runTz1 := TruncatedNow()
	dsl.given_a_test_run("run-B",&pb.TestRun{
		Id:           "run-1",
		ChangeListId: clID,
		OutputUrl:    "run-1",
		Tz:           runTz1,
		TestResult: map[string]*pb.TestResult{
			"test-1": {NumFails: 1, NumRuns: 10},
		},
	})
	runTz2 := TruncatedNow()
	dsl.given_a_test_run("run-C",&pb.TestRun{
		Id:           "run-1",
		ChangeListId: clID,
		OutputUrl:    "run-1",
		Tz:           runTz2,
		TestResult: map[string]*pb.TestResult{
			"test-1": {NumFails: 1, NumRuns: 20},
		},
	})

	dsl.when_the_change_is_saved("change-A")
	dsl.when_the_test_run_is_saved("run-B")
	dsl.when_the_test_run_is_saved("run-C")
	dsl.when_the_change_is_fetched(clID)
	dsl.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           changeListTz,
		TestRun: []*pb.TestRun{
			{
				Id:           "run-1",
				ChangeListId: clID,
				OutputUrl:    "run-1",
				Tz:           runTz2,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 20},
				},
			},
		},
	})
}
