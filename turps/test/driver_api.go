package test

import (
	"context"
	"google.golang.org/protobuf/proto"
	"testing"
	pb "turps/api"
	"turps/internal/grpc"
	"turps/internal/postgres"
	pg_testing "turps/internal/postgres/testing"
)

type apiWorld struct {
	ctx               context.Context
	server            *grpc.Server
	Testing           *testing.T
	changesMap        map[string]*pb.ChangeList
	testRunsMap       map[string]*pb.TestRun
	lastFetchedChange *pb.ChangeList
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

func (w *apiWorld) given_a_change(ref string, c *pb.ChangeList) {
	w.changesMap[ref] = c
}
func (w *apiWorld) given_a_test_run(ref string, r *pb.TestRun) {
	w.testRunsMap[ref] = r
}

func (w *apiWorld) when_the_change_is_saved(ref string) {
	c := w.changesMap[ref]

	var err error
	_, err = w.server.UpsertChangeList(w.ctx, &pb.UpsertChangeListRequest{ChangeList: c})
	if err != nil {
		w.Testing.Fatalf("failed to store change list: %v", err)
	}
}
func (w *apiWorld) when_the_test_run_is_saved(ref string) {
	run := w.testRunsMap[ref]

	_, err := w.server.UpsertTestResult(w.ctx, &pb.UpsertTestRunRequest{TestRun: run})
	if err != nil {
		w.Testing.Fatalf("failed to store test run: %v", err)
	}
}
func (w *apiWorld) when_the_change_is_fetched(changeListId string) {
	var err error
	fetchResponse, err := w.server.GetChangeList(w.ctx, &pb.GetChangeListRequest{
		ChangeListId: changeListId,
	})
	if err != nil {
		w.Testing.Fatalf("Failed reading change %v", err)
	}
	w.lastFetchedChange = fetchResponse.ChangeList
}

func (w *apiWorld) then_the_change_should_be(expected *pb.ChangeList) {
	if !proto.Equal(expected, w.lastFetchedChange) {
		w.Testing.Fatalf("Value fetched does not match value stored.\nexpected=%s\n  actual=%s", expected, w.lastFetchedChange)
	}
}
