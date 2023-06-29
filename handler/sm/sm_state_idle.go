package sm

import (
	"context"
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	"github.com/ciaolink-game-platform/cgp-common/presenter"
)

type StateIdle struct {
	lib.StateBase
	timeout time.Duration
}

func NewIdleState(fn lib.FireFn) lib.StateHandler {
	return &StateIdle{
		StateBase: lib.NewStateBase(fn),
		timeout:   10 * time.Second,
	}
}

func (s *StateIdle) Enter(ctx context.Context, _ ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	state.SetUpCountDown(s.timeout)
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
