package tarzan

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &tarzanEngine{}

type tarzanEngine struct {
	engines map[pb.SiXiangGame]lib.Engine
}

func NewEngine() lib.Engine {
	e := &tarzanEngine{
		engines: make(map[pb.SiXiangGame]lib.Engine),
	}
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_NORMAL] = NewNormal(nil)
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9] = NewFreeSpinX9(nil)
	e.engines[pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE] = NewJungTreasure(nil)
	return e
}

// NewGame implements lib.Engine
func (e *tarzanEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	engine.NewGame(matchState)
	return matchState, nil
}

// Process implements lib.Engine
func (e *tarzanEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	s.PerlGreenForest++
	s.PerlGreenForestChips += s.Bet().GetChips() / 2
	return engine.Process(matchState)
}

// Random implements lib.Engine
func (e *tarzanEngine) Random(min int, max int) int {
	panic("unimplemented")
}

// Finish implements lib.Engine
func (e *tarzanEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	result, err := engine.Finish(matchState)
	if err != nil {
		return result, err
	}
	slotDesk := result.(*pb.SlotDesk)
	if slotDesk == nil {
		return result, err
	}
	slotDesk.PerlGreenForest = int32(s.PerlGreenForest)
	slotDesk.GameReward.ChipsBonus = s.PerlGreenForestChips
	slotDesk.GameReward.UpdateChipsBonus = false
	if s.PerlGreenForest >= 100 {
		slotDesk.GameReward.UpdateChipsBonus = true
		s.PerlGreenForestChips = 0
		s.PerlGreenForest = 0
	}

	// slotDesk.BigWin = e.transformLineWinToBigWin(s.LineWinByGame[s.CurrentSiXiangGame])
	slotDesk.BigWin = e.transformLineWinToBigWin(int(s.ChipStat.LineWin(s.CurrentSiXiangGame)))
	// slotDesk.CollectionSymbols = s.CollectionSymbolToSlice(s.CurrentSiXiangGame, 0)
	return slotDesk, err
}

func (e *tarzanEngine) Loop(matchState interface{}) (interface{}, error) {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		fmt.Println("Recovered. Error:\n", r)
	// 	}
	// }()
	s := matchState.(*entity.SlotsMatchState)
	engine := e.engines[s.CurrentSiXiangGame]
	return engine.Loop(s)
}

func (e *tarzanEngine) transformLineWinToBigWin(lineWin int) pb.BigWin {
	if lineWin > 10000 {
		return pb.BigWin_BIG_WIN_MEGA
	}
	if lineWin > 2000 {
		return pb.BigWin_BIG_WIN_HUGE
	}
	if lineWin > 1000 {
		return pb.BigWin_BIG_WIN_BIG
	}
	return pb.BigWin_BIG_WIN_UNSPECIFIED
}
