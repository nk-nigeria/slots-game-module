package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/cgbdb"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/heroiclabs/nakama-common/runtime"

	"google.golang.org/grpc/codes"
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

	s.SetBetInfo(&pb.InfoBet{
		Chips: -1,
	})
	_, err := p.engine.NewGame(matchState)
	if err != nil {
		logger.WithField("err", err).Error("Engine new game failed")
		return
	}
	// FIXME: remove after test
	// {
	// 	s.CurrentSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS
	// 	s.NextSiXiangGame = s.CurrentSiXiangGame
	// 	p.engine.NewGame(matchState)
	// }

	presence := s.GetPresences()[0]
	saveGame := p.loadSaveGame(ctx, logger, nk, db, dispatcher, s.GetPlayingPresences()[0].GetUserId(), s.Label.Code)
	s.LoadSaveGame(saveGame, func(mcbInSaveGame int64) int64 {
		return p.suggestMcb(ctx, logger, nk, presence.GetUserId(), mcbInSaveGame)
	})
	// p.suggestMcb(ctx, logger, nk, presence.GetUserId(), s)
	s.Bet().EmitNewgameEvent = false

	p.InitSpecialGameDesk(ctx, logger, nk, db, dispatcher, matchState)

	for _, player := range s.GetPlayingPresences() {
		p.getInfoTable(ctx,
			logger, nk, db,
			dispatcher,
			player.GetUserId(),
			s)
	}
}

func (p *processor) ProcessGame(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	messages []runtime.MatchData,
	matchState interface{},
) {
	// logger.Info("tic")
	if p.delayTime.After(time.Now()) {
		return
	}
	s := matchState.(*entity.SlotsMatchState)
	s.InitNewRound()
	defer s.SetAllowSpin(true)
	p.InitSpecialGameDesk(ctx, logger, nk, db, dispatcher, matchState)
	s.Bet().EmitNewgameEvent = false
	// auto run in some game
	for _, message := range messages {
		// if s.CurrentSiXiangGame != s.NextSiXiangGame {
		// 	s.CurrentSiXiangGame = s.NextSiXiangGame
		// 	if s.CurrentSiXiangGame != pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		// 		logger.
		// 			WithField("game", s.CurrentSiXiangGame.String()).
		// 			WithField("next game", s.NextSiXiangGame.String()).
		// 			Info("InitSpecialGameDesk")
		// 		p.InitSpecialGameDesk(ctx, logger, nk, db, dispatcher, matchState)
		// 	} else {
		// 		logger.
		// 			WithField("game", s.CurrentSiXiangGame.String()).
		// 			WithField("next game", s.NextSiXiangGame.String()).
		// 			Info("Ignore InitSpecialGameDesk")
		// 	}
		// 	if s.Bet().EmitNewgameEvent {
		// 		logger.Info("emit handlerRequestGetInfoTable by new game state")
		// 		for _, player := range s.GetPlayingPresences() {
		// 			p.getInfoTable(ctx,
		// 				logger, nk, db,
		// 				dispatcher, player.GetUserId(), s)
		// 		}
		// 	}
		// }
		p.InitSpecialGameDesk(ctx, logger, nk, db, dispatcher, matchState)

		switch message.GetOpCode() {
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_SPIN):
			p.doSpin(ctx, logger, nk, db, dispatcher, message, s)
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_INFO_TABLE):
			logger.Info("handlerRequestGetInfoTable by user request")
			p.getInfoTable(ctx, logger, nk, db, dispatcher, message.GetUserId(), s)
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_BUY_SIXIANG_GEM):
			p.buySixiangGem(ctx, logger, nk, db, dispatcher, message, s)
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_BET):
			p.doChangeBet(ctx, logger, nk, db, dispatcher, message, s)
		}
	}
	{
		res, err := p.engine.Loop(s)
		if err != nil {
			logger.WithField("err", err).Error("loop with error")
		} else if res != nil {
			if slotDesk, ok := res.(*pb.SlotDesk); ok {
				logger.Info("game summary in loop auto")
				slotDesk.InfoBet = s.Bet()
				p.gameSummary(ctx, logger, nk, dispatcher, s.GetPlayingPresences()[0].GetUserId(), s, slotDesk, 0)
			}
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
	p.saveGame(ctx, logger, nk, db, dispatcher, pendingLeaves[0].GetUserId(),
		s.SaveGameJson(), s.Label.Code)
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
		p.getInfoTable(ctx, logger, nk, db, dispatcher, newuser.GetUserId(), s)
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
	if len(listUserId) == 0 {
		return
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

func (p *processor) doSpin(ctx context.Context,
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
	logger.Debug("Recv request bet user %s , payload %s",
		message.GetUserId(), message.GetData())
	if err != nil {
		logger.WithField("err", err.Error()).
			WithField("msg", message.GetData()).
			WithField("user id", message.GetUserId()).
			Error("unmarshal bet info failed")
		p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR),
			&pb.Error{
				Code:  int64(pb.OpCodeUpdate_OPCODE_ERROR),
				Error: entity.ErrorInfoBetInvalid.Error(),
			},
			[]runtime.Presence{s.GetPresence(message.GetUserId())}, nil, false)
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
	chipBetFee := int64(0)
	s.Bet().Id = bet.Id
	// only update new bet in normal game
	if s.CurrentSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		s.SetBetInfo(bet)
		// sub chip fee in wallet
		err := p.checkEnoughChipFromWallet(ctx, logger, nk, message.GetUserId(), s.Bet().Chips)
		if err != nil {
			logger.Error("deduce chip bet failed %s", err.Error())
			p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR), &pb.Error{
				Code:  int64(codes.Aborted),
				Error: entity.ErrorChipNotEnough.Error(),
			}, []runtime.Presence{s.GetPresence(message.GetUserId())}, nil, false)
			return
		}
		chipBetFee = s.Bet().Chips
	} else {
		bet.Chips = s.Bet().GetChips()
		s.SetBetInfo(bet)
	}
	s.Bet().ReqSpecGame = bet.ReqSpecGame
	_, err = p.engine.Process(s)
	if err != nil {
		logger.WithField("error", err.Error()).
			WithField("bet info", s.Bet()).
			WithField("game", s.Label.Code).WithField("state", s.CurrentSiXiangGame).
			Error("engine process failed")
		return
	}
	result, err := p.engine.Finish(s)
	if err != nil {
		logger.WithField("error", err.Error()).
			Error("engine finish failed")
		return
	}
	slotDesk := result.(*pb.SlotDesk)
	slotDesk.InfoBet = s.Bet()
	// logger.WithField("###### @@@@@@ ", s.GameEyePlayed()).
	// WithField("bet", s.Bet()).
	// Info("########")
	p.gameSummary(ctx, logger, nk, dispatcher, message.GetUserId(), s, slotDesk, chipBetFee)
}

func (p *processor) getInfoTable(
	ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	userID string,
	s *entity.SlotsMatchState,
) {
	info, err := p.engine.Info(s)
	if err != nil {
		logger.WithField("user", userID).WithField("err", err.Error()).Error("request info table")
		return
	}
	slotdesk, ok := info.(*pb.SlotDesk)
	if !ok {
		logger.WithField("user", userID).WithField("error", "result can not convert to pb.SlotDesk").Error("request info table")
		return
	}
	logger.WithField("user", userID).Info("request info table")
	if s.LastResult != nil {
		slotdesk.GameReward = s.LastResult.GameReward
	}
	gameReward := slotdesk.GameReward
	if gameReward == nil {
		gameReward = &pb.GameReward{}
	}
	gameReward.UpdateWallet = false
	wallet, err := entity.ReadWalletUser(ctx, nk, logger, s.GetPlayingPresences()[0].GetUserId())
	if err != nil {
		logger.WithField("error", err.Error()).
			WithField("user id", s.GetPlayingPresences()[0].GetUserId()).
			Error("get profile user failed")
	} else {
		gameReward.BalanceChipsWalletBefore = wallet.Chips
		gameReward.BalanceChipsWalletAfter = wallet.Chips
	}
	slotdesk.GameReward = gameReward
	p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
		slotdesk, []runtime.Presence{s.GetPresence(userID)}, nil, true)
}

func (p *processor) buySixiangGem(
	ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	// userID string,
	message runtime.MatchData,
	s *entity.SlotsMatchState,
) {
	// if s.Label.Code != define.SixiangGameName {
	// 	p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR), &pb.Error{
	// 		Code:  int64(codes.Aborted),
	// 		Error: entity.ErrorInvalidRequestGame.Error(),
	// 	}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
	// 	return
	// }
	// numGemCollect := s.NumGameEyePlayed()
	// if numGemCollect < 0 || numGemCollect > 4 {
	// 	p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR), &pb.Error{
	// 		Code:  int64(codes.Aborted),
	// 		Error: entity.ErrorInternal.Error(),
	// 	}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
	// 	return
	// }
	// ratio := entity.PriceBuySixiangGem[numGemCollect]
	// chips := ratio * int(s.Bet().Chips)
	if s.CurrentSiXiangGame != pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		logger.WithField("payload", message.GetData()).Error("buy gem failed, only allow buy in normal game")
		return
	}
	userID := message.GetUserId()
	if userID == "" {
		logger.Error("user id is empty. ignore req buy gem")
		return
	}
	request := &pb.InfoBet{}
	err := p.unmarshaler.Unmarshal(message.GetData(), request)
	if err != nil {
		logger.WithField("payload", message.GetData()).Error("buy gem failed, invalid payload")
		return
	}
	if s.Bet().Chips == 0 {
		logger.Error("buy gem failed, chip mcb is zero")
		return
	}
	if s.NumGameEyePlayed() >= 4 {
		logger.Error("buy gem failed, num gem >=4")
		return
	}
	gemWantBuy := pb.SiXiangGame_SI_XIANG_GAME_UNSPECIFIED
	listSymbol := []pb.SiXiangGame{
		pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY,
		pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
	}
	for _, v := range listSymbol {
		if request.GetId() == int32(v) {
			gemWantBuy = v
			break
		}
	}
	if gemWantBuy == pb.SiXiangGame_SI_XIANG_GAME_UNSPECIFIED {
		logger.WithField("payload", message.GetData()).Error("invalid type gem buy")
		return
	}
	if len(s.GameEyePlayed()) > 0 {
		if _, exist := s.GameEyePlayed()[gemWantBuy]; exist {
			logger.WithField("gem", gemWantBuy.String()).Error("gem alreay collect")
			return
		}
	}
	chips, err := s.PriceBuySixiangGem()
	if err != nil {
		p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR), &pb.Error{
			Code:  int64(codes.Aborted),
			Error: err.Error(),
		}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
		return
	}
	if chips == 0 {
		p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR), &pb.Error{
			Code:  int64(codes.Aborted),
			Error: entity.ErrorInternal.Error(),
		}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
	}
	err = p.checkEnoughChipFromWallet(ctx, logger, nk, userID, int64(chips))
	if err != nil {
		logger.WithField("err", err.Error()).Error("chip not enough for buy gem")
		p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR), &pb.Error{
			Code:  int64(codes.Aborted),
			Error: entity.ErrorChipNotEnough.Error(),
		}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
		return
	}
	err = p.updateChipUser(ctx, logger, nk,
		s.GetPlayingPresences()[0].GetUserId(), s.Label.Code,
		-int64(chips), map[string]interface{}{"action": "buy_gem"},
	)
	if err != nil {
		logger.WithField("err", err.Error()).Error("update chip buy gem failed")
		p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR), &pb.Error{
			Code:  int64(codes.Aborted),
			Error: entity.ErrorInternal.Error(),
		}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
		return
	}
	// add SI XIANG GEMS
	// gamePlayed := s.GameEyePlayed()

	// for _, sym := range listSymbol {
	// 	if _, ok := gamePlayed[sym]; !ok {
	// 		s.AddGameEyePlayed(sym)
	// 		break
	// 	}
	// }
	logger.WithField("payload", request).Info("buy gem success")
	s.AddGameEyePlayed(gemWantBuy)
	// if s.CurrentSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL && s.NumGameEyePlayed() >= 4 {
	// 	s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS
	// }
	if s.CurrentSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		s.NextSiXiangGame = gemWantBuy
	}
	p.saveGame(ctx, logger, nk, db, dispatcher, userID, s.SaveGameJson(), s.Label.Code)
	// p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_BUY_SIXIANG_GEM),
	// 	&pb.InfoBet{}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
	s.Bet().EmitNewgameEvent = true
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

func (m *processor) updateChipUser(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule,
	userId string, gameCode string, amountChipAdd int64, metadata map[string]interface{}) error {
	// logger.Info("updateChipByResultGameFinish %v", balanceResult)
	// logger.WithField("data", balanceResult).Info("update game reward to wallet ")
	if metadata == nil {
		metadata = map[string]interface{}{}
	}
	metadata["game_reward"] = gameCode
	walletUpdates := make([]*runtime.WalletUpdate, 0)
	// for _, result := range balanceResult.Updates {
	amountChip := amountChipAdd
	changeset := map[string]int64{
		"chips": amountChip, // Substract amountChip coins to the user's wallet.
	}
	walletUpdates = append(walletUpdates, &runtime.WalletUpdate{
		UserID:    userId,
		Changeset: changeset,
		Metadata:  metadata,
	})
	// }
	_, err := nk.WalletsUpdate(ctx, walletUpdates, true)
	if err != nil {
		payload, _ := json.Marshal(walletUpdates)
		logger.
			WithField("payload", string(payload)).
			WithField("err", err).
			Error("Wallets update error.")
		return err
	}
	return err
}

func (p *processor) InitSpecialGameDesk(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	matchState interface{}) {
	s, ok := matchState.(*entity.SlotsMatchState)
	if !ok {
		logger.
			WithField("s is nil", s).
			Error("InitSpecialGameDesk failed")
		return
	}
	if s.CurrentSiXiangGame != s.NextSiXiangGame {
		prevGame := s.CurrentSiXiangGame
		s.CurrentSiXiangGame = s.NextSiXiangGame
		p.engine.NewGame(s)
		// if s.CurrentSiXiangGame != pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		logger.
			WithField("prev", prevGame.String()).
			WithField("new game", s.NextSiXiangGame.String()).
			Info("InitSpecialGameDesk success")
		// p.engine.NewGame(s)
		// } else {
		// 	logger.
		// 		WithField("prev", s.CurrentSiXiangGame.String()).
		// 		WithField("new game", s.NextSiXiangGame.String()).
		// 		Info("Ignore InitSpecialGameDesk")

		// }
		if s.Bet().EmitNewgameEvent {
			logger.Info("emit handlerRequestGetInfoTable by new game state")
			for _, player := range s.GetPlayingPresences() {
				p.getInfoTable(ctx,
					logger, nk, db,
					dispatcher, player.GetUserId(), s)
			}
		}
	}

}

func (p *processor) checkValidBetInfo(s *entity.SlotsMatchState, bet *pb.InfoBet) bool {

	switch s.CurrentSiXiangGame {
	case pb.SiXiangGame_SI_XIANG_GAME_NORMAL:
		if bet.Chips <= 0 {
			return false
		}
		for _, betLv := range entity.BetLevels {
			if bet.Chips == betLv {
				return true
			}
		}
		return false
	default:
		return true
	}
}

func (p *processor) reportStatistic(logger runtime.Logger, userId string, slotDesk *pb.SlotDesk, s *entity.SlotsMatchState) {
	// send to statistic
	if slotDesk.IsFinishGame && slotDesk.GameReward != nil {
		// report to operation module
		report := lib.NewReportGame()
		// report.AddFee(totalFee)
		report.AddMatch(&pb.MatchData{
			GameId:   0,
			GameCode: s.Label.Code,
			Mcb:      int64(s.Bet().Chips),
			ChipFee:  slotDesk.GameReward.ChipFee,
		})
		report.AddPlayerData(&pb.PlayerData{
			UserId:  userId,
			Chip:    slotDesk.GameReward.BalanceChipsWalletBefore,
			ChipAdd: slotDesk.GameReward.BalanceChipsWalletAfter - slotDesk.GameReward.BalanceChipsWalletBefore,
		})
		// reportUrl := "http://103.226.250.195:8350"
		data, status, err := report.Commit()
		if err != nil || status > 300 {
			if err != nil {
				logger.Error("Report game (%s) operation -> url %s failed, response %s status %d err %s",
					lib.HostReport, s.Label.Code, string(data), status, err.Error())
			} else {
				logger.Info("Report game (%s) operatio -> %s successful", s.Label.Code)
			}
		}
	}
}

func (p *processor) checkEnoughChipFromWallet(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, userId string, chipRequired int64) error {
	wallet, err := entity.ReadWalletUser(ctx, nk, logger, userId)
	if err != nil {
		return err
	}
	if wallet.Chips < chipRequired {
		return entity.ErrorChipNotEnough
	}
	return nil
}
func (p *processor) gameSummary(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher, userId string, s *entity.SlotsMatchState,
	slotDesk *pb.SlotDesk, chipBetFee int64,
) {
	wallet, err := entity.ReadWalletUser(ctx, nk, logger, s.GetPlayingPresences()[0].GetUserId())
	if err != nil {
		logger.WithField("error", err.Error()).
			WithField("user id", s.GetPlayingPresences()[0].GetUserId()).
			Error("get profile user failed")
		return
	}
	if slotDesk.GameReward == nil {
		slotDesk.GameReward = &pb.GameReward{}
	}
	slotDesk.GameReward.BalanceChipsWalletBefore = wallet.Chips
	slotDesk.GameReward.BalanceChipsWalletAfter = wallet.Chips
	if slotDesk.IsFinishGame {
		if chipBetFee <= 0 && (slotDesk.GameReward != nil && slotDesk.GameReward.ChipsWin <= 0) {
			logger.WithField("user", s.GetPlayingPresences()[0].GetUserId()).
				WithField("current game", slotDesk.CurrentSixiangGame.String()).
				WithField("next game", slotDesk.NextSixiangGame.String()).
				Info("no need update wallet, because chip win <= 0")
		} else {

			gameReward := slotDesk.GameReward
			gameReward.UpdateWallet = true
			gameReward.BalanceChipsWalletBefore = wallet.Chips
			gameReward.ChipBetFee = chipBetFee
			// FIXME: hard code 10%,
			gameReward.ChipFee = gameReward.TotalChipsWinByGame / 10
			chipWinGame := gameReward.TotalChipsWinByGame -
				gameReward.GetChipBetFee() - slotDesk.GameReward.ChipFee
			gameReward.BalanceChipsWalletAfter = gameReward.BalanceChipsWalletBefore + chipWinGame
			// update chip win/loose by game
			p.updateChipUser(ctx, logger, nk,
				s.GetPlayingPresences()[0].GetUserId(),
				s.Label.Code, chipWinGame, nil)
			chipBonus := gameReward.GetPerlGreenForestChips()
			// update bonus chip
			if gameReward.UpdateChipsBonus && chipBonus > 0 {
				gameReward.BalanceChipsWalletAfter += chipBonus
				p.updateChipUser(ctx, logger, nk,
					s.GetPlayingPresences()[0].GetUserId(),
					s.Label.Code, chipBonus, map[string]interface{}{"action": "bonus_perl_green_forest"},
				)
			}
			slotDesk.GameReward = gameReward
		}
	}
	slotDesk.ChipsBuyGem, _ = s.PriceBuySixiangGem()
	slotDesk.BetLevels = make([]int64, 0)
	slotDesk.BetLevels = append(slotDesk.BetLevels, 100, 200, 500, 1000)
	slotDesk.TsUnix = time.Now().Unix()
	// sixiang bonus
	if slotDesk.IsFinishGame && s.NumGameEyePlayed() >= 4 {
		logger.WithField("game", s).Info("collect full 4 gem -> goto sixiang bonus")
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS
	}
	// if slotDesk.CurrentSixiangGame != slotDesk.NextSixiangGame && s.Bet().EmitNewgameEvent {
	// 	p.delayTime = time.Now().Add(2 * time.Second)
	// }
	p.broadcastMessage(logger, dispatcher,
		int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
		slotDesk,
		s.GetPlayingPresences(),
		nil, false)
	p.reportStatistic(logger, userId, slotDesk, s)
}

func (p *processor) saveGame(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule, db *sql.DB,
	dispatcher runtime.MatchDispatcher, userId string, saveGameJson, gameCode string) {
	// save game
	saveGame := &pb.SaveGame{
		Data:           saveGameJson,
		LastUpdateUnix: time.Now().Unix(),
	}
	data, err := json.Marshal(saveGame)
	if err != nil {
		logger.WithField("err", err.Error()).Error("masharl save game sixiang failed")
	} else {
		cgbdb.UpdateUsersSaveGame(ctx, logger, db, userId, gameCode,
			string(data))
	}
}
func (p *processor) loadSaveGame(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule,
	db *sql.DB, dispatcher runtime.MatchDispatcher, userId string, gameCode string) *pb.SaveGame {
	account, err := nk.AccountGetId(ctx, userId)
	if err != nil {
		logger.WithField("err", err.Error()).WithField("user id", userId).Error("get account failed")
		return &pb.SaveGame{}
	}
	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(account.User.GetMetadata()), &metadata)
	if err != nil {
		logger.WithField("err", err.Error()).WithField("user id", userId).Error("marshal account metadata failed")
		return &pb.SaveGame{}
	}
	data, ok := metadata[fmt.Sprintf("savegame.%s", gameCode)].(string)
	if !ok {
		return &pb.SaveGame{}
	}
	saveGame := &pb.SaveGame{}
	json.Unmarshal([]byte(data), &saveGame)
	// expire save game
	if time.Now().Unix()-saveGame.LastUpdateUnix > 30*86400 {
		return &pb.SaveGame{}
	}
	return saveGame
}

// "Gợi ý đưa vào bàn :
// TH1 : user mới chưa chơi bao giờ  -> đưa vào MCB dựa theo số chips mang vào
// TH2 : user đã chơi -> quay lại chơi -> số chips mang vào >= mức bet đã chơi
// -> sever đưa vào lại MCB cũ.
// TH3 : user đã chơi -> quay lại chơi -> số chips mang vào  < mức bet đã chơi
// -> sever đưa vào MCB dựa theo số chips mang vào."
func (p *processor) suggestMcb(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule,
	userId string, mcbInSaveGame int64) int64 {

	//load wallet
	wallet, err := entity.ReadWalletUser(ctx, nk, logger, userId)
	if err != nil {
		logger.WithField("err", err.Error()).
			WithField("user id", userId).
			Error("load wallet user failed")
	}
	// TH2 : user đã chơi -> quay lại chơi -> số chips mang vào >= mức bet đã chơi
	// -> sever đưa vào lại MCB cũ.
	if mcbInSaveGame > 0 && wallet.Chips > mcbInSaveGame {
		// if mcb in savegame not equal any value
		// in bet level -> mcb in savegame invalid
		for _, val := range entity.BetLevels {
			if mcbInSaveGame == val {
				return mcbInSaveGame
			}
		}
	}
	//TH1 : user mới chưa chơi bao giờ  -> đưa vào MCB dựa theo số chips mang vào
	// TH3 : user đã chơi -> quay lại chơi -> số chips mang vào  < mức bet đã chơi
	// -> sever đưa vào MCB dựa theo số chips mang vào."
	betsLevel := make([]int64, len(entity.BetLevels))
	copy(betsLevel, entity.BetLevels)
	// sort mcb desc
	sort.Slice(betsLevel, func(i, j int) bool {
		x := betsLevel[i]
		y := betsLevel[j]
		return x > y
	})
	mcbSuggest := entity.BetLevels[0]
	for _, betLv := range betsLevel {
		if betLv < wallet.Chips {
			mcbSuggest = betLv
			break
		}
	}
	return mcbSuggest
}

func (p *processor) doChangeBet(
	ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	message runtime.MatchData,
	s *entity.SlotsMatchState,
) {
	bet := &pb.InfoBet{}
	err := p.unmarshaler.Unmarshal(message.GetData(), bet)
	// logger.Debug("Recv request add bet user %s , payload %s with parse error %v",
	// 	message.GetUserId(), message.GetData(), err)
	if err != nil {
		logger.WithField("err", err.Error()).
			WithField("msg", message.GetData()).
			WithField("user id", message.GetUserId()).
			Error("unmarshal bet info failed")
		p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR),
			&pb.Error{
				Code:  int64(pb.OpCodeUpdate_OPCODE_ERROR),
				Error: entity.ErrorInfoBetInvalid.Error(),
			},
			[]runtime.Presence{s.GetPresence(message.GetUserId())}, nil, false)
		return
	}
	if !p.checkValidBetInfo(s, bet) {
		logger.WithField("user id", message.GetUserId()).
			WithField("game", s.CurrentSiXiangGame.String()).
			WithField("bet", bet).
			Error("invalid bet ")
		return
	}
	s.SetBetInfo(bet)
	p.getInfoTable(ctx, logger, nk, db, dispatcher, message.GetUserId(), s)
}
