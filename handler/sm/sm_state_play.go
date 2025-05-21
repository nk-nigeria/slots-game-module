package sm

import (
	"context"
	"math"

	pb "github.com/nakamaFramework/cgp-common/proto"

	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	"github.com/nakamaFramework/cgp-common/lib"
)

type StatePlay struct {
	lib.StateBase
}

func NewStatePlay(fn lib.FireFn) lib.StateHandler {
	return &StatePlay{
		StateBase: lib.NewStateBase(fn),
	}
}

func (s *StatePlay) Enter(ctx context.Context, _ ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	procPkg.GetLogger().Info("[playing] enter")
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	// Setup count down
	// state.SetUpCountDown(playTimeout)
	procPkg.GetProcessor().NotifyUpdateGameState(

		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GameStatePlay,
			CountDown: int64(math.Round(state.GetRemainCountDown())),
		},
		state,
	)
	procPkg.GetProcessor().ProcessNewGame(
		procPkg.GetContext(),
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		procPkg.GetMatchState(),
	)
	return nil
}

func (s *StatePlay) Exit(_ context.Context, _ ...interface{}) error {
	return nil
}

func (s *StatePlay) Process(ctx context.Context, args ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)

	if state.GetPresenceSize() <= 0 {
		procPkg.GetLogger().Info("no user in game")
		s.Trigger(ctx, lib.TriggerStateFinishFailed)
		return nil
	}

	message := procPkg.GetMessages()
	procPkg.GetProcessor().ProcessGame(procPkg.GetContext(),
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		message,
		state)
	return nil
}
