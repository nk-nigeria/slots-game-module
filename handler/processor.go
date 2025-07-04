package handler

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/nk-nigeria/slots-game-module/cgbdb"

	"github.com/nk-nigeria/slots-game-module/entity"

	"github.com/heroiclabs/nakama-common/runtime"
	pb "github.com/nk-nigeria/cgp-common/proto"
	"google.golang.org/protobuf/proto"

	"github.com/nk-nigeria/cgp-common/lib"
)

var _ lib.Processor = &processor{}

type processor struct {
	engine      lib.Engine
	marshaler   *proto.MarshalOptions
	unmarshaler *proto.UnmarshalOptions
	// turnBaseEngine *lib.TurnBaseEngine
}

func NewMatchProcessor(marshaler *proto.MarshalOptions,
	unmarshaler *proto.UnmarshalOptions,
	engine lib.Engine) lib.Processor {
	p := processor{
		marshaler:   marshaler,
		unmarshaler: unmarshaler,
		engine:      engine,
		// turnBaseEngine: lib.NewTurnBaseEngine(),
	}
	return &p
}
func (p *processor) ProcessNewGame(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	matchState interface{}) {
	s := matchState.(*entity.SlotsMatchState)
	s.InitNewMatch()
	s.SetupMatchPresence()
	s.SetAllowSpin(true)
	p.engine.NewGame(matchState)
	if len(s.GetPlayingPresences()) == 0 {
		return
	}
	logger.Info("List matrix new game %v", string(s.GetMatrix().List))
	slotDesk := &pb.SlotDesk{}
	presence := s.GetPlayingPresences()[0]
	// wallet, err := entity.ReadWalletUser(ctx, nk, logger, presence.GetUserId())
	// if err != nil {
	// 	logger.WithField("error", err.Error()).
	// 		WithField("user id", presence.GetUserId()).
	// 		Error("get profile user failed")
	// 	return
	// }
	matrix := s.GetMatrix()
	slotDesk.Matrix = matrix.ToPbSlotMatrix()
	// slotDesk.BalanceChipsWalletBefore = wallet.Chips
	// slotDesk.BalanceChipsWalletAfter = wallet.Chips

	p.broadcastMessage(logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
		slotDesk,
		[]runtime.Presence{presence},
		nil, false)
}

func (p *processor) ProcessGame(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	messages []runtime.MatchData,
	matchState interface{},
) {
	s := matchState.(*entity.SlotsMatchState)
	defer s.SetAllowSpin(true)

	for _, message := range messages {
		if message.GetOpCode() == int64(pb.OpCodeRequest_OPCODE_REQUEST_SPIN) {
			if !s.IsAllowSpin() {
				continue
			}
			bet := &pb.InfoBet{}
			err := p.unmarshaler.Unmarshal(message.GetData(), bet)
			logger.Debug("Recv request add bet user %s , payload %s with parse error %v",
				message.GetUserId(), message.GetData(), err)
			if err != nil {
				continue
			}
			if bet.Chips <= 0 {
				logger.WithField("user id", message.GetUserId()).
					Error("bet with mcb <= 0 ")
				continue
			}
			s.SetAllowSpin(false)
			s.SetBetInfo(bet)
			p.engine.Process(matchState)
			// wallet, err := entity.ReadWalletUser(ctx, nk, logger, s.GetPlayingPresences()[0].GetUserId())
			// if err != nil {
			// 	logger.WithField("error", err.Error()).
			// 		WithField("user id", s.GetPlayingPresences()[0].GetUserId()).
			// 		Error("get profile user failed")
			// 	continue
			// }
			result, err := p.engine.Finish(matchState)
			if err != nil {
				logger.WithField("error", err.Error()).
					Error("engine finish failed")
				continue
			}
			slotDesk := result.(*pb.SlotDesk)
			// slotDesk.BalanceChipsWalletBefore = wallet.Chips
			// slotDesk.BalanceChipsWalletAfter = wallet.Chips + slotDesk.GetChipsWinInSpin() - bet.Chips
			// p.updateChipByResultGameFinish(ctx, logger, nk, &pb.BalanceResult{
			// 	Updates: []*pb.BalanceUpdate{
			// 		{
			// 			UserId:            s.GetPlayingPresences()[0].GetUserId(),
			// 			AmountChipBefore:  slotDesk.BalanceChipsWalletBefore,
			// 			AmountChipCurrent: slotDesk.BalanceChipsWalletAfter,
			// 			AmountChipAdd:     slotDesk.BalanceChipsWalletAfter - slotDesk.BalanceChipsWalletBefore,
			// 		},
			// 	},
			// })
			p.broadcastMessage(logger, dispatcher,
				int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
				slotDesk,
				s.GetPlayingPresences(),
				nil, false)
			continue
		}
	}
}

func (p *processor) NotifyUpdateGameState(
	logger runtime.Logger,
	dispatcher runtime.MatchDispatcher,
	updateState proto.Message,
	matchState interface{},
) {
	p.broadcastMessage(
		logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE),
		updateState, nil, nil, true)
}

func (p *processor) ProcessApplyPresencesLeave(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	matchState interface{},
) {
	s := matchState.(*entity.SlotsMatchState)
	pendingLeaves := s.GetLeavePresences()
	if len(pendingLeaves) == 0 {
		return
	}
	logger.Info("process apply presences")
	s.RemovePresence(pendingLeaves...)
	s.ApplyLeavePresence()
}

func (p *processor) ProcessFinishGame(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	matchState interface{},
) {
	// s := matchState.(*entity.SlotsMatchState)
	// if s.WaitSpinMatrix {
	// 	p.engine.Process(s)
	// 	s.WaitSpinMatrix = false
	// }
	// result, _ := p.engine.Finish(s)
	// updateFinish, ok := result.(*pb.SlotDesk)
	// if !ok {
	// 	logger.WithField("result", result).
	// 		Error("Result from engine can not convert to SlotDesk")
	// 	return
	// }
	// p.broadcastMessage(
	// 	logger, dispatcher,
	// 	int64(pb.OpCodeUpdate_OPCODE_UPDATE_FINISH),
	// 	updateFinish, nil, nil, true)

	// logger.Info("process finish game done %v", updateFinish)
}

func (p *processor) ProcessMessageFromUser(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	messages []runtime.MatchData,
	matchState interface{},
) {
	// s := matchState.(*entity.SlotsMatchState)
	// for _, message := range messages {
	// 	switch pb.OpCodeRequest(message.GetOpCode()) {
	// 	case pb.OpCodeRequest_OPCODE_REQUEST_BET:
	// 		if s.IsAllowBet() == false {
	// 			return
	// 		}
	// 		if s.WaitSpinMatrix {
	// 			return
	// 		}
	// 		bet := &pb.InfoBet{}
	// 		err := p.unmarshaler.Unmarshal(message.GetData(), bet)
	// 		logger.Debug("Recv request add bet user %s , payload %s with parse error %v",
	// 			message.GetUserId(), message.GetData(), err)
	// 		s.ResetUserNotInteract(message.GetUserId())
	// 		s.WaitSpinMatrix = true
	// 	}
	// }
}

func (p *processor) ProcessPresencesJoin(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	presences []runtime.Presence,
	matchState interface{},
) {
	s := matchState.(*entity.SlotsMatchState)
	logger.Info("process presences join %v", presences)
	// update new presence
	newJoins := make([]runtime.Presence, 0)
	for _, presence := range presences {
		// check in list leave pending
		{
			_, found := s.LeavePresences.Get(presence.GetUserId())
			if found {
				s.LeavePresences.Remove(presence.GetUserId())
			} else {
				newJoins = append(newJoins, presence)
			}
		}
	}
	s.AddPresence(ctx, nk, newJoins)
	s.JoinsInProgress -= len(newJoins)
	// for _, presence := range newJoins {
	// 	slotDesk := &pb.SlotDesk{}
	// 	wallet, err := entity.ReadWalletUser(ctx, nk, logger, presence.GetUserId())
	// 	if err != nil {
	// 		logger.WithField("error", err.Error()).
	// 			WithField("user id", presence.GetUserId()).
	// 			Error("get profile user failed")
	// 		return
	// 	}
	// 	matrix := s.GetMatrix()
	// 	slotDesk.Matrix = matrix.ToPbSlotMatrix()
	// 	slotDesk.BalanceChipsWalletBefore = wallet.Chips
	// 	slotDesk.BalanceChipsWalletAfter = wallet.Chips

	// 	p.broadcastMessage(logger, dispatcher,
	// 		int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
	// 		slotDesk,
	// 		[]runtime.Presence{presence},
	// 		nil, false)
	// }
}

func (p *processor) ProcessPresencesLeave(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	presences []runtime.Presence,
	matchState interface{},
) {
	s := matchState.(*entity.SlotsMatchState)
	s.RemovePresence(presences...)
	var listUserId []string
	for _, p := range presences {
		listUserId = append(listUserId, p.GetUserId())
	}
	cgbdb.UpdateUsersPlayingInMatch(ctx, logger, db, listUserId, "")
}

func (p *processor) ProcessPresencesLeavePending(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher,
	presences []runtime.Presence,
	matchState interface{},
) {
	s := matchState.(*entity.SlotsMatchState)
	logger.Info("process presences leave pending %v", presences)
	for _, presence := range presences {
		_, found := s.PlayingPresences.Get(presence.GetUserId())
		if found {
			s.AddLeavePresence(presence)
		} else {
			s.RemovePresence(presence)
		}
	}
}

func (p *processor) broadcastMessage(logger runtime.Logger,
	dispatcher runtime.MatchDispatcher,
	opCode int64,
	data proto.Message,
	presences []runtime.Presence,
	sender runtime.Presence,
	reliable bool,
) error {
	dataJson, err := p.marshaler.Marshal(data)
	if err != nil {
		logger.Error("Error when marshaler data for broadcastMessage")
		return err
	}
	err = dispatcher.BroadcastMessage(opCode, dataJson, presences, sender, true)
	if opCode == int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE) {
		return nil
	}
	logger.Info("broadcast message opcode %v, to %v, data %v", opCode, presences, string(dataJson))
	if err != nil {
		logger.Error("Error BroadcastMessage, message: %s", string(dataJson))
		return err
	}
	return nil
}

func (m *processor) updateChipByResultGameFinish(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, balanceResult *pb.BalanceResult) {
	logger.Info("updateChipByResultGameFinish %v", balanceResult)
	walletUpdates := make([]*runtime.WalletUpdate, 0, len(balanceResult.Updates))
	for _, result := range balanceResult.Updates {
		amountChip := result.AmountChipCurrent - result.AmountChipBefore
		changeset := map[string]int64{
			"chips": amountChip, // Substract amountChip coins to the user's wallet.
		}
		metadata := map[string]interface{}{
			"game_reward": entity.ModuleName,
		}
		walletUpdates = append(walletUpdates, &runtime.WalletUpdate{
			UserID:    result.UserId,
			Changeset: changeset,
			Metadata:  metadata,
		})
	}
	logger.Info("wallet update ctx %v, walletUpdates %v", ctx, walletUpdates)
	_, err := nk.WalletsUpdate(ctx, walletUpdates, true)
	if err != nil {
		payload, _ := json.Marshal(walletUpdates)
		logger.
			WithField("payload", string(payload)).
			WithField("err", err).
			Error("Wallets update error.")
		return
	}
}
