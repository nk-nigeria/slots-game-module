package juicy

import (
	"testing"

	"github.com/nakamaFramework/cgb-slots-game-module/entity"
	api "github.com/nakamaFramework/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_fruitBasket_NewGame(t *testing.T) {
	name := "Test_fruitBasket_NewGame"
	t.Run(name, func(t *testing.T) {
		e := NewFruitBaseket()
		s := entity.NewSlotsMathState(nil)
		e.NewGame(s)
		assert.NotNil(t, s)
		assert.Equal(t, int(2), len(s.MatrixSpecial.List))
	})
}

func Test_fruitBasket_Finish(t *testing.T) {
	name := "Test_fruitBasket_Finish"
	t.Run(name, func(t *testing.T) {
		e := NewFruitBaseket()

		for i := 0; i < 2; i++ {
			s := entity.NewSlotsMathState(nil)
			s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET
			e.NewGame(s)
			selectGame := s.MatrixSpecial.Flip(0)
			result, err := e.Finish(s)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			slotDesk := result.(*api.SlotDesk)
			assert.Equal(t, 2, len(slotDesk.Matrix.Lists))
			if selectGame == api.SiXiangSymbol_SI_XIANG_SYMBOL_JUICE_FUIT_SELECT_FREE_GAME {
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_JUICE_FREE_GAME, slotDesk.NextSixiangGame)
			} else {
				assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_JUICE_FRUIT_BASKET, s.CurrentSiXiangGame)
			}
			mapSymbol := make(map[api.SiXiangSymbol]bool, 0)
			for _, symbol := range slotDesk.Matrix.Lists {
				_, exist := mapSymbol[symbol]
				t.Logf("symbol %s", symbol.String())
				assert.Equal(t, false, exist)
				mapSymbol[symbol] = true
			}
			assert.Equal(t, true, slotDesk.IsFinishGame)
		}
	})
}
