package grpc

import (
	"github.com/golang/protobuf/ptypes"
	pb "turps/api"
	"turps/internal"
)

func NewChangeListProto(list *internal.ChangeList) (*pb.ChangeList, error) {

	protoTz, err := ptypes.TimestampProto(list.Timestamp)
	if err != nil {
		return nil, err
	}

	var runs []*pb.TestRun
	if list.TestRuns != nil {
		runs, err = NewTestRunSliceProto(list.TestRuns)
		if err != nil {
			return nil, err
		}
	}

	return &pb.ChangeList{
		ChangeListId: string(list.Id),
		Tz:           protoTz,
		TestIds:      list.DependentTestIds,
		TestRun:      runs,
	}, nil
}

func NewTestRunSliceProto(run []internal.TestRun) ([]*pb.TestRun, error) {
	to := []*pb.TestRun{}
	for _, v := range run {
		protoRun, err := NewTestRunProto(&v)
		if err != nil {
			return nil, err
		}
		to = append(to, protoRun)
	}
	return to, nil
}

func NewTestRunProto(run *internal.TestRun) (*pb.TestRun, error) {
	testResultMap := NewTestResultMapProto(run.TestResults)

	protoTz, err := ptypes.TimestampProto(run.Timestamp)
	if err != nil {
		return nil, err
	}

	return &pb.TestRun{
		Id:           string(run.BuildId),
		ChangeListId: string(run.ChangeListId),
		OutputUrl:    run.OutputUrl,
		Tz:           protoTz,
		TestResult:   testResultMap,
	}, nil
}

func NewTestResultMapProto(m map[string]internal.TestResult) map[string]*pb.TestResult {
	to := map[string]*pb.TestResult{}
	for k, v := range m {
		to[k] = NewTestResultProto(&v)
	}
	return to
}

func NewTestResultProto(result *internal.TestResult) *pb.TestResult {
	return &pb.TestResult{NumFails: uint64(result.Fails), NumRuns: uint64(result.Runs)}
}
func NewTestResultMapInternal(protoResultMap map[string]*pb.TestResult) map[string]internal.TestResult {
	to := map[string]internal.TestResult{}
	for k, v := range protoResultMap {
		to[k] = NewTestResultInternal(v)
	}
	return to
}

func NewTestRunInternal(run *pb.TestRun) (*internal.TestRun, error) {

	tz, err := ptypes.Timestamp(run.Tz)
	if err != nil {
		return nil, err
	}

	return &internal.TestRun{
		ChangeListId: internal.ChangeListId(run.ChangeListId),
		BuildId:      internal.BuildId(run.Id),
		TestResults:  NewTestResultMapInternal(run.TestResult),
		OutputUrl:    run.OutputUrl,
		Timestamp:    tz,
	}, nil
}

func NewTestResultInternal(result *pb.TestResult) internal.TestResult {
	return internal.TestResult{Runs: result.NumRuns, Fails: result.NumFails}
}
