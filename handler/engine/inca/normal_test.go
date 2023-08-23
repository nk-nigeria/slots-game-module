package inca

import (
	"testing"

	"github.com/ciaolink-game-platform/cgb-slots-game-module/entity"
	api "github.com/ciaolink-game-platform/cgp-common/proto"
	"github.com/stretchr/testify/assert"
)

func Test_normal_NewGame(t *testing.T) {
	name := "Test_normal_NewGame"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(nil)
		s := entity.NewSlotsMathState(nil)
		e.NewGame(s)
		assert.NotEqual(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED, s.Matrix.List[0])
	})

}

func Test_normal_Process(t *testing.T) {
	name := "Test_normal_Process"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(nil)
		s := entity.NewSlotsMathState(nil)
		e.NewGame(s)
		e.Process(s)
		assert.NotEqual(t, api.SiXiangSymbol_SI_XIANG_SYMBOL_UNSPECIFIED, s.Matrix.List[0])
		assert.Equal(t, entity.ColsIncaMatrix*entity.RowsIncaMatrix, len(s.WildMatrix.List))
	})
}

func Test_normal_Finish(t *testing.T) {
	name := "Test_normal_Finish"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(nil)
		s := entity.NewSlotsMathState(nil)
		s.Bet().Chips = 100
		e.NewGame(s)
		e.Process(s)
		result, err := e.Finish(s)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		slotDesk := result.(*api.SlotDesk)
		assert.Equal(t, entity.ColsIncaMatrix*entity.RowsIncaMatrix, len(slotDesk.Matrix.Lists))
		assert.Equal(t, true, slotDesk.IsFinishGame)
		assert.NotNil(t, slotDesk.GameReward)
		assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, slotDesk.CurrentSixiangGame)
		assert.Equal(t, api.SiXiangGame_SI_XIANG_GAME_NORMAL, slotDesk.NextSixiangGame)
		assert.Equal(t, s.Bet().Chips, slotDesk.ChipsMcb)
	})
}
