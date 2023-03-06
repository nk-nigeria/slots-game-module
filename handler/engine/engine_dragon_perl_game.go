package engine

import (
	"errors"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &dragonPearlEngine{}

const (
	defaultGemSpin = 3
	bonusGemSpin   = 3
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
		engine.randomIntFn = RandomInt
	}
	if randomFloat64 != nil {
		engine.randomFloat64 = randomFloat64
	} else {
		engine.randomFloat64 = RandomFloat64
	}
	return &engine
}

func (e *dragonPearlEngine) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	matrix := entity.NewSiXiangMatrixDragonPearl()
	s.MatrixSpecial = ShuffleMatrix(matrix)
	// s.ChipsWinInSpecialGame = 0
	s.SpinSymbols = []*pb.SpinSymbol{
		{Symbol: pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED},
	}
	s.EyeSiXiangRemain = ShuffleSlice(entity.ListEyeSiXiang[:])
	s.GemSpin = defaultGemSpin
	s.EyeSiXiangSpined = make([]pb.SiXiangSymbol, 0)
	s.RatioBonus = 1
	return s, nil
}

func (e *dragonPearlEngine) Random(min, max int) int {
	return e.randomIntFn(min, max)
}

func (e *dragonPearlEngine) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.GemSpin <= 0 {
		return s, errors.New("gem spin not enough")
	}
	if len(s.MatrixSpecial.TrackFlip) == 15 {
		return s, errors.New("Spin all")
	}
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
		s.GemSpin--
		return true
	})
	switch s.SpinSymbols[0].Symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_BIRD:
		// add spin
		s.GemSpin += bonusGemSpin
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER:
		// x2 money in gem
		s.RatioBonus = 2
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_WARRIOR:
		// add more than gen money
		for i := 0; i < 3; i++ {
			for {
				// todo if gen money not enough, not add more gem money
				if len(s.MatrixSpecial.List)-len(s.MatrixSpecial.TrackFlip) == len(s.EyeSiXiangRemain) {
					break
				}
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
		// add jp pearl
		// listSymbolJP := []pb.SiXiangSymbol{
		// 	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MINOR,
		// 	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MAJOR,
		// 	pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MEGA,
		// }
		// randomJp := ShuffleSlice(listSymbolJP)[0]
		// s.SpinSymbols = append(s.SpinSymbols, &pb.SpinSymbol{
		// 	Symbol: randomJp,
		// })
	}
	return s, nil
}

func (e *dragonPearlEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotDesk := &pb.SlotDesk{}
	if s.GemSpin <= 0 || len(s.MatrixSpecial.TrackFlip) == 15 {
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
		if s.MatrixSpecial.TrackFlip[idx] == false {
			slotDesk.Matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
		} else {
			v := entity.ListSymbolDragonPearl[symbol].Value
			totalMcb += float64(e.randomFloat64(float64(v.Min), float64(v.Max)))
		}
	})

	//todo jacpot reward
	for _, spin := range s.SpinSymbols {
		v := entity.ListSymbolDragonPearl[spin.Symbol].Value
		totalMcb += float64(e.randomFloat64(float64(v.Min), float64(v.Max)))
	}
	chips := float64(s.RatioBonus) * float64(totalMcb*float64(s.GetBetInfo().Chips))
	slotDesk.ChipsWin = int64(chips)

	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.ChipsMcb = s.GetBetInfo().Chips
	return slotDesk, nil
}

func (e *dragonPearlEngine) checkJackpot(s *entity.SlotsMatchState) bool {
	if len(s.MatrixSpecial.TrackFlip) >= 15 {
		return true
	}
	return false
}

func (e *dragonPearlEngine) randomPearl(
	s *entity.SlotsMatchState,
	fn func(symbolRand, eyeRand pb.SiXiangSymbol, row, col int) bool,
) bool {
	idRandom, symbolRandom := s.MatrixSpecial.RandomSymbolNotFlip(e.randomIntFn)
	eyeRandom := pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
	if symbolRandom == pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_LUCKMONEY {
		eyeRandom = s.EyeSiXiangRemain[0]
	}
	row, col := s.MatrixSpecial.RowCol(idRandom)
	acceptSymbol := fn(symbolRandom, eyeRandom, row, col)
	if acceptSymbol {
		if eyeRandom != pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED {
			s.EyeSiXiangRemain = s.EyeSiXiangRemain[1:]
			s.EyeSiXiangSpined = append(s.EyeSiXiangSpined, symbolRandom)
		}
		s.MatrixSpecial.TrackFlip[idRandom] = true
	}
	return acceptSymbol
}
