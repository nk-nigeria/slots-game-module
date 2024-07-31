package api

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	"github.com/nakamaFramework/cgb-slots-game-module/handler"
	"github.com/nakamaFramework/cgb-slots-game-module/handler/engine/inca"
	incaclone "github.com/nakamaFramework/cgb-slots-game-module/handler/engine/inca_clone"
	"github.com/nakamaFramework/cgb-slots-game-module/handler/engine/juicy"
	"github.com/nakamaFramework/cgb-slots-game-module/handler/engine/sixiang"
	"github.com/nakamaFramework/cgb-slots-game-module/handler/engine/tarzan"
	"github.com/nakamaFramework/cgb-slots-game-module/handler/sm"
	"github.com/nakamaFramework/cgp-common/define"
	"github.com/nakamaFramework/cgp-common/lib"
	pb "github.com/nakamaFramework/cgp-common/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

var _ runtime.Match = &MatchHandler{}

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

	var processor lib.Processor
	switch moduleName {
	case define.SixiangGameName.String(),
		define.JourneyToTheWest.String():
		processor = handler.NewMatchProcessor(marshaler, unmarshaler, sixiang.NewEngine())
	case define.TarzanGameName.String(),
		define.FortuneFoundFortune.String():
		processor = handler.NewMatchProcessor(marshaler, unmarshaler, tarzan.NewEngine())
	case define.JuicyGardenName.String(),
		define.CryptoRush.String():
		processor = handler.NewMatchProcessor(marshaler, unmarshaler, juicy.NewEngine())
	case define.IncaGameName.String():
		processor = handler.NewMatchProcessor(marshaler, unmarshaler, inca.NewEngine())
	case define.NoelGameName.String(),
		define.FruitGameName.String():
		processor = handler.NewMatchProcessor(marshaler, unmarshaler, incaclone.NewEngine(define.GameName(moduleName)))
	}

	return &MatchHandler{
		processor: processor,
		machine:   lib.NewGameStateMachine(sm.NewSlotsStateMachineState()),
	}
}

func (m *MatchHandler) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	logger.Info("match init: %v", params)
	label, ok := params["data"].(string)
	if !ok {
		logger.WithField("params", params).Error("invalid match init parameter \"data\"")
		return nil, entity.TickRate, ""
	}
	matchInfo := &pb.Match{}
	err := entity.DefaulUnmarshaler.Unmarshal([]byte(label), matchInfo)
	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, entity.TickRate, ""
	}
	matchInfo.MatchId, _ = ctx.Value(runtime.RUNTIME_CTX_MATCH_ID).(string)
	labelJSON, err := entity.DefaultMarshaler.Marshal(matchInfo)

	if err != nil {
		logger.Error("match init json label failed ", err)
		return nil, entity.TickRate, ""
	}

	logger.Info("match init label= %s", string(labelJSON))
	matchState := entity.NewSlotsMathState(matchInfo)
	if matchState == nil {
		return nil, entity.TickRate, string(labelJSON)
	}
	procPkg := lib.NewProcessorPackage(matchState, m.processor, logger, nil, nil, nil, nil, nil)
	m.machine.TriggerIdle(lib.GetContextWithProcessorPackager(procPkg))
	return matchState, entity.TickRate, string(labelJSON)
}
