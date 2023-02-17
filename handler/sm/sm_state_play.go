package sm

import (
	"context"
	"math"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
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
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	// Setup count down
	state.SetUpCountDown(playTimeout)
	procPkg.GetProcessor().NotifyUpdateGameState(

		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GameStatePlay,
			CountDown: int64(math.Round(state.GetRemainCountDown())),
		},
		state,
	)
	state.SetupMatchPresence()
	state.SetAllowBet(true)
	return nil
}

func (s *StatePlay) Exit(_ context.Context, _ ...interface{}) error {
	return nil
}

func (s *StatePlay) Process(ctx context.Context, args ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	remain := state.GetRemainCountDown()
	if remain <= 0 {
		procPkg.GetLogger().Info("[play] timeout reach %v", remain)
		s.Trigger(ctx, lib.TriggerPlayTimeout)
		return nil
	}

	message := procPkg.GetMessages()
	procPkg.GetProcessor().ProcessGame(ctx,
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state)

	if len(message) > 0 {
		procPkg.GetProcessor().ProcessMessageFromUser(ctx,
			procPkg.GetLogger(),
			procPkg.GetNK(),
			procPkg.GetDb(),
			procPkg.GetDispatcher(),
			message,
			state)
	}

	if state.IsNeedNotifyCountDown() {
		remainCountDown := int(math.Round(state.GetRemainCountDown()))
		procPkg.GetProcessor().NotifyUpdateGameState(
			procPkg.GetLogger(),
			procPkg.GetDispatcher(),
			&pb.UpdateGameState{
				State:     pb.GameState_GameStatePlay,
				CountDown: int64(remainCountDown),
			},
			state,
		)
		state.SetLastCountDown(remainCountDown)
	}
	return nil
}
