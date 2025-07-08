package sm

import (
	"context"
	"time"

	"github.com/nk-nigeria/slots-game-module/entity"
	"github.com/nk-nigeria/cgp-common/lib"
	"github.com/nk-nigeria/cgp-common/presenter"
)

type StateIdle struct {
	lib.StateBase
	timeout time.Duration
	Count   int
}

func NewIdleState(fn lib.FireFn) lib.StateHandler {
	return &StateIdle{
		StateBase: lib.NewStateBase(fn),
		timeout:   100 * time.Second,
		Count:     0,
	}
}

func (s *StateIdle) Enter(ctx context.Context, _ ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	if s.Count == 0 {
		s.timeout = 500 * time.Second
	} else {
		s.timeout = 10 * time.Second
	}
	state.SetUpCountDown(s.timeout)
	dispatcher := procPkg.GetDispatcher()
	if dispatcher == nil {
		procPkg.GetLogger().Warn("missing dispatcher don't broadcast")
		return nil
	}
	return nil
}

func (s *StateIdle) Exit(_ context.Context, _ ...interface{}) error {
	s.Count++
	return nil
}

func (s *StateIdle) Process(ctx context.Context, args ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	if state.GetPresenceSize() > 0 {
		s.Trigger(ctx, lib.TriggerStateFinishSuccess)
		s.timeout = 0
		return nil
	}
	if remain := state.GetRemainCountDown(); remain < 0 {
		// Do finish here
		procPkg.GetLogger().Info("[idle] idle timeout => exit")
		s.Trigger(ctx, lib.TriggerExit)
		return presenter.ErrGameFinish
	}
	return nil
}
