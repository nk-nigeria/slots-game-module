package juicy

import (
	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	"github.com/ciaolink-game-platform/cgp-common/lib"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
)

var _ lib.Engine = &freeGame{}

type freeGame struct {
	ratioWild ratioWild
	normal
}

func NewFreeGame(randomIntFn func(int, int) int) lib.Engine {
	e := &freeGame{}
	if randomIntFn != nil {
		e.randomFn = randomIntFn
	} else {
		e.randomFn = entity.RandomInt
	}
	return e
}

// NewGame implements lib.Engine
func (e *freeGame) NewGame(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		s.NumSpinLeft = entity.NumSpinByScatterSeq(int(s.GameConfig.NumScatterSeq))
		matrixSpecial := entity.NewJuicyMatrix()
		matrix := e.SpinMatrix(matrixSpecial, e.ratioWild)
		s.MatrixSpecial = &matrix
		s.ChipStat.Reset(s.CurrentSiXiangGame)
	}
	s.IsSpinChange = false
	return matchState, nil
}

// Process implements lib.Engine
func (e *freeGame) Process(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if s.NumSpinLeft <= 0 {
		return matchState, entity.ErrorSpinReachMax
	}
	s.IsSpinChange = true
	s.NumSpinLeft--
	var matrix entity.SlotMatrix
	for {
		matrix = e.SpinMatrix(*s.MatrixSpecial, ratioWild1_0)
		// không cho phép xuất hiện các loại bonus khác (Free tiếp, 6 giỏ, Scatter 345, hoặc JP)
		numScatterSeq := e.countScattersSequent(&matrix)
		if numScatterSeq >= 3 {
			continue
		}
		numFruitBasket := e.countFruitBasket(&matrix)
		if numFruitBasket >= 6 {
			continue
		}
		// numJpSymbol := 0
		// matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		// 	if entity.IsFruitJPSymbol(symbol) {
		// 		numJpSymbol++
		// 	}
		// })
		// if numJpSymbol > 0 {
		// 	continue
		// }
		break
	}
	s.MatrixSpecial = &matrix
	s.SetWildMatrix(e.WildMatrix(matrix))
	s.SetPaylines(e.Paylines(matrix))
	s.SpinList = make([]*pb.SpinSymbol, 0)
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		spinSymbol := &pb.SpinSymbol{
			Symbol: symbol,
			Index:  int32(idx),
			Row:    int32(row),
			Col:    int32(col),
		}
		if entity.IsFruitBasketSymbol(symbol) && !entity.IsFruitJPSymbol(symbol) {
			for {
				randSym := entity.ShuffleSlice(entity.JuicyFruitRainSybol[1:])[0]
				if randSym != pb.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_SPIN {
					spinSymbol.Symbol = randSym
					break
				}
			}
			val := entity.JuicyBasketSymbol[spinSymbol.Symbol]
			spinSymbol.Ratio = float32(e.Random(int(val.Value.Min), int(val.Value.Max)))
			spinSymbol.WinJp = entity.JuicySpinSymbolToJp(spinSymbol.Symbol)
		}
		s.SpinList = append(s.SpinList, spinSymbol)
	})
	return matchState, nil
}

// Random implements lib.Engine
func (e *freeGame) Random(min int, max int) int {
	return e.randomFn(min, max)
}

// Finish implements lib.Engine
func (e *freeGame) Finish(matchState interface{}) (interface{}, error) {
	s := matchState.(*entity.SlotsMatchState)
	if !s.IsSpinChange {
		return s.LastResult, nil
	}
	s.IsSpinChange = false
	lineWin := 0
	paylines := s.Paylines()
	for _, payline := range paylines {
		lineWin += int(payline.GetRate())
		payline.Chips = int64(payline.GetRate()) * s.Bet().Chips / 20
	}
	// s.RatioFruitBasket = 1
	// scatter x3 x4 x5 tính điểm tương ứng 3 4 5 x line bet
	if s.GameConfig.NumScatterSeq >= 3 {
		lineWin *= int(s.GameConfig.NumScatterSeq)
	}
	for _, spin := range s.SpinList {
		spin.WinAmount = int64(spin.Ratio) * s.Bet().Chips / 20
	}
	// s.RatioFruitBasket = e.transformNumScaterSeqToRationFruitBasket(s.GameConfig.NumScatterSeq)
	s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
		if entity.IsFruitBasketSymbol(symbol) && !entity.IsFruitJPSymbol(symbol) {
			val := entity.JuicyBasketSymbol[symbol]
			lineWin += int(float64(val.Value.Min) * float64(s.GameConfig.GetRatioBasket()))
		}
	})
	chipWin := int64(lineWin) * s.Bet().Chips / 20
	s.ChipStat.AddChipWin(s.CurrentSiXiangGame, chipWin)
	s.ChipStat.AddLineWin(s.CurrentSiXiangGame, int64(lineWin))
	s.NextSiXiangGame = e.GetNextSiXiangGame(s)
	slotDesk := &pb.SlotDesk{
		GameReward: &pb.GameReward{
			ChipsWin:            chipWin,
			TotalChipsWinByGame: s.ChipStat.ChipWin(s.CurrentSiXiangGame),
		},
		ChipsMcb:           s.Bet().Chips,
		Matrix:             s.MatrixSpecial.ToPbSlotMatrix(),
		Paylines:           paylines,
		CurrentSixiangGame: s.CurrentSiXiangGame,
		NextSixiangGame:    s.NextSiXiangGame,
		IsFinishGame:       true,
		NumSpinLeft:        int64(s.NumSpinLeft),
	}
	slotDesk.Matrix.SpinLists = s.SpinList
	return slotDesk, nil
}

func (e *freeGame) GetNextSiXiangGame(s *entity.SlotsMatchState) pb.SiXiangGame {
	if s.NumSpinLeft <= 0 {
		if s.NumFruitBasket >= 6 {
			return pb.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN
		}
		return pb.SiXiangGame_SI_XIANG_GAME_NORMAL
	}
	return pb.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME
}

func (e *freeGame) transformNumScaterSeqToRationFruitBasket(numScatterSeq int) int {
	switch numScatterSeq {
	case 4:
		return 2
	case 5:
		return 4
	}
	return 1
}
