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
	given_a_change(changeList *pb.ChangeList)
	when_tests_are_run(runs []*pb.TestRun)
	then_the_change_should_be(changeList *pb.ChangeList)
}

func TruncatedNow() *timestamp.Timestamp {
	tz := time.Now().Truncate(time.Microsecond)
	pbTz, _ := ptypes.TimestampProto(tz)
	return pbTz
}

func ShouldCreateAndFetchChangeList(dsl DSL) {
	tz := TruncatedNow()
	clID := internal.RandomStringOfLength(10)
	dsl.given_a_turps_server()
	dsl.given_a_change(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           tz,
	})
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
	dsl.given_a_change(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           tz1,
		TestIds:      []string{"test-1"},
	})
	dsl.given_a_change(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           tz2,
		TestIds:      []string{"test-2"},
	})
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
	dsl.given_a_change(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           changeListTz,
	})
	dsl.when_tests_are_run([]*pb.TestRun{
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
	)
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
	dsl.given_a_change(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           changeListTz,
	})

	runTz1 := TruncatedNow()
	runTz2 := TruncatedNow()

	dsl.when_tests_are_run([]*pb.TestRun{
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
	)
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
	dsl.given_a_change(&pb.ChangeList{
		ChangeListId: clID,
		Tz:           changeListTz,
	})
	runTz1 := TruncatedNow()
	runTz2 := TruncatedNow()
	dsl.when_tests_are_run([]*pb.TestRun{
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
			Id:           "run-1",
			ChangeListId: clID,
			OutputUrl:    "run-1",
			Tz:           runTz2,
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 20},
			},
		},
	},
	)
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
