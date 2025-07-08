package sixiang

import (
	"github.com/nk-nigeria/slots-game-module/entity"
	"github.com/nk-nigeria/cgp-common/lib"
	pb "github.com/nk-nigeria/cgp-common/proto"
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
	engine.enginesGame[pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL] = NewDragonPearlEngine(1, nil, nil)
	engine.enginesGame[pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW] = NewLuckyDrawEngine(1, nil, nil)
	engine.enginesGame[pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK] = NewGoldPickEngine(1, nil, nil)
	engine.enginesGame[pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY] = NewRapidPayEngine(1, nil, nil)

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
	if s.LastResult != nil && s.LastResult.GameReward != nil {
		s.LastResult.GameReward.TotalChipsWinByGame *= int64(e.ratioBonus)
		s.LastResult.GameReward.ChipsWin *= int64(e.ratioBonus)
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
	if s.CurrentSiXiangGame != pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY {
		return s, nil
	}
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

func (e *sixiangBonusIngameEngine) Info(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *sixiangBonusIngameEngine) processResult(s *entity.SlotsMatchState, slotDesk *pb.SlotDesk) *pb.SlotDesk {
	// s.ChipStat.AddChipWin(s.CurrentSiXiangGame, -slotDesk.GameReward.ChipsWin)

	// slotDesk.GameReward.ChipsWin *= int64(e.ratioBonus)
	// // s.ChipStat.AddChipWin(s.CurrentSiXiangGame, slotDesk.GameReward.ChipsWin)
	// slotDesk.GameReward.TotalChipsWinByGame *= int64(e.ratioBonus)
	slotDesk.GameReward.ChipsWin *= int64(e.ratioBonus)
	slotDesk.GameReward.TotalChipsWinByGame *= int64(e.ratioBonus)
	slotDesk.IsInSixiangBonus = true
	// slotDesk.SixiangGems = make([]pb.SiXiangGame, 0)
	// slotDesk.GameReward.RatioBonus = float32(e.ratioBonus)
	s.ClearGameEyePlayed()
	s.LastResult = slotDesk
	return slotDesk
}
