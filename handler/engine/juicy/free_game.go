package juicy

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &freeGame{}

type freeGame struct {
	ratioWild ratioWild
	normal
}

func NewFreeGame(randomIntFn func(int, int) int) lib.Engine {
	e := &freeGame{}
	if randomIntFn != nil {
		e.randomFn = randomIntFn
	} else {
		e.randomFn = RandomInt
	}
	return e
}

// NewGame implements lib.Engine
func (e *freeGame) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	switch s.NumScatterSeq {
	case 3:
		s.RatioFruitBasket = 1
		e.ratioWild = ratioWild1_2
		s.GemSpin = 6
	case 4:
		s.RatioFruitBasket = 2
		e.ratioWild = ratioWild1_5
		s.GemSpin = 9
	case 5:
		s.RatioFruitBasket = 4
		e.ratioWild = ratioWild2_0
		s.GemSpin = 15
	default:
		s.RatioFruitBasket = 1
	}
	matrix := e.SpinMatrix(s.MatrixSpecial, e.ratioWild)
	s.MatrixSpecial = matrix
	s.ChipWinByGame[s.CurrentSiXiangGame] = 0
	return matchState, nil
}

// Process implements lib.Engine
func (e *freeGame) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.GemSpin <= 0 {
		return matchState, ErrorSpinReadMax
	}
	matrix := e.SpinMatrix(s.Matrix, ratioWild1_0)
	s.MatrixSpecial = (matrix)
	s.SetWildMatrix(e.WildMatrix(matrix))
	s.SetPaylines(e.Paylines(s.WildMatrix))
	s.GemSpin--
	return matchState, nil
}

// Random implements lib.Engine
func (e *freeGame) Random(min int, max int) int {
	return e.randomFn(min, max)
}

// Finish implements lib.Engine
func (e *freeGame) Finish(matchState interface{}) (interface{}, error) {
	slotDesk := &pb.SlotDesk{}
	s := matchState.(*entity.SlotsMatchState)
	lineWin := 0
	for _, payline := range s.Paylines() {
		lineWin += int(payline.GetRate())
	}

	slotDesk.ChipsWin = int64(lineWin) * s.Bet().Chips / 100
	slotDesk.ChipsMcb = s.Bet().Chips
	s.ChipWinByGame[s.CurrentSiXiangGame] += slotDesk.ChipsWin
	slotDesk.TotalChipsWinByGame = s.ChipWinByGame[s.CurrentSiXiangGame]
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	slotDesk.Paylines = s.Paylines()

	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.IsFinishGame = s.GemSpin <= 0
	if slotDesk.IsFinishGame {
		s.RatioFruitBasket = 1
	}
	return slotDesk, nil
}
func (e *freeGame) GetNextSiXiangGame(s *entity.SlotsMatchState) pb.SiXiangGame {
	if s.GemSpin <= 0 {
		return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	return pb.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME
}
