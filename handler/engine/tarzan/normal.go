package tarzan

import (
	"fmt"
	"math"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &normal{}

type normal struct {
	maxDropTarzanSymbol int
	maxDropFreeSpin     int
	allowDropFreeSpinx9 bool
	maxDropLetterSymbol int
	randomIntFn         func(int, int) int
}

// todo save JUNGLE when exit
func NewNormal(randomIntFn func(int, int) int) lib.Engine {
	e := &normal{
		maxDropTarzanSymbol: 1,
		maxDropLetterSymbol: 1,
		maxDropFreeSpin:     math.MaxInt,
		allowDropFreeSpinx9: true,
	}
	if randomIntFn == nil {
		e.randomIntFn = entity.RandomInt
	} else {
		e.randomIntFn = randomIntFn
	}
	return e

}

// NewGame implements lib.Engine
func (e *normal) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewTarzanMatrix()
	s.SetMatrix(e.SpinMatrix(matrix))
	s.SetWildMatrix(s.Matrix)
	s.TrackIndexFreeSpinSymbol = make(map[int]bool)
	s.NumSpinLeft = -1
	s.SpinList = make([]*pb.SpinSymbol, 0)
	return s, nil
}

// Process implements lib.Engine
func (e *normal) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.IsSpinChange = true
	s.TrackIndexFreeSpinSymbol = make(map[int]bool)
	matrix := e.SpinMatrix(s.Matrix)
	s.NumSpinRemain6thLetter++
	// spin letter symbol
	numLetterPerSpin := 0
	s.SpinSymbols = make([]*pb.SpinSymbol, 0)
	// make sure num spin for 6th reach before appear 6th letter in matrix
	for numLetterPerSpin < e.maxDropLetterSymbol {
		if len(s.LetterSymbol) >= 5 && s.NumSpinRemain6thLetter <= entity.MinNumSpinLetter6th {
			fmt.Printf("reject spin letter symbol cause by num s.LetterSymbol = %d and  s.NumSpinRemain6thLetter(%d) < s.NumSpinRemain6thLetter(%d) not meet \r\n",
				len(s.LetterSymbol), s.NumSpinRemain6thLetter, entity.MinNumSpinLetter6th)
			break
		}
		numLetterPerSpin++
		rIdx := e.randomIntFn(0, 600)
		if rIdx > 100 {
			fmt.Printf("letter symbol not drop, ridx %d \r\n", rIdx)
			continue
		}
		rIdx = e.randomIntFn(int(pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_J), int(pb.SiXiangSymbol_SI_XIANG_SYMBOL_LETTER_E)+1)
		letterSymbol := pb.SiXiangSymbol(rIdx)
		for {
			rIdx = e.randomIntFn(0, matrix.Size)
			sym := matrix.List[rIdx]
			if entity.TarzanLowSymbol[sym] || entity.TarzanMidSymbol[sym] {
				break
			}
		}
		row, col := matrix.RowCol(rIdx)
		fmt.Printf("letter symbol drop %s\r\n", letterSymbol.String())
		s.SpinSymbols = append(s.SpinSymbols, &pb.SpinSymbol{
			Index:  int32(rIdx),
			Symbol: letterSymbol,
			Row:    int32(row),
			Col:    int32(col),
		})
	}
	s.SetMatrix(matrix)
	s.SetWildMatrix(e.TarzanSwing(matrix))
	// cheat custom game
	{
		switch s.Bet().ReqSpecGame {
		case int32(pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9):
			s.Matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
				if col >= entity.Col_3 {
					s.Matrix.List[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN
				}
			})
		case int32(pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE):
			for sym := range entity.TarzanLetterSymbol {
				// s.AddCollectionSymbol(s.CurrentSiXiangGame, 0, sym)
				s.LetterSymbol[sym] = true
			}
		}
	}
	// end set custom game
	newLetterSymbolAppear := false
	// matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
	// 	if entity.TarzanLetterSymbol[symbol] {
	// 		if s.LetterSymbol[symbol] {
	// 			newLetterSymbolAppear = true
	// 		}
	// 		s.LetterSymbol[symbol] = true
	// 	}
	// })
	for _, sym := range s.SpinSymbols {
		symbol := sym.Symbol
		if entity.TarzanLetterSymbol[symbol] {
			if s.LetterSymbol[symbol] {
				newLetterSymbolAppear = true
			}
			s.LetterSymbol[symbol] = true
		}
	}
	// save if spin new letter symbol
	if newLetterSymbolAppear {
		s.SaveGameJson()
	}
	s.PerlGreenForest++
	s.PerlGreenForestChips += s.Bet().GetChips() / 2
	return matchState, nil
}

// Finish implements lib.Engine
func (e *normal) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return s.LastResult, nil
	}
	s.IsSpinChange = false
	paylines := e.Paylines(s.WildMatrix)
	lineWin := int64(len(paylines))
	matrix := s.Matrix
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN {
			lineWin += 200
		}
	})
	chipWin := int64(lineWin * s.Bet().Chips / 100)
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	// if next game is freex9, save index freespin symbol
	if s.NextSiXiangGame == pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9 {
		matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
				s.TrackIndexFreeSpinSymbol[idx] = true
			}
		})
	}
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{
			ChipsWin:            chipWin,
			TotalChipsWinByGame: chipWin,
			RatioWin:            float32(lineWin) / 100,
			LineWin:             lineWin,
			TotalRatioWin:       float32(lineWin),
		},
		ChipsMcb:           s.Bet().Chips,
		Paylines:           paylines,
		Matrix:             s.Matrix.ToPbSlotMatrix(),
		SpreadMatrix:       s.WildMatrix.ToPbSlotMatrix(),
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		IsFinishGame:       true,
		NumSpinLeft:        -1,
		SpinSymbols:        s.SpinSymbols,
	}
	for k := range s.LetterSymbol {
		slotDesk.LetterSymbols = append(slotDesk.LetterSymbols, k)
	}
	s.LastResult = slotDesk
	return slotDesk, nil
}

// Random implements lib.Engine
func (e *normal) Random(min int, max int) int {
	return e.randomIntFn(min, max)
}

func (e *normal) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *normal) Info(matchState interface{}) (interface{}, error) {
	return nil, nil
}

func (e *normal) SpinMatrix(m entity.SlotMatrix) entity.SlotMatrix {
	numTarzanSymbolSpin := 0
	// numLetterSymbolSpin := 0
	numFreeSpinSymbolSpin := 0
	matrix := entity.NewSlotMatrix(m.Rows, m.Cols)
	matrix.List = make([]pb.SiXiangSymbol, m.Size)
	listSymbol := entity.ShuffleSlice(entity.TarzanSymbols)
	lenSymbols := len(listSymbol)
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
	loop:
		for {
			numRandom := e.Random(0, lenSymbols-1)
			randSymbol := listSymbol[numRandom]
			switch randSymbol {
			// Tarzan symbol chỉ xuất hiện ở col 5 và chỉ xuất hiện 1 lần
			case pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN:
				if col != entity.Col_5 || numTarzanSymbolSpin >= e.maxDropTarzanSymbol {
					continue loop
				}
				numTarzanSymbolSpin++
			// chỉ xuất hiện free spin ở col 3, 4, 5
			case pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN:
				if row < entity.Col_3 || numFreeSpinSymbolSpin >= e.maxDropLetterSymbol {
					continue loop
				}
				// kiểm tra điều kiện cho phép ra freespin symbol
				// nhưng không cho phép ra freespinx9 game
				if !e.allowDropFreeSpinx9 && e.countFreeSpinSymbolByCol(matrix) >= 2 {
					continue loop
				}
				numFreeSpinSymbolSpin++
			}
			// // Letter symbol only one per spin
			// if entity.TarzanLetterSymbol[randSymbol] {
			// 	if numLetterSymbolSpin >= e.maxDropLetterSymbol {
			// 		continue loop
			// 	}
			// 	numLetterSymbolSpin++
			// }
			matrix.List[idx] = randSymbol
			break
		}
	})
	return matrix
}

func (e *normal) TarzanSwing(matrix entity.SlotMatrix) entity.SlotMatrix {
	swingMatrix := entity.SlotMatrix{
		List: make([]pb.SiXiangSymbol, matrix.Size),
		Cols: matrix.Cols,
		Rows: matrix.Rows,
		Size: matrix.Size,
	}
	copy(swingMatrix.List, matrix.List)
	hasTarzanSymbol := false
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN {
			hasTarzanSymbol = true
		}
	})
	if !hasTarzanSymbol {
		return swingMatrix
	}
	swingMatrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		isMidSymbol := entity.TarzanMidSymbol[symbol]
		if isMidSymbol {
			swingMatrix.List[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_WILD
		}
	})
	return swingMatrix
}

func (e *normal) GetNextSiXiangGame(s *entity.SlotsMatchState) pb.SiXiangGame {
	if len(s.LetterSymbol) == len(entity.TarzanLetterSymbol) {
		return pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE
	}
	matrix := s.Matrix
	nummFreeSpinSymbolPerCol := 0
	matrix.ForEachCol(func(col int, symbols []pb.SiXiangSymbol) {
		if col < entity.Col_3 {
			return
		}
		for _, sym := range symbols {
			if sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
				nummFreeSpinSymbolPerCol++
				break
			}
		}
	})
	if nummFreeSpinSymbolPerCol >= 3 {
		return pb.SiXiangGame_SI_XIANG_GAME_TARZAN_FREESPINX9
	}
	return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
}

func (e *normal) Paylines(matrix entity.SlotMatrix) []*pb.Payline {
	paylines := make([]*pb.Payline, 0)
	for pair := entity.PaylineTarzanMapping.Oldest(); pair != nil; pair = pair.Next() {
		paylineIndexs, isPayline := matrix.IsPayline(pair.Value)
		if !isPayline {
			continue
		}
		payline := &pb.Payline{
			Indices: make([]int32, 0),
		}
		payline.Id = int32(pair.Key)
		payline.Symbol = matrix.List[paylineIndexs[0]]
		payline.NumOccur = int32(len(paylineIndexs))
		for _, val := range paylineIndexs {
			payline.Indices = append(payline.GetIndices(), int32(val))
		}
		paylines = append(paylines, payline)
	}
	return paylines
}

func (e *normal) countFreeSpinSymbolByCol(matrix entity.SlotMatrix) int {
	count := 0
	matrix.ForEachCol(func(col int, symbols []pb.SiXiangSymbol) {
		for _, symbol := range symbols {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
				count++
				break
			}
		}
	})
	return count
}
