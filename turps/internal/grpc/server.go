package grpc

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
	"turps/internal"

	pb "turps/api"
)

type Server struct {
	pb.UnimplementedTurpsServer
	internal.ChangeListRepository
}

func NewServer(changeListRepository internal.ChangeListRepository) *Server {
	return &Server{ChangeListRepository: changeListRepository}
}

func (s *Server) UpsertChangeList(ctx context.Context, in *pb.UpsertChangeListRequest) (*pb.UpsertChangeListResponse, error) {

	if in.ChangeList == nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"ChangeList must be supplied")
	}
	if in.ChangeList.ChangeListId == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"ChangeListId is mandatory")
	}
	if in.ChangeList.Tz == nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"Timestamp is mandatory")
	}

	tz, err := ptypes.Timestamp(in.ChangeList.Tz)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"Timestamp cannot be converted %s", err)
	}

	changeList := internal.ChangeList{
		Id:               internal.ChangeListId(in.ChangeList.ChangeListId),
		DependentTestIds: in.ChangeList.TestIds,
		Timestamp:        tz,
	}

	err = s.ChangeListRepository.Save(ctx, &changeList)
	if err != nil {
		return nil, err
	}
	return &pb.UpsertChangeListResponse{ChangeList: in.ChangeList}, nil
}

func (s *Server) GetChangeList(ctx context.Context, in *pb.GetChangeListRequest) (*pb.GetChangeListResponse, error) {
	if in.ChangeListId == "" {
		return nil, status.Errorf(codes.InvalidArgument,
			"ChangeListId is mandatory")
	}
	changeList, err := s.ChangeListRepository.ChangeList(ctx, internal.ChangeListId(in.ChangeListId))
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"Unable to fetch change list %s", err)
	}

	protoChangeList, err := NewChangeListProto(changeList)
	if err != nil {
		return nil, status.Errorf(codes.Internal,
			"Unable to marshall change list %s", err)
	}
	return &pb.GetChangeListResponse{ChangeList: protoChangeList}, nil
}

func (s *Server) UpsertTestResult(ctx context.Context, in *pb.UpsertTestRunRequest) (*pb.UpsertTestRunResponse, error) {
	internal, err := NewTestRunInternal(in.TestRun)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument,
			"cannot understand test run %s", err)
	}
	s.ChangeListRepository.SaveTestRun(ctx, internal)

	return &pb.UpsertTestRunResponse{TestRun: in.TestRun}, nil
}

func TruncatedNow() *timestamp.Timestamp {
	tz := time.Now().Truncate(time.Microsecond)
	pbTz, _ := ptypes.TimestampProto(tz)
	return pbTz
}
