package sm

import (
	"context"
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
)

type StateMatching struct {
	lib.StateBase
}

func NewStateMatching(fn lib.FireFn) lib.StateHandler {
	return &StateMatching{
		StateBase: lib.NewStateBase(fn),
	}
}

func (s *StateMatching) Enter(ctx context.Context, _ ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	procPkg.GetLogger().Info("[matching] enter")
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	state.SetUpCountDown(1 * time.Second)
	procPkg.GetLogger().Info("apply leave presence")
	procPkg.GetProcessor().ProcessApplyPresencesLeave(
		procPkg.GetContext(),
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state)
	procPkg.GetProcessor().ProcessNewGame(procPkg.GetLogger(), procPkg.GetDispatcher(), state)
	return nil
}
func (s *StateMatching) Exit(_ context.Context, _ ...interface{}) error {
	return nil
}

func (s *StateMatching) Process(ctx context.Context, args ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	message := procPkg.GetMessages()
	if len(message) > 0 {
		procPkg.GetProcessor().ProcessMessageFromUser(ctx,
			procPkg.GetLogger(),
			procPkg.GetNK(),
			procPkg.GetDb(),
			procPkg.GetDispatcher(),
			message, state)
	}
	remain := state.GetRemainCountDown()
	if remain > 0 {
		return nil
	}
	if state.IsReadyToPlay() {
		s.Trigger(ctx, lib.TriggerPresenceReady)
	} else {
		s.Trigger(ctx, lib.TriggerIdle)
	}
	return nil
}
