package sixiang

import (
	"errors"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &dragonPearlEngine{}

const (
	// default remain spin remain at start game
	defaultDragonPearlGemSpin = 3
	// add more spin when appear symbol eye bird (chu tuoc)
	bonusDragonPearlGemSpin = 3
)

type dragonPearlEngine struct {
	randomIntFn   func(min, max int) int
	randomFloat64 func(min, max float64) float64
	// ration bonus when draw DRAGONPEARL_GEM, default 1
	// will x2 when draw eye dragon
	ratioGem int
}

func NewDragonPearlEngine(randomIntFn func(min, max int) int, randomFloat64 func(min, max float64) float64) lib.Engine {
	engine := dragonPearlEngine{
		ratioGem: 1,
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

func (e *dragonPearlEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewMatrixDragonPearl()
	s.MatrixSpecial = entity.ShuffleMatrix(matrix)
	// s.ChipsWinInSpecialGame = 0
	s.SpinSymbols = []*pb.SpinSymbol{}
	// init spin list
	s.SpinList = make([]*pb.SpinSymbol, 0)
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		s.SpinList = append(s.SpinList, &pb.SpinSymbol{
			Symbol:    symbol,
			Row:       int32(row),
			Col:       int32(col),
			Index:     int32(idx),
			Ratio:     0,
			WinJp:     pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED,
			WinAmount: 0,
		})
	})
	{
		list := make([]pb.SiXiangSymbol, 0)
		for k := range entity.ListEyeSiXiang {
			list = append(list, k)
		}
		s.CollectionSymbolRemain = entity.ShuffleSlice(list)
	}

	s.NumSpinLeft = defaultDragonPearlGemSpin
	// s.CollectionSymbol = make(map[int]map[pb.SiXiangSymbol]int)
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	s.TurnSureSpin = e.randomIntFn(1, s.NumSpinLeft)
	s.ChipStat.Reset(s.CurrentSiXiangGame)
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
	if s.NumSpinLeft == s.TurnSureSpin && s.SizeCollectionSymbol(s.CurrentSiXiangGame, int(s.Bet().Chips)) == 0 {
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
			Index:  int32(idxRandom),
		}
		s.AddCollectionSymbol(s.CurrentSiXiangGame, int(s.Bet().GetChips()), eyeRandom)
		s.SpinSymbols = []*pb.SpinSymbol{spinSymbol}
		s.NumSpinLeft--
	} else {
		e.randomPearl(s, func(symbolRand, eyeRand pb.SiXiangSymbol, row, col int) bool {
			spinSymbol := &pb.SpinSymbol{
				Symbol: symbolRand,
				Row:    int32(row),
				Col:    int32(col),
				Index:  int32(row*s.Matrix.Cols + col),
			}
			if eyeRand != pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
				spinSymbol.Symbol = eyeRand
			}
			s.SpinSymbols = []*pb.SpinSymbol{spinSymbol}
			s.NumSpinLeft--
			return true
		})
	}
	s.SpinList[s.SpinSymbols[0].Index] = s.SpinSymbols[0]
	switch s.SpinSymbols[0].Symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_BIRD:
		// add spin
		s.NumSpinLeft += bonusDragonPearlGemSpin
		s.AddCollectionSymbol(s.CurrentSiXiangGame, int(s.Bet().GetChips()), pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_BIRD)
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER:
		// x2 money in gem
		s.NumSpinLeft += 1
		e.ratioGem *= 2
		s.AddCollectionSymbol(s.CurrentSiXiangGame, int(s.Bet().GetChips()), pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER)
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
						index := row*s.Matrix.Cols + col
						spinSymbol := &pb.SpinSymbol{
							Symbol: symbolRand,
							Row:    int32(row),
							Col:    int32(col),
							Index:  int32(index),
						}
						s.SpinSymbols = append(s.SpinSymbols, spinSymbol)
						s.SpinList[index] = spinSymbol
						return true
					}
					return false
				})
				if valid {
					break
				}
			}
		}
		s.AddCollectionSymbol(
			s.CurrentSiXiangGame,
			int(s.Bet().GetChips()),
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_WARRIOR,
		)
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
		index := s.SpinSymbols[0].Index
		spinSymbol := &pb.SpinSymbol{
			Symbol: randomJp,
			Row:    s.SpinSymbols[0].Row,
			Col:    s.SpinSymbols[0].Col,
			Index:  index,
		}
		s.SpinSymbols = append(s.SpinSymbols, spinSymbol)
		s.SpinList[index].WinJp = pb.WinJackpot(randomJp)
		s.AddCollectionSymbol(
			s.CurrentSiXiangGame,
			int(s.Bet().GetChips()),
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_DRAGON,
		)
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
	ratioWin := float32(0)
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if s.MatrixSpecial.TrackFlip[idx] {
			// v := entity.ListSymbolDragonPearl[symbol].Value
			// ratioWin += e.randomFloat64(float64(v.Min), float64(v.Max))
		} else {
			slotDesk.Matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		}
	})
	for _, sym := range s.SpinSymbols {
		v := entity.ListSymbolDragonPearl[sym.Symbol].Value
		sym.Ratio = float32(e.randomFloat64(float64(v.Min), float64(v.Max)))
		ratioWin += sym.Ratio
		s.SpinSymbols[sym.Index].Ratio = sym.Ratio
		s.SpinSymbols[sym.Index].WinAmount = int64(float64(float64(ratioWin) * float64(s.Bet().Chips)))
	}
	ratioWin *= float32(e.ratioGem)

	ratioJPBonus := float32(1)
	for _, eyeSym := range s.CollectionSymbolToSlice(s.CurrentSiXiangGame, int(s.Bet().Chips)) {
		r := entity.ListEyeSiXiang[eyeSym.Symbol].Value.Min
		if (r) > ratioJPBonus {
			ratioJPBonus = float32(r)
		}
	}
	slotDesk.GameReward.ChipsWin = int64(float64(ratioJPBonus) * float64(float64(ratioWin)*float64(s.Bet().Chips)))
	// totalChipsWin := float64(ratioJPBonus) * float64(ratioWin*float64(s.Bet().Chips))
	// slotDesk.GameReward.ChipsWin = int64(totalChipsWin)
	if s.WinJp != pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED {
		slotDesk.GameReward.ChipsWin = s.Bet().Chips * int64(s.WinJp)
	}

	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.Matrix.SpinLists = s.SpinList
	slotDesk.ChipsMcb = s.Bet().Chips
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	// s.ChipStat.ResetChipWin(s.CurrentSiXiangGame)
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, int64(slotDesk.GameReward.ChipsWin))
	slotDesk.GameReward.TotalChipsWinByGame = s.ChipStat.TotalChipWin(s.CurrentSiXiangGame)
	slotDesk.GameReward.RatioWin = float32(ratioWin)
	slotDesk.CollectionSymbols = s.CollectionSymbolToSlice(s.CurrentSiXiangGame, int(s.Bet().Chips))
	return slotDesk, nil
}

func (e *dragonPearlEngine) Loop(s interface{}) (interface{}, error) {
	return s, nil
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
			s.AddCollectionSymbol(s.CurrentSiXiangGame, int(s.Bet().Chips), eyeRandom)
		}

		// s.MatrixSpecial.TrackFlip[idRandom] = true
		s.MatrixSpecial.Flip(idRandom)
	}
	return acceptSymbol
}
