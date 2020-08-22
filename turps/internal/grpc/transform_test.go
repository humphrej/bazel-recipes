package grpc

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"reflect"
	"testing"
	"time"
	pb "turps/api"
)

func genTestResult() gopter.Gen {
	return gen.StructPtr(reflect.TypeOf(&pb.TestResult{}), map[string]gopter.Gen{
		"NumRuns":  gen.UInt64(),
		"NumFails": gen.UInt64(),
	})
}

func genTimestamp() gopter.Gen {
	maxTz, _ := time.Parse(time.RFC3339, "9999-12-31T23:59:59Z")
	maxSeconds := maxTz.Unix()
	return gen.StructPtr(reflect.TypeOf(&timestamp.Timestamp{}), map[string]gopter.Gen{
		"Seconds": gen.Int64Range(0, maxSeconds),
		"Nanos":   gen.Int32Range(0, 999999999),
	})
}

func Test_ShouldRoundTripTimestamp(t *testing.T) {
	properties := gopter.NewProperties(gopter.DefaultTestParameters())

	properties.Property("Roundtrip time <-> proto", prop.ForAll(
		func(time time.Time) bool {
			pbTimestamp, err := ptypes.TimestampProto(time)
			if err != nil {
				return false
			}

			newTime, err := ptypes.Timestamp(pbTimestamp)
			if err != nil {
				return false
			}
			fmt.Printf("%v %v", time, newTime)
			return time.Equal(newTime)
		},
		gen.Time(),
	))

	properties.TestingRun(t)
}

func Test_ShouldRoundTripPbTimestamp(t *testing.T) {

	properties := gopter.NewProperties(gopter.DefaultTestParameters())

	properties.Property("Roundtrip proto <-> time", prop.ForAll(
		func(pbTimestamp *timestamp.Timestamp) (bool, error) {
			time, err := ptypes.Timestamp(pbTimestamp)
			if err != nil {
				return false, err
			}
			newTimestamp, err := ptypes.TimestampProto(time)
			return pbTimestamp.Seconds == newTimestamp.Seconds &&
				pbTimestamp.Nanos == newTimestamp.Nanos, nil
		},
		genTimestamp(),
	))

	properties.TestingRun(t)
}

func genTestRun() gopter.Gen {
	return gen.StructPtr(reflect.TypeOf(&pb.TestRun{}), map[string]gopter.Gen{
		"Id":           gen.AlphaString(),
		"ChangeListId": gen.AlphaString(),
		"OutputUrl":    gen.AlphaString(),
		"Tz":           genTimestamp(),
		"TestResult":   gen.MapOf(gen.AlphaString(), genTestResult()),
	})
}
func Test_ShouldTransformTestRun(t *testing.T) {

	properties := gopter.NewProperties(gopter.DefaultTestParameters())

	properties.Property("proto -> internal -> proto", prop.ForAll(
		func(v *pb.TestRun) (bool, error) {
			internal, err := NewTestRunInternal(v)
			if err != nil {
				return false, err
			}
			newPb, err := NewTestRunProto(internal)
			return reflect.DeepEqual(newPb, v), nil
		},
		genTestRun(),
	))

	properties.TestingRun(t)
}
