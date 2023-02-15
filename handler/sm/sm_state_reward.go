package sm

import (
	"context"
	"math"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

type StateReward struct {
	lib.StateBase
}

func NewStateReward(fn lib.FireFn) lib.StateHandler {
	return &StateReward{
		StateBase: lib.NewStateBase(fn),
	}
}
func (s *StateReward) Enter(ctx context.Context, _ ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	procPkg.GetLogger().Info("[reward] enter")
	// setup reward timeout
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	state.SetUpCountDown(rewardTimeout)
	procPkg.GetProcessor().NotifyUpdateGameState(
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GameStateReward,
			CountDown: int64(math.Round(state.GetRemainCountDown())),
		},
		state,
	)
	// process finish
	procPkg.GetProcessor().ProcessFinishGame(
		procPkg.GetContext(),
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state)

	return nil
}

func (s *StateReward) Exit(ctx context.Context, _ ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	state.ResetBalanceResult()
	return nil
}

func (s *StateReward) Process(ctx context.Context, args ...interface{}) error {
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
	if remain := state.GetRemainCountDown(); remain <= 0 {
		s.Trigger(ctx, lib.TriggerRewardTimeout)
	}
	return nil
}
