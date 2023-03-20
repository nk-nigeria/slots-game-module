package sixiang

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &bonusEngine{}

type bonusEngine struct {
	randomIntFn func(min, max int) int
}

func NewBonusEngine(randomIntFn func(min, max int) int) lib.Engine {
	engine := bonusEngine{}
	if randomIntFn != nil {
		engine.randomIntFn = randomIntFn
	} else {
		engine.randomIntFn = RandomInt
	}
	return &engine
}

func (e *bonusEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	matrix := entity.NewMatrixBonusGame()
	s.MatrixSpecial = ShuffleMatrix(matrix)
	// s.ChipsWinInSpecialGame = 0
	s.SpinSymbols = []*pb.SpinSymbol{}
	return s, nil
}

func (e *bonusEngine) Random(min, max int) int {
	return e.randomIntFn(min, max)
}

func (e *bonusEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	id := e.Random(0, len(s.MatrixSpecial.List))
	if s.GetBetInfo().ReqSpecGame > 0 {
		reqBonusGame := pb.SiXiangSymbol(s.GetBetInfo().ReqSpecGame)
		for idx, symbol := range s.MatrixSpecial.List {
			if symbol == reqBonusGame {
				id = idx
				break
			}
		}
	}
	s.MatrixSpecial.TrackFlip[id] = true
	spinSymbol := &pb.SpinSymbol{
		Symbol: s.MatrixSpecial.List[id],
	}
	row, col := s.MatrixSpecial.RowCol(id)
	spinSymbol.Row = int32(row)
	spinSymbol.Col = int32(col)
	s.SpinSymbols = []*pb.SpinSymbol{spinSymbol}
	return s, nil
}

func (e *bonusEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	matrix := s.MatrixSpecial
	slotDesk := &pb.SlotDesk{
		Matrix: &pb.SlotMatrix{
			Rows: int32(matrix.Rows),
			Cols: int32(matrix.Cols),
		},
	}
	slotDesk.Matrix.Lists = make([]pb.SiXiangSymbol, slotDesk.Matrix.Cols*slotDesk.Matrix.Rows)
	// cacl ratio chips by symbol, only goldx10,20,30,50 has ratio > 0
	{
		ratio := entity.ListSymbolBonusGame[s.SpinSymbols[0].Symbol].Value.Min
		slotDesk.ChipsWin = int64(float64(ratio) * float64(s.GetBetInfo().GetChips()))
	}
	// slotDesk.ChipsWinInSpecialGame = s.ChipsWinInSpecialGame
	// slotDesk.ChipsWin = slotDesk.ChipsWinInSpecialGame
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.IsFinishGame = true
	slotDesk.ChipsMcb = s.GetBetInfo().Chips
	return slotDesk, nil
}

func (e *bonusEngine) GetNextSiXiangGame(s *entity.SixiangMatchState) pb.SiXiangGame {
	switch s.SpinSymbols[0].Symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX10,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX20,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX30,
		pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDX50:
		return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_DRAGONBALL:
		return pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_LUCKYDRAW:
		return pb.SiXiangGame_SI_XIANG_GAME_LUCKDRAW
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_GOLDPICK:
		return pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_BONUS_RAPIDPAY:
		return pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY
	default:
		return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
}
