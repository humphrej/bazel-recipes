package grpc

import (
	"context"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"reflect"
	"testing"
	pb "turps/api"
	"turps/internal"
	"turps/internal/postgres"
	pg_testing "turps/internal/postgres/testing"
)

type world struct {
	ctx              context.Context
	server           *Server
	upsertResponse   *pb.UpsertChangeListResponse
	changeListId     string
	Testing          *testing.T
	testRunResponses []*pb.UpsertTestRunResponse
}

func Test_ShouldCreateAndFetchChangeList(t *testing.T) {
	world := world{Testing: t}
	given_a_turps_server(&world)
	and_a_change(&world, &pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           ptypes.TimestampNow(),
	})
	then_the_change_should_be(&world, &pb.ChangeList{
		ChangeListId: world.changeListId,
		Tz:           world.upsertResponse.ChangeList.Tz,
	})
}

func Test_ShouldUpdateOnDoubleUpsert(t *testing.T) {
	world := world{Testing: t}
	given_a_turps_server(&world)
	and_a_change(&world, &pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           ptypes.TimestampNow(),
		TestIds:      []string{"test-1"},
	})
	and_a_change(&world, &pb.ChangeList{
		ChangeListId: world.changeListId,
		Tz:           ptypes.TimestampNow(),
		TestIds:      []string{"test-2"},
	})
	then_the_change_should_be(&world, &pb.ChangeList{
		ChangeListId: world.changeListId,
		Tz:           world.upsertResponse.ChangeList.Tz,
		TestIds:      []string{"test-2"}, //changed
	})
}

func Test_ShouldUpdateChangeListWithSingleTestRun(t *testing.T) {
	world := world{Testing: t}
	given_a_turps_server(&world)
	and_a_change(&world, &pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           ptypes.TimestampNow(),
	})
	when_tests_are_run(&world, []*pb.TestRun{
		{
			Id:           "run-1",
			ChangeListId: world.changeListId,
			OutputUrl:    "run-1",
			Tz:           ptypes.TimestampNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 10},
			},
		},
	},
	)
	then_the_change_should_be(&world, &pb.ChangeList{
		ChangeListId: world.changeListId,
		Tz:           world.upsertResponse.ChangeList.Tz,
		TestRun: []*pb.TestRun{
			{
				Id:           "run-1",
				ChangeListId: world.changeListId,
				OutputUrl:    "run-1",
				Tz:           world.testRunResponses[0].TestRun.Tz,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 10},
				},
			},
		},
	})
}

func Test_ShouldUpdateChangeListWithDoubleTestRun(t *testing.T) {
	world := world{Testing: t}
	given_a_turps_server(&world)
	and_a_change(&world, &pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           ptypes.TimestampNow(),
	})
	when_tests_are_run(&world, []*pb.TestRun{
		{
			Id:           "run-1",
			ChangeListId: world.changeListId,
			OutputUrl:    "run-1",
			Tz:           ptypes.TimestampNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 10},
			},
		},
		{
			Id:           "run-2",
			ChangeListId: world.changeListId,
			OutputUrl:    "run-2",
			Tz:           ptypes.TimestampNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 10},
				"test-2": {NumFails: 1, NumRuns: 1},
			},
		},
	},
	)
	then_the_change_should_be(&world, &pb.ChangeList{
		ChangeListId: world.changeListId,
		Tz:           world.upsertResponse.ChangeList.Tz,
		TestRun: []*pb.TestRun{
			{
				Id:           "run-1",
				ChangeListId: world.changeListId,
				OutputUrl:    "run-1",
				Tz:           world.testRunResponses[0].TestRun.Tz,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 10},
				},
			},
			{
				Id:           "run-2",
				ChangeListId: world.changeListId,
				OutputUrl:    "run-2",
				Tz:           world.testRunResponses[1].TestRun.Tz,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 10},
					"test-2": {NumFails: 1, NumRuns: 1},
				},
			},
		},
	})
}
func Test_ShouldUpsertChangeList(t *testing.T) {
	world := world{Testing: t}
	given_a_turps_server(&world)
	and_a_change(&world, &pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           ptypes.TimestampNow(),
	})
	when_tests_are_run(&world, []*pb.TestRun{
		{
			Id:           "run-1",
			ChangeListId: world.changeListId,
			OutputUrl:    "run-1",
			Tz:           ptypes.TimestampNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 10},
			},
		},
		{
			Id:           "run-1",
			ChangeListId: world.changeListId,
			OutputUrl:    "run-1",
			Tz:           ptypes.TimestampNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 20},
			},
		},
	},
	)
	then_the_change_should_be(&world, &pb.ChangeList{
		ChangeListId: world.changeListId,
		Tz:           world.upsertResponse.ChangeList.Tz,
		TestRun: []*pb.TestRun{
			{
				Id:           "run-1",
				ChangeListId: world.changeListId,
				OutputUrl:    "run-1",
				Tz:           world.testRunResponses[1].TestRun.Tz,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 20},
				},
			},
		},
	})
}

func given_a_turps_server(w *world) {
	w.ctx = context.Background()

	pool, err := pg_testing.NewPool(w.ctx)
	if err != nil {
		glog.Fatalf("failed to create connection pool: %v", err)
	}

	repo := postgres.ChangeListStorage{
		Pool: pool,
	}
	w.server = NewServer(repo)
}

func and_a_change(w *world, c *pb.ChangeList) {
	var err error
	w.upsertResponse, err = w.server.UpsertChangeList(w.ctx, &pb.UpsertChangeListRequest{ChangeList: c})
	if err != nil {
		glog.Fatalf("failed to store change list: %v", err)
	}
	w.changeListId = w.upsertResponse.ChangeList.ChangeListId
}
func when_tests_are_run(w *world, runs []*pb.TestRun) {
	for _, run := range runs {
		upsertTestRunResponse, err := w.server.UpsertTestResult(w.ctx, &pb.UpsertTestRunRequest{TestRun: run})
		if err != nil {
			w.Testing.Fatalf("failed to store test run: %v", err)
		}
		w.testRunResponses = append(w.testRunResponses, upsertTestRunResponse)
	}
}
func then_the_change_should_be(w *world, expected *pb.ChangeList) {
	var err error
	fetchResponse, err := w.server.GetChangeList(w.ctx, &pb.GetChangeListRequest{
		ChangeListId: w.changeListId,
	})
	if err != nil {
		w.Testing.Fatalf("Failed reading change %v", err)
	}

	if !reflect.DeepEqual(expected, fetchResponse.ChangeList) {
		w.Testing.Fatalf("Value fetched does not match value stored.\nexpected=%s\n  actual=%s", expected, fetchResponse.ChangeList)
	}
}
