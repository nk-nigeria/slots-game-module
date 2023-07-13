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
	s := matchState.(*entity.SlotsMatchState)
	engine := e.enginesGame[s.CurrentSiXiangGame]
	if engine == nil {
		return matchState, entity.ErrorNoGameEngine
	}
	_, err := engine.NewGame(matchState)
	if err != nil {
		return s, err
	}
	// s.CurrentSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_BONUS
	// s.ChipStat.ResetChipWin(s.CurrentSiXiangGame)
	return s, nil
}

func (e *sixiangBonusIngameEngine) Random(min, max int) int {
	return entity.RandomInt(min, max)
}

func (e *sixiangBonusIngameEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.enginesGame[s.CurrentSiXiangGame]
	if engine == nil {
		return matchState, entity.ErrorNoGameEngine
	}
	// clear after spin
	if s.IsSpinChange {
		s.ClearGameEyePlayed()
	}
	return engine.Process(matchState)
}

func (e *sixiangBonusIngameEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.enginesGame[s.CurrentSiXiangGame]
	if engine == nil {
		return matchState, entity.ErrorNoGameEngine
	}
	result, err := engine.Finish(matchState)
	if err != nil {
		return result, err
	}
	slotDesk, ok := result.(*pb.SlotDesk)
	if !ok {
		return result, nil
	}
	slotDesk = e.processResult(s, slotDesk)
	return slotDesk, nil
}

func (e *sixiangBonusIngameEngine) Loop(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.enginesGame[s.CurrentSiXiangGame]
	if engine == nil {
		return matchState, entity.ErrorNoGameEngine
	}
	result, err := engine.Loop(s)
	if err != nil {
		return nil, err
	}
	if s.IsSpinChange {
		s.ClearGameEyePlayed()
	}
	slotDesk, ok := result.(*pb.SlotDesk)
	if !ok {
		return result, nil
	}
	slotDesk = e.processResult(s, slotDesk)
	return slotDesk, nil
}

func (e *sixiangBonusIngameEngine) processResult(s *entity.SlotsMatchState, slotDesk *pb.SlotDesk) *pb.SlotDesk {
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, -slotDesk.GameReward.ChipsWin)
	slotDesk.GameReward.ChipsWin *= int64(e.ratioBonus)
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, slotDesk.GameReward.ChipsWin)
	slotDesk.GameReward.TotalChipsWinByGame = s.ChipStat.TotalChipWin(s.CurrentSiXiangGame)
	slotDesk.IsInSixiangBonus = true
	slotDesk.SixiangGems = make([]pb.SiXiangGame, 0)
	s.ClearGameEyePlayed()
	s.LastResult = slotDesk
	return slotDesk
}
