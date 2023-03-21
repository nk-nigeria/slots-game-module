package tarzan

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &normal{}

type normal struct {
	randomIntFn func(int, int) int
}

func NewNormal(randomIntFn func(int, int) int) lib.Engine {
	e := &normal{}
	if randomIntFn == nil {
		e.randomIntFn = RandomInt
	} else {
		e.randomIntFn = randomIntFn
	}
	return e
}

// NewGame implements lib.Engine
func (e *normal) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.TarzanMatchState)
	s.Matrix = entity.NewTarzanMatrix()
	s.Matrix = e.SpinMatrix(s.Matrix)
	s.FreeSpinSymbolIndexs = make([]int, 0)
	return s, nil
}

// Process implements lib.Engine
func (e *normal) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.TarzanMatchState)
	s.Matrix = e.SpinMatrix(s.Matrix)
	s.SwingMatrix = e.TarzanSwing(s.Matrix)
	s.FreeSpinSymbolIndexs = make([]int, 0)
	s.Matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if entity.TarzanLetterSymbol[symbol] {
			s.AddCollectionSymbol(symbol)
		}
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN {
			s.FreeSpinSymbolIndexs = append(s.FreeSpinSymbolIndexs, idx)
		}
	})
	return matchState, nil
}

// Finish implements lib.Engine
func (e *normal) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.TarzanMatchState)
	slotDesk := &pb.SlotDesk{}
	slotDesk.Paylines = e.Paylines(s.SwingMatrix)
	slotDesk.ChipsMcb = s.Bet.Chips
	ratio := len(slotDesk.Paylines)
	s.Matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if entity.TarzanLetterSymbol[symbol] {
			ratio += 500
		}
	})

	slotDesk.ChipsWin = int64(ratio) * slotDesk.ChipsMcb

	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	return nil, nil
}

// Random implements lib.Engine
func (e *normal) Random(min int, max int) int {
	return e.randomIntFn(min, max)
}

func (e *normal) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	tarzanSymbolAppear := false
	letterSymbolAppear := false
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
	loop:
		for {
			numRandom := e.Random(0, len(entity.TarzanSymbols)-1)
			randSymbol := entity.TarzanSymbols[numRandom]
			switch randSymbol {
			// Tarzan symbol chỉ xuất hiện ở row 5 và chỉ xuất hiện 1 lần
			case pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN:
				if col != entity.Row_5 || tarzanSymbolAppear {
					continue loop
				}
			// chỉ xuất hiện free spin ở row 3, 4, 5
			case pb.SiXiangSymbol_SI_XIANG_SYMBOL_FREE_SPIN:
				if row < entity.Row_3 {
					continue loop
				}
			}
			// Letter symbol only one per spin
			if entity.TarzanLetterSymbol[randSymbol] {
				if letterSymbolAppear {
					continue
				}
				letterSymbolAppear = true
			}
			matrix.List[idx] = randSymbol
			break
		}
	})
	return matrix
}

func (e *normal) TarzanSwing(matrix entity.SlotMatrix) entity.SlotMatrix {
	swingMatrix := entity.SlotMatrix{
		List: make([]pb.SiXiangSymbol, 0, matrix.Size),
		Cols: matrix.Cols,
		Rows: matrix.Rows,
		Size: matrix.Size,
	}
	hasTarzanSymbol := false
	matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_TARZAN {
			hasTarzanSymbol = true
		}
	})
	if !hasTarzanSymbol {
		copy(swingMatrix.List, matrix.List)
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

func (e *normal) PaylineMatrix(matrix entity.SlotMatrix) []*pb.Payline {

	return nil
}

func (e *normal) GetNextSiXiangGame(s *entity.TarzanMatchState) pb.SiXiangGame {
	if s.SizeCollectionSymbol() == len(entity.TarzanLetterSymbol) {
		return pb.SiXiangGame_SI_XIANG_GAME_TARZAN_JUNGLE_TREASURE
	}
	nummFreeSpinSymbolPerCol := 0
	s.Matrix.ForEachCol(func(col int, symbols []pb.SiXiangSymbol) {
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
	return pb.SiXiangGame_SI_XIANG_GAME_TARZAN_NORMAL
}

func (e *normal) Paylines(matrix entity.SlotMatrix) []*pb.Payline {
	paylines := make([]*pb.Payline, 0)
	for pair := entity.PaylineTarzanMapping.Oldest(); pair != nil; pair = pair.Next() {
		paylineIndexs, isPayline := matrix.IsPayline(pair.Value)
		if !isPayline {
			continue
		}
		payline := &pb.Payline{
			// Id: int32(idx),
		}
		payline.Id = int32(pair.Key)
		payline.Symbol = matrix.List[paylineIndexs[0]]
		payline.NumOccur = int32(len(paylineIndexs))
		paylines = append(paylines, payline)
	}
	return paylines
}
