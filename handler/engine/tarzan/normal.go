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
	e := &normal{
	}
	if randomIntFn == nil {
		e.randomIntFn=RandomInt 
	} else {
		e.randomIntFn = randomIntFn
	}
	return e
}

// Finish implements lib.Engine
func (*normal) Finish(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// NewGame implements lib.Engine
func (*normal) NewGame(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Process implements lib.Engine
func (*normal) Process(matchState interface{}) (interface{}, error) {
	panic("unimplemented")
}

// Random implements lib.Engine
func (e *normal) Random(min int, max int) int {
	return e.randomIntFn(min, max)
}

func (e *normal) SpinMatrix(matrix entity.SlotMatrix) entity.SlotMatrix {
	tarzanSymbolAppear := false
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
