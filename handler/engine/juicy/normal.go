package juicy

import (
	"fmt"

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
		e.randomFn = entity.RandomInt
	}
	return e
}

// NewGame implements lib.Engine
func (e *normal) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewSlotMatrix(entity.RowsJuicynMatrix, entity.ColsJuicyMatrix)
	matrix = e.SpinMatrix(matrix, ratioWild1_0)
	s.SetMatrix(matrix)
	s.SetWildMatrix(matrix)
	// s.ChipStat.Reset(s.CurrentSiXiangGame)
	s.NumSpinLeft = -1
	return matchState, nil
}

// Process implements lib.Engine
func (e *normal) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.IsSpinChange = true
	// s.ChipStat.Reset(s.CurrentSiXiangGame)
	matrix := e.SpinMatrix(s.Matrix, ratioWild1_0)
	// cheat game
	switch s.Bet().ReqSpecGame {
	case int32(pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET),
		int32(pb.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME):
		for i := 0; i < 3; i++ {
			matrix.List[i] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		}
	case int32(pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN):
		for i := 0; i < 6; i++ {
			matrix.List[i] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_MINI
		}
	}
	// end cheat game
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
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return s.LastResult, nil
	}
	s.IsSpinChange = false
	s.AddChipAccum(s.Bet().GetChips())
	s.NumScatterSeq = e.countScattersSequent(&s.Matrix)
	lineWin := 0
	paylines := s.Paylines()
	for _, payline := range paylines {
		lineWin += int(payline.GetRate())
		payline.Chips = int64(payline.GetRate()) * s.Bet().Chips / 20
	}

	s.RatioFruitBasket = e.transformNumScaterSeqToRationFruitBasket(s.NumScatterSeq)
	s.NumFruitBasket = e.countFruitBasket(&s.Matrix)
	// scatter x3 x4 x5 tính điểm tương ứng 3 4 5 x line bet
	if s.NumScatterSeq >= 3 {
		lineWin *= s.NumScatterSeq
	}
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	chipWin := int64(lineWin) * s.Bet().Chips / 20
	s.PerlGreenForestChipsCollect += s.Bet().Chips
	slotDesk := &pb.SlotDesk{
		ChipsMcb: s.Bet().Chips,
		GameReward: &pb.GameReward{
			ChipsWin:            chipWin,
			TotalChipsWinByGame: chipWin,
			LineWin:             int64(lineWin),
			TotalLineWin:        int64(lineWin),
			RatioWin:            float32(lineWin / 100.0),
			TotalRatioWin:       float32(lineWin / 100.0),
			RatioBonus:          float32(s.NumScatterSeq),
		},
		Matrix:             s.Matrix.ToPbSlotMatrix(),
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		Paylines:           paylines,
		RatioFruitBasket:   int64(s.RatioFruitBasket),
		IsFinishGame:       true,
		NumSpinLeft:        int64(s.NumSpinLeft),
	}
	slotDesk.SpreadMatrix = slotDesk.Matrix
	return slotDesk, nil
}

func (e *normal) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *normal) Info(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *normal) SpinMatrix(matrix entity.SlotMatrix, ratioWild ratioWild) entity.SlotMatrix {
	var list []pb.SiXiangSymbol
	switch ratioWild {
	case ratioWild1_0:
		list = entity.ShuffleSlice(entity.JuiceAllSymbols)
	case ratioWild1_2:
		list = entity.ShuffleSlice(entity.JuiceAllSymbolsWildRatio1_2)
	case ratioWild1_5:
		list = entity.ShuffleSlice(entity.JuiceAllSymbolsWildRatio1_5)
	case ratioWild2_0:
		list = entity.ShuffleSlice(entity.JuiceAllSymbolsWildRatio2_0)
	default:
		list = entity.ShuffleSlice(entity.JuiceAllSymbols)
	}
	// list = ShuffleSlice(entity.JuiceAllSymbols)
	spinMatrix := entity.NewSlotMatrix(matrix.Rows, matrix.Cols)
	spinMatrix.List = make([]pb.SiXiangSymbol, spinMatrix.Size)
	for i := 0; i < spinMatrix.Size; i++ {
		// for {
		randSymbol := entity.JuicySpinSymbol(e.randomFn, list)
		// if randSymbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_SPIN {
		// 	continue
		// }
		spinMatrix.List[i] = randSymbol
		// break
		// }
	}
	//  Wild (reel 2 3 4 5)
	spinMatrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD && col == entity.Col_1 {
			for {
				randSymbol := entity.JuicySpinSymbol(e.randomFn, list)
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
		symbols := matrix.ListFromIndexs(pair.Value)
		if len(symbols) == 0 {
			continue
		}
		payline := &pb.Payline{}
		// get payline win largest chips
		for idx, symbol := range symbols {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD ||
				symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER ||
				symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
				continue
			}
			if entity.IsFruitBasketSymbol(symbol) {
				continue
			}
			if pair.Key == 8 {
				fmt.Println("")
			}
			symbolCompare, indecies := e.countSymbolSeq(matrix, idx, pair.Value)
			newPayline := e.buildPaylineFromSymbol(symbolCompare, len(indecies))
			newPayline.Indices = make([]int32, 0, len(indecies))
			for _, ii := range indecies {
				newPayline.Indices = append(newPayline.Indices, int32(ii))
			}
			if newPayline.Rate > payline.Rate {
				payline = newPayline
			}
		}
		if payline.Rate > 0 {
			payline.Id = int32(pair.Key)
			paylines = append(paylines, payline)
		}
	}
	return paylines
}

func (e *normal) countSymbolSeq(matrix entity.SlotMatrix, startIndex int, indexs []int) (pb.SiXiangSymbol, []int) {
	defer func() {
		if r := recover(); r != nil {
			matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
				fmt.Printf("%d, ", symbol.Number())
			})
			fmt.Println("")
			fmt.Printf("%v \r\n %v", startIndex, indexs)
			panic("fiss")
		}
	}()
	numSameSymbol := 0
	startCount := false
	/*
	 find wild symbol before start index
	 0  1	2		3			4  5
	 x wild wild start_index wild  x
	 -> start_index wild set = 1
	*/
	symbolCount := matrix.List[indexs[startIndex]]
	for i := startIndex - 1; i >= 0; i-- {
		sym := matrix.List[indexs[i]]
		if sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
			startIndex = i
			continue
		}
		break
	}
	newIndexs := indexs[startIndex:]
	paylineIndexs := make([]int, 0, len(newIndexs))
	for idx, symbol := range matrix.ListFromIndexs(newIndexs) {
		if symbol == symbolCount || symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD {
			numSameSymbol++
			startCount = true
			paylineIndexs = append(paylineIndexs, newIndexs[idx])
			continue
		}
		if startCount {
			break
		}
	}
	return symbolCount, paylineIndexs
}

func (e *normal) buildPaylineFromSymbol(symbol pb.SiXiangSymbol, occur int) *pb.Payline {
	rate := entity.RatioJuicyPaylineMap[symbol][int32(occur)]
	payline := &pb.Payline{
		Symbol:   symbol,
		NumOccur: int32(occur),
		Rate:     rate,
	}
	return payline
}

func (e *normal) countScattersSequent(matrix *entity.SlotMatrix) int {
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

func (e *normal) countFruitBasket(matrix *entity.SlotMatrix) int {
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
