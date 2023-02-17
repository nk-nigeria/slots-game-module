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
	engine         lib.Engine
	marshaler      *protojson.MarshalOptions
	unmarshaler    *protojson.UnmarshalOptions
	turnBaseEngine *lib.TurnBaseEngine
}

func NewMatchProcessor(marshaler *protojson.MarshalOptions,
	unmarshaler *protojson.UnmarshalOptions,
	engine lib.Engine) lib.Processor {
	p := processor{
		marshaler:      marshaler,
		unmarshaler:    unmarshaler,
		engine:         engine,
		turnBaseEngine: lib.NewTurnBaseEngine(),
	}
	return &p
}
func (p *processor) ProcessNewGame(logger runtime.Logger,
	dispatcher runtime.MatchDispatcher,
	matchState interface{}) {
	listUserId := make([]string, 0)
	s := matchState.(*entity.SlotsMatchState)
	for _, p := range s.GetPlayingPresences() {
		listUserId = append(listUserId, p.GetUserId())
	}
	p.turnBaseEngine.Config(listUserId, 5*time.Second)
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
	s := matchState.(*entity.SlotsMatchState)
	updateFinish, _ := p.engine.Finish(s)
	// p.broadcastMessage(
	// 	logger, dispatcher,
	// 	int64(pb.OpCodeUpdate_OPCODE_UPDATE_FINISH),
	// 	updateFinish, nil, nil, true)
	logger.Info("process finish game done %v", updateFinish)
}

func (p *processor) ProcessMessageFromUser(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	messages []runtime.MatchData,
	matchState interface{},
) {
	s := matchState.(*entity.SlotsMatchState)
	for _, message := range messages {
		switch pb.OpCodeRequest(message.GetOpCode()) {
		case pb.OpCodeRequest_OPCODE_REQUEST_BET:
			if s.IsAllowBet() == false {
				return
			}
			if s.WaitSpinMatrix {
				return
			}
			bet := &pb.InfoBet{}
			err := p.unmarshaler.Unmarshal(message.GetData(), bet)
			logger.Debug("Recv request add bet user %s , payload %s with parse error %v",
				message.GetUserId(), message.GetData(), err)
			s.ResetUserNotInteract(message.GetUserId())
			s.WaitSpinMatrix = true
		}
	}
}

func (p *processor) ProcessGame(ctx context.Context,
	logger runtime.Logger,
	nk runtime.NakamaModule,
	db *sql.DB,
	dispatcher runtime.MatchDispatcher,
	matchState interface{},
) {
	s := matchState.(*entity.SlotsMatchState)
	if !s.WaitSpinMatrix {
		return
	}
	p.turnBaseEngine.Loop()
	if p.turnBaseEngine.IsNewTurn() {
		p.engine.Process(s)
		s.WaitSpinMatrix = false
		slotsDesk := pb.SlotDesk{
			Matrices: make([]int32, 0),
		}
		matrix, cols, rows := s.GetMatrix()
		for row := 0; row < rows; row++ {
			for col := 0; col < cols; col++ {
				slotsDesk.Matrices = append(slotsDesk.Matrices, int32(matrix[row][col]))
			}
		}
		p.broadcastMessage(logger, dispatcher, int64(pb.OpCodeUpdate_OPCODE_UPDATE_TABLE),
			&slotsDesk, s.GetPlayingPresences(), nil, false)
		return
	}
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
	{
		var listUserId []string
		for _, p := range newJoins {
			listUserId = append(listUserId, p.GetUserId())
		}
		matchId, _ := ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
		playingMatch := &pb.PlayingMatch{
			Code:    entity.ModuleName,
			MatchId: matchId,
		}
		playingMatchJson, _ := json.Marshal(playingMatch)
		cgbdb.UpdateUsersPlayingInMatch(ctx, logger, db, listUserId, string(playingMatchJson))
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
