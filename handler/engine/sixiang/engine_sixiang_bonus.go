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
	engine := &sixiangBonusEngine{}
	return engine
}

func (e *sixiangBonusEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewMatrixSiXiangBonus()
	s.MatrixSpecial = matrix
	s.SpinSymbols = nil
	s.NumSpinLeft = 1
	return s, nil
}

func (e *sixiangBonusEngine) Random(min, max int) int {
	return entity.RandomInt(min, max)
}

func (e *sixiangBonusEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return s, entity.ErrorSpinReachMax
	}
	gameSymbolSelect := pb.SiXiangSymbol(s.Bet().Id)
	if s.Bet().Id < 0 || gameSymbolSelect == pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
		return s, entity.ErrorInfoBetInvalid
	}
	s.IsSpinChange = true
	indexFlip := -1
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if gameSymbolSelect == symbol {
			indexFlip = idx
		}
	})
	if indexFlip < 0 {
		return s, entity.ErrorInfoBetInvalid
	}
	// _, sym := s.MatrixSpecial.RandomSymbolNotFlip(e.Random)
	sym := s.MatrixSpecial.Flip(indexFlip)
	row, col := s.MatrixSpecial.RowCol(indexFlip)
	s.SpinSymbols = []*pb.SpinSymbol{
		{
			Symbol: sym,
			Row:    int32(row),
			Col:    int32(col),
		},
	}
	s.NumSpinLeft--
	s.ClearGameEyePlayed()
	return s, nil
}

func (e *sixiangBonusEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{},
	}
	if !s.IsSpinChange {
		return slotDesk, entity.ErrorSpinNotChange
	}
	s.IsSpinChange = false
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
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.IsFinishGame = true
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	return slotDesk, nil
}

func (e *sixiangBonusEngine) Loop(s interface{}) (interface{}, error) {
	return nil, nil
}
