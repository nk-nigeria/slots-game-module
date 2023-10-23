package sm

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

type StatePreparing struct {
	lib.StateBase
}

func NewStatePreparing(fn lib.FireFn) lib.StateHandler {
	return &StatePreparing{
		StateBase: lib.NewStateBase(fn),
	}
}
func (s *StatePreparing) Enter(ctx context.Context, _ ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	procPkg.GetLogger().Info("state %v", state.Presences)
	state.SetUpCountDown(0 * time.Second)
	// remove all user not interact 2 game conti
	listPrecense := state.GetPresenceNotInteract(2)
	if len(listPrecense) > 0 {
		listUserId := make([]string, len(listPrecense))
		for _, p := range listPrecense {
			listUserId = append(listUserId, p.GetUserId())
		}
		procPkg.GetLogger().Info("Kick %d user from math %s",
			len(listPrecense), strings.Join(listUserId, ","))
		state.AddLeavePresence(listPrecense...)
	}
	procPkg.GetProcessor().ProcessApplyPresencesLeave(
		procPkg.GetContext(),
		procPkg.GetLogger(),
		procPkg.GetNK(),
		procPkg.GetDb(),
		procPkg.GetDispatcher(),
		state,
	)
	procPkg.GetProcessor().NotifyUpdateGameState(
		procPkg.GetLogger(),
		procPkg.GetDispatcher(),
		&pb.UpdateGameState{
			State:     pb.GameState_GameStatePreparing,
			CountDown: int64(math.Round(float64(state.GetRemainCountDown()))),
		},
		state,
	)
	return nil
}

func (s *StatePreparing) Exit(_ context.Context, _ ...interface{}) error {
	return nil
}

func (s *StatePreparing) Process(ctx context.Context, args ...interface{}) error {
	procPkg := lib.GetProcessorPackagerFromContext(ctx)
	state := procPkg.GetMatchState().(*entity.SlotsMatchState)
	remain := state.GetRemainCountDown()
	message := procPkg.GetMessages()
	if len(message) > 0 {
		procPkg.GetProcessor().ProcessMessageFromUser(
			procPkg.GetContext(),
			procPkg.GetLogger(),
			procPkg.GetNK(),
			procPkg.GetDb(),
			procPkg.GetDispatcher(),
			message, state)
	}
	if remain <= 0 {
		if state.IsReadyToPlay() {
			s.Trigger(ctx, lib.TriggerStateFinishSuccess)
		} else {
			// change to wait
			s.Trigger(ctx, lib.TriggerStateFinishFailed)
		}
		return nil
	}
	return nil
}
