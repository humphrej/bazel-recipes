package internal

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type ChangeListId string
type BuildId string

// ChangeList represents a discrete SCM change.
type ChangeList struct {
	Id               ChangeListId
	PrevId           ChangeListId
	DependentTestIds []string
	Timestamp        time.Time
	TestRuns         []TestRun
}

func (c ChangeList) String() string {
	return fmt.Sprintf("id: %s, test_ids=[%s], timestamp=%s, runs=%s",
		c.Id, strings.Join(c.DependentTestIds, ","), c.Timestamp, c.TestRuns)
}

// TestRun represents an execution of one or more tests
type TestRun struct {
	ChangeListId
	BuildId
	TestResults map[string]TestResult
	OutputUrl   string
	Timestamp   time.Time
}

func (r TestRun) String() string {
	var results []string

	for k, v := range r.TestResults {
		results = append(results, k+":"+v.String())
	}
	return fmt.Sprintf("buildId=%s, buildUrl=%s, timestamp=%s, results=%s",
		r.BuildId, r.OutputUrl, r.Timestamp,
		results,
	)
}

// TestResult holds the results of a test execution
type TestResult struct {
	Runs  uint64
	Fails uint64
}

func (r TestResult) String() string {
	return fmt.Sprintf("%d/%d", r.Runs, r.Fails)
}
func (r TestResult) add(other TestResult) TestResult {
	var result TestResult
	result.Fails = r.Fails + other.Fails
	result.Runs = r.Runs + other.Runs
	return result
}

type ChangeListRepository interface {
	ChangeList(ctx context.Context, id ChangeListId) (*ChangeList, error)
	Save(ctx context.Context, c *ChangeList) error
	SaveTestRun(ctx context.Context, run *TestRun) error
}

func (cl ChangeList) AggregateByTest() map[string]TestResult {

	agg := map[string]TestResult{}

	for _, testRun := range cl.TestRuns {
		for testName, summary := range testRun.TestResults {
			agg[testName] = agg[testName].add(summary)
		}
	}
	return agg
}
