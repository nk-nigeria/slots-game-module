package sixiang

import (
	"github.com/nk-nigeria/slots-game-module/entity"
	"github.com/nk-nigeria/cgp-common/lib"
	pb "github.com/nk-nigeria/cgp-common/proto"
)

var _ lib.Engine = &goldPickEngine{}

type goldPickEngine struct {
	randomIntFn         func(min, max int) int
	randomFloat64       func(min, max float64) float64
	ratioInSixiangBonus int
}

const (
	defaultGoldPickGemSpin = 20
)

func NewGoldPickEngine(ratioInSixiangBonus int, randomIntFn func(min, max int) int, randomFloat64 func(min, max float64) float64) lib.Engine {
	engine := goldPickEngine{
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

func (e *goldPickEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewMatrixGoldPick()
	matrixSpecial := entity.ShuffleMatrix(matrix)
	s.MatrixSpecial = &matrixSpecial
	s.SpinSymbols = []*pb.SpinSymbol{}
	s.NumSpinLeft = defaultGoldPickGemSpin
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
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
	spin.Index = int32(idRandom)
	s.SpinSymbols = []*pb.SpinSymbol{spin}
	_, spin.WinJp = entity.GoldPickSymbolToReward(spin.Symbol)

	s.SpinList[idRandom] = spin
	return s, nil
}

func (e *goldPickEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return s.LastResult, entity.ErrorSpinNotChange
	}
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{},
		ChipsMcb:   s.Bet().Chips,
	}
	s.IsSpinChange = false
	if s.NumSpinLeft <= 0 {
		slotDesk.IsFinishGame = true
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}

	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	for id := range s.MatrixSpecial.List {
		if !s.MatrixSpecial.IsFlip(id) {
			slotDesk.Matrix.Lists[id] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		}
	}
	for _, spin := range s.SpinSymbols {
		sym := spin.Symbol
		if sym >= pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MINOR {
			// s.AddCollectionSymbol(s.CurrentSiXiangGame, int(s.Bet().Chips), sym)
			_, s.WinJp = entity.GoldPickSymbolToReward(sym)
		}
		ratio, chips := e.calcRewardBySymbol(spin, s.Bet().Chips)
		spin.Ratio = ratio
		spin.WinAmount = chips
		s.SpinList[spin.Index].Ratio = spin.Ratio
		s.SpinList[spin.Index].WinAmount = chips
		slotDesk.GameReward.ChipsWin += chips
	}
	for _, spin := range s.SpinList {
		// _, chips := e.calcRewardBySymbol(spin, s.Bet().Chips)
		slotDesk.GameReward.TotalChipsWinByGame += spin.WinAmount
	}
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	// s.ChipStat.AddChipWin(s.CurrentSiXiangGame, slotDesk.GameReward.ChipsWin)
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	if slotDesk.IsFinishGame {
		s.AddGameEyePlayed(pb.SiXiangGame_SI_XIANG_GAME_GOLDPICK)
	}
	slotDesk.WinJp = s.WinJp
	slotDesk.Matrix.SpinLists = s.SpinList
	s.LastResult = slotDesk
	return slotDesk, nil
}

func (e *goldPickEngine) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *goldPickEngine) Info(s interface{}) (interface{}, error) {
	return s, nil
}

// return ratio, chip
func (e *goldPickEngine) calcRewardBySymbol(spin *pb.SpinSymbol, mcb int64) (float32, int64) {
	sym := spin.Symbol
	if sym == pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_TRYAGAIN {
		return 0, 0
	}
	ratio := float64(spin.Ratio)
	if ratio <= 0 {
		ratio = e.randomFloat64(
			float64(entity.ListSymbolGoldPick[sym].Value.Min),
			float64(entity.ListSymbolGoldPick[sym].Value.Max),
		)
	}
	ratio *= float64(e.ratioInSixiangBonus)
	chips := ratio * 100 * float64(mcb) / 100
	return float32(ratio), int64(chips)
}
