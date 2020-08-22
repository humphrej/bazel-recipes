package grpc

import (
	"context"
	"github.com/golang/glog"
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
	w := world{Testing: t}
	w.given_a_turps_server()
	w.given_a_change(&pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           TruncatedNow(),
	})
	w.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: w.changeListId,
		Tz:           w.upsertResponse.ChangeList.Tz,
	})
}

func Test_ShouldUpdateOnDoubleUpsert(t *testing.T) {
	w := world{Testing: t}
	w.given_a_turps_server()
	w.given_a_change(&pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           TruncatedNow(),
		TestIds:      []string{"test-1"},
	})
	w.given_a_change(&pb.ChangeList{
		ChangeListId: w.changeListId,
		Tz:           TruncatedNow(),
		TestIds:      []string{"test-2"},
	})
	w.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: w.changeListId,
		Tz:           w.upsertResponse.ChangeList.Tz,
		TestIds:      []string{"test-2"}, //changed
	})
}

func Test_ShouldUpdateChangeListWithSingleTestRun(t *testing.T) {
	w := world{Testing: t}
	w.given_a_turps_server()
	w.given_a_change(&pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           TruncatedNow(),
	})
	w.when_tests_are_run([]*pb.TestRun{
		{
			Id:           "run-1",
			ChangeListId: w.changeListId,
			OutputUrl:    "run-1",
			Tz:           TruncatedNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 10},
			},
		},
	},
	)
	w.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: w.changeListId,
		Tz:           w.upsertResponse.ChangeList.Tz,
		TestRun: []*pb.TestRun{
			{
				Id:           "run-1",
				ChangeListId: w.changeListId,
				OutputUrl:    "run-1",
				Tz:           w.testRunResponses[0].TestRun.Tz,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 10},
				},
			},
		},
	})
}

func Test_ShouldUpdateChangeListWithDoubleTestRun(t *testing.T) {
	w := world{Testing: t}
	w.given_a_turps_server()
	w.given_a_change(&pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           TruncatedNow(),
	})
	w.when_tests_are_run([]*pb.TestRun{
		{
			Id:           "run-1",
			ChangeListId: w.changeListId,
			OutputUrl:    "run-1",
			Tz:           TruncatedNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 10},
			},
		},
		{
			Id:           "run-2",
			ChangeListId: w.changeListId,
			OutputUrl:    "run-2",
			Tz:           TruncatedNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 10},
				"test-2": {NumFails: 1, NumRuns: 1},
			},
		},
	},
	)
	w.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: w.changeListId,
		Tz:           w.upsertResponse.ChangeList.Tz,
		TestRun: []*pb.TestRun{
			{
				Id:           "run-1",
				ChangeListId: w.changeListId,
				OutputUrl:    "run-1",
				Tz:           w.testRunResponses[0].TestRun.Tz,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 10},
				},
			},
			{
				Id:           "run-2",
				ChangeListId: w.changeListId,
				OutputUrl:    "run-2",
				Tz:           w.testRunResponses[1].TestRun.Tz,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 10},
					"test-2": {NumFails: 1, NumRuns: 1},
				},
			},
		},
	})
}
func Test_ShouldUpsertChangeList(t *testing.T) {
	w := world{Testing: t}
	w.given_a_turps_server()
	w.given_a_change(&pb.ChangeList{
		ChangeListId: internal.RandomStringOfLength(10),
		Tz:           TruncatedNow(),
	})
	w.when_tests_are_run([]*pb.TestRun{
		{
			Id:           "run-1",
			ChangeListId: w.changeListId,
			OutputUrl:    "run-1",
			Tz:           TruncatedNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 10},
			},
		},
		{
			Id:           "run-1",
			ChangeListId: w.changeListId,
			OutputUrl:    "run-1",
			Tz:           TruncatedNow(),
			TestResult: map[string]*pb.TestResult{
				"test-1": {NumFails: 1, NumRuns: 20},
			},
		},
	},
	)
	w.then_the_change_should_be(&pb.ChangeList{
		ChangeListId: w.changeListId,
		Tz:           w.upsertResponse.ChangeList.Tz,
		TestRun: []*pb.TestRun{
			{
				Id:           "run-1",
				ChangeListId: w.changeListId,
				OutputUrl:    "run-1",
				Tz:           w.testRunResponses[1].TestRun.Tz,
				TestResult: map[string]*pb.TestResult{
					"test-1": {NumFails: 1, NumRuns: 20},
				},
			},
		},
	})
}

func (w *world) given_a_turps_server() {
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

func (w *world) given_a_change(c *pb.ChangeList) {
	var err error
	w.upsertResponse, err = w.server.UpsertChangeList(w.ctx, &pb.UpsertChangeListRequest{ChangeList: c})
	if err != nil {
		glog.Fatalf("failed to store change list: %v", err)
	}
	w.changeListId = w.upsertResponse.ChangeList.ChangeListId
}
func (w *world) when_tests_are_run(runs []*pb.TestRun) {
	for _, run := range runs {
		upsertTestRunResponse, err := w.server.UpsertTestResult(w.ctx, &pb.UpsertTestRunRequest{TestRun: run})
		if err != nil {
			w.Testing.Fatalf("failed to store test run: %v", err)
		}
		w.testRunResponses = append(w.testRunResponses, upsertTestRunResponse)
	}
}
func (w *world) then_the_change_should_be(expected *pb.ChangeList) {
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
