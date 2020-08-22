package main

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
	pb "turps/api"
	"turps/internal/grpc"
)

func NewChangeList() pb.ChangeList {
	return pb.ChangeList{
		ChangeListId: "new-change-list-1",
		Tz:           grpc.TruncatedNow(),
	}
}
func writePbMessage(m proto.Message) (string, error) {

	bs, err := protojson.Marshal(m)
	if err != nil {
		return "", err
	}

	tmpDir, ok := os.LookupEnv("TEST_TMPDIR")
	if !ok {
		tmpDir = "/tmp"
	}
	tmpfile, err := ioutil.TempFile(tmpDir, "turps_acceptance")
	if err != nil {
		return "", err
	}

	err = ioutil.WriteFile(tmpfile.Name(), bs, 0644)
	if err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func TestCreateAndFetchChange(t *testing.T) {

	changeList := NewChangeList()
	name, err := writePbMessage(&changeList)
	if err != nil {
		t.Errorf("error generating testdata %v", err)
	}
	defer os.Remove(name)

	createResult, err := runTurps([]byte{}, "create", "--from_file="+name, "changelist")
	if err != nil {
		t.Fatalf("error running turps %v", err)
	}

	if createResult.exitCode != 0 {
		t.Fatalf("Exit code in error %d", createResult.exitCode)
	}

	if len(createResult.stderr) > 0 {
		t.Fatalf("unexpected stderr %v", string(createResult.stderr))
	}

	getResult, err := runTurps([]byte{}, "get", changeList.ChangeListId)
	if err != nil {
		t.Fatalf("error running turps %v", err)
	}

	if getResult.exitCode != 0 {
		t.Fatalf("Exit code in error %d", getResult.exitCode)
	}

	if len(getResult.stderr) > 0 {
		t.Fatalf("unexpected stderr %v", string(getResult.stderr))
	}

	var fetchedChangeList pb.ChangeList

	err = protojson.Unmarshal(getResult.stdout, &fetchedChangeList)
	if err != nil {
		t.Fatalf("Cannot unmarshal fetched changelist %v", err)
	}

	if !reflect.DeepEqual(fetchedChangeList, changeList) {
		t.Fatalf("actual:\n%v\nexpected:\n%v", fetchedChangeList, changeList)
	}

}

type TurpsResult struct {
	stdout   []byte
	stderr   []byte
	exitCode int
}

func NewTurpsResult(exitCode int, stdout []byte, stderr []byte) *TurpsResult {
	return &TurpsResult{
		stdout:   stdout,
		stderr:   stderr,
		exitCode: exitCode,
	}
}

func runTurps(stdinData []byte, args ...string) (*TurpsResult, error) {

	turpsBinary, ok := os.LookupEnv("TURPS_BINARY")
	if !ok {
		return nil, errors.New("TURPS_BINARY must be set")
	}

	cmd := exec.Command(turpsBinary, args...)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	stdin.Write(stdinData) // TODO error check / retry?

	slurpedStderr, _ := ioutil.ReadAll(stderr)
	fmt.Printf("stderr: %s\n", slurpedStderr)
	slurpedStdout, _ := ioutil.ReadAll(stdout)
	fmt.Printf("stdout: %s\n", slurpedStdout)
	if err := cmd.Wait(); err != nil {
		return nil, err
	}
	return NewTurpsResult(cmd.ProcessState.ExitCode(), slurpedStdout, slurpedStderr), nil
}
