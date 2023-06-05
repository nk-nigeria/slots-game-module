package juicy

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"

	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &fruitBasket{}

type fruitBasket struct {
}

func NewFruitBaseket() lib.Engine {
	return &fruitBasket{}
}

// NewGame implements lib.Engine
func (*fruitBasket) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	s.MatrixSpecial = entity.NewSlotMatrix(1, 2)
	s.MatrixSpecial.List = append(s.MatrixSpecial.List, pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FUIT_SELECT_FREE_GAME)
	s.MatrixSpecial.List = append(s.MatrixSpecial.List, pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FUIT_SELECT_FRUIT_RAIN)
	s.MatrixSpecial = entity.ShuffleMatrix(s.MatrixSpecial)
	s.NumSpinLeft = 1
	return matchState, nil
}

// Process implements lib.Engine
func (e *fruitBasket) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	randIdx := e.Random(0, s.MatrixSpecial.Size)
	s.MatrixSpecial.Flip(randIdx)
	s.NumSpinLeft--
	return matchState, nil
}

// Random implements lib.Engine
func (*fruitBasket) Random(min int, max int) int {
	return entity.RandomInt(min, max)
}

// Finish implements lib.Engine
func (*fruitBasket) Finish(matchState interface{}) (interface{}, error) {
	slotDesk := &pb.SlotDesk{}
	s := matchState.(*entity.SlotsMatchState)
	slotDesk.ChipsMcb = s.Bet().Chips
	for idx := range s.MatrixSpecial.TrackFlip {
		sym := s.MatrixSpecial.List[idx]
		switch sym {
		case pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FUIT_SELECT_FRUIT_RAIN:
			s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN
		default:
			s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME
		}
	}
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.IsFinishGame = true
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	return slotDesk, nil
}
