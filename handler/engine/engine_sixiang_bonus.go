package engine

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
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewSiXiangMatrixBonusGame()
	s.MatrixSpecial = ShuffleMatrix(matrix)
	return s, nil
}

func (e *sixiangBonusEngine) Random(min, max int) int {
	return RandomInt(min, max)
}

func (e *sixiangBonusEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	bet := s.GetBetInfo()
	s.MatrixSpecial.TrackFlip[int(bet.GetId())] = true
	return s, nil
}

func (e *sixiangBonusEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := s.MatrixSpecial
	slotDesk := &pb.SlotDesk{
		Matrix: &pb.SlotMatrix{
			Rows: int32(matrix.Rows),
			Cols: int32(matrix.Cols),
		},
	}
	slotDesk.Matrix.Lists = make([]pb.SiXiangSymbol, slotDesk.Matrix.Cols*slotDesk.Matrix.Rows)
	drawSymbol := pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
	for k := range matrix.TrackFlip {
		slotDesk.Matrix.Lists[k] = s.MatrixSpecial.List[k]
		drawSymbol = matrix.List[k]
	}

	switch drawSymbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_DRAGONBALL:
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_LUCKYDRAW:
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDPICK:
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_RAPIDPAY:
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY
	default:
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
		ratio := entity.ListSymbolBonusGame[drawSymbol].Value.Min
		slotDesk.ChipsWin = int64(float64(ratio) * float64(s.GetBetInfo().GetChips()))
		// slotDesk.ChipsWin = slotDesk.ChipsWinInSpecialGame
	}
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.IsFinishGame = true
	return slotDesk, nil
}
