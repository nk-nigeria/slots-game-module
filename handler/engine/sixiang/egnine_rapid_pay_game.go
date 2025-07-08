package sixiang

import (
	"time"

	"github.com/nk-nigeria/slots-game-module/entity"
	"github.com/nk-nigeria/cgp-common/lib"
	pb "github.com/nk-nigeria/cgp-common/proto"
)

var _ lib.Engine = &rapidPayEngine{}

const (
	defaultRapidPayGemSpin = entity.Row_5 + 1
	defaultAddRatioMcb     = float64(0.5)
	// duration auto spin if no interract after first countdown
	durationAutoSpinNoInteract = 10 * time.Second
	// duration auto spin if no interract first
	durationAutoSpin       = 2 * time.Second
	delayAnimationNewGame  = 8 * time.Second // animation at begin game
	delayAnimationSpinGame = 200 * time.Millisecond
)

type rapidPayEngine struct {
	randomIntFn         func(min, max int) int
	randomFloat64       func(min, max float64) float64
	ratioInSixiangBonus int
}

func NewRapidPayEngine(ratioInSixiangBonus int, randomIntFn func(min, max int) int, randomFloat64 func(min, max float64) float64) lib.Engine {
	engine := rapidPayEngine{
		ratioInSixiangBonus: ratioInSixiangBonus,
	}
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
	s.MatrixSpecial = &matrix
	s.SpinSymbols = []*pb.SpinSymbol{}
	s.NumSpinLeft = defaultRapidPayGemSpin
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	s.LastSpinTime = time.Now()
	s.DurationTriggerAutoSpin = durationAutoSpinNoInteract
	s.DurationTriggerAutoSpin += delayAnimationNewGame //(8s animation in client)
	// s.ResetCollection(s.CurrentSiXiangGame, int(s.Bet().Chips))
	s.ChipStat.Reset(s.CurrentSiXiangGame)
	s.SpinList = make([]*pb.SpinSymbol, 0)
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		s.SpinList = append(s.SpinList, &pb.SpinSymbol{
			Symbol:    pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED,
			Row:       int32(row),
			Col:       int32(col),
			Index:     int32(idx),
			Ratio:     0,
			WinJp:     pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED,
			WinAmount: 0,
		})
	})
	return s, nil
}

func (e *rapidPayEngine) Random(min, max int) int {
	return e.randomIntFn(min, max)
}

func (e *rapidPayEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return s, entity.ErrorSpinReachMax
	}
	defer func() {
		s.LastSpinTime = time.Now()
		s.DurationTriggerAutoSpin = durationAutoSpinNoInteract + delayAnimationSpinGame
	}()
	s.SpinSymbols = make([]*pb.SpinSymbol, 0)
	s.IsSpinChange = true
	indexStart := (s.NumSpinLeft - 1) * s.MatrixSpecial.Cols
	{
		beginId := 0
		for _, sym := range entity.ShuffleSlice(e.symbolsAtRow(s.NumSpinLeft - 1)) {
			for {
				if beginId >= s.MatrixSpecial.Cols {
					break
				}
				if s.MatrixSpecial.List[indexStart+beginId] != pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_LUCKYBOX {
					beginId++
					continue
				}
				s.MatrixSpecial.List[indexStart+beginId] = sym
				beginId++
				break
			}
		}
	}
	arrSpin := s.MatrixSpecial.List[indexStart : indexStart+s.MatrixSpecial.Cols]
	var idRandom int = int(s.Bet().Id)
	var symRandom pb.SiXiangSymbol
	if idRandom <= 0 || idRandom >= len(arrSpin) {
		idRandom = e.randomIntFn(0, len(arrSpin))
	}
	for {
		symRandom = arrSpin[idRandom]
		if symRandom != pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
			break
		}
		idRandom = e.randomIntFn(0, len(arrSpin))
	}
	row, col := s.MatrixSpecial.RowCol(int(indexStart) + idRandom)
	s.MatrixSpecial.Flip(int(indexStart) + idRandom)
	spin := &pb.SpinSymbol{
		Symbol: symRandom,
		Row:    int32(row),
		Col:    int32(col),
		Index:  int32(indexStart) + int32(idRandom),
	}
	s.SpinSymbols = []*pb.SpinSymbol{spin}
	s.SpinList[spin.Index] = spin
	s.NumSpinLeft--
	return nil, nil
}

func (e *rapidPayEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{},
		ChipsMcb:   s.Bet().Chips,
	}
	if len(s.SpinSymbols) == 0 {
		return s.LastResult, entity.ErrorMissingSpinSymbol
	}
	if !s.IsSpinChange {
		return s.LastResult, entity.ErrorSpinNotChange
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
			v := entity.ListSymbolRapidPay[sym].Value.Min
			if v > 0 {
				ratioTotal *= float64(v)
			}
		}
	}
	ratioTotal *= float64(e.ratioInSixiangBonus)
	ratio := float64(defaultAddRatioMcb)
	for _, sym := range s.SpinSymbols {
		sym.Ratio = entity.ListSymbolRapidPay[sym.GetSymbol()].Value.Min * float32(e.ratioInSixiangBonus)
		s.SpinList[sym.Index].Ratio = sym.Ratio
		s.SpinList[sym.Index].WinAmount = int64(sym.Ratio*10) * int64(slotDesk.ChipsMcb) / 10
		ratio *= float64(sym.Ratio)
	}
	slotDesk.SpreadMatrix = s.MatrixSpecial.ToPbSlotMatrix()
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.Matrix.SpinLists = s.SpinList
	slotDesk.GameReward.TotalChipsWinByGame = int64(ratioTotal*10) * int64(slotDesk.ChipsMcb) / 10
	// s.ChipStat.ResetChipWin(s.CurrentSiXiangGame)
	slotDesk.GameReward.ChipsWin = int64(ratio * float64(slotDesk.ChipsMcb))
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	slotDesk.GameReward.RatioWin = float32(ratioTotal)
	if slotDesk.IsFinishGame {
		s.AddGameEyePlayed(pb.SiXiangGame_SI_XIANG_GAME_RAPIDPAY)
	}
	s.LastResult = slotDesk
	return slotDesk, nil
}

func (e *rapidPayEngine) Loop(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	delay := time.Since(s.LastSpinTime)
	// fmt.Printf("loop delay %d game %s \r\n", delay.Milliseconds(), s.CurrentSiXiangGame.String())
	// fmt.Println("sss")
	if delay > s.DurationTriggerAutoSpin {
		e.Process(s)
		s.DurationTriggerAutoSpin = durationAutoSpin + delayAnimationSpinGame
		return e.Finish(s)
	}
	return nil, nil
}

func (e *rapidPayEngine) Info(s interface{}) (interface{}, error) {
	return s, nil
}

// row 0: x4 END
// row 1:x3 x4 END
// row 2: x2 x3 x4 END
// row 3: x2 x3 X4 END
// row 4: x2 x2 x3 x3 x4
func (e *rapidPayEngine) symbolsAtRow(row int) []pb.SiXiangSymbol {
	switch row {
	case 0:
		return []pb.SiXiangSymbol{
			//
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END,
		}
	case 1:
		return []pb.SiXiangSymbol{
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END,
		}
	case 2:
		return []pb.SiXiangSymbol{
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END,
		}
	case 3:
		return []pb.SiXiangSymbol{
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_END,
		}
	case 4:
		return []pb.SiXiangSymbol{
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X2,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X3,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_RAPIDPAY_X4,
		}
	}
	return []pb.SiXiangSymbol{}
}
