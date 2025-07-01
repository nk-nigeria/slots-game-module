package handler

import (
	"context"

	"github.com/nk-nigeria/cgp-common/lib"
	"github.com/nk-nigeria/cgp-common/presenter"
	"github.com/nk-nigeria/slots-game-module/entity"
)

type StateIdle struct {
	lib.StateBase
}

func NewIdleState(fn lib.FireFn) lib.StateHandler {
	return &StateIdle{
		StateBase: lib.NewStateBase(fn),
	}
}

func (s *StateIdle) Enter(ctx context.Context, _ ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	state.SetUpCountDown(idleTimeout)
	dispatcher := procPkg.GetDispatcher()
	if dispatcher == nil {
		procPkg.GetLogger().Warn("missing dispatcher don't broadcast")
		return nil
	}
	// procPkg.GetProcessor().NotifyUpdateGameState(
	// 	state,
	// 	procPkg.GetLogger(),
	// 	procPkg.GetDispatcher(),
	// 	&pb.UpdateGameState{
	// 		State: pb.GameState_GameStateIdle,
	// 	},
	// )
	return nil
}

func (s *StateIdle) Exit(_ context.Context, _ ...interface{}) error {
	return nil
}

func (s *StateIdle) Process(ctx context.Context, args ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	if state.GetPresenceSize() > 0 {
		s.Trigger(ctx, triggerMatching)
		return nil
	}
	if remain := state.GetRemainCountDown(); remain < 0 {
		// Do finish here
		procPkg.GetLogger().Info("[idle] idle timeout => exit")
		s.Trigger(ctx, triggerNoOne)
		return presenter.ErrGameFinish
	}
	return nil
}
