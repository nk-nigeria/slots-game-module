package sixiang

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &goldPickEngine{}

type goldPickEngine struct {
	randomIntFn   func(min, max int) int
	randomFloat64 func(min, max float64) float64
}

const (
	defaultGoldPickGemSpin = 20
)

func NewGoldPickEngine(randomIntFn func(min, max int) int, randomFloat64 func(min, max float64) float64) lib.Engine {
	engine := goldPickEngine{}
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

func (e *goldPickEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewMatrixGoldPick()
	s.MatrixSpecial = entity.ShuffleMatrix(matrix)
	s.SpinSymbols = []*pb.SpinSymbol{}
	s.NumSpinLeft = defaultGoldPickGemSpin
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	s.ChipStat.ResetChipWin(s.CurrentSiXiangGame)
	s.ResetCollection(s.CurrentSiXiangGame, int(s.Bet().Chips))
	return s, nil
}

func (e *goldPickEngine) Random(min, max int) int {
	return e.randomIntFn(min, max)
}

func (e *goldPickEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return s, entity.ErrorSpinReachMax
	}
	s.IsSpinChange = true
	idRandom, symbolRandom := s.MatrixSpecial.RandomSymbolNotFlip(e.randomIntFn)
	row, col := s.MatrixSpecial.RowCol(idRandom)
	spin := &pb.SpinSymbol{
		Symbol: symbolRandom,
		Row:    int32(row),
		Col:    int32(col),
	}
	s.NumSpinLeft--
	if symbolRandom == pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD5 {
		arr := []pb.SiXiangSymbol{
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD5,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MINOR,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MAJOR,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MEGA,
		}
		symbolRandom = entity.ShuffleSlice(arr)[e.randomIntFn(0, len(arr))]
		spin.Symbol = symbolRandom
	}
	s.MatrixSpecial.Flip(idRandom)
	s.MatrixSpecial.List[idRandom] = symbolRandom
	s.SpinSymbols = []*pb.SpinSymbol{spin}
	return s, nil
}

func (e *goldPickEngine) Finish(matchState interface{}) (interface{}, error) {
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{},
	}
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return slotDesk, entity.ErrorSpinNotChange
	}

	s.IsSpinChange = false
	if s.NumSpinLeft <= 0 {
		slotDesk.IsFinishGame = true
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	// for id, sym := range s.MatrixSpecial.List {
	// 	if s.MatrixSpecial.TrackFlip[id] {
	// 		value := entity.ListSymbolGoldPick[sym].Value
	// 		ratio += e.randomFloat64(float64(value.Min), float64(value.Max))
	// 	} else {
	// 		slotDesk.Matrix.Lists[id] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
	// 	}
	// }
	ratioWin := float32(0)
	for _, spin := range s.SpinSymbols {
		sym := spin.Symbol
		if sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_TRYAGAIN {
			continue
		}
		if sym >= pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MINOR {
			s.AddCollectionSymbol(s.CurrentSiXiangGame, int(s.Bet().Chips), sym)
		}
		spin.Ratio = float32(e.randomFloat64(
			float64(entity.ListSymbolGoldPick[sym].Value.Min),
			float64(entity.ListSymbolGoldPick[sym].Value.Max),
		))
		ratioWin += spin.Ratio
	}
	slotDesk.ChipsMcb = s.Bet().GetChips()
	slotDesk.GameReward.ChipsWin = int64(float64(ratioWin) * float64(slotDesk.ChipsMcb))
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, slotDesk.GameReward.ChipsWin)
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	slotDesk.GameReward.TotalChipsWinByGame = s.ChipStat.TotalChipWin(s.CurrentSiXiangGame)
	slotDesk.CollectionSymbols = s.CollectionSymbolToSlice(s.CurrentSiXiangGame, int(s.Bet().Chips))
	slotDesk.GameReward.RatioWin = ratioWin
	return slotDesk, nil
}

func (e *goldPickEngine) Loop(s interface{}) (interface{}, error) {
	return s, nil
}
