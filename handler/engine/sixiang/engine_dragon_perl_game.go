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
	//
	ratioInSixiangBonus int
}

func NewDragonPearlEngine(ratioInSixiangBonus int, randomIntFn func(min, max int) int, randomFloat64 func(min, max float64) float64) lib.Engine {
	engine := dragonPearlEngine{
		ratioGem:            0,
		ratioInSixiangBonus: 1,
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
	matrixSpecial := entity.ShuffleMatrix(matrix)
	s.MatrixSpecial = &matrixSpecial
	// s.ChipsWinInSpecialGame = 0
	s.SpinSymbols = []*pb.SpinSymbol{}
	// init spin list
	s.SpinList = make([]*pb.SpinSymbol, 0)
	var eyeSymbols []pb.SiXiangSymbol
	{
		list := make([]pb.SiXiangSymbol, 0)
		for k := range entity.ListEyeSiXiang {
			list = append(list, k)
		}
		eyeSymbols = entity.ShuffleSlice(list)
	}
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
		if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_LUCKMONEY {
			s.MatrixSpecial.List[idx] = eyeSymbols[0]
			eyeSymbols = eyeSymbols[1:]
		}
	})
	// {
	// 	list := make([]pb.SiXiangSymbol, 0)
	// 	for k := range entity.ListEyeSiXiang {
	// 		list = append(list, k)
	// 	}
	// 	s.EyeSymbolRemains = entity.ShuffleSlice(list)
	// }

	s.NumSpinLeft = defaultDragonPearlGemSpin
	// s.CollectionSymbol = make(map[int]map[pb.SiXiangSymbol]int)
	s.WinJp = pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED
	s.NumSpinLeft += 6
	s.ChipStat.Reset(s.CurrentSiXiangGame)
	// Rồng thần nhả 6 ngọc có sẵn cho user.
	s.NotDropEyeSymbol = true
	var res interface{}
	for i := 0; i < 6; i++ {
		_, _ = e.Process(s)
		res, _ = e.Finish(s)
		// fmt.Println(s.MatrixSpecial.TrackFlip)
	}
	s.LastResult = res.(*pb.SlotDesk)
	s.NotDropEyeSymbol = false
	s.NumSpinLeft = 3
	s.TurnSureSpinSpecial = e.randomIntFn(1, 3)
	// s.TurnSureSpinEye = 3
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
	// nên đầu game random ra lân quay chắc chắn sẽ ra eye nếu tới lượt đó
	// nhưng chưa quay ra eye
	var spinSymbol *pb.SpinSymbol
	if s.NumSpinLeft == s.TurnSureSpinSpecial && !s.NotDropEyeSymbol {
		s.MatrixSpecial.ForEeachNotFlip(func(idx, row, col int, symbol pb.SiXiangSymbol) {
			if entity.IsSixiangEyeSymbol(symbol) && spinSymbol == nil {
				// eyeSymbolSpin = symbol
				spinSymbol = &pb.SpinSymbol{
					Symbol: symbol,
					Row:    int32(row),
					Col:    int32(col),
					Index:  int32(idx),
				}
			}
		})
		s.TurnSureSpinSpecial = -1
	} else {
		// normal spin
		for {
			idxRandom, symbol := s.MatrixSpecial.RandomSymbolNotFlip(e.randomIntFn)
			if entity.IsSixiangEyeSymbol(symbol) {
				if s.NotDropEyeSymbol {
					continue
				}
				s.TurnSureSpinSpecial = 0
				// pay in turn 100% spin eye
				if s.TurnSureSpinSpecial < 0 {
					continue
				}
				s.TurnSureSpinSpecial = 0
			}
			// check condition when spin eye warrior
			// Mắt huyền vũ: làm rơi thêm random 3 viên ngọc tiền lên bảng (không nhả ngọc JP)
			if symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_WARRIOR {
				numGemNotFlip := s.MatrixSpecial.CountSymbolCond(func(idx int, symbol pb.SiXiangSymbol) bool {
					return !entity.IsSixiangEyeSymbol(symbol) && !s.MatrixSpecial.IsFlip(idx)
				})
				if numGemNotFlip <= 3 {
					// change warrior to another eye symbol
					mapEye := make(map[pb.SiXiangSymbol]int)
					for k := range entity.ListEyeSiXiang {
						mapEye[k] = 1
					}
					s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
						mapEye[symbol] = 0
					})
					for eyeSymbol, inList := range mapEye {
						if inList == 1 {
							s.MatrixSpecial.List[idxRandom] = eyeSymbol
							symbol = eyeSymbol
							break
						}
					}
				}
			}
			row, col := s.MatrixSpecial.RowCol(idxRandom)
			spinSymbol = &pb.SpinSymbol{
				Symbol: symbol,
				Row:    int32(row),
				Col:    int32(col),
				Index:  int32(idxRandom),
				// Ratio:  1.0,
			}
			break
		}
	}
	if spinSymbol == nil {
		return s, entity.ErrorInternal
	}
	// cheat alway drop eye tiger
	// symbolSpec := pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER
	// if entity.IsSixiangEyeSymbol(spinSymbol.Symbol) &&
	// 	spinSymbol.Symbol != symbolSpec {
	// 	// check index eye
	// 	idxEyeSpec := -1
	// 	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
	// 		if symbol == symbolSpec {
	// 			idxEyeSpec = idx
	// 		}
	// 	})
	// 	if idxEyeSpec == -1 {
	// 		// in matrix init not contain eye tiger
	// 		spinSymbol.Symbol = symbolSpec
	// 		s.MatrixSpecial.List[spinSymbol.Index] = spinSymbol.Symbol
	// 	} else if !s.MatrixSpecial.IsFlip(idxEyeSpec) {
	// 		// swap to another eye
	// 		s.MatrixSpecial.List[idxEyeSpec] = spinSymbol.Symbol
	// 		spinSymbol.Symbol = symbolSpec
	// 		s.MatrixSpecial.List[spinSymbol.Index] = spinSymbol.Symbol
	// 	}
	// }
	s.MatrixSpecial.Flip(int(spinSymbol.GetIndex()))
	s.SpinSymbols = []*pb.SpinSymbol{spinSymbol}
	s.NumSpinLeft--

	s.SpinList[s.SpinSymbols[0].Index] = s.SpinSymbols[0]
	switch s.SpinSymbols[0].Symbol {
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_BIRD:
		// add spin
		s.NumSpinLeft += bonusDragonPearlGemSpin
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER:
		// x2 money in gem
		s.NumSpinLeft += 1
		// s.SpinSymbols[0].Ratio = 2.0
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_WARRIOR:
		s.NumSpinLeft += 1
		// add 3 gen money
		for {
			idx, symbol := s.MatrixSpecial.RandomSymbolNotFlip(e.randomIntFn)
			if entity.IsSixiangEyeSymbol(symbol) {
				continue
			}
			if idx < 0 {
				break
			}
			if len(s.SpinSymbols) >= 4 {
				break
			}
			row, col := s.MatrixSpecial.RowCol(idx)
			newSpin := &pb.SpinSymbol{
				Symbol: symbol,
				Row:    int32(row),
				Col:    int32(col),
				Index:  int32(idx),
			}
			s.MatrixSpecial.Flip(int(newSpin.GetIndex()))
			s.SpinSymbols = append(s.SpinSymbols, newSpin)
		}
		for _, spin := range s.SpinSymbols {
			s.MatrixSpecial.Flip(int(spin.GetIndex()))
			s.SpinList[spin.Index] = spin
		}
	case pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_DRAGON:
		s.NumSpinLeft += 1
		// roll jackpot
		listSymbolJP := []pb.SiXiangSymbol{
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MINOR,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MAJOR,
			pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_JP_MEGA,
		}
		randomJp := entity.ShuffleSlice(listSymbolJP)[e.randomIntFn(0, len(listSymbolJP))]
		spinSymbol := s.SpinSymbols[0]
		_, spinSymbol.WinJp = entity.DragonPearlSymbolToReward(randomJp)
		s.SpinList[spinSymbol.GetIndex()].WinJp = spinSymbol.WinJp
	default:
		s.NumSpinLeft++
	}
	for _, sym := range s.SpinSymbols {
		v := entity.ListSymbolDragonPearl[sym.Symbol].Value
		ratio := e.randomFloat64(float64(v.Min), float64(v.Max))
		sym.Ratio = float32(ratio)
		sym.RatioBonus = 0
		s.SpinList[sym.GetIndex()].Ratio = sym.Ratio
		s.SpinList[sym.GetIndex()].RatioBonus = float32(e.ratioGem)
	}
	return s, nil
}

/*
Cách tính tiền khi xuất hiện cả mắt bạch hổ(x2 giá trị ngọc tiền) và thanh long(random 1 ngọc) trong trường hợp si xiang bonus (x4)
Mắt bạch hổ:
x2 đối với ngọc tiền
chỉ tính tiền tổng bên dưới, không x2 giá trị ngọc tiền hiển thị trên bảng
không x2 giá trị ngọc random mở ra từ mắt thanh long
+ Tính tiền cuối = (tiền thắng được trong game) x4
*/
func (e *dragonPearlEngine) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{},
	}
	if !s.IsSpinChange {
		return s.LastResult, entity.ErrorSpinNotChange
	}
	s.IsSpinChange = false
	if s.NumSpinLeft <= 0 || len(s.MatrixSpecial.TrackFlip) == 15 {
		slotDesk.IsFinishGame = true
	}
	if slotDesk.IsFinishGame {
		if len(s.MatrixSpecial.TrackFlip) == 15 {
			s.WinJp = pb.WinJackpot_WIN_JACKPOT_GRAND
		}
		s.NextSiXiangGame = pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	slotDesk.Matrix = s.MatrixSpecial.ToPbSlotMatrix()
	s.MatrixSpecial.ForEeachNotFlip(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		slotDesk.Matrix.Lists[idx] = pb.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
	})

	if s.SpinSymbols[0].Symbol == pb.SiXiangSymbol_SI_XIANG_SYMBOL_DRAGONPEARL_EYE_TIGER {
		e.ratioGem = 1
		// update old result
		for _, sym := range s.SpinList {
			sym.RatioBonus = float32(e.ratioGem)
		}
		// không x2 giá trị ngọc random mở ra từ mắt thanh long
		// for _, sym := range s.SpinSymbols {
		// 	s.SpinList[sym.Index].RatioBonus = 0
		// }
	}

	for _, sym := range s.SpinSymbols {
		// sym.Ratio, sym.WinAmount = e.calcRewardBySymbol(sym, s.Bet().Chips)
		sym.WinAmount = int64(float64(sym.Ratio)*100*float64(s.Bet().Chips)) / 100
		sym.WinAmount += int64(sym.GetWinJp()) * s.Bet().Chips
		slotDesk.GameReward.ChipsWin += sym.WinAmount
		slotDesk.GameReward.RatioWin += sym.Ratio
		s.SpinList[sym.Index].WinAmount = sym.WinAmount
	}
	for _, sym := range s.SpinList {
		slotDesk.GameReward.TotalRatioWin += sym.Ratio + sym.RatioBonus
		chipsWin := int64(float64(sym.Ratio)*100*float64(s.Bet().Chips)) / 100
		slotDesk.GameReward.TotalChipsWinByGame += chipsWin
		// add more chip by win EYE_TIGER
		for i := 0; i < e.ratioGem; i++ {
			slotDesk.GameReward.TotalChipsWinByGame += chipsWin
		}
		slotDesk.GameReward.TotalChipsWinByGame += int64(sym.GetWinJp()) * s.Bet().Chips
	}
	if s.WinJp != pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED {
		chipsJp := s.Bet().Chips * int64(s.WinJp)
		slotDesk.GameReward.ChipsWin += chipsJp
		slotDesk.GameReward.TotalChipsWinByGame += chipsJp
	}
	slotDesk.WinJp = s.WinJp
	slotDesk.CurrentSixiangGame = s.CurrentSiXiangGame
	slotDesk.NextSixiangGame = s.NextSiXiangGame
	slotDesk.SpinSymbols = s.SpinSymbols
	slotDesk.Matrix.SpinLists = s.SpinList
	slotDesk.ChipsMcb = s.Bet().Chips
	slotDesk.NumSpinLeft = int64(s.NumSpinLeft)
	if slotDesk.IsFinishGame {
		s.AddGameEyePlayed(pb.SiXiangGame_SI_XIANG_GAME_DRAGON_PEARL)
	}
	s.LastResult = slotDesk
	return slotDesk, nil
}

func (e *dragonPearlEngine) Loop(s interface{}) (interface{}, error) {
	return s, nil
}

func (e *dragonPearlEngine) Info(s interface{}) (interface{}, error) {
	return s, nil
}

// return ratio, chip
// func (e *dragonPearlEngine) calcRewardBySymbol(sym *pb.SpinSymbol, mcb int64) (float32, int64) {
// 	symbol := sym.Symbol
// 	v := entity.ListSymbolDragonPearl[symbol].Value
// 	ratio := float64(sym.Ratio)
// 	if ratio <= 0 {
// 		ratio = e.randomFloat64(float64(v.Min), float64(v.Max))
// 	}
// 	ratio *= float64(e.ratioGem)
// 	winJpRatio := sym.GetWinJp()
// 	if winJpRatio != pb.WinJackpot_WIN_JACKPOT_UNSPECIFIED {
// 		ratio += float64(winJpRatio)
// 	}
// 	ratio *= float64(e.ratioInSixiangBonus)
// 	chips := ratio * 100 * float64(mcb) / 100
// 	return float32(ratio), int64(chips)
// }
