package sm

import (
	"context"
	"time"

	"github.com/qmuntal/stateless"

	"github.com/ciaolink-game-platform/cgp-common/lib"
)

const (
	idleTimeout      = time.Second * 15
	preparingTimeout = time.Second * 5
	playTimeout      = time.Second * 15
	//playTimeout      = time.Second * 10
	rewardTimeout = time.Second * 10
	//rewardTimeout    = time.Second * 10
)

var _ lib.StateMachineState = &SlotsStateMachineState{}

type SlotsStateMachineState struct {
	lib.StateMachineState
}

func NewSlotsStateMachineState() lib.StateMachineState {
	s := SlotsStateMachineState{}
	return &s
}

func (s *SlotsStateMachineState) NewIdleState(fn lib.FireFn) lib.StateHandler {
	return NewIdleState(fn)
}

func (s *SlotsStateMachineState) NewStateMatching(fn lib.FireFn) lib.StateHandler {
	return NewStateMatching(fn)
}

func (s *SlotsStateMachineState) NewStatePreparing(fn lib.FireFn) lib.StateHandler {
	return NewStatePreparing(fn)
}

func (s *SlotsStateMachineState) NewStatePlay(fn lib.FireFn) lib.StateHandler {
	return NewStatePlay(fn)
}
func (s *SlotsStateMachineState) NewStateReward(fn lib.FireFn) lib.StateHandler {
	return NewStateReward(fn)
}

func (s *SlotsStateMachineState) OnTransitioning(
	ctx context.Context,
	t stateless.Transition,
) {
}
