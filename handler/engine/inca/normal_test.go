package inca

import (
	"testing"

	"github.com/nk-nigeria/slots-game-module/entity"
	api "github.com/nk-nigeria/cgp-common/proto"
	pb "github.com/nk-nigeria/cgp-common/proto"
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

func Test_normal_Finish_Win_Chips(t *testing.T) {
	name := "Test_normal_Finish_Win_Chips"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(nil)
		s := entity.NewSlotsMathState(nil)
		s.Bet().Chips = 1000
		e.NewGame(s)
		e.Process(s)
		s.Matrix.List[0] = api.SiXiangSymbol_SI_XIANG_SYMBOL_EAGLE_GARUDA
		s.Matrix.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUN
		s.Matrix.List[2] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_SPADES
		s.Matrix.List[3] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
		s.Matrix.List[4] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q

		s.Matrix.List[5] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUN
		s.Matrix.List[6] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_DIAMONDS
		s.Matrix.List[7] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		s.Matrix.List[8] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[9] = api.SiXiangSymbol_SI_XIANG_SYMBOL_EAGLE_GARUDA

		s.Matrix.List[10] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_DIAMONDS
		s.Matrix.List[11] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER
		s.Matrix.List[12] = api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD
		s.Matrix.List[13] = api.SiXiangSymbol_SI_XIANG_SYMBOL_EAGLE_GARUDA
		s.Matrix.List[14] = api.SiXiangSymbol_SI_XIANG_SYMBOL_Q
		s.WildMatrix = s.Matrix
		engine := e.(*normal)
		s.SetPaylines(engine.Paylines(s.Matrix))
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
		assert.Equal(t, int64(8700), slotDesk.GameReward.TotalChipsWinByGame)
	})
}

func Test_normal_SpinMatrix_ScatterOccur(t *testing.T) {
	name := "Test_normal_SpinMatrix_ScatterOccur"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(nil)
		s := entity.NewSlotsMathState(nil)
		s.GameConfig = &entity.GameConfig{
			GameConfig: &pb.GameConfig{},
		}
		s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_NORMAL
		s.Bet().Chips = 1000
		e.NewGame(s)

		for i := 0; i < 10000; i++ {
			e.Process(s)
			result, err := e.Finish(s)
			assert.NoError(t, err)
			assert.NotNil(t, result)
			slotDesk := result.(*api.SlotDesk)
			sm := &entity.SlotMatrix{}
			sm.FromPbMatrix(slotDesk.Matrix)
			// sm.ForEachCol(func(col int, symbols []api.SiXiangSymbol) {
			// 	if col == entity.Col_1 || col == entity.Col_5 {
			// 		for row, sym := range symbols {
			// 			if sym == api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
			// 				t.Logf("row %d col %d", row, col)
			// 			}
			// 			assert.NotEqual(t, pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER, sym)
			// 		}
			// 	}
			// })
			sm.ForEeach(func(idx, row, col int, symbol pb.SiXiangSymbol) {
				if col == entity.Col_1 || col == entity.Col_5 {
					if symbol == api.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER {
						t.Logf("row %d col %d", row, col)
					}
					assert.NotEqual(t, pb.SiXiangSymbol_SI_XIANG_SYMBOL_SCATTER, symbol)

				}
			})
		}
	})
}

func Test_normal_Payline(t *testing.T) {
	name := "Test_normal_Payline"
	t.Run(name, func(t *testing.T) {
		e := NewNormal(nil)
		s := entity.NewSlotsMathState(nil)
		s.GameConfig = &entity.GameConfig{
			GameConfig: &pb.GameConfig{},
		}
		s.CurrentSiXiangGame = api.SiXiangGame_SI_XIANG_GAME_NORMAL
		s.Bet().Chips = 1000
		e.NewGame(s)
		e.Process(s)
		// s.Matrix.List[0] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		// s.Matrix.List[1] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		// s.Matrix.List[2] = api.SiXiangSymbol_SI_XIANG_SYMBOL_EAGLE_GARUDA
		// s.Matrix.List[3] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_HEARTS
		// s.Matrix.List[4] = api.SiXiangSymbol_SI_XIANG_SYMBOL_EAGLE_GARUDA

		// s.Matrix.List[5] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUN
		// s.Matrix.List[6] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUN
		// s.Matrix.List[7] = api.SiXiangSymbol_SI_XIANG_SYMBOL_WILD
		// s.Matrix.List[8] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUN
		// s.Matrix.List[9] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_HEARTS

		// s.Matrix.List[10] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		// s.Matrix.List[11] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_HEARTS
		// s.Matrix.List[12] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUIT_SPADES
		// s.Matrix.List[13] = api.SiXiangSymbol_SI_XIANG_SYMBOL_SUN
		// s.Matrix.List[14] = api.SiXiangSymbol_SI_XIANG_SYMBOL_K
		// engine := e.(*normal)
		// paylines := engine.Paylines(s.Matrix)
		// assert.NotNil(t, paylines)
		// assert.Equal(t, 2, len(paylines))
		finish, _ := e.Finish(s)
		assert.NotNil(t, finish)
	})
}
