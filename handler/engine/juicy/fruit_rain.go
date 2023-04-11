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
		e.randomIntFn = RandomInt
	} else {
		e.randomIntFn = randomIntFn
	}
	return e
}

// NewGame implements lib.Engine
func (e *fruitRain) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.MatrixSpecial = entity.NewJuicyMatrix()
	s.GemSpin = 3
	e.autoRefillGemSpin = true
	m := entity.NewJuicyFruitRainMaxtrix()
	e.matrixFruitRainBasket = m.List
	return s, nil
}

// Process implements lib.Engine
func (e *fruitRain) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := s.MatrixSpecial
	matrix = e.SpinMatrix(matrix)
	s.GemSpin--
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if entity.IsFruitBasketSymbol(symbol) {
			return
		}
		newSymbol := matrix.List[idx]
		if entity.IsFruitBasketSymbol(newSymbol) {
			newSymbol = e.matrixFruitRainBasket[idx]
			if newSymbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_7 {
				arr := []pb.SiXiangSymbol{
					pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_RANDOM_7,
					pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINI,
					pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINOR,
					pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MAJOR,
				}
				newSymbol = ShuffleSlice(arr)[e.randomIntFn(0, len(arr))]
			}
		}
		s.MatrixSpecial.List[idx] = newSymbol
		s.MatrixSpecial.Flip(idx)
		if entity.IsFruitBasketSymbol(newSymbol) {
			s.GemSpin = 3
			e.autoRefillGemSpin = false
		}
	})
	if s.GemSpin == 0 && e.autoRefillGemSpin {
		s.GemSpin = 3
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
		ChipsMcb: s.Bet().Chips,
		Matrix:   s.MatrixSpecial.ToPbSlotMatrix(),
	}
	isFinish := s.GemSpin == 0
	if !isFinish {
		numFruitBasket := 0
		s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if entity.IsFruitBasketSymbol(symbol) {
				numFruitBasket++
			}
		})
		if numFruitBasket == len(s.MatrixSpecial.List) {
			isFinish = true
		}
	}
	if isFinish {
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
		chipWin := int64(0)
		lineWin := 0
		s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			val := entity.JuicyBasketSymbol[symbol]
			lineWin += e.randomIntFn(int(val.Value.Min), int(val.Value.Max))
		})
		chipWin = int64(lineWin) * s.Bet().Chips / 100
		slotDesk.ChipsWin = chipWin
		slotDesk.TotalChipsWinByGame = chipWin
	}
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	return slotDesk, nil
}

func (e *fruitRain) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	spinMatrix := entity.NewSlotMatrix(matrix.Rows, matrix.Cols)
	spinMatrix.List = make([]pb.SiXiangSymbol, spinMatrix.Size)
	spinMatrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		for {
			randomIdx := e.randomIntFn(0, len(entity.JuiceAllSymbols))
			randomSymbol := entity.JuiceAllSymbols[randomIdx]
			spinMatrix.List[idx] = randomSymbol
		}
	})
	return spinMatrix
}
