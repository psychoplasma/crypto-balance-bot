package concurrency_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/psychoplasma/crypto-balance-bot/infrastructure/concurrency"
)

var errTest = errors.New("test error")
var msgTest = "hello world!"
var r = concurrency.Retrial{
	Limit: 5,
	Delay: time.Millisecond * 100,
}

func TestCallRecurcively_SuccedingFunction(t *testing.T) {
	output, err := concurrency.Retry(r, func() (interface{}, error) {
		return succedingFunc(msgTest)
	})
	if err != nil {
		t.Fatal(err)
	}

	if output != msgTest {
		t.Fatalf("\nreturn value: %s\nwant: %s\n", output, msgTest)
	}
}

func TestCallRecurcively_FailingFunction(t *testing.T) {
	errFail := fmt.Sprintf("retry has reached to limit: %d", r.Limit)

	_, err := concurrency.Retry(r, failingFunc)
	if err == nil || err.Error() != errFail {
		t.Fatalf("\nerror: %s\nwant: %s\n", err, errFail)
	}
}

func TestCallRecurcively_SuccedingFunctionAfterSomeTrials(t *testing.T) {
	trials := 0
	maxTrials := 3

	output, err := concurrency.Retry(r, func() (interface{}, error) {
		return succedingFuncAfterTrial(msgTest, &trials, maxTrials)
	})
	if err != nil {
		t.Fatal(err)
	}

	if output != msgTest {
		t.Fatalf("\nreturn value: %s\nwant: %s\n", output, msgTest)
	}

	if trials != maxTrials {
		t.Fatalf("\nnumber of trials before success: %d\nwant: %d\n", trials, maxTrials)
	}
}

func succedingFunc(a string) (string, error) {
	time.Sleep(time.Millisecond * 100)
	return a, nil
}

func failingFunc() (interface{}, error) {
	time.Sleep(time.Millisecond * 100)
	return "", errTest
}

func succedingFuncAfterTrial(a string, trial *int, maxTrials int) (string, error) {
	time.Sleep(time.Millisecond * 100)
	if *trial < maxTrials {
		*trial++
		return "", errTest
	}
	return a, nil
}

func cleanUp(filename string) error {
	return nil
}
