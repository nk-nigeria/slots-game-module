package sixiang

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &rapidPayEngine{}

const (
	defaultRapidPayGemSpin = entity.Row_5 + 1
	defaultAddRatioMcb     = float64(0.1)
)

type rapidPayEngine struct {
	randomIntFn   func(min, max int) int
	randomFloat64 func(min, max float64) float64
}

func NewRapidPayEngine(randomIntFn func(min, max int) int, randomFloat64 func(min, max float64) float64) lib.Engine {
	engine := rapidPayEngine{}
	if randomIntFn != nil {
		engine.randomIntFn = randomIntFn
	} else {
		engine.randomIntFn = entity.RandomInt
	}
	if randomFloat64 != nil {
		engine.randomFloat64 = randomFloat64
	} else {
		engine.randomFloat64 = entity.RandomFloat64
	}
	return &engine
}

func (e *rapidPayEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewMatrixRapidPay()
	s.MatrixSpecial = matrix
	s.SpinSymbols = []*pb.SpinSymbol{}
	s.NumSpinLeft = defaultRapidPayGemSpin
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	return s, nil
}

func (e *rapidPayEngine) Random(min, max int) int {
	return e.randomIntFn(min, max)
}

func (e *rapidPayEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return s, entity.ErrorSpinReadMax
	}
	indexStart := (s.NumSpinLeft - 1) * s.MatrixSpecial.Cols
	arrSpin := s.MatrixSpecial.List[indexStart : indexStart+s.MatrixSpecial.Cols]
	var idRandom int
	var symRandom pb.SiXiangSymbol
	for {
		idRandom = e.randomIntFn(0, len(arrSpin))
		symRandom = entity.ShuffleSlice(arrSpin)[idRandom]
		if symRandom != pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
			break
		}
	}
	row, col := s.MatrixSpecial.RowCol(int(indexStart) + idRandom)
	s.MatrixSpecial.Flip(int(indexStart) + idRandom)
	spin := &pb.SpinSymbol{
		Symbol: symRandom,
		Row:    int32(row),
		Col:    int32(col),
	}
	s.SpinSymbols = []*pb.SpinSymbol{spin}
	s.NumSpinLeft--
	return nil, nil
}

func (e *rapidPayEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotDesk := &pb.SlotDesk{}
	if len(s.SpinSymbols) == 0 {
		return slotDesk, entity.ErrorMissingSpinSymbol
	}
	if s.NumSpinLeft <= 0 || s.SpinSymbols[0].Symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END {
		slotDesk.IsFinishGame = true
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	} else {
		s.NextSiXiangGame = s.CurrentSiXiangGame
	}
	ratio := defaultAddRatioMcb
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	for idx, sym := range s.MatrixSpecial.List {
		if s.MatrixSpecial.TrackFlip[idx] {
			ratio += float64(entity.ListSymbolRapidPay[sym].Value.Min)
		}
	}
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.ChipsMcb = s.Bet().Chips
	slotDesk.ChipsWin = int64(ratio * float64(slotDesk.ChipsMcb))
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	return slotDesk, nil
}
