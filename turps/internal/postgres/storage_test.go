package postgres

import (
	"context"
	"reflect"
	"testing"
	"time"
	"turps/internal"
	pg_testing "turps/internal/postgres/testing"
)

func randomChangeListId() internal.ChangeListId {
	return internal.ChangeListId(internal.RandomStringOfLength(10))
}

func randomBuildId() internal.BuildId {
	return internal.BuildId(internal.RandomStringOfLength(15))
}

func truncatedNow() time.Time {
	return time.Now().Truncate(time.Microsecond)
}

func newChangeList(id internal.ChangeListId) *internal.ChangeList {
	return &internal.ChangeList{
		Id:               id,
		DependentTestIds: []string{"test-1"},
		Timestamp:        truncatedNow(),
	}
}

func TestCreateAndFetchChangeList(t *testing.T) {

	ctx := context.Background()

	pool, err := pg_testing.NewPool(ctx)
	if err != nil {
		t.Fatalf("error %s", err)
	}
	defer pool.Close()
	storage := ChangeListStorage{Pool: pool}

	changeListId := randomChangeListId()
	expected := newChangeList(changeListId)
	err = storage.Save(ctx, expected)
	if err != nil {
		t.Errorf("error %s", err)
	}

	var actual *internal.ChangeList
	actual, err = storage.ChangeList(ctx, changeListId)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Value fetched does not match value stored.  Create=%s, Get=%s", expected, actual)
	}
}
func Test_ShouldSaveTestRuns(t *testing.T) {

	ctx := context.Background()

	pool, err := pg_testing.NewPool(ctx)
	if err != nil {
		t.Fatalf("error %s", err)
	}
	defer pool.Close()
	storage := ChangeListStorage{Pool: pool}

	changeListId := randomChangeListId()
	buildId := randomBuildId()
	expected := newChangeList(changeListId)
	const EXPECTED_OUTPUT_URL = "EXPECTED_OUTPUT_URL"
	expected.TestRuns = append(expected.TestRuns, internal.TestRun{
		ChangeListId: changeListId,
		BuildId:      buildId,
		OutputUrl:    EXPECTED_OUTPUT_URL,
		Timestamp:    truncatedNow(),
		TestResults: map[string]internal.TestResult{
			"test-1": {Runs: 1, Fails: 1},
		},
	})

	err = storage.Save(ctx, expected)
	if err != nil {
		t.Errorf("error %s", err)
	}

	err = storage.SaveTestRun(ctx, &expected.TestRuns[0])
	if err != nil {
		t.Errorf("error %s", err)
	}

	var actual *internal.ChangeList
	actual, err = storage.ChangeList(ctx, changeListId)
	if err != nil {
		t.Errorf("error %s", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Value fetched does not match value stored.  Create=%s, Get=%s", expected, actual)
	}
}
