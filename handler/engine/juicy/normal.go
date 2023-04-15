package juicy

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &normal{}

type ratioWild string

const (
	ratioWild1_0 = "1.0"
	ratioWild1_2 = "1.2"
	ratioWild1_5 = "1.5"
	ratioWild2_0 = "2.0"
)

type normal struct {
	randomFn func(min, max int) int
}

func NewNormal(randomIntFn func(int, int) int) lib.Engine {
	e := &normal{}
	if randomIntFn != nil {
		e.randomFn = randomIntFn
	} else {
		e.randomFn = RandomInt
	}
	return e
}

// NewGame implements lib.Engine
func (e *normal) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewSlotMatrix(entity.RowsJuicynMatrix, entity.ColsJuicyMatrix)
	matrix = e.SpinMatrix(matrix, ratioWild1_0)
	s.SetMatrix(matrix)
	return matchState, nil
}

// Process implements lib.Engine
func (e *normal) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := e.SpinMatrix(s.Matrix, ratioWild1_0)
	s.SetMatrix(matrix)
	s.SetWildMatrix(e.WildMatrix(matrix))
	s.SetPaylines(e.Paylines(s.WildMatrix))
	return matchState, nil
}

// Random implements lib.Engine
func (e *normal) Random(min int, max int) int {
	return e.randomFn(min, max)
}

// Finish implements lib.Engine
func (e *normal) Finish(matchState interface{}) (interface{}, error) {
	slotDesk := &pb.SlotDesk{}
	s := matchState.(*entity.SlotsMatchState)
	s.NumScatterSeq = e.countScattersSequent(s.Matrix)
	lineWin := 0
	for _, payline := range s.Paylines() {
		lineWin += int(payline.GetRate())
	}
	// scatter x3 x4 x5 tính điểm tương ứng 3 4 5 x line bet
	if s.NumScatterSeq >= 3 {
		lineWin *= s.NumScatterSeq
	}
	s.RatioFruitBasket = e.transformNumScaterSeqToRationFruitBasket(s.NumScatterSeq)
	s.Matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if entity.IsFruitBasketSymbol(symbol) {
			val := entity.JuicyBasketSymbol[symbol]
			lineWin += int(val.Value.Min) * s.RatioFruitBasket
		}
	})
	s.ChipWinByGame[s.CurrentSiXiangGame] = int64(lineWin) * s.Bet().Chips / 100
	s.LineWinByGame[s.CurrentSiXiangGame] = lineWin
	slotDesk.ChipsMcb = s.Bet().Chips
	slotDesk.ChipsWin = s.ChipWinByGame[s.CurrentSiXiangGame]
	slotDesk.TotalChipsWinByGame = slotDesk.ChipsWin
	slotDesk.Matrix = s.Matrix.ToPbSlotMatrix()
	slotDesk.Paylines = s.Paylines()

	s.NumFruitBasket = e.countFruitBasket(s.Matrix)
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.RatioFruitBasket = int64(s.RatioFruitBasket)
	slotDesk.IsFinishGame = true
	return slotDesk, nil
}

func (e *normal) SpinMatrix(matrix entity.SlotMatrix, ratioWild ratioWild) entity.SlotMatrix {
	var list []pb.SiXiangSymbol
	switch ratioWild {
	case ratioWild1_0:
		list = ShuffleSlice(entity.JuiceAllSymbols)
	case ratioWild1_2:
		list = ShuffleSlice(entity.JuiceAllSymbolsWildRatio1_2)
	case ratioWild1_5:
		list = ShuffleSlice(entity.JuiceAllSymbolsWildRatio1_5)
	case ratioWild2_0:
		list = ShuffleSlice(entity.JuiceAllSymbolsWildRatio2_0)
	default:
		list = ShuffleSlice(entity.JuiceAllSymbols)
	}
	// list = ShuffleSlice(entity.JuiceAllSymbols)
	lenList := len(list)
	spinMatrix := entity.NewSlotMatrix(matrix.Rows, matrix.Cols)
	spinMatrix.List = make([]pb.SiXiangSymbol, spinMatrix.Size)
	for i := 0; i < spinMatrix.Size; i++ {
		idx := e.randomFn(0, lenList)
		randSymbol := list[idx]
		spinMatrix.List[i] = randSymbol
	}
	//  Wild (reel 2 3 4 5)
	spinMatrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD && col == entity.Col_1 {
			for {
				randomId := e.randomFn(0, lenList)
				randSymbol := list[randomId]
				if randSymbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
					continue
				}
				spinMatrix.List[idx] = randSymbol
				return
			}
		}
	})
	return spinMatrix
}

func (e *normal) WildMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	return matrix
}

func (e *normal) GetNextSiXiangGame(s *entity.SlotsMatchState) pb.SiXiangGame {
	if s.NumScatterSeq >= 3 {
		return pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET
	}
	if s.NumFruitBasket >= 6 {
		return pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN
	}
	return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
}

func (e *normal) Paylines(matrix entity.SlotMatrix) []*pb.Payline {
	paylines := make([]*pb.Payline, 0)
	for pair := entity.MapJuicyPaylineIdx.Oldest(); pair != nil; pair = pair.Next() {
		payline := &pb.Payline{
			// Id: int32(idx),
		}
		payline.Id = int32(pair.Key)
		// idx++
		symbols := matrix.ListFromIndexs(pair.Value)
		if len(symbols) == 0 {
			continue
		}
		firstSymbol := symbols[0]
		if entity.IsFruitBasketSymbol(firstSymbol) {
			continue
		}
		numSameSymbol := 0
		for _, symbol := range symbols {
			if symbol == firstSymbol || symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
				numSameSymbol++
				continue
			}
			break
		}
		if numSameSymbol >= 2 {
			rate := entity.RatioJuicyPaylineMap[firstSymbol][int32(numSameSymbol)]
			if rate == 0 {
				continue
			}
			payline.Symbol = firstSymbol
			payline.NumOccur = int32(numSameSymbol)
			payline.Rate = rate
			paylines = append(paylines, payline)
		}
	}
	return paylines
}

func (e *normal) countScattersSequent(matrix entity.SlotMatrix) int {
	numScaterSeq := 0
	matrix.ForEachCol(func(col int, symbols []pb.SiXiangSymbol) {
		for _, symbol := range symbols {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
				numScaterSeq++
				return
			}
		}
	})
	return numScaterSeq
}

func (e *normal) countFruitBasket(matrix entity.SlotMatrix) int {
	numFruitBasket := 0
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if entity.IsFruitBasketSymbol(symbol) {
			numFruitBasket++
		}
	})
	return numFruitBasket
}

func (e *normal) transformNumScaterSeqToRationFruitBasket(numScatterSeq int) int {
	if numScatterSeq < 3 {
		return 1
	}
	return numScatterSeq
}
