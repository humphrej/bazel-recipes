package test

import (
	"context"
	"reflect"
	"testing"
	pb "turps/api"
	"turps/internal/grpc"
	"turps/internal/postgres"
	pg_testing "turps/internal/postgres/testing"
)

type apiWorld struct {
	ctx              context.Context
	server           *grpc.Server
	upsertResponse   *pb.UpsertChangeListResponse
	changeListId     string
	Testing          *testing.T
	testRunResponses []*pb.UpsertTestRunResponse
}

func (w *apiWorld) given_a_turps_server() {
	w.ctx = context.Background()

	pool, err := pg_testing.NewPool(w.ctx)
	if err != nil {
		w.Testing.Fatalf("failed to create connection pool: %v", err)
	}

	repo := postgres.ChangeListStorage{
		Pool: pool,
	}
	w.server = grpc.NewServer(repo)
}

func (w *apiWorld) given_a_change(c *pb.ChangeList) {
	var err error
	w.upsertResponse, err = w.server.UpsertChangeList(w.ctx, &pb.UpsertChangeListRequest{ChangeList: c})
	if err != nil {
		w.Testing.Fatalf("failed to store change list: %v", err)
	}
	w.changeListId = w.upsertResponse.ChangeList.ChangeListId
}
func (w *apiWorld) when_tests_are_run(runs []*pb.TestRun) {
	for _, run := range runs {
		upsertTestRunResponse, err := w.server.UpsertTestResult(w.ctx, &pb.UpsertTestRunRequest{TestRun: run})
		if err != nil {
			w.Testing.Fatalf("failed to store test run: %v", err)
		}
		w.testRunResponses = append(w.testRunResponses, upsertTestRunResponse)
	}
}
func (w *apiWorld) then_the_change_should_be(expected *pb.ChangeList) {
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
