package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/cgbdb"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/proto"

	"github.com/ciaolink-game-platform/cgp-common/lib"
	"google.golang.org/protobuf/encoding/protojson"
)

var _ lib.Processor = &processor{}

type processor struct {
	engine      lib.Engine
	marshaler   *protojson.MarshalOptions
	unmarshaler *protojson.UnmarshalOptions

	delayTime time.Time
}

func NewMatchProcessor(marshaler *protojson.MarshalOptions,
	unmarshaler *protojson.UnmarshalOptions,
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
	matchState interface{},
) {
	s := matchState.(*entity.SlotsMatchState)
	logger.WithField("size player", s.GetPresenceSize()).Info("ProcessNewGame")
	s.InitNewRound()
	s.SetupMatchPresence()
	s.SetAllowSpin(true)
	s.CurrentSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	s.NextSiXiangGame = s.CurrentSiXiangGame
	_, err := p.engine.NewGame(matchState)
	if err != nil {
		logger.WithField("err", err).Error("Engine new game failed")
		return
	}
	if s.GetPresenceSize() <= 0 {
		logger.
			WithField("game", s.Label.Code).
			Info("no player broadcast new game data")
		return
	}
	logger.
		WithField("game", s.Label.Code).
		WithField("data", s.Matrix).
		Info("new game")

	// slotDesk := &pb.SlotDesk{}
	// presence := s.GetPresences()[0]
	// wallet, err := entity.ReadWalletUser(ctx, nk, logger, presence.GetUserId())
	// if err != nil {
	// 	logger.WithField("error", err.Error()).
	// 		WithField("user id", presence.GetUserId()).
	// 		Error("get profile user failed")
	// 	return
	// }
	// matrix := s.GetMatrix()
	// slotDesk.Matrix = matrix.ToPbSlotMatrix()
	// slotDesk.UpdateWallet = true
	// slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	// slotDesk.NextSixiangGame = s.NextSiXiangGame
	// slotDesk.BalanceChipsWalletBefore = wallet.Chips
	// slotDesk.BalanceChipsWalletAfter = wallet.Chips

	// p.broadcastMessage(logger, dispatcher,
	// 	int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
	// 	slotDesk,
	// 	[]runtime.Presence{presence},
	// 	nil, false)
	presence := s.GetPresences()[0]
	p.handlerRequestGetInfoTable(ctx, logger, nk, db, dispatcher, presence.GetUserId(), s)
}

func (p *processor) ProcessGame(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	messages []runtime.MatchData,
	matchState interface{},
) {
	if p.delayTime.After(time.Now()) {
		return
	}
	s := matchState.(*entity.SlotsMatchState)
	s.InitNewRound()
	defer s.SetAllowSpin(true)
	if s.CurrentSiXiangGame != s.NextSiXiangGame {
		s.CurrentSiXiangGame = s.NextSiXiangGame
		if s.CurrentSiXiangGame != pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
			p.InitSpecialGameDesk(ctx, logger, nk, db, dispatcher, matchState)
		}
		for _, player := range s.GetPlayingPresences() {
			p.handlerRequestGetInfoTable(ctx,
				logger, nk, db,
				dispatcher,
				player.GetUserId(),
				s)
		}
	}

	for _, message := range messages {
		switch message.GetOpCode() {
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_SPIN):
			p.handlerRequestBet(ctx, logger, nk, db, dispatcher, message, s)
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_INFO_TABLE):
			p.handlerRequestGetInfoTable(ctx, logger, nk, db, dispatcher, message.GetUserId(), s)
		}
	}
	p.ProcessApplyPresencesLeave(ctx, logger, nk, db, dispatcher, s)
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
	logger.WithField("size", len(pendingLeaves)).Info("process apply presences")
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
	logger.WithField("users", presences).Info("presences join")
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
	if len(s.Matrix.List) == 0 {
		logger.Debug("game state not init")
		return
	}
	for _, newuser := range newJoins {
		p.handlerRequestGetInfoTable(ctx, logger, nk, db, dispatcher, newuser.GetUserId(), s)
	}
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
	logger.WithField("user", presences).Info("process presences leave pending")
	for _, presence := range presences {
		_, found := s.PlayingPresences.Get(presence.GetUserId())
		if found {
			s.AddLeavePresence(presence)
		} else {
			s.RemovePresence(presence)
		}
	}
}

func (p *processor) handlerRequestBet(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	message runtime.MatchData,
	s *entity.SlotsMatchState,
) {
	if !s.IsAllowSpin() {
		return
	}
	bet := &pb.InfoBet{}
	err := p.unmarshaler.Unmarshal(message.GetData(), bet)
	logger.Debug("Recv request add bet user %s , payload %s with parse error %v",
		message.GetUserId(), message.GetData(), err)
	if err != nil {
		return
	}
	if !p.checkValidBetInfo(s, bet) {
		logger.WithField("user id", message.GetUserId()).
			WithField("game", s.CurrentSiXiangGame.String()).
			WithField("bet", bet).
			Error("invalid bet ")
		return
	}
	s.SetAllowSpin(false)
	// only update new bet in normal game
	if s.CurrentSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		s.SetBetInfo(bet)
	} else {
		bet.Chips = s.Bet().GetChips()
		s.SetBetInfo(bet)
	}
	p.engine.Process(s)

	result, err := p.engine.Finish(s)
	if err != nil {
		logger.WithField("error", err.Error()).
			Error("engine finish failed")
		return
	}
	slotDesk := result.(*pb.SlotDesk)
	if slotDesk.IsFinishGame {
		if slotDesk.ChipsWin <= 0 {
			logger.WithField("user", s.GetPlayingPresences()[0].GetUserId()).
				WithField("current game", slotDesk.CurrentSixiangGame.String()).
				WithField("next game", slotDesk.NextSixiangGame.String()).
				Info("no need update wallet, because chip win <= 0")
		}
		wallet, err := entity.ReadWalletUser(ctx, nk, logger, s.GetPlayingPresences()[0].GetUserId())
		if err != nil {
			logger.WithField("error", err.Error()).
				WithField("user id", s.GetPlayingPresences()[0].GetUserId()).
				Error("get profile user failed")
			return
		}
		slotDesk.UpdateWallet = true
		slotDesk.BalanceChipsWalletBefore = wallet.Chips
		slotDesk.BalanceChipsWalletAfter = wallet.Chips + slotDesk.GetChipsWin() - bet.Chips
		p.updateChipByResultGameFinish(ctx, logger, nk, &pb.BalanceResult{
			Updates: []*pb.BalanceUpdate{
				{
					UserId:            s.GetPlayingPresences()[0].GetUserId(),
					AmountChipBefore:  slotDesk.BalanceChipsWalletBefore,
					AmountChipCurrent: slotDesk.BalanceChipsWalletAfter,
					AmountChipAdd:     slotDesk.BalanceChipsWalletAfter - slotDesk.BalanceChipsWalletBefore,
				},
			},
		})
	}
	slotDesk.TsUnix = time.Now().Unix()
	p.broadcastMessage(logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
		slotDesk,
		s.GetPlayingPresences(),
		nil, false)
	if slotDesk.CurrentSixiangGame != slotDesk.NextSixiangGame {
		p.delayTime = time.Now().Add(2 * time.Second)
	}
}

func (p *processor) handlerRequestGetInfoTable(
	ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	userID string,
	s *entity.SlotsMatchState,
) {
	logger.WithField("user", userID).Info("request info table")
	slotdesk := &pb.SlotDesk{
		ChipsMcb:           s.Bet().Chips,
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		TsUnix:             time.Now().Unix(),
	}
	switch s.CurrentSiXiangGame {
	case pb.SiXiangGame_SI_XIANG_GAME_NORMAL,
		pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9,
		pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN,
		pb.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME:

		matrix := s.Matrix
		slotdesk.Matrix = matrix.ToPbSlotMatrix()
	case pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
		pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE,
		pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET:
		matrix := s.MatrixSpecial
		slotdesk.Matrix = matrix.ToPbSlotMatrix()
		for idx, symbol := range matrix.List {
			if matrix.TrackFlip[idx] {
				slotdesk.Matrix.Lists[idx] = symbol
			} else {
				slotdesk.Matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
			}
		}
	default:
		matrix := s.MatrixSpecial
		slotdesk.Matrix = matrix.ToPbSlotMatrix()

	}
	slotdesk.NextSixiangGame = s.NextSiXiangGame
	wallet, err := entity.ReadWalletUser(ctx, nk, logger, s.GetPlayingPresences()[0].GetUserId())
	if err != nil {
		logger.WithField("error", err.Error()).
			WithField("user id", s.GetPlayingPresences()[0].GetUserId()).
			Error("get profile user failed")
	} else {
		slotdesk.UpdateWallet = true
		slotdesk.BalanceChipsWalletBefore = wallet.Chips
		slotdesk.BalanceChipsWalletAfter = slotdesk.BalanceChipsWalletBefore
	}
	p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
		slotdesk, []runtime.Presence{s.GetPresence(userID)}, nil, true)
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
		logger.WithField("err", err).
			Error("Error when marshaler data for broadcastMessage")
		return err
	}
	err = dispatcher.BroadcastMessage(opCode, dataJson, presences, sender, true)
	if opCode == int64(pb.OpCodeUpdate_OPCODE_UPDATE_GAME_STATE) {
		return nil
	}
	logger.Info("broadcast message opcode %v, to %v, data %v", opCode, presences, string(dataJson))
	if err != nil {
		logger.
			WithField("message", string(dataJson)).
			Error("Error BroadcastMessage")
		return err
	}
	return nil
}

func (m *processor) updateChipByResultGameFinish(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, balanceResult *pb.BalanceResult) {
	// logger.Info("updateChipByResultGameFinish %v", balanceResult)
	logger.WithField("data", balanceResult).Info("update game reward to wallet ")

	walletUpdates := make([]*runtime.WalletUpdate, 0, len(balanceResult.Updates))
	for _, result := range balanceResult.Updates {
		amountChip := result.AmountChipCurrent - result.AmountChipBefore
		changeset := map[string]int64{
			"chips": amountChip, // Substract amountChip coins to the user's wallet.
		}
		metadata := map[string]interface{}{
			"game_reward": entity.SixiangGameName,
		}
		walletUpdates = append(walletUpdates, &runtime.WalletUpdate{
			UserID:    result.UserId,
			Changeset: changeset,
			Metadata:  metadata,
		})
	}
	// logger.Info("wallet update ctx %v, walletUpdates %v", ctx, walletUpdates)
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

func (p *processor) InitSpecialGameDesk(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	matchState interface{}) {
	// s := matchState.(*entity.SlotsMatchState)
	p.engine.NewGame(matchState)

	// slotdesk := &pb.SlotDesk{
	// 	Matrix:             s.MatrixSpecial.ToPbSlotMatrix(),
	// 	ChipsMcb:           s.GetBetInfo().Chips,
	// 	CurrentSixiangGame: s.CurrentSiXiangGame,
	// 	NextSixiangGame:    s.NextSiXiangGame,
	// }
	// p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_INIT_SPECIAL_TABLE),
	// 	slotdesk, s.GetPlayingPresences(), nil, true)
}

func (p *processor) checkValidBetInfo(s *entity.SlotsMatchState, bet *pb.InfoBet) bool {

	switch s.CurrentSiXiangGame {
	case pb.SiXiangGame_SI_XIANG_GAME_NORMAL:
		if bet.Chips <= 0 {
			return false
		}
		return true
	default:
		return true
	}
}
