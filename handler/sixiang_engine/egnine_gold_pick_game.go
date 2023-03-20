package sixiangengine

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
		engine.randomIntFn = RandomInt
	}
	if randomFloat64 != nil {
		engine.randomFloat64 = randomFloat64
	} else {
		engine.randomFloat64 = RandomFloat64
	}
	return &engine
}

func (e *goldPickEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	matrix := entity.NewMatrixGoldPick()
	s.MatrixSpecial = ShuffleMatrix(matrix)
	s.SpinSymbols = []*pb.SpinSymbol{}
	s.GemSpin = defaultGoldPickGemSpin
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	return s, nil
}

func (e *goldPickEngine) Random(min, max int) int {
	return e.randomIntFn(min, max)
}

func (e *goldPickEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SixiangMatchState)
	if s.GemSpin <= 0 {
		return s, ErrorSpinReadMax
	}
	idRandom, symbolRandom := s.MatrixSpecial.RandomSymbolNotFlip(e.randomIntFn)
	row, col := s.MatrixSpecial.RowCol(idRandom)
	spin := &pb.SpinSymbol{
		Symbol: symbolRandom,
		Row:    int32(row),
		Col:    int32(col),
	}
	s.GemSpin--
	if symbolRandom == pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD5 {
		arr := []pb.SiXiangSymbol{
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_GOLD5,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MINOR,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MAJOR,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_GOLD_PICK_JP_MEGA,
		}
		symbolRandom = ShuffleSlice(arr)[e.randomIntFn(0, len(arr))]
		spin.Symbol = symbolRandom
	}
	s.MatrixSpecial.Flip(idRandom)
	s.MatrixSpecial.List[idRandom] = symbolRandom
	s.SpinSymbols = []*pb.SpinSymbol{spin}
	return s, nil
}

func (e *goldPickEngine) Finish(matchState interface{}) (interface{}, error) {
	slotDesk := &pb.SlotDesk{}
	s := matchState.(*entity.SixiangMatchState)
	if s.GemSpin <= 0 {
		slotDesk.IsFinishGame = true
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	ratio := float64(0)
	for id, sym := range s.MatrixSpecial.List {
		if s.MatrixSpecial.TrackFlip[id] {
			value := entity.ListSymbolGoldPick[sym].Value
			ratio += e.randomFloat64(float64(value.Min), float64(value.Max))
		} else {
			slotDesk.Matrix.Lists[id] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		}
	}
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.ChipsMcb = s.GetBetInfo().GetChips()
	slotDesk.ChipsWin = int64(ratio * float64(slotDesk.ChipsMcb))
	return slotDesk, nil
}
