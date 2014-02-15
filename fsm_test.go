package fsm

import (
	"fmt"
	"reflect"
	"testing"
)

var (
	MyFsm       *Fsm
	green       = State("green")
	yellow      = State("yellow")
	red         = State("red")
	transitions = map[State]map[State]bool{
		green:  {yellow: true},
		yellow: {green: true, red: true},
		red:    {yellow: true, green: true},
	}
)

var testData = []struct {
	cur      State
	possible map[State]bool
	avail    State
	wrong    State
}{
	{yellow, map[State]bool{green: true, red: true}, green, yellow},
	{green, map[State]bool{yellow: true}, yellow, red},
	{yellow, map[State]bool{green: true, red: true}, red, yellow},
	{red, map[State]bool{green: true, yellow: true}, green, red},
}

func init() {
	MyFsm = NewFsm(green, transitions)
}

func resetState() {
	MyFsm.state = green
}

func ExampleState() {
	norm := State("normal")
	warn := State("warning")
	crit := State("critical")
	fmt.Println(norm)
	fmt.Println(warn)
	fmt.Println(crit)
	// Output:

	// normal
	// warning
	// critical
}

func ExampleNewFsm() {
	// Imagine monitoring states: green for "ok", yellow for "warning", red
	// for "critical".
	var (
		green  = State("green")
		yellow = State("yellow")
		red    = State("red")
	)
	// Transitions map: keys are source states, values are sets (boolean
	// maps) of possible destination states.
	transitions = map[State]map[State]bool{
		green:  {yellow: true},
		yellow: {green: true, red: true},
		red:    {yellow: true, green: true},
	}
	myFsm := NewFsm(green, transitions)
	fmt.Printf("%s", myFsm)
	// Output:

	// Fsm:green
}

func TestFull(t *testing.T) {
	resetState()
	for _, testCase := range testData {
		currentState, err := MyFsm.To(testCase.cur)
		if err != nil {
			t.Errorf("Failed to change state to %s: %s", testCase.cur, err)
		}
		if currentState != testCase.cur {
			t.Errorf("Failed to change state: got %s, want %s", currentState, testCase.cur)
		}
		possibleStates := MyFsm.Possible()
		if !reflect.DeepEqual(possibleStates, testCase.possible) {
			t.Errorf("Invalid possible states: got %s, want %s", possibleStates, testCase.possible)
		}
		if !MyFsm.Avail(testCase.avail) {
			t.Errorf("Expected state unavailable: %s", testCase.avail)
		}
		if MyFsm.Avail(testCase.wrong) {
			t.Errorf("State available, but should not be: %s", testCase.wrong)
		}
		if MyFsm.State() != testCase.cur {
			t.Errorf("Invalid current state: %s, want %s", MyFsm.State, testCase.cur)
		}
	}
}
