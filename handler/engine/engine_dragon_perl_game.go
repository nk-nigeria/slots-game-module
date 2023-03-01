package engine

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &dragonPearlEngine{}

type dragonPearlEngine struct {
	randomIntFn func(min, max int) int
}

func NewDragonPearlEngine(randomIntFn func(min, max int) int) lib.Engine {
	engine := dragonPearlEngine{}
	if randomIntFn != nil {
		engine.randomIntFn = randomIntFn
	} else {
		engine.randomIntFn = RandomInt
	}
	return &engine
}

func (e *dragonPearlEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewSiXiangMatrixDragonPearl()
	s.MatrixSpecial = ShuffleMatrix(matrix)
	s.ChipsWinInSpecialGame = 0
	s.SpinSymbol = &pb.SpinSymbol{
		Symbol: pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED,
	}
	s.EyeSiXiangRemain = ShuffleSlice(entity.ListEyeSiXiang[:])
	return s, nil
}

func (e *dragonPearlEngine) Random(min, max int) int {
	return e.randomIntFn(min, max)
}

func (e *dragonPearlEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	// bet := s.GetBetInfo()
	idRandom, symbolRandom := s.MatrixSpecial.RandomSymbolNotFlip(e.randomIntFn)
	if symbolRandom == pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_LUCKMONEY {
		symbolRandom = s.EyeSiXiangRemain[0]
		s.EyeSiXiangRemain = s.EyeSiXiangRemain[1:]
	}
	s.MatrixSpecial.TrackFlip[idRandom] = true
	s.SpinSymbol = &pb.SpinSymbol{
		Symbol: symbolRandom,
	}
	row, col := s.MatrixSpecial.RowCol(idRandom)
	s.SpinSymbol.Row = int32(row)
	s.SpinSymbol.Col = int32(col)
	return s, nil
}

func (e *dragonPearlEngine) Finish(matchState interface{}) (interface{}, error) {
	return nil, nil
}
