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

func NewJungTreasure(randomIntFn func(int, int) int) lib.Engine {
	e := &jungleTreasure{}
	if randomIntFn == nil {
		e.randomIntFn = entity.RandomInt
	} else {
		e.randomIntFn = randomIntFn
	}
	return e
}

// NewGame implements lib.Engine
func (e *jungleTreasure) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.MatrixSpecial = entity.NewTarzanJungleTreasureMatrix()
	s.MatrixSpecial = entity.ShuffleMatrix(s.MatrixSpecial)
	s.SpinSymbols = nil
	s.NumSpinLeft = 5
	// s.ChipWinByGame[s.CurrentSiXiangGame] = 0
	// s.LineWinByGame[s.CurrentSiXiangGame] = 0
	s.ChipStat.ResetChipWin(0)
	s.ChipStat.ResetLineWin(0)
	e.sureTurnSpinSymboTurnX3 = e.randomIntFn(1, s.NumSpinLeft+1)
	return s, nil
}

// Process implements lib.Engine
func (e *jungleTreasure) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft == 0 {
		return s, entity.ErrorSpinReachMax
	}
	s.IsSpinChange = true
	var randIdx int
	var randSymbol pb.SiXiangSymbol
	if s.NumSpinLeft != e.sureTurnSpinSymboTurnX3 {
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
		Index:  int32(randIdx),
	}
	s.SpinSymbols = []*pb.SpinSymbol{spin}
	s.NumSpinLeft--
	return matchState, nil
}

// Finish implements lib.Engine
func (e *jungleTreasure) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return nil, entity.ErrorSpinNotChange
	}
	lineWin := 0
	for _, spin := range s.SpinSymbols {
		switch spin.Symbol {
		case pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX2:
			s.NumSpinLeft += 2
		case pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN_MORE_TURNX3:
			s.NumSpinLeft += 3
		default:
			symInfo := entity.TarzanJungleTreasureSymbol[spin.Symbol]
			lineWin += e.randomIntFn(int(symInfo.Value.Min), int(symInfo.Value.Max))
		}
	}

	if s.NumSpinLeft <= 0 {
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	slotDesk := &pb.SlotDesk{
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		ChipsMcb:           s.Bet().GetChips(),
		IsFinishGame:       s.NumSpinLeft == 0,
		GameReward:         &pb.GameReward{},
	}
	chipsWin := int64(lineWin * int(slotDesk.ChipsMcb) / 100)
	// s.ChipWinByGame[s.CurrentSiXiangGame] = s.ChipWinByGame[s.CurrentSiXiangGame] + chipsWin
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, chipsWin)
	// s.LineWinByGame[s.CurrentSiXiangGame] = s.LineWinByGame[s.CurrentSiXiangGame] + lineWin
	s.ChipStat.AddLineWin(s.CurrentSiXiangGame, int64(lineWin))
	slotDesk.GameReward.ChipsWin = chipsWin
	// slotDesk.TotalChipsWinByGame = s.ChipWinByGame[s.CurrentSiXiangGame]
	slotDesk.GameReward.TotalChipsWinByGame = s.ChipStat.ChipWin(s.CurrentSiXiangGame)
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if !s.MatrixSpecial.TrackFlip[idx] {
			slotDesk.Matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		}
	})
	slotDesk.GameReward.RatioWin = float32(lineWin) / 100.0
	slotDesk.GameReward.LineWin = int64(lineWin)
	slotDesk.GameReward.TotalLineWin = s.ChipStat.TotalLineWin(s.CurrentSiXiangGame)
	slotDesk.GameReward.TotalRatioWin = float32(s.ChipStat.TotalLineWin(s.CurrentSiXiangGame)) / 100.0
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	return slotDesk, nil
}

// Random implements lib.Engine
func (e *jungleTreasure) Random(min int, max int) int {
	return e.randomIntFn(min, max)
}

func (e *jungleTreasure) Loop(s interface{}) (interface{}, error) {
	return s, nil
}
