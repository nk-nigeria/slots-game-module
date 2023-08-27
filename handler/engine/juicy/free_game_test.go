package juicy

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	pb "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func TestNewFreeGame(t *testing.T) {
	name := "TestNewFreeGame"
	t.Run(name, func(t *testing.T) {
		e := NewFreeGame(nil)
		assert.NotNil(t, e)
	})

}

func Test_freeGame_NewGame(t *testing.T) {
	name := "Test_freeGame_NewGame"
	t.Run(name, func(t *testing.T) {
		e := NewFreeGame(nil)
		engine := e.(*freeGame)
		s := entity.NewSlotsMathState(nil)
		s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME
		for i := 0; i <= 5; i++ {
			s.GameConfig.NumScatterSeq = int64(i)
			e.NewGame(s)
			ratioFruitBasket := 1
			var ratioWild ratioWild = ratioWild1_0
			gemSpin := 3
			// switch i {
			// case 3:
			// 	ratioFruitBasket = 1
			// 	ratioWild = ratioWild1_2
			// 	gemSpin = 6
			// case 4:
			// 	ratioFruitBasket = 2
			// 	ratioWild = ratioWild1_5
			// 	gemSpin = 9
			// case 5:
			// 	ratioFruitBasket = 4
			// 	ratioWild = ratioWild2_0
			// 	gemSpin = 15
			// }
			s.GameConfig = entity.GameConfigFreeGame(i)
			assert.Equal(t, float32(ratioFruitBasket), float32(s.GameConfig.RatioBasket))
			assert.Equal(t, int(gemSpin), s.NumSpinLeft)
			assert.Equal(t, ratioWild, engine.ratioWild)
			assert.Equal(t, int(entity.RowsJuicynMatrix*entity.ColsJuicyMatrix), int(len(s.MatrixSpecial.List)))
			// assert.Equal(t, int(0), s.ChipWinByGame[s.CurrentSiXiangGame])
			// assert.Equal(t, int(0), int(s.ChipStat.ChipWin(s.CurrentSiXiangGame)))
		}
	})
}

func Test_freeGame_Process(t *testing.T) {
	name := "Test_freeGame_Process"
	t.Run(name, func(t *testing.T) {
		e := NewFreeGame(nil)
		engine := e.(*freeGame)
		s := entity.NewSlotsMathState(nil)
		e.NewGame(s)
		// make sure num scatter < 3, num fruit basket < 6, no grand JP
		for i := 0; i < 10000; i++ {
			e.NewGame(s)
			_, err := e.Process(s)
			assert.NoError(t, err)
			assert.NotNil(t, s)
			assert.Equal(t, int(entity.RowsJuicynMatrix*entity.ColsJuicyMatrix), int(len(s.MatrixSpecial.List)))
			numScatterSeq := engine.countScattersSequent(s.MatrixSpecial)
			assert.Greater(t, int(3), int(numScatterSeq))
			numFruitBasket := engine.countFruitBasket(s.MatrixSpecial)
			assert.Greater(t, int(6), int(numFruitBasket))
			s.MatrixSpecial.ForEeach(func(idx, row, col int, symbol api.SiXiangSymbol) {
				assert.NotEqual(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FRUITBASKET_GRAND, symbol)
			})
		}
	})
}

func Test_freeGame_maxSpin_Process(t *testing.T) {
	name := "Test_freeGame_maxSpin_Process"
	t.Run(name, func(t *testing.T) {
		e := NewFreeGame(nil)
		s := entity.NewSlotsMathState(nil)
		e.NewGame(s)
		s.NumSpinLeft = 3
		for i := 0; i < 3; i++ {
			_, err := e.Process(s)
			assert.NoError(t, err)
		}
		_, err := e.Process(s)
		assert.Equal(t, entity.ErrorSpinReachMax, err)
	})

}

func Test_freeGame_Random(t *testing.T) {
	name := "Test_freeGame_Random"
	t.Run(name, func(t *testing.T) {
		e := NewFreeGame(nil)
		mapNum := make(map[int]bool)
		for i := 0; i < 1000; i++ {
			randNum := e.Random(0, 1000)
			mapNum[randNum] = true
		}
		assert.Less(t, int(100), int(len(mapNum)))
	})
}

func Test_freeGame_GetNextSiXiangGame(t *testing.T) {
	name := "Test_freeGame_GetNextSiXiangGame"
	t.Run(name, func(t *testing.T) {
		e := NewFreeGame(nil)
		engine := e.(*freeGame)
		s := entity.NewSlotsMathState(nil)
		e.NewGame(s)
		s.NumFruitBasket = 0
		s.NumSpinLeft = 3
		nextGame := engine.GetNextSiXiangGame(s)
		assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME, nextGame)
		s.NumSpinLeft = 0
		for i := 0; i < 6; i++ {
			s.NumFruitBasket = i
			nextGame = engine.GetNextSiXiangGame(s)
			assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, nextGame)
		}
		for i := 6; i < 8; i++ {
			s.NumFruitBasket = i
			nextGame = engine.GetNextSiXiangGame(s)
			assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_RAIN, nextGame)
		}
	})
}

func Test_freeGame_transformNumScaterSeqToRationFruitBasket(t *testing.T) {
	name := "Test_freeGame_transformNumScaterSeqToRationFruitBasket"
	t.Run(name, func(t *testing.T) {
		e := NewFreeGame(nil)
		enginre := e.(*freeGame)
		var mapRatioFruitBasket = map[int]int{0: 1, 1: 1, 2: 1, 3: 1, 4: 2, 5: 4}
		for k, v := range mapRatioFruitBasket {
			ratio := enginre.transformNumScaterSeqToRationFruitBasket(k)
			assert.Equal(t, int(v), int(ratio))
		}
	})
}

func Test_freeGame_only_payline_Finish(t *testing.T) {
	name := "Test_freeGame_only_payline_Finish"
	t.Run(name, func(t *testing.T) {
		e := NewFreeGame(nil)
		// engine := e.(*freeGame)
		// s := entity.NewSlotsMathState(nil)
		// e.NewGame(s)
		listNumScatterSeq := []int{1, 2, 3, 4, 5}
		for _, numScatterSeq := range listNumScatterSeq {
			s := entity.NewSlotsMathState(nil)
			chipMcb := 100
			s.SetBetInfo(&api.InfoBet{
				Chips: int64(chipMcb),
			})
			s.Matrix = entity.NewJuicyMatrix()
			s.NumFruitBasket = 0
			s.NumSpinLeft = 2
			s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME
			s.GameConfig.NumScatterSeq = int64(numScatterSeq)
			s.Matrix.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
				s.Matrix.List[idx] = api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED
			})
			for i := 0; i < numScatterSeq; i++ {
				s.Matrix.List[i] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
			}
			matrixSpecial := s.Matrix
			s.MatrixSpecial = &matrixSpecial
			// paylineSymbols := s.Matrix.ListFromIndexs(ids)
			lineWin := 100
			payline := &pb.Payline{
				Symbol:   api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED,
				NumOccur: 0,
				Id:       1,
				Rate:     float64(lineWin),
			}
			s.SetPaylines([]*pb.Payline{payline})
			s.IsSpinChange = true
			result, err := e.Finish(s)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			slotDesk := result.(*api.SlotDesk)
			assert.Equal(t, false, slotDesk.IsFinishGame)
			assert.Less(t, int(0), int(s.NumSpinLeft))
			assert.NotNil(t, slotDesk)
			assert.Equal(t, int(100), int(slotDesk.ChipsMcb))
			// ratioFruitBasket := engine.transformNumScaterSeqToRationFruitBasket(numScatterSeq)
			// t.Logf("scatter seq %d ration fruitbasket %d", numScatterSeq, ratioFruitBasket)
			// assert.Equal(t, int(ratioFruitBasket), int(s.RatioFruitBasket))
			chipWin := int(lineWin * numScatterSeq * int(slotDesk.ChipsMcb) / 20)
			if numScatterSeq < 3 {
				// assert.Equal(t, int(lineWin), int(s.LineWinByGame[s.CurrentSiXiangGame]))
				// assert.Equal(t, int(lineWin), int(s.ChipStat.LineWin(s.CurrentSiXiangGame)))
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, slotDesk.NextSixiangGame)
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, s.NextSiXiangGame)
				assert.Equal(t, int(lineWin*int(slotDesk.ChipsMcb)/20), int(slotDesk.GameReward.ChipsWin))

			} else {
				// assert.Equal(t, int(lineWin*numScatterSeq), int(s.LineWinByGame[s.CurrentSiXiangGame]))
				// assert.Equal(t, int(lineWin*numScatterSeq), int(s.ChipStat.LineWin(s.CurrentSiXiangGame)))
				assert.Equal(t, int(chipWin), int(slotDesk.GameReward.ChipsWin))
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, slotDesk.NextSixiangGame)
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, s.NextSiXiangGame)
			}
			assert.Equal(t, 0, s.NumFruitBasket)
		}

	})
}

func Test_freeGame_only_payline_Finish_Flow(t *testing.T) {
	name := "Test_freeGame_only_payline_Finish"
	t.Run(name, func(t *testing.T) {
		e := NewFreeGame(nil)
		s := entity.NewSlotsMathState(nil)
		s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME
		s.NextSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME
		s.Bet().Chips = 100
		e.NewGame(s)
		for i := 0; i < 1000; i++ {
			e.Process(s)
			e.Finish(s)
			if s.NumSpinLeft <= 0 {
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, s.NextSiXiangGame)
				s = entity.NewSlotsMathState(nil)
				s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME
				s.Bet().Chips = 100
				e.NewGame(s)
			}
		}
	})
}
