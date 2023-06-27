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

	"github.com/ciaolink-game-platform/cgp-common/define"
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
		Chips: 10000,
	})
	_, err := p.engine.NewGame(matchState)
	if err != nil {
		logger.WithField("err", err).Error("Engine new game failed")
		return
	}
	// FIXME: remove after test
	// {
	// 	s.CurrentSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY
	// 	s.NextSiXiangGame = s.CurrentSiXiangGame
	// 	p.engine.NewGame(matchState)
	// }
	if s.GetPresenceSize() <= 0 {
		logger.
			WithField("game", s.Label.Code).
			Info("no player broadcast new game data")
		return
	}
	logger.
		WithField("game", s.Label.Code).
		// WithField("data", s.Matrix).
		Info("new game")
	presence := s.GetPresences()[0]
	saveGame := p.loadSaveGame(ctx, logger, nk, db, dispatcher, s.GetPlayingPresences()[0].GetUserId(), s.Label.Code)
	s.LoadSaveGame(saveGame)
	p.suggestMcb(ctx, logger, nk, presence.GetUserId(), s)
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
		if s.Bet().EmitNewgameEvent {
			for _, player := range s.GetPlayingPresences() {
				p.handlerRequestGetInfoTable(ctx,
					logger, nk, db,
					dispatcher,
					player.GetUserId(),
					s)
			}
		}
	}
	// auto run in some game
	{
		res, _ := p.engine.Loop(s)
		// logger.Info("loop")
		if res != nil {
			if slotDesk, ok := res.(*pb.SlotDesk); ok {
				p.handlerResult(ctx, logger, nk, dispatcher, s.GetPlayingPresences()[0].GetUserId(), s, slotDesk, 0)
			}
		}
	}

	for _, message := range messages {
		switch message.GetOpCode() {
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_SPIN):
			p.handlerRequestSpin(ctx, logger, nk, db, dispatcher, message, s)
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_INFO_TABLE):
			p.handlerRequestGetInfoTable(ctx, logger, nk, db, dispatcher, message.GetUserId(), s)
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_BUY_SIXIANG_GEM):
			p.handlerBuySixiangGemInfoTable(ctx, logger, nk, db, dispatcher, message.GetUserId(), s)
		case int64(pb.OpCodeRequest_OPCODE_REQUEST_BET):
			p.handlerChangeBet(ctx, logger, nk, db, dispatcher, message, s)

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

func (p *processor) handlerRequestSpin(ctx context.Context,
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
	p.handlerResult(ctx, logger, nk, dispatcher, message.GetUserId(), s, slotDesk, chipBetFee)
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
	case pb.SiXiangGame_SI_XIANG_GAME_NORMAL, pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9,
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
	slotdesk.SpreadMatrix = s.MatrixSpecial.ToPbSlotMatrix()
	slotdesk.Matrix.SpinLists = s.SpinList
	slotdesk.NextSixiangGame = s.NextSiXiangGame
	for gem := range s.GameEyePlayed() {
		slotdesk.SixiangGems = append(slotdesk.SixiangGems, gem)
	}
	wallet, err := entity.ReadWalletUser(ctx, nk, logger, s.GetPlayingPresences()[0].GetUserId())
	if err != nil {
		logger.WithField("error", err.Error()).
			WithField("user id", s.GetPlayingPresences()[0].GetUserId()).
			Error("get profile user failed")
	} else {
		gameReward := &pb.GameReward{}
		gameReward.UpdateWallet = true
		gameReward.BalanceChipsWalletBefore = wallet.Chips
		gameReward.BalanceChipsWalletAfter = gameReward.BalanceChipsWalletBefore
		// gameReward.TotalChipsWinByGame = slotdesk.GameReward.TotalChipsWinByGame
		slotdesk.GameReward = gameReward
	}
	slotdesk.NumSpinLeft = int64(s.NumSpinLeft)
	slotdesk.BetLevels = append(slotdesk.BetLevels, entity.BetLevels...)
	p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
		slotdesk, []runtime.Presence{s.GetPresence(userID)}, nil, true)
}

func (p *processor) handlerBuySixiangGemInfoTable(
	ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	userID string,
	s *entity.SlotsMatchState,
) {
	if s.Label.Code != define.SixiangGameName {
		p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR), &pb.Error{
			Code:  int64(codes.Aborted),
			Error: entity.ErrorInvalidRequestGame.Error(),
		}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
		return
	}
	numGemCollect := s.NumGameEyePlayed()
	if numGemCollect < 0 || numGemCollect > 4 {
		p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_ERROR), &pb.Error{
			Code:  int64(codes.Aborted),
			Error: entity.ErrorInternal.Error(),
		}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
		return
	}
	ratio := entity.PriceBuySixiangGem[numGemCollect]
	chips := ratio * int(s.Bet().Chips)
	err := p.checkEnoughChipFromWallet(ctx, logger, nk, userID, int64(chips))
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
	gamePlayed := s.GameEyePlayed()
	listSymbol := []pb.SiXiangGame{
		pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK,
		pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY,
		pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL,
		pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW,
	}
	for _, sym := range listSymbol {
		if _, ok := gamePlayed[sym]; !ok {
			s.AddGameEyePlayed(sym)
			break
		}
	}
	if s.CurrentSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS
	}
	p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_BUY_SIXIANG_GEM),
		&pb.InfoBet{}, []runtime.Presence{s.GetPresence(userID)}, nil, false)
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
			Mcb:      int64(s.Label.Bet),
			ChipFee:  slotDesk.GameReward.ChipFee,
		})
		report.AddPlayerData(&pb.PlayerData{
			UserId:  userId,
			Chip:    slotDesk.GameReward.BalanceChipsWalletAfter,
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
	// chipBet := chipRequired
	if wallet.Chips < chipRequired {
		return entity.ErrorChipNotEnough
	}
	return nil
}
func (p *processor) handlerResult(ctx context.Context, logger runtime.Logger, nk runtime.NakamaModule,
	dispatcher runtime.MatchDispatcher, userId string, s *entity.SlotsMatchState, slotDesk *pb.SlotDesk, chipBetFee int64) {
	if slotDesk.IsFinishGame {
		if chipBetFee <= 0 && (slotDesk.GameReward != nil && slotDesk.GameReward.ChipsWin <= 0) {
			logger.WithField("user", s.GetPlayingPresences()[0].GetUserId()).
				WithField("current game", slotDesk.CurrentSixiangGame.String()).
				WithField("next game", slotDesk.NextSixiangGame.String()).
				Info("no need update wallet, because chip win <= 0")
		} else {
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
			gameReward := slotDesk.GameReward
			gameReward.UpdateWallet = true
			gameReward.BalanceChipsWalletBefore = wallet.Chips
			gameReward.ChipBetFee = chipBetFee
			// FIXME: hard code 10%,
			gameReward.ChipFee = gameReward.TotalChipsWinByGame / 10
			chipWinGame := gameReward.TotalChipsWinByGame -
				gameReward.GetChipBetFee() - slotDesk.GameReward.ChipFee
			gameReward.BalanceChipsWalletAfter = gameReward.BalanceChipsWalletBefore + chipWinGame + gameReward.ChipsBonus
			// update chip win/loose by game
			p.updateChipUser(ctx, logger, nk,
				s.GetPlayingPresences()[0].GetUserId(),
				s.Label.Code, chipWinGame, nil)
			// update bonus chip
			if gameReward.UpdateChipsBonus && gameReward.ChipsBonus > 0 {
				p.updateChipUser(ctx, logger, nk,
					s.GetPlayingPresences()[0].GetUserId(),
					s.Label.Code, gameReward.ChipsBonus, map[string]interface{}{"action": "bonus"},
				)
			}
			slotDesk.GameReward = gameReward
		}
	}
	slotDesk.BetLevels = append(slotDesk.BetLevels, 100, 200, 500, 1000)
	slotDesk.TsUnix = time.Now().Unix()
	if slotDesk.CurrentSixiangGame != slotDesk.NextSixiangGame {
		p.delayTime = time.Now().Add(2 * time.Second)
	} else {
		// sixiang bonus
		if s.NumGameEyePlayed() >= 4 {
			s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS
			p.delayTime = time.Now().Add(2 * time.Second)
		}
	}
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
	userId string, s *entity.SlotsMatchState) {
	//TH1 : user mới chưa chơi bao giờ  -> đưa vào MCB dựa theo số chips mang vào
	if s.Bet().Chips > 0 {
		return
	}
	//load wallet
	wallet, err := entity.ReadWalletUser(ctx, nk, logger, userId)
	if err != nil {
		logger.WithField("err", err.Error()).
			WithField("user id", userId).
			Error("load wallet user failed")
	}
	s.Bet().Chips = entity.BetLevels[0]
	// TH2 : user đã chơi -> quay lại chơi -> số chips mang vào >= mức bet đã chơi
	// -> sever đưa vào lại MCB cũ.
	if wallet.Chips > s.Bet().Chips {
		return
	}
	// TH3 : user đã chơi -> quay lại chơi -> số chips mang vào  < mức bet đã chơi
	// -> sever đưa vào MCB dựa theo số chips mang vào."
	betsLevel := make([]int64, 0, len(entity.BetLevels))
	copy(betsLevel, entity.BetLevels)
	// sort mcb desc
	sort.Slice(betsLevel, func(i, j int) bool {
		x := betsLevel[i]
		y := betsLevel[j]
		return x < y
	})
	for _, betLv := range betsLevel {
		if betLv < wallet.Chips {
			s.Bet().Chips = betLv
			break
		}
	}
}

func (p *processor) handlerChangeBet(
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
	p.handlerRequestGetInfoTable(ctx, logger, nk, db, dispatcher, message.GetUserId(), s)
}
