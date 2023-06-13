package sixiang

import (
	"errors"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &dragonPearlEngine{}

const (
	defaultDragonPearlGemSpin = 3
	bonusDragonPearlGemSpin   = 3
)

type dragonPearlEngine struct {
	randomIntFn   func(min, max int) int
	randomFloat64 func(min, max float64) float64
}

func NewDragonPearlEngine(randomIntFn func(min, max int) int, randomFloat64 func(min, max float64) float64) lib.Engine {
	engine := dragonPearlEngine{}
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

func (e *dragonPearlEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewMatrixDragonPearl()
	s.MatrixSpecial = entity.ShuffleMatrix(matrix)
	// s.ChipsWinInSpecialGame = 0
	s.SpinSymbols = []*pb.SpinSymbol{}
	{
		list := make([]pb.SiXiangSymbol, 0)
		for k := range entity.ListEyeSiXiang {
			list = append(list, k)
		}
		s.CollectionSymbolRemain = entity.ShuffleSlice(list)
	}

	s.NumSpinLeft = defaultDragonPearlGemSpin
	s.CollectionSymbol = make(map[int]map[pb.SiXiangSymbol]int)
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	s.TurnSureSpin = e.randomIntFn(1, s.NumSpinLeft)
	s.ChipStat.ResetChipWin(s.CurrentSiXiangGame)
	return s, nil
}

func (e *dragonPearlEngine) Random(min, max int) int {
	return e.randomIntFn(min, max)
}

func (e *dragonPearlEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return s, errors.New("gem spin not enough")
	}
	if len(s.MatrixSpecial.TrackFlip) >= 15 {
		return s, entity.ErrorSpinReachMax
	}
	s.IsSpinChange = true
	// Setup sao cho số lượt spins của user ít nhất được 8 ngọc và 1 phong bao
	// nên đầu game random ra lân quay chắc chắn sẽ ra ngọc nếu tới lượt đó
	// nhưng chưaquay ra ngọc
	if s.NumSpinLeft == s.TurnSureSpin && len(s.CollectionSymbol[int(s.Bet().Chips)]) == 0 {
		listIdEyeSymbol := make([]int, 0)
		s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_LUCKMONEY && !s.MatrixSpecial.TrackFlip[idx] {
				listIdEyeSymbol = append(listIdEyeSymbol, idx)
			}
		})
		eyeRandom := entity.ShuffleSlice(s.CollectionSymbolRemain)[0]
		idxRandom := entity.ShuffleSlice(listIdEyeSymbol)[0]
		row, col := s.MatrixSpecial.RowCol(idxRandom)
		spinSymbol := &pb.SpinSymbol{
			Symbol: eyeRandom,
			Row:    int32(row),
			Col:    int32(col),
		}
		s.AddCollectionSymbol(int(s.Bet().GetChips()), eyeRandom)
		s.SpinSymbols = []*pb.SpinSymbol{spinSymbol}
		s.NumSpinLeft--
	} else {
		e.randomPearl(s, func(symbolRand, eyeRand pb.SiXiangSymbol, row, col int) bool {
			spinSymbol := &pb.SpinSymbol{
				Symbol: symbolRand,
				Row:    int32(row),
				Col:    int32(col),
			}
			if eyeRand != pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
				spinSymbol.Symbol = eyeRand
			}
			s.SpinSymbols = []*pb.SpinSymbol{spinSymbol}
			s.NumSpinLeft--
			return true
		})
	}
	switch s.SpinSymbols[0].Symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_BIRD:
		// add spin
		s.NumSpinLeft += bonusDragonPearlGemSpin
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER:
		// x2 money in gem
		s.NumSpinLeft += 1
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_WARRIOR:
		s.NumSpinLeft += 1
		// add 3 gen money
		for i := 0; i < 3; i++ {
			for {
				// todo if gen money not enough, not add more gem money
				// if len(s.MatrixSpecial.List)-len(s.MatrixSpecial.TrackFlip) == len(s.EyeSiXiangRemain) {
				// 	break
				// }
				valid := e.randomPearl(s, func(symbolRand, eyeRand pb.SiXiangSymbol, row, col int) bool {
					if eyeRand == pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
						spinSymbol := &pb.SpinSymbol{
							Symbol: symbolRand,
							Row:    int32(row),
							Col:    int32(col),
						}
						s.SpinSymbols = append(s.SpinSymbols, spinSymbol)
						return true
					}
					return false
				})
				if valid {
					break
				}
			}
		}
		// todo
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_DRAGON:
		s.NumSpinLeft += 1
		// roll jackpot
		listSymbolJP := []pb.SiXiangSymbol{
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MINOR,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MAJOR,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MEGA,
		}
		randomJp := entity.ShuffleSlice(listSymbolJP)[e.randomIntFn(0, len(listSymbolJP))]
		spinSymbol := &pb.SpinSymbol{
			Symbol: randomJp,
			Row:    s.SpinSymbols[0].Row,
			Col:    s.SpinSymbols[0].Col,
		}
		s.SpinSymbols = append(s.SpinSymbols, spinSymbol)
	}
	return s, nil
}

func (e *dragonPearlEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{},
	}
	if !s.IsSpinChange {
		return slotDesk, entity.ErrorSpinNotChange
	}
	s.IsSpinChange = false
	if s.NumSpinLeft <= 0 || len(s.MatrixSpecial.TrackFlip) == 15 {
		slotDesk.IsFinishGame = true
	}
	if slotDesk.IsFinishGame {
		if len(s.MatrixSpecial.List) == 15 {
			slotDesk.WinJp = pb.WinJackpot_WIN_JACKPOT_GRAND
		}
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	totalMcb := float64(0)
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if s.MatrixSpecial.TrackFlip[idx] {
			v := entity.ListSymbolDragonPearl[symbol].Value
			totalMcb += e.randomFloat64(float64(v.Min), float64(v.Max))
		} else {
			slotDesk.Matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		}
	})

	//todo jacpot reward
	// for _, spin := range s.SpinSymbols {
	// 	v := entity.ListSymbolDragonPearl[spin.Symbol].Value
	// 	totalMcb += float64(e.randomFloat64(float64(v.Min), float64(v.Max)))
	// }
	ratioBonus := float64(1)
	for _, eyeSym := range s.CollectionSymbolToSlice(int(s.Bet().Chips)) {
		r := entity.ListEyeSiXiang[eyeSym].Value.Min
		if float64(r) > ratioBonus {
			ratioBonus = float64(r)
		}
	}
	chips := ratioBonus * float64(totalMcb*float64(s.Bet().Chips))
	slotDesk.GameReward.ChipsWin = int64(chips)
	if s.WinJp != pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED {
		slotDesk.GameReward.ChipsWin = s.Bet().Chips * int64(s.WinJp)
	}

	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.ChipsMcb = s.Bet().Chips
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, slotDesk.GameReward.ChipsWin)
	slotDesk.GameReward.TotalChipsWinByGame = s.ChipStat.TotalChipWin(s.CurrentSiXiangGame)
	return slotDesk, nil
}

// func (e *dragonPearlEngine) checkJackpot(s *entity.SlotsMatchState) bool {
// 	return len(s.MatrixSpecial.TrackFlip) >= 15
// }

func (e *dragonPearlEngine) randomPearl(
	s *entity.SlotsMatchState,
	fn func(symbolRand, eyeRand pb.SiXiangSymbol, row, col int) bool,
) bool {
	idRandom, symbolRandom := s.MatrixSpecial.RandomSymbolNotFlip(e.randomIntFn)
	eyeRandom := pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
	if symbolRandom == pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_LUCKMONEY {
		eyeRandom = s.CollectionSymbolRemain[0]
	}
	row, col := s.MatrixSpecial.RowCol(idRandom)
	acceptSymbol := fn(symbolRandom, eyeRandom, row, col)
	if acceptSymbol {
		if eyeRandom != pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
			s.CollectionSymbolRemain = s.CollectionSymbolRemain[1:]
			s.AddCollectionSymbol(int(s.Bet().Chips), eyeRandom)
		}

		// s.MatrixSpecial.TrackFlip[idRandom] = true
		s.MatrixSpecial.Flip(idRandom)
	}
	return acceptSymbol
}
