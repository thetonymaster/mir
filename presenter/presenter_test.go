package presenter

import (
	"errors"
	"testing"
)

func TestNewPresenter(t *testing.T) {
	p := NewPresenter(nil)
	if p == nil {
		t.Error("Presenter is not built correctly")
	}
}

type fakeRepo struct{}

func (f fakeRepo) Save(table string, data map[string]interface{}) error {
	return nil
}

func ExamplePresenter_PrintResult_success() {
	p := NewPresenter(fakeRepo{})
	res := []Result{
		{
			Task: "Task A",
			Time: 10.0,
		},
		{
			Task: "Task B",
			Time: 10.0,
		},
		{
			Task: "Task C",
			Time: 10.0,
		},
		{
			Task: "D",
			Time: 10.0,
		},
	}

	p.PrintResult(res, 10.0)

	// Output:
	//	SUCCESS
	// Task A took 10.000000
	// Task B took 10.000000
	// Task C took 10.000000
	// D took 10.000000
	// TOTAL: 40.000000
	// AVERAGE: 10.000000
	// REAL TIME: 10.000000

}

func ExamplePresenter_PrintResult_fail() {
	p := NewPresenter(fakeRepo{})
	resf := []Result{
		{
			Task: "Task A",
			Time: 10.0,
		},
		{
			Task:   "Task B",
			Time:   10.0,
			Error:  errors.New("Failed Test"),
			Output: "Failed tests:\nTest asd: failed\nTests run:",
		},
		{
			Task: "Task C",
			Time: 10.0,
		},
		{
			Task: "D",
			Time: 10.0,
		},
	}

	p.PrintResult(resf, 10.0)

	// Output:
	// FAIL
	// Task A took 10.000000
	// Task B took 10.000000
	// Test asd: failed
	// Task C took 10.000000
	// D took 10.000000
	// TOTAL: 40.000000
	// AVERAGE: 10.000000
	// REAL TIME: 10.000000

}

func TestPresenter_Save(t *testing.T) {
	p := NewPresenter(fakeRepo{})

	err := p.Save(10.0, 10.0, 10.0)
	if err != nil {
		t.Fatal("Error should hot have happened")
	}

}
