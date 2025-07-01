package handler

import (
	"context"
	"math"

	pb "github.com/nk-nigeria/cgp-common/proto"

	"github.com/nk-nigeria/cgp-common/lib"
	"github.com/nk-nigeria/slots-game-module/entity"
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
	return nil
}

func (s *StatePlay) Exit(_ context.Context, _ ...interface{}) error {
	return nil
}

func (s *StatePlay) Process(ctx context.Context, args ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	// remain := state.GetRemainCountDown()
	// if remain <= 0 {
	// 	procPkg.GetLogger().Info("[play] timeout reach %v", remain)
	// 	s.Trigger(ctx, lib.TriggerPlayTimeout)
	// 	return nil
	// }
	if state.GetPresenceSize() <= 0 {
		if state.CountDownReachTime.Unix() <= 0 {
			state.SetUpCountDown(playTimeout)
		}
		return nil
	}

	message := procPkg.GetMessages()
	procPkg.GetProcessor().ProcessGame(ctx,
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		message,
		state)

	if state.CountDownReachTime.Unix() <= 0 {
		s.Trigger(ctx, triggerNoOne)
	}

	// if len(message) > 0 {

	// 	procPkg.GetProcessor().ProcessMessageFromUser(ctx,
	// 		procPkg.GetLogger(),
	// 		procPkg.GetNK(),
	// 		procPkg.GetDb(),
	// 		procPkg.GetDispatcher(),
	// 		message,
	// 		state)
	// }

	// if state.IsNeedNotifyCountDown() {
	// 	remainCountDown := int(math.Round(state.GetRemainCountDown()))
	// 	procPkg.GetProcessor().NotifyUpdateGameState(
	// 		procPkg.GetLogger(),
	// 		procPkg.GetDispatcher(),
	// 		&pb.UpdateGameState{
	// 			State:     pb.GameState_GameStatePlay,
	// 			CountDown: int64(remainCountDown),
	// 		},
	// 		state,
	// 	)
	// 	state.SetLastCountDown(remainCountDown)
	// }
	return nil
}
