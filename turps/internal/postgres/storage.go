package postgres

import (
	"context"
	"github.com/golang/glog"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
	"turps/internal"
)

type ChangeListStorage struct {
	Pool *pgxpool.Pool
}

// ChangeList performs a full query by change_list_id
func (pg ChangeListStorage) ChangeList(ctx context.Context, id internal.ChangeListId) (*internal.ChangeList, error) {
	var err error
	var timestamp time.Time
	var testIds []string
	changeListRow := pg.Pool.QueryRow(ctx,
		"SELECT id,timestamp, test_ids FROM change_list where id=$1", id)
	err = changeListRow.Scan(&id, &timestamp, &testIds)
	if err != nil {
		return nil, err
	}

	cl := &internal.ChangeList{
		Id:               id,
		DependentTestIds: testIds,
		Timestamp:        timestamp,
	}

	var buildId internal.BuildId
	var buildTz time.Time
	var buildOutput string
	testRunRows, err := pg.Pool.Query(ctx,
		"SELECT id,output_url,timestamp FROM test_run where change_list_id=$1", id)
	defer testRunRows.Close()
	if err != nil {
		return nil, err
	}
	for testRunRows.Next() {
		err = testRunRows.Scan(&buildId, &buildOutput, &buildTz)
		if err != nil {
			return nil, err
		}
		cl.TestRuns = append(cl.TestRuns, internal.TestRun{
			ChangeListId: id,
			BuildId:      buildId,
			OutputUrl:    buildOutput,
			Timestamp:    buildTz,
			TestResults:  map[string]internal.TestResult{},
		})
	}

	for _, run := range cl.TestRuns {
		testResultRows, err := pg.Pool.Query(ctx,
			"SELECT test_id, count_runs, count_fails from test_run_result where run_id = $1", run.BuildId)
		defer testResultRows.Close()
		if err != nil {
			return nil, err
		}
		var testId string
		var countRuns uint64
		var countFails uint64
		for testResultRows.Next() {
			err = testResultRows.Scan(&testId, &countRuns, &countFails)
			if err != nil {
				return nil, err
			}
			run.TestResults[testId] = internal.TestResult{
				Runs:  countRuns,
				Fails: countFails,
			}
		}
	}

	return cl, nil
}

func (pg ChangeListStorage) Save(ctx context.Context, c *internal.ChangeList) error {

	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := pg.Pool.Exec(ctx,
		"INSERT INTO change_list(id, timestamp, test_ids) values ($1,$2,$3) on conflict on constraint change_list_pkey do update set timestamp=$2, test_ids=$3",
		c.Id, c.Timestamp, c.DependentTestIds)
	if err != nil {
		return err
	}

	glog.Infof("insert change_list rows_affected %d", commandTag.RowsAffected())

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (pg ChangeListStorage) SaveTestRun(ctx context.Context, run *internal.TestRun) error {

	tx, err := pg.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	commandTag, err := pg.Pool.Exec(ctx,
		`INSERT INTO test_run(id, change_list_id, output_url, timestamp) values ($1,$2,$3,$4)
on conflict on constraint test_run_pkey do update set change_list_id = $2, output_url=$3, timestamp=$4`,
		run.BuildId, run.ChangeListId, run.OutputUrl, run.Timestamp)
	if err != nil {
		return err
	}
	glog.Infof("insert test_run rows_affected %d", commandTag.RowsAffected())

	for testId, summary := range run.TestResults {
		commandTag, err := pg.Pool.Exec(ctx,
			`INSERT INTO test_run_result(run_id, test_id, count_runs, count_fails) values ($1,$2,$3,$4)
on conflict on constraint test_run_result_run_id_test_id_key do update set count_runs=$3,count_fails=$4`,
			run.BuildId, testId, summary.Runs, summary.Fails)
		if err != nil {
			return err
		}
		glog.Infof("insert test_run_result rows_affected %d", commandTag.RowsAffected())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
