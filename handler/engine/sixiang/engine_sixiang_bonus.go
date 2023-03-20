package sixiang

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &sixiangBonusEngine{}

type sixiangBonusEngine struct {
}

func NewSixiangBonusEngine() lib.Engine {
	engine := bonusEngine{}
	return &engine
}

func (e *sixiangBonusEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	matrix := entity.NewMatrixSiXiangBonus()
	s.MatrixSpecial = ShuffleMatrix(matrix)
	s.SpinSymbols = nil
	return s, nil
}

func (e *sixiangBonusEngine) Random(min, max int) int {
	return RandomInt(min, max)
}

func (e *sixiangBonusEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	idRamdom, sym := s.MatrixSpecial.RandomSymbolNotFlip(e.Random)
	row, col := s.MatrixSpecial.RowCol(idRamdom)
	s.SpinSymbols = []*pb.SpinSymbol{
		{
			Symbol: sym,
			Row:    int32(row),
			Col:    int32(col),
		},
	}
	s.MatrixSpecial.Flip(idRamdom)
	return s, nil
}

func (e *sixiangBonusEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	slotDesk := &pb.SlotDesk{}
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	switch s.SpinSymbols[0].Symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_SIXANGBONUS_DRAGONPEARL_GAME:
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_DRAGON_PEARL
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_SIXANGBONUS_LUCKYDRAW_GAME:
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_LUCKDRAW
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_SIXANGBONUS_GOLDPICK_GAME:
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_GOLDPICK
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_SIXANGBONUS_RAPIDPAY_GAME:
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_SIXANGBONUS_RAPIDPAY
	default:
	}
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.IsFinishGame = true
	return slotDesk, nil
}
