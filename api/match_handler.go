package api

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler/engine/juicy"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler/engine/sixiang"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler/engine/tarzan"
	"github.com/ciaolink-game-platform/cgb-slots-game-module/handler/sm"

	"github.com/ciaolink-game-platform/cgp-common/define"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

var _ runtime.Match = &MatchHandler{}

const (
	tickRate = 4
)

type MatchHandler struct {
	processor lib.Processor
	machine   lib.UseCase
}

func (m *MatchHandler) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	//panic("implement me")
	return state, ""
}

func NewMatchHandler(
	moduleName string,
	marshaler *protojson.MarshalOptions,
	unmarshaler *protojson.UnmarshalOptions,
) *MatchHandler {
	var matchHandler *MatchHandler
	switch moduleName {
	case define.SixiangGameName:
		matchHandler = &MatchHandler{
			processor: handler.NewMatchProcessor(marshaler, unmarshaler,
				sixiang.NewEngine()),
			machine: lib.NewGameStateMachine(sm.NewSlotsStateMachineState()),
		}
	case define.TarzanGameName:
		matchHandler = &MatchHandler{
			processor: handler.NewMatchProcessor(marshaler, unmarshaler,
				tarzan.NewEngine()),
			machine: lib.NewGameStateMachine(sm.NewSlotsStateMachineState()),
		}
	case define.JuicyGardenName:
		{
			matchHandler = &MatchHandler{
				processor: handler.NewMatchProcessor(marshaler, unmarshaler, juicy.NewEngine()),
				machine:   lib.NewGameStateMachine(sm.NewSlotsStateMachineState()),
			}
		}
	}
	return matchHandler
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
	}
	code, ok := params["code"].(string)
	if !ok {
		logger.Warn("invalid match init parameter \"code\"")
	}
	password, ok := params["password"].(string)
	if !ok {
		logger.Warn("invalid match init parameter \"password\"")
	}
	open := int32(1)
	if password != "" {
		open = 0
	}
	label := &lib.MatchLabel{
		Open:     open,
		Bet:      bet,
		Code:     code,
		Name:     name,
		Password: password,
		MaxSize:  entity.MaxPresences,
	}

	labelJSON, err := json.Marshal(label)
	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, tickRate, ""
	}
	logger.WithField("label", string(labelJSON)).Info("match init")
	matchState := entity.NewSlotsMathState(label)
	if matchState == nil {
		return nil, tickRate, string(labelJSON)
	}
	procPkg := lib.NewProcessorPackage(matchState, m.processor, logger, nil, nil, nil, nil, nil)
	m.machine.TriggerIdle(lib.GetContextWithProcessorPackager(procPkg))
	return matchState, tickRate, string(labelJSON)
}
