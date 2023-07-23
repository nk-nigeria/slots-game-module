package juicy

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

// todo implement chip win with % total chip win
var _ lib.Engine = &fruitRain{}

type fruitRain struct {
	randomIntFn           func(min, max int) int
	autoRefillGemSpin     bool
	matrixFruitRainBasket []pb.SiXiangSymbol
}

func NewFruitRain(randomIntFn func(min, max int) int) lib.Engine {
	e := &fruitRain{}
	if randomIntFn == nil {
		e.randomIntFn = entity.RandomInt
	} else {
		e.randomIntFn = randomIntFn
	}
	return e
}

// NewGame implements lib.Engine
func (e *fruitRain) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.MatrixSpecial = entity.NewJuicyMatrix()
	s.NumSpinLeft = 3
	e.autoRefillGemSpin = true
	m := entity.NewJuicyFruitRainMaxtrix()
	e.matrixFruitRainBasket = m.List
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	switch s.NumScatterSeq {
	case 3:
		s.RatioFruitBasket = 1
	case 4:
		s.RatioFruitBasket = 2
	case 5:
		s.RatioFruitBasket = 4
	default:
		s.RatioFruitBasket = 1
	}
	s.ChipStat.Reset(s.CurrentSiXiangGame)
	return s, nil
}

// Process implements lib.Engine
func (e *fruitRain) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return s, entity.ErrorSpinReachMax
	}
	// matrix := s.MatrixSpecial
	matrix := e.SpinMatrix(s.MatrixSpecial)
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		// keep symbol if fruitbasket
		if entity.IsFruitBasketSymbol(symbol) {
			return
		}
		newSymbol := matrix.List[idx]
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_SPIN {
			newSymbol = e.matrixFruitRainBasket[idx]
		}
		s.MatrixSpecial.List[idx] = newSymbol
		s.MatrixSpecial.Flip(idx)
		if entity.IsFruitBasketSymbol(newSymbol) {
			s.NumSpinLeft = 3
			e.autoRefillGemSpin = false
		}
	})
	s.NumSpinLeft--
	if s.NumSpinLeft == 0 && e.autoRefillGemSpin {
		s.NumSpinLeft = 3
		e.autoRefillGemSpin = false
	}
	return s, nil
}

// Random implements lib.Engine
func (e *fruitRain) Random(min int, max int) int {
	return e.randomIntFn(min, max)
}

// Finish implements lib.Engine
func (e *fruitRain) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotDesk := &pb.SlotDesk{
		ChipsMcb:   s.Bet().Chips,
		Matrix:     s.MatrixSpecial.ToPbSlotMatrix(),
		GameReward: &pb.GameReward{},
	}
	isFinish := s.NumSpinLeft == 0
	if !isFinish {
		numFruitBasket := 0
		s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if entity.IsFruitBasketSymbol(symbol) {
				numFruitBasket++
			}
			if entity.IsFruitJPSymbol(symbol) {
				switch symbol {
				case pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINI:
					s.WinJp = pb.WinJackpot_WIN_JACKPOT_MINOR
				case pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINOR:
					s.WinJp = pb.WinJackpot_WIN_JACKPOT_MAJOR
				case pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MAJOR:
					s.WinJp = pb.WinJackpot_WIN_JACKPOT_MEGA
				}
			}
		})
		if numFruitBasket == len(s.MatrixSpecial.List) {
			isFinish = true
			s.WinJp = pb.WinJackpot_WIN_JACKPOT_GRAND
		}
	}
	if s.WinJp != pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED {
		isFinish = true
	}
	if isFinish {
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
		lineWin := 0
		s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			val := entity.JuicyBasketSymbol[symbol]
			lineWin += e.randomIntFn(int(val.Value.Min), int(val.Value.Max))

		})
		// TODO:  add calc chip win jackpot
		chipWin := int64(lineWin) * s.Bet().Chips / 100
		slotDesk.GameReward.ChipsWin = chipWin
		slotDesk.GameReward.TotalChipsWinByGame = chipWin
		slotDesk.GameReward.LineWin = int64(lineWin)
		slotDesk.GameReward.TotalLineWin = int64(lineWin)
		s.NumFruitBasket = 0
	}
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	return slotDesk, nil
}

func (e *fruitRain) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *fruitRain) Info(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *fruitRain) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	spinMatrix := entity.NewSlotMatrix(matrix.Rows, matrix.Cols)
	spinMatrix.List = make([]pb.SiXiangSymbol, spinMatrix.Size)
	spinMatrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		randomSymbol := entity.JuicySpintSymbol(entity.JuiceAllSymbols)
		spinMatrix.Flip(idx)
		spinMatrix.List[idx] = randomSymbol
	})
	return spinMatrix
}
