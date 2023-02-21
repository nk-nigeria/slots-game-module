package api

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler/sm"

	"github.com/ciaolink-game-platform/cgp-common/lib"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

var _ runtime.Match = &MatchHandler{}

const (
	tickRate = 2
)

type MatchHandler struct {
	processor lib.Processor
	machine   lib.UseCase
}

func (m *MatchHandler) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	//panic("implement me")
	s := state.(*entity.SlotsMatchState)
	return s, ""
}

func NewMatchHandler(
	marshaler *protojson.MarshalOptions,
	unmarshaler *protojson.UnmarshalOptions,
) *MatchHandler {
	return &MatchHandler{
		processor: handler.NewMatchProcessor(marshaler, unmarshaler,
			handler.NewSlotsEngine()),
		machine: lib.NewGameStateMachine(sm.NewSlotsStateMachineState()),
	}
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	logger.Info("match init: %v", params)
	bet, ok := params["bet"].(int32)
	if !ok {
		logger.Error("invalid match init parameter \"bet\"")
		return nil, 0, ""
	}

	name, ok := params["name"].(string)
	if !ok {
		logger.Warn("invalid match init parameter \"name\"")
		//return nil, 0, ""
	}

	password, ok := params["password"].(string)
	if !ok {
		logger.Warn("invalid match init parameter \"password\"")
		//return nil, 0, ""
	}

	open := int32(1)
	if password != "" {
		open = 0
	}

	// mockCodeCard, _ := params["mock_code_card"].(int32)

	label := &lib.MatchLabel{
		Open:     open,
		Bet:      bet,
		Code:     entity.ModuleName,
		Name:     name,
		Password: password,
		MaxSize:  entity.MaxPresences,
	}

	labelJSON, err := json.Marshal(label)
	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, tickRate, ""
	}

	logger.Info("match init label= %s", string(labelJSON))

	matchState := entity.NewSlotsMathState(label)

	procPkg := lib.NewProcessorPackage(matchState, m.processor, logger, nil, nil, nil, nil, nil)
	m.machine.TriggerIdle(lib.GetContextWithProcessorPackager(procPkg))

	return matchState, tickRate, string(labelJSON)
}
