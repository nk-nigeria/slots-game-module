package juicy

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
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
			s.NumScatterSeq = i
			e.NewGame(s)
			ratioFruitBasket := 1
			var ratioWild ratioWild = ratioWild1_0
			gemSpin := 3
			switch i {
			case 3:
				ratioFruitBasket = 1
				ratioWild = ratioWild1_2
				gemSpin = 6
			case 4:
				ratioFruitBasket = 2
				ratioWild = ratioWild1_5
				gemSpin = 9
			case 5:
				ratioFruitBasket = 4
				ratioWild = ratioWild2_0
				gemSpin = 15
			}
			assert.Equal(t, int(ratioFruitBasket), s.RatioFruitBasket)
			assert.Equal(t, int(gemSpin), s.GemSpin)
			assert.Equal(t, ratioWild, engine.ratioWild)
			assert.Equal(t, int(entity.RowsJuicynMatrix*entity.ColsJuicyMatrix), int(len(s.MatrixSpecial.List)))
			assert.Equal(t, int(0), s.ChipWinByGame[s.CurrentSiXiangGame])
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
		s.GemSpin = 3
		for i := 0; i < 3; i++ {
			_, err := e.Process(s)
			assert.NoError(t, err)
		}
		_, err := e.Process(s)
		assert.Equal(t, ErrorSpinReadMax, err)
	})

}
