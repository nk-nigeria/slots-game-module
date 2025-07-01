package api

import (
	"context"
	"database/sql"

	pb "github.com/nk-nigeria/cgp-common/proto"
	"github.com/nk-nigeria/slots-game-module/entity"
	"github.com/nk-nigeria/slots-game-module/handler"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nk-nigeria/cgp-common/lib"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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
		machine: lib.NewGameStateMachine(handler.NewSlotsStateMachineState()),
	}
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	rawLabel, ok := params["label"].(string) // đọc label từ param
	if !ok {
		logger.Error("match init label not found")
		return nil, 0, ""
	}

	matchInfo := &pb.Match{}
	err := proto.Unmarshal([]byte(rawLabel), matchInfo)
	if err != nil {
		logger.Error("failed to unmarshal match label: %v", err)
		return nil, 0, ""
	}

	matchInfo.Name = entity.ModuleName
	matchInfo.MaxSize = entity.MaxPresences
	matchInfo.MockCodeCard = 0

	labelJSON, err := protojson.Marshal(matchInfo)
	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, tickRate, ""
	}

	logger.Info("match init label= %s", string(labelJSON))

	matchState := entity.NewSlotsMathState(matchInfo)

	procPkg := lib.NewProcessorPackage(matchState, m.processor, logger, nil, nil, nil, nil, nil)
	m.machine.TriggerIdle(lib.GetContextWithProcessorPackager(procPkg))

	return matchState, tickRate, string(labelJSON)
}
