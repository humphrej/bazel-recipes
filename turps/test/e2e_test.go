package test

import (
	"reflect"
	"runtime"
	"testing"
)

func newCLIDSL(t *testing.T) DSL {
	return &cliWorld{Testing: t}
}

func newAPIDSL(t *testing.T) DSL {
	return &apiWorld{Testing: t}
}

var (
	allTests = []func(DSL){
		ShouldCreateAndFetchChangeList,
		ShouldUpdateOnDoubleUpsert,
		ShouldUpdateChangeListWithSingleTestRun,
		ShouldUpdateChangeListWithDoubleTestRun,
		ShouldUpsertTestRun}
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
func TestAPI(t *testing.T) {
	for _, testFunc := range allTests {
		t.Run("API "+GetFunctionName(testFunc), partial(newAPIDSL(t), testFunc))
	}
}

func TestCLI(t *testing.T) {
	for _, testFunc := range allTests {
		t.Run("CLI "+GetFunctionName(testFunc), partial(newCLIDSL(t), testFunc))
	}
}

func partial(dsl DSL, testFunc func(dsl DSL)) func(t *testing.T) {
	return func(t *testing.T) {
		testFunc(dsl)
	}
}
