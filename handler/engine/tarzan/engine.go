package tarzan

import (
	"time"

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
	slotDesk.GameReward.PerlGreenForest = int32(s.PerlGreenForest)
	slotDesk.GameReward.PerlGreenForestChips = 0
	slotDesk.GameReward.UpdateChipsBonus = false
	if s.PerlGreenForest >= 100 {
		slotDesk.GameReward.UpdateChipsBonus = true
		perlGreenForestRemain := s.PerlGreenForest - 100
		slotDesk.GameReward.PerlGreenForest = 0
		slotDesk.GameReward.PerlGreenForestChips = s.PerlGreenForestChipsCollect - int64(perlGreenForestRemain)*s.Bet().Chips/2
		s.PerlGreenForest = perlGreenForestRemain
		s.PerlGreenForestChipsCollect = int64(s.PerlGreenForest) * s.Bet().Chips / 2
	}
	slotDesk.GameReward.PerlGreenForestChipsCollect = s.PerlGreenForestChipsCollect
	slotDesk.BigWin = e.transformLineWinToBigWin(int(s.ChipStat.TotalLineWin(s.CurrentSiXiangGame)))
	return slotDesk, err
}

func (e *tarzanEngine) Loop(matchState interface{}) (interface{}, error) {
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

func (e *tarzanEngine) Info(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	var matrix *pb.SlotMatrix
	var spreadMatrix *pb.SlotMatrix
	switch s.CurrentSiXiangGame {
	case pb.SiXiangGame_SI_XIANG_GAME_NORMAL:
		spreadMatrix = s.WildMatrix.ToPbSlotMatrix()
		matrix = spreadMatrix
	case pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9:
		matrix = s.Matrix.ToPbSlotMatrix()
		spreadMatrix = s.WildMatrix.ToPbSlotMatrix()
	case pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE:
		matrix = s.MatrixSpecial.ToPbSlotMatrix()
		for idx, symbol := range s.MatrixSpecial.List {
			if s.MatrixSpecial.IsFlip(idx) {
				matrix.Lists[idx] = symbol
			} else {
				matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
			}
		}
	default:
		matrix = s.MatrixSpecial.ToPbSlotMatrix()
		spreadMatrix = s.MatrixSpecial.ToPbSlotMatrix()
	}
	matrix.SpinLists = s.SpinList
	slotdesk := &pb.SlotDesk{
		Matrix:             matrix,
		SpreadMatrix:       spreadMatrix,
		ChipsMcb:           s.Bet().Chips,
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		TsUnix:             time.Now().Unix(),
		SpinSymbols:        s.SpinSymbols,
		NumSpinLeft:        int64(s.NumSpinLeft),
		InfoBet:            s.Bet(),
		WinJpHistory:       s.WinJPHistory(),
		BetLevels:          entity.BetLevels[:],
		GameReward: &pb.GameReward{
			PerlGreenForest:             int32(s.PerlGreenForest),
			PerlGreenForestChipsCollect: s.PerlGreenForestChipsCollect,
			RatioBonus:                  float32(s.CountLineCrossFreeSpinSymbol),
			UpdateWallet:                false,
			TotalChipsWinByGame:         s.ChipStat.TotalChipWin(s.CurrentSiXiangGame),
			TotalLineWin:                s.ChipStat.TotalLineWin(s.CurrentSiXiangGame),
		},
	}
	if slotdesk.GameReward.RatioBonus < 1 {
		slotdesk.GameReward.RatioBonus = 1
	}
	// slotdesk.ChipsBuyGem, _ = s.PriceBuySixiangGem()
	slotdesk.LetterSymbols = make([]pb.SiXiangSymbol, 0)
	for k := range s.LetterSymbol {
		slotdesk.LetterSymbols = append(slotdesk.LetterSymbols, k)
	}
	// slotdesk.SixiangGems = make([]pb.SiXiangGame, 0)
	// for gem := range s.GameEyePlayed() {
	// 	slotdesk.SixiangGems = append(slotdesk.SixiangGems, gem)
	// }
	return slotdesk, nil
}
