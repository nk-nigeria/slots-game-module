package tarzan

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &jungleTreasure{}

type jungleTreasure struct {
	randomIntFn func(int, int) int
	// đảm bảo 100% sẽ có 1 lần ra turn x3
	sureTurnSpinSymboTurnX3 int
}

func NewJungleTrease(randomIntFn func(int, int) int) lib.Engine {
	e := &jungleTreasure{}
	if randomIntFn == nil {
		e.randomIntFn = RandomInt
	} else {
		e.randomIntFn = randomIntFn
	}
	return e
}

// NewGame implements lib.Engine
func (e *jungleTreasure) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.MatrixSpecial = entity.NewTarzanJungleTreasureMatrix()
	s.MatrixSpecial = ShuffleMatrix(s.MatrixSpecial)
	s.SpinSymbols = nil
	s.GemSpin = 5
	s.ChipWinByGame[s.CurrentSiXiangGame] = 0
	e.sureTurnSpinSymboTurnX3 = e.randomIntFn(1, s.GemSpin+1)
	return s, nil
}

// Process implements lib.Engine
func (e *jungleTreasure) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.GemSpin == 0 {
		return s, ErrorSpinReachMaximum
	}
	var randIdx int
	var randSymbol pb.SiXiangSymbol
	if s.GemSpin != e.sureTurnSpinSymboTurnX3 {
		randIdx, randSymbol = s.MatrixSpecial.RandomSymbolNotFlip(e.randomIntFn)
	} else {
		randSymbol = pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3
		s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if symbol == randSymbol {
				randIdx = idx
			}
		})
	}
	if randSymbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3 {
		// nếu đã quay được turnx3,
		// ko cần force ra turnx3 nữa
		e.sureTurnSpinSymboTurnX3 = -1
	}
	s.MatrixSpecial.Flip(randIdx)
	row, col := s.MatrixSpecial.RowCol(randIdx)
	spin := &pb.SpinSymbol{
		Symbol: randSymbol,
		Row:    int32(row),
		Col:    int32(col),
	}
	s.SpinSymbols = []*pb.SpinSymbol{spin}
	s.GemSpin--
	return matchState, nil
}

// Finish implements lib.Engine
func (e *jungleTreasure) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	lineWin := 0
	for _, spin := range s.SpinSymbols {
		switch spin.Symbol {
		case pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX2:
			s.GemSpin += 2
		case pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3:
			s.GemSpin += 3
		default:
			symInfo := entity.TarzanJungleTreasureSymbol[spin.Symbol]
			lineWin += e.randomIntFn(int(symInfo.Value.Min), int(symInfo.Value.Max))
		}
	}
	if s.GemSpin == 0 {
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	slotDesk := &pb.SlotDesk{
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		ChipsMcb:           s.Bet().GetChips(),
		IsFinishGame:       s.GemSpin == 0,
	}
	chipsWin := int64(lineWin/100) * slotDesk.ChipsMcb
	s.ChipWinByGame[s.CurrentSiXiangGame] += chipsWin
	slotDesk.ChipsWin = chipsWin
	slotDesk.TotalChipsWinByGame = s.ChipWinByGame[s.CurrentSiXiangGame]
	return slotDesk, nil
}

// Random implements lib.Engine
func (e *jungleTreasure) Random(min int, max int) int {
	return e.randomIntFn(min, max)
}
