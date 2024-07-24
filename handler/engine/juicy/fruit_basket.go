package juicy

import (
	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	"github.com/nakamaFramework/cgp-common/lib"

	pb "github.com/nakamaFramework/cgp-common/proto"
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
	matrixSpecial := entity.NewSlotMatrix(1, 2)
	matrixSpecial.List = append(matrixSpecial.List, pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FUIT_SELECT_FREE_GAME)
	matrixSpecial.List = append(matrixSpecial.List, pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FUIT_SELECT_FRUIT_RAIN)
	matrixSpecial = entity.ShuffleMatrix(matrixSpecial)
	s.MatrixSpecial = &matrixSpecial
	s.NumSpinLeft = 1
	s.ChipStat.Reset(s.CurrentSiXiangGame)
	return matchState, nil
}

// Process implements lib.Engine
func (e *fruitBasket) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	// randIdx := e.Random(0, s.MatrixSpecial.Size)
	if s.Bet().Id < 0 {
		return nil, entity.ErrorInfoBetInvalid
	}
	if s.NumSpinLeft <= 0 {
		return nil, entity.ErrorSpinReachMax
	}
	// cheat
	if s.Bet().ReqSpecGame == int32(pb.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME) {
		s.MatrixSpecial.List[s.Bet().Id] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FUIT_SELECT_FREE_GAME
		s.MatrixSpecial.List[(s.Bet().Id+1)%2] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FUIT_SELECT_FRUIT_RAIN
	}
	// end cheat
	s.IsSpinChange = true
	idx := s.Bet().GetId()
	symbol := s.MatrixSpecial.Flip(int(idx))
	row, col := s.MatrixSpecial.RowCol(int(idx))
	s.SpinSymbols = append(s.SpinSymbols, &pb.SpinSymbol{
		Symbol: symbol,
		Index:  idx,
		Row:    int32(row),
		Col:    int32(col),
	})
	s.NumSpinLeft--
	return matchState, nil
}

// Random implements lib.Engine
func (*fruitBasket) Random(min int, max int) int {
	return entity.RandomInt(min, max)
}

// Finish implements lib.Engine
func (*fruitBasket) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return s.LastResult, nil
	}
	s.IsSpinChange = false
	for idx := range s.MatrixSpecial.TrackFlip {
		sym := s.MatrixSpecial.List[idx]
		switch sym {
		case pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FUIT_SELECT_FRUIT_RAIN:
			s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN
		default:
			s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME
		}
	}
	if s.NextSiXiangGame != pb.SiXiangGame_SI_XIANG_GAME_NORMAL {
		s.GameConfig = entity.GameConfigFreeGame(int(s.GameConfig.NumScatterSeq))
	}
	slotDesk := &pb.SlotDesk{
		ChipsMcb:           s.Bet().Chips,
		Matrix:             s.MatrixSpecial.ToPbSlotMatrix(),
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		IsFinishGame:       true,
		NumSpinLeft:        int64(s.NumSpinLeft),
		SpinSymbols:        s.SpinSymbols,
		GameConfig:         s.GameConfig.GameConfig,
	}

	return slotDesk, nil
}

func (e *fruitBasket) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *fruitBasket) Info(s interface{}) (interface{}, error) {
	return s, nil
}
