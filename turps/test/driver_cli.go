package test

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/encoding/protojson"
	gproto "google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
	pb "turps/api"
)

type cliWorld struct {
	Testing           *testing.T
	changesMap        map[string]*pb.ChangeList
	testRunsMap       map[string]*pb.TestRun
	lastFetchedChange *pb.ChangeList
	lastError         string
}

func (w *cliWorld) given_a_turps_server() {
}

func (w *cliWorld) given_a_change(ref string, c *pb.ChangeList) {
	w.changesMap[ref] = c
}
func (w *cliWorld) given_a_test_run(ref string, r *pb.TestRun) {
	w.testRunsMap[ref] = r
}

func (w *cliWorld) when_the_change_is_saved(ref string) {
	c := w.changesMap[ref]

	name, err := writePbMessage(c)
	if err != nil {
		w.Testing.Errorf("error generating testdata %v", err)
	}
	defer os.Remove(name)

	createResult, err := runTurps([]byte{}, "create", "--from_file="+name, "changelist")
	if err != nil {
		w.Testing.Fatalf("error running turps %v", err)
	}

	if createResult.exitCode != 0 {
		w.Testing.Fatalf("Exit code in error %d", createResult.exitCode)
	}

	if len(createResult.stderr) > 0 {
		w.Testing.Fatalf("unexpected stderr %v", string(createResult.stderr))
	}

}
func (w *cliWorld) when_the_test_run_is_saved(ref string) {
	run := w.testRunsMap[ref]

	name, err := writePbMessage(run)
	if err != nil {
		w.Testing.Errorf("error generating testdata %v", err)
	}
	defer os.Remove(name)

	createResult, err := runTurps([]byte{}, "create", "--from_file="+name, "testrun")
	if err != nil {
		w.Testing.Fatalf("error running turps %v", err)
	}

	if createResult.exitCode != 0 {
		w.Testing.Fatalf("Exit code in error %d", createResult.exitCode)
	}

	if len(createResult.stderr) > 0 {
		w.Testing.Fatalf("unexpected stderr %v", string(createResult.stderr))
	}
}

//func (w *cliWorld) then_the_change_should_be(expected *pb.ChangeList) {
func (w *cliWorld) when_the_change_is_fetched(changeListId string) {

	getResult, err := runTurps([]byte{}, "get", changeListId)
	if err != nil {
		w.Testing.Fatalf("error running turps %v", err)
	}

	if getResult.exitCode != 0 {
		w.Testing.Fatalf("Exit code in error %d", getResult.exitCode)
	}

	if len(getResult.stderr) > 0 {
		w.lastError = string(getResult.stderr)
		return
		//w.Testing.Fatalf("unexpected stderr %v", string(getResult.stderr))
	}

	var fetchedChangeList pb.ChangeList

	err = protojson.Unmarshal(getResult.stdout, &fetchedChangeList)
	if err != nil {
		w.Testing.Fatalf("Cannot unmarshal fetched changelist %v", err)
	}

	w.lastFetchedChange = &fetchedChangeList
}

func (w *cliWorld) then_the_change_should_be(expected *pb.ChangeList) {

	if !proto.Equal(expected, w.lastFetchedChange) {
		w.Testing.Fatalf("Value fetched does not match value stored.\nexpected=%s\n  actual=%s", expected, w.lastFetchedChange)
	}
}
func (w *cliWorld) then_an_error(matchingString string) {
	if !strings.Contains(w.lastError, matchingString) {
		w.Testing.Fatalf("error message does not match error:%s match:%s", w.lastError, matchingString)
	}
}

func writePbMessage(m gproto.Message) (string, error) {

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
