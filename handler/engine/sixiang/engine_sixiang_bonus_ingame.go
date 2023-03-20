package sixiang

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &sixiangBonusIngameEngine{}

type sixiangBonusIngameEngine struct {
	ratioBonus  int
	enginesGame map[pb.SiXiangGame]lib.Engine
}

func NewSixiangBonusInGameEngine(
	ratioBonus int) lib.Engine {
	engine := sixiangBonusIngameEngine{
		ratioBonus:  ratioBonus,
		enginesGame: make(map[pb.SiXiangGame]lib.Engine),
	}
	engine.enginesGame[pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL] = NewDragonPearlEngine(nil, nil)
	engine.enginesGame[pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW] = NewLuckyDrawEngine(nil, nil)
	engine.enginesGame[pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK] = NewGoldPickEngine(nil, nil)
	engine.enginesGame[pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY] = NewRapidPayEngine(nil, nil)

	return &engine
}

func (e *sixiangBonusIngameEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	engine := e.enginesGame[s.CurrentSiXiangGame]
	if engine == nil {
		return matchState, ErrorNoGameEngine
	}
	return engine.NewGame(matchState)
}

func (e *sixiangBonusIngameEngine) Random(min, max int) int {
	return RandomInt(min, max)
}

func (e *sixiangBonusIngameEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	engine := e.enginesGame[s.CurrentSiXiangGame]
	if engine == nil {
		return matchState, ErrorNoGameEngine
	}
	return engine.Process(matchState)
}

func (e *sixiangBonusIngameEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	engine := e.enginesGame[s.CurrentSiXiangGame]
	if engine == nil {
		return matchState, ErrorNoGameEngine
	}
	result, err := engine.Finish(matchState)
	if err != nil {
		return result, err
	}
	slotDesk := result.(*pb.SlotDesk)
	slotDesk.ChipsWin *= int64(e.ratioBonus)
	slotDesk.IsInSixiangBonus = true
	return slotDesk, nil
}
