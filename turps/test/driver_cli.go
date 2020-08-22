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
	"testing"
	pb "turps/api"
)

type cliWorld struct {
	Testing *testing.T
}

func (w *cliWorld) given_a_turps_server() {
}

func (w *cliWorld) given_a_change(c *pb.ChangeList) {

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
func (w *cliWorld) when_tests_are_run(runs []*pb.TestRun) {

	for _, run := range runs {
		invokeTestRun(w, run)
	}
}

func invokeTestRun(w *cliWorld, testRun *pb.TestRun) {

	name, err := writePbMessage(testRun)
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
func (w *cliWorld) then_the_change_should_be(expected *pb.ChangeList) {

	getResult, err := runTurps([]byte{}, "get", expected.ChangeListId)
	if err != nil {
		w.Testing.Fatalf("error running turps %v", err)
	}

	if getResult.exitCode != 0 {
		w.Testing.Fatalf("Exit code in error %d", getResult.exitCode)
	}

	if len(getResult.stderr) > 0 {
		w.Testing.Fatalf("unexpected stderr %v", string(getResult.stderr))
	}

	var fetchedChangeList pb.ChangeList

	err = protojson.Unmarshal(getResult.stdout, &fetchedChangeList)
	if err != nil {
		w.Testing.Fatalf("Cannot unmarshal fetched changelist %v", err)
	}

	if !proto.Equal(&fetchedChangeList, expected) {
		w.Testing.Fatalf("actual:\n%v\nexpected:\n%v", fetchedChangeList, expected)
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
