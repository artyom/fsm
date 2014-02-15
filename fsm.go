/*
Package fsm provides a convenient struct to track complex finite state
machine-like state.
*/
package fsm

import (
	"fmt"
	"sync"
)

// Possible state
type State string

var InvalidState = State("")

// State change error: transition to invalid state requested
type FsmTransitionError struct {
	State State
}

func (e FsmTransitionError) Error() string {
	return fmt.Sprintf("Cannot switch from current state %s", e.State)
}

type Fsm struct {
	// transitions map
	transitions map[State]map[State]bool
	// Current state
	state State
	m     sync.RWMutex
}

func (fsm Fsm) String() string {
	return fmt.Sprintf("Fsm:%s", fsm.state)
}

// NewFsm returns initialized Fsm struct with initial state set and transitions map.
func NewFsm(initial State, transitions map[State]map[State]bool) *Fsm {
	return &Fsm{
		transitions: transitions,
		state:       initial,
	}
}

// State returns current fsm state
func (fsm *Fsm) State() State {
	return fsm.state
}

// Possible returns set (map) of states, transition to which can be made from current state
func (fsm *Fsm) Possible() map[State]bool {
	return fsm.transitions[fsm.state]
}

// Avail checks whether transition to given state can be made
func (fsm *Fsm) Avail(dst State) bool {
	_, present := fsm.transitions[fsm.state][dst]
	return present
}

// To changes state to given one, returning FsmTransitionError if transition cannot be made
func (fsm *Fsm) To(dst State) (State, error) {
	fsm.m.Lock()
	defer fsm.m.Unlock()
	if !fsm.Avail(dst) {
		return InvalidState, FsmTransitionError{fsm.state}
	}
	fsm.state = dst
	return fsm.state, nil
}
