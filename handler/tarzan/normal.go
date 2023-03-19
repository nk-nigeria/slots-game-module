package tarzan

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &normal{}

type normal struct {
}

func NewNormal() lib.Engine {
	return &normal{}
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
func (*normal) Random(min int, max int) int {
	panic("unimplemented")
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
