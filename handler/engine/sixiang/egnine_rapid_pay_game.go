package sixiang

import (
	"time"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &rapidPayEngine{}

const (
	defaultRapidPayGemSpin  = entity.Row_5 + 1
	defaultAddRatioMcb      = float64(0.1)
	durationTriggerAutoSpin = 2 * time.Second
)

type rapidPayEngine struct {
	randomIntFn   func(min, max int) int
	randomFloat64 func(min, max float64) float64
	lastSpinTime  time.Time
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
	e.lastSpinTime = time.Now()
	return s, nil
}

func (e *rapidPayEngine) Random(min, max int) int {
	return e.randomIntFn(min, max)
}

func (e *rapidPayEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	defer func() {
		e.lastSpinTime = time.Now()
	}()
	if s.NumSpinLeft <= 0 {
		return s, entity.ErrorSpinReachMax
	}
	s.SpinSymbols = make([]*pb.SpinSymbol, 0)
	s.IsSpinChange = true
	indexStart := (s.NumSpinLeft - 1) * s.MatrixSpecial.Cols
	arrSpin := s.MatrixSpecial.List[indexStart : indexStart+s.MatrixSpecial.Cols]
	var idRandom int
	var symRandom pb.SiXiangSymbol
	for {
		idRandom = e.randomIntFn(0, len(arrSpin))
		symRandom = arrSpin[idRandom]
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
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{},
	}
	if len(s.SpinSymbols) == 0 {
		return slotDesk, entity.ErrorMissingSpinSymbol
	}
	if !s.IsSpinChange {
		return slotDesk, entity.ErrorSpinNotChange
	}
	s.IsSpinChange = false
	if s.NumSpinLeft <= 0 || s.SpinSymbols[0].Symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END {
		slotDesk.IsFinishGame = true
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
		s.NumSpinLeft = 0
	} else {
		s.NextSiXiangGame = s.CurrentSiXiangGame
	}
	ratioTotal := defaultAddRatioMcb
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	for idx, sym := range s.MatrixSpecial.List {
		if s.MatrixSpecial.TrackFlip[idx] {
			ratioTotal += float64(entity.ListSymbolRapidPay[sym].Value.Min)
		}
	}
	ratio := float64(0)
	for _, sym := range s.SpinSymbols {
		ratio += float64(entity.ListSymbolRapidPay[sym.GetSymbol()].Value.Min)
	}
	slotDesk.SpreadMatrix = s.MatrixSpecial.ToPbSlotMatrix()
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.ChipsMcb = s.Bet().Chips
	slotDesk.GameReward.TotalChipsWinByGame = int64(ratioTotal * float64(slotDesk.ChipsMcb))
	// s.ChipStat.ResetChipWin(s.CurrentSiXiangGame)
	slotDesk.GameReward.ChipsWin = int64(ratio * float64(slotDesk.ChipsMcb))
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	slotDesk.GameReward.RatioWin = float32(ratioTotal)
	return slotDesk, nil
}

func (e *rapidPayEngine) Loop(s interface{}) (interface{}, error) {
	delay := time.Since(e.lastSpinTime)
	if delay > durationTriggerAutoSpin {
		e.Process(s)
		return e.Finish(s)
	}
	return s, nil
}
